package emu

import (
	//NODISPLAY//"github.com/strickyak/gomar/display"
	. "github.com/strickyak/gomar/gu"
	"github.com/strickyak/gomar/sym"

	//"bufio"
	"bytes"
	"flag"
	//"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	//"sort"
	//"strconv"
	//"strings"
)

var FlagTraceVerbosity = flag.String("vv", "", "Trace verbosity chars") // Trace Verbosity
var FlagTraceAfter = flag.Int64("t", MaxInt64, "Tracing starts after this many steps")

func Main() {
	CompileWatches()
	SetVerbosityBits(*FlagInitialVerbosity)
	InitHardware()
	keystrokes := make(chan byte, 0)
	go InputRoutine(keystrokes)

	//NODISPLAY// CocodChan := make(chan *display.CocoDisplayParams, 50)
	//NODISPLAY// Disp = display.NewDisplay(mem[:], 80, 25, CocodChan, keystrokes, &sam, PeekBWithInt)

	Logd("(begin roms)")
	if *FlagBootImageFilename != "" {
		{
			// TODO: this code is duplicated????? Search for FlagBootImageFilename and find the other one.
			// Open disk image.
			fd, err := os.OpenFile(*FlagDiskImageFilename, os.O_RDWR, 0644)
			if err != nil {
				log.Fatalf("Cannot open disk image: %q: %v", *FlagBootImageFilename, err)
			}
			disk_fd = fd
		}

		{
			// Read disk_sector_0.
			n, err := disk_fd.Read(disk_sector_0[:])
			if err != nil {
				log.Panicf("Bad disk sector read: err=%v", err)
			}
			if n != 256 {
				log.Panicf("Short disk sector read: n=%d", n)
			}

			disk_dd_fmt = disk_sector_0[16]

			tracks_per_sector := int(disk_sector_0[17])*256 + int(disk_sector_0[18])
			if tracks_per_sector != 18 {
				log.Panicf("Not 18 sectors per track: %d.", tracks_per_sector)
			}
		}
	}

	if *FlagRomA000Filename != "" {
		rom, err := ioutil.ReadFile(*FlagRomA000Filename)
		if err != nil {
			log.Fatalf("Cannot read rom image: %q: %v", *FlagRomA000Filename, err)
		}
		Logd("Loading Rom %q at %04x", *FlagRomA000Filename, 0xA000)
		LoadRom(0xA000, rom)
		//for i := Word(0); i < 16; i++ {
		//PokeB(0xFFF0+i, PeekB(0xbff0+i)) // Install interrupt vectors.
		//}
		usedRom = true
	}

	if *FlagRom8000Filename != "" {
		rom, err := ioutil.ReadFile(*FlagRom8000Filename)
		if err != nil {
			log.Fatalf("Cannot read rom image: %q: %v", *FlagRom8000Filename, err)
		}
		Logd("Loading Rom %q at %04x", *FlagRom8000Filename, 0x8000)
		LoadRom(0x8000, rom)
		usedRom = true
	}

	if *FlagCartFilename != "" {
		rom, err := ioutil.ReadFile(*FlagCartFilename)
		if err != nil {
			log.Fatalf("Cannot read rom image: %q: %v", *FlagCartFilename, err)
		}
		Logd("Loading Cart %q", *FlagCartFilename)
		LoadCart(rom)
	}
	Logd("(end roms)")

	if *FlagLoadmFilename != "" {
		loadm, err := ioutil.ReadFile(*FlagLoadmFilename)
		if err != nil {
			log.Fatalf("Cannot read loadm image: %q: %v", *FlagLoadmFilename, err)
		}
		pcreg = Loadm(loadm)
		sreg = 0x8000
	}

	if *FlagBootImageFilename != "" {
		// Loading a binary file skipping the first 256 bytes of RAM
		// and starting pcreg of 0x100 was a convention from the sbc09.c code.
		// This is probably not the right thing for a coco emulator,
		// but I started using it in the early days, and haven't switched away yet.
		boot, err := ioutil.ReadFile(*FlagBootImageFilename)
		if err != nil {
			log.Fatalf("Cannot read boot image: %q: %v", *FlagDiskImageFilename, err)
		}
		L("boot mem size: %x", len(boot))
		for i, b := range boot {
			PokeB(Word(i+0x100), b)
		}
		pcreg = 0x100
		DumpAllMemory()
	} else if *FlagKernelFilename != "" {
		kernel, err := ioutil.ReadFile(*FlagKernelFilename)
		if err != nil {
			log.Fatalf("Cannot read kernel image: %q: %v", *FlagKernelFilename, err)
		}
		if kernel[0] != 'O' || kernel[1] != 'S' {
			log.Fatalf("--kernel does not begin with OS")
		}
		L("kernel mem size: %x", len(kernel))
		for i, b := range kernel {
			PokeB(Word(i+0x2600), b)
		}
		PutW(0xFFF2, 0xFEEE) // SWI3
		PutW(0xFFF4, 0xFEF1) // SWI2
		PutW(0xFFFA, 0xFEFA) // SWI
		PutW(0xFFFC, 0xFEFD) // NMI
		PutW(0xFFF8, 0xFEF7) // IRQ
		PutW(0xFFF6, 0xFEF4) // FIRQ
		pcreg = 0x2602
		DumpAllMemory()
	}

	if *FlagUserResetVector {
		pcreg = PeekW(0xFFFE)
	}

	if usedRom {
		enableRom = true
		pcreg = PeekW(0xFFFE)
		pcreg = HiLo(internalRom[0x7Ffe], internalRom[0x7Fff])
		pcreg = HiLo(internalRom[0x3Ffe], internalRom[0x3Fff])
	}
	if pcreg == 0 {
		log.Fatalf("Before run, pcreg is still 0")
	}

	sreg = 0
	dpreg = 0
	iflag = 0

	Dis_len(0)

	defer func() {
		Finish()
	}()

	max := int64(MaxInt64)
	if *FlagMaxSteps > 0 {
		max = *FlagMaxSteps
	}
	stepsUntilTimer := *FlagClock
	early := true

	Cycles = int64(0)
	for Cycles < max {
		if early {
			early = EarlyAction()
		}

		pcreg_prev = pcreg

		if stepsUntilTimer == 0 {
			DoMemoryDumps()
			FireTimerInterrupt()
			stepsUntilTimer = *FlagClock
		} else {
			stepsUntilTimer--
		}

		if Waiting {
			continue
		}

		if (irqs_pending) != 0 {
			if (irqs_pending & NMI_PENDING) != 0 {
				nmi()
				continue
			}
			if (irqs_pending&IRQ_PENDING) != 0 && (ccreg&CC_INHIBIT_IRQ) == 0 {

				irq(keystrokes)
				continue
			}
		}

		ireg = B(pcreg)
		if pcreg == Word(*FlagTriggerPc) && ireg == byte(*FlagTriggerOp) {
			*FlagTraceAfter = 1
			SetVerbosityBits(*FlagTraceVerbosity)
			log.Printf("TRIGGERED")
			// MemoryModules()
			// DoDumpAllMemory()
		}
		pcreg++

		// Process instruction
		instructionTable[ireg]()

		if true || Cycles >= *FlagTraceAfter {
			Trace()
		}

		if paranoid && !early {
			ParanoidAsserts()
		}
	} /* next step */

	if max > 0 {
		if Cycles >= max {
			log.Fatalf("MAX CYCLES REACHED: %d.", Cycles)
		}
	}
}

func ParanoidAsserts() {
	if pcreg < 0x005E /* D.BtDbg */ {
		log.Panicf("PC in page 0: 0x%x", pcreg)
	}
	if pcreg >= 0xFF00 {
		log.Panicf("PC in page FF: 0x%x", pcreg)
	}
	if pcreg >= 0x0200 && pcreg < 0x04FF {
		log.Panicf("PC in sys data: 0x%x", pcreg)
	}
	if Level == 1 {
		if sreg < 256 {
			log.Panicf("S in page 0: 0x%x", sreg)
		}
	}
	if sreg >= 0xFF00 {
		log.Panicf("S in page FF: 0x%x", sreg)
	}
	if sreg >= 0x0140 && sreg < 0x0400 {
		log.Panicf("S in sys data: 0x%x", sreg)
	}
}

func interrupt(vector_addr Word) {
	PushWord(pcreg)
	if vector_addr == VECTOR_FIRQ {
		// Fast IRQ.
		ccreg &= ^byte(CC_ENTIRE)
	} else {
		// Other IRQs.
		PushWord(ureg)
		PushWord(yreg)
		PushWord(xreg)
		PushByte(dpreg)
		PushWord(dreg)
	}
	PushByte(ccreg)
	if vector_addr == VECTOR_FIRQ {
		// Fast IRQ.
		ccreg &= ^byte(CC_ENTIRE)
	} else {
		// Other IRQs.
		ccreg |= byte(CC_ENTIRE)
	}
	// All IRQs.
	ccreg |= (CC_INHIBIT_FIRQ | CC_INHIBIT_IRQ)
	pcreg = W(vector_addr)
}

func irq(keystrokes <-chan byte) {
	kbd_cycle++
	L("INTERRUPTING with IRQ (kbd_cycle = %d)", kbd_cycle)
	Assert(0 == (ccreg&CC_INHIBIT_IRQ), ccreg)

	if (kbd_cycle & 1) == 0 {
		ch := inkey(keystrokes)
		kbd_ch = ch
		if kbd_ch != 0 {
			log.Printf("key/irq $%x=%d.", kbd_ch, kbd_ch)
		}

		L("getchar -> ch %x %q kbd_ch %x %q (kbd_cycle = %d)\n", ch, string(rune((ch))), kbd_ch, string(rune((kbd_ch))), kbd_cycle)
	} else {
		kbd_ch = 0
	}
	L("irq -> kbd_ch %x %q (kbd_cycle = %d)\n", kbd_ch, string(rune(kbd_ch)), kbd_cycle)

	interrupt(VECTOR_IRQ)
	irqs_pending &^= IRQ_PENDING
}

var swi_name = []string{"swi", "swi2", "swi3"}

func swi() {
	Dis_inst(swi_name[iflag], "", 5)
	Dis_len(3 /* Often an extra byte after the SWI opcode */)

	ccregOrig, sregOrig := ccreg, sreg

	ccreg |= 0x80
	PushWord(pcreg)
	PushWord(ureg)
	PushWord(yreg)
	PushWord(xreg)
	PushByte(dpreg)
	PushWord(dreg)
	PushByte(ccreg)

	var handler Word
	switch iflag {
	case 0: /* SWI */
		L("SWI")
		if false { // Nothing should still be using this.
			// Intercept HyperOp on SWI
			op := PeekB(pcreg)
			pcreg++
			L("HyperOp %d.", op)
			HyperOp(op)
			ccreg, sreg = ccregOrig, sregOrig
		} else {
			// Normal SWI.
			ccreg |= 0xd0
			handler = W(0xfffa)
		}
		return
	case 1: /* SWI2 */
		describe, returns := DecodeOs9Opcode(B(pcreg))
		proc := W0(sym.D_Proc)
		pid := B0(proc + sym.P_ID)
		pmodul := W0(proc + sym.P_PModul)
		moduleName := Os9String(pmodul + W(pmodul+4))

		luser := 0
		if Level == 1 && dpreg != 0 {
			luser = 1
		}
		if Level == 2 && MmuTask != 0 {
			luser = 1
		}

		L("{proc=%x%q} OS9KERNEL%d: %s", pid, moduleName, luser, describe)
		L("\tregs: %s", Regs())
		L("\t%s", ExplainMMU())

		stack := MapAddr(sreg, true /*quiet*/)
		if returns {
			Os9Description[stack] = describe
		} else {
			Os9Description[stack] = ""
		}

		handler = W(0xfff4)
	case 2: /* SWI3 */
		handler = W(0xfff2)
	default:
		log.Panicf("bad swi iflag=: %d", iflag)
	}

	if paranoid {
		if handler < 256 {
			log.Panicf("FATAL: Attempted SWI%d with small handler: 0x%04x", handler)
		}
		if handler >= 0xFF00 {
			log.Panicf("FATAL: Attempted SWI%d with large handler: 0x%04x", handler)
		}
	}

	syscall := B(pcreg)
	handled := false

	if hyp && iflag == 1 {
		handled = Os9HypervisorCall(syscall)
	}

	if !handled {
		pcreg = handler
	}
}

func rti() {
	if Cycles >= *FlagTraceAfter {
		DoDumpSysMap()
	}

	stack := MapAddr(sreg, true /*quiet*/)
	describe := Os9Description[stack]

	if *FlagTraceOnOS9 != "" && describe != "" {
		if RegexpTraceOnOS9 == nil {
			RegexpTraceOnOS9 = regexp.MustCompile(*FlagTraceOnOS9)
		}
		if RegexpTraceOnOS9.MatchString(describe) {
			*FlagTraceAfter = 1
		}
	}

	entire := ccreg & CC_ENTIRE
	if entire == 0 {
		Dis_inst("rti", "", 6)
	} else {
		Dis_inst("rti", "", 15)
	}
	Dis_len(1)
	PullByte(&ccreg)
	if entire != 0 {
		PullWord(&dreg)
		PullByte(&dpreg)
		PullWord(&xreg)
		PullWord(&yreg)
		PullWord(&ureg)
	}
	PullWord(&pcreg)

	back3 := B(pcreg - 3)
	back2 := B(pcreg - 2)
	back1 := B(pcreg - 1)
	if back3 == 0x10 && back2 == 0x3f && describe != "" {
		if (ccreg & 1 /* carry bit indicates error */) != 0 {
			errcode := GetBReg()

			luser := 0
			if Level == 1 && dpreg != 0 {
				luser = 1
			}
			if Level == 2 && MmuTask != 0 {
				luser = 1
			}

			/*
				PrettyDumpHex64(0, 0xFFFF)
			*/

			L("RETURN ERROR: $%x(%v): OS9KERNEL%d %s #%d", errcode, DecodeOs9Error(errcode), luser, describe, Cycles)
			L("\tregs: %s  #%d", Regs(), Cycles)
			L("\t%s", ExplainMMU())
			DoDumpAllMemory() // yak
		} else {
			switch back1 {
			case 0x82, 0x83, 0x84: // I$Dup, I$Create, I$Open
				describe += F(" -> path $%x", GetAReg())
			case 0x28: // F$SRqMem
				describe += F(" -> size $%x addr $%04x", dreg, ureg)
			case 0x30:
				describe += F(" -> base $%x blocknum $%x addr $%x", xreg, GetAReg(), yreg)
			case 0x00:
				describe += F(" -> addr $%x entry $%x", ureg, yreg)
			}

			luser := 0
			if Level == 1 && dpreg != 0 {
				luser = 1
			}
			if Level == 2 && MmuTask != 0 {
				luser = 1
			}

			L("RETURN OKAY: OS9KERNEL%d %s #%d", luser, describe, Cycles)
			L("\tregs: %s  #%d", Regs(), Cycles)
			L("\t%s", ExplainMMU())
			DoDumpAllMemory() // yak

			if back1 == 0x8B {
				var buf bytes.Buffer
				for i := Word(0); i < yreg; i++ {
					buf.WriteRune(rune(PeekB(xreg + i)))
				}
				L("ReadLn returns: [$%x] %q", buf.Len(), buf.String())
			}
		}

		// Os9Description[stack] = "" // Clear description
		delete(Os9Description, stack)

	}
}
