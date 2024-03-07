package emu

import (
	. "github.com/strickyak/gomar/gu"
	"github.com/strickyak/gomar/sym"

	//"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"regexp"
	"sync/atomic"
	"syscall"
	//"sort"
	//"strconv"
	//"strings"
)

var FlagTraceVerbosity = flag.String("vv", "", "Trace verbosity chars") // Trace Verbosity
var FlagTraceAfter = flag.Int64("t", MaxInt64, "Tracing starts after this many steps")

var FlagVdgRate = flag.Int("vdg_rate", 10003, "how often to print text screen")

// doubles in Os9 //  var FlagKbdRate = flag.Int("kbd_rate", 100031, "how often to frob keyboard")
var FlagKbdRate = flag.Int("kbd_rate", 50031, "how often to frob keyboard")

func Main() {
	LoadRomListings()
	InitExpectations()
	CompileWatches()
	SetVerbosityBits(*FlagInitialVerbosity)
	InitHardware()

	//NODISPLAY// CocodChan := make(chan *display.CocoDisplayParams, 50)
	//NODISPLAY// Disp = display.NewDisplay(mem[:], 80, 25, CocodChan, keystrokes, &sam, PeekBWithInt)

	// LoadBootImage()

	Logd("(begin roms)")

	if *FlagDiskImageFilename != "" {
		{
			// TODO: this code is duplicated????? Search for FlagBootImageFilename and find the other one.
			// Open disk image.
			fd, err := os.OpenFile(*FlagDiskImageFilename, os.O_RDWR, 0644)
			if err != nil {
				log.Fatalf("Cannot open disk image: %q: %v", *FlagDiskImageFilename, err)
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

			/*
				tracks_per_sector := int(disk_sector_0[17])*256 + int(disk_sector_0[18])
				if tracks_per_sector != 18 {
					log.Panicf("Not 18 sectors per track: %d.", tracks_per_sector)
				}
			*/
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

	if *FlagUserResetVector { // TODO ???
		pcreg = PeekW(0xFFFE)
	}

	if usedRom {
		enableRom = true
		pcreg = PeekW(0xFFFE)
		pcreg = HiLo(internalRom[0x7Ffe], internalRom[0x7Fff])
		pcreg = HiLo(internalRom[0x3Ffe], internalRom[0x3Fff])
	}
	if pcreg == 0 {
		pcreg = W(0xFFFE)
		log.Printf("Using reset vector for pcreg: $%04x", pcreg)
	}

	sreg = 0x8000
	dpreg = 0
	iflag = 0

	Dis_len(0)

	displayCount := *FlagVdgRate

	kbdCount := *FlagKbdRate
	keystrokes := make(chan byte, 0)
	go InputRoutine(keystrokes)

	defer func() {
		display.Tick(0)
		Finish()
	}()

	max := int64(MaxInt64)
	if *FlagMaxSteps > 0 {
		max = *FlagMaxSteps
	}
	// stepsUntilTimer := *FlagClock
	early := true

	////////////////////////////////////////
	var haltDumpAndExit int32
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGQUIT)
	go func() {
		_ = <-sigChan
		atomic.StoreInt32(&haltDumpAndExit, 1)
	}()
	////////////////////////////////////////

	Cycles = int64(0)
	for Cycles < max {
		L("[[[ fnord #%d ]]]", Cycles)
		if atomic.LoadInt32(&haltDumpAndExit) > 0 {
			V['d'] = true
			V['p'] = true
			Logd("haltDumpAndExit ...")
			DoDumpAllMemoryPhys()
			JustDoDumpAllMemory()
			Logd("... haltDumpAndExit.")
			fmt.Printf("\n... haltDumpAndExit.\n")
			os.Exit(99)
		}
		if early {
			early = EarlyAction()
		}

		pcreg_prev = pcreg

		{
			kbdCount--
			if kbdCount <= 0 {
				kbdService(keystrokes)
				kbdCount = *FlagKbdRate
			}

			displayCount--
			if displayCount <= 0 {
				display.Tick(Cycles)
				displayCount = *FlagVdgRate
			}

			if Pia0FrameSyncInterruptEnable {
				if Cycles > frameCycles {
					framePending = true
					incr := FastCyclesPerVertical
					frameCycles = Cycles + int64(incr)
				}
			}
			if Pia0HorzSyncInterruptEnable {
				if Cycles > horzCycles {
					horzPending = true
					incr := FastCyclesPerHorizontal
					horzCycles = Cycles + int64(incr)
				}
			}

			if GimeVirtSyncInterruptEnable {
				if Cycles > gimeVirtCycles {
					gimeVirtPending = true
					incr := FastCyclesPerVertical
					gimeVirtCycles = Cycles + int64(incr)
				}
			}
			if GimeHorzSyncInterruptEnable {
				if Cycles > gimeHorzCycles {
					gimeHorzPending = true
					incr := FastCyclesPerHorizontal
					gimeHorzCycles = Cycles + int64(incr)
				}
			}
		}

		if nmiPending {
			nmiPending = false
			nmi()
			continue
		}

		if (gimeVirtPending || gimeHorzPending || framePending || horzPending) && (ccreg&CC_INHIBIT_IRQ) == 0 {
			if gimeVirtPending {
				L("interrupting due to gimeVirtPending...")
			}
			if gimeHorzPending {
				L("interrupting due to gimeHorzPending...")
			}
			if framePending {
				L("interrupting due to framePending...")
			}
			if horzPending {
				L("interrupting due to horzPending...")
			}
			for p := 0xFFF0; p < 0xFFFF; p += 2 {
				L("  [%04x]  peek=%04x  int=%02x%02x  ext=%02x%02x", p, PeekW(Word(p)), internalRom[0x3FFF&p], internalRom[0x3FFF&(p+1)], cartRom[0x3FFF&p], cartRom[0x3FFF&(p+1)])
			}
			Waiting = false
			irq()
			continue
		}
		if Waiting {
			Cycles += 10 // move along
			continue
		}

		ireg = B(pcreg)
		L("[[[ fnord #%d: %02x at %04x ]]]", Cycles, ireg, pcreg)
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

		if Cycles >= *FlagTraceAfter {
			Trace()
		}

		if paranoid && !early {
			ParanoidAsserts()
		}

	} /* next step */

	if Expectations != nil {
		log.Fatalf("\n===@=== UNMET EXPECTATIONS: %#v\n", Expectations)
	}

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
		ccreg &^= byte(CC_ENTIRE)
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
		ccreg &^= byte(CC_ENTIRE)
	} else {
		// Other IRQs.
		ccreg |= byte(CC_ENTIRE)
	}
	// All IRQs.
	ccreg |= (CC_INHIBIT_FIRQ | CC_INHIBIT_IRQ)
	// pcreg = W(vector_addr)
	pcreg = HiLo(internalRom[0x3fff&vector_addr], internalRom[0x3fff&(vector_addr+1)])
	//panic("---------------interrupt--------------")
}

func kbdService(keystrokes <-chan byte) {
	kbd_cycle++

	if (kbd_cycle & 1) == 0 {
		// On Odd cycles, do a keystroke.
		ch := inkey(keystrokes)
		kbd_ch = ch

		L("kbdService: getchar -> ch %x %q kbd_ch %x %q (kbd_cycle = %d)\n", ch, string(rune((ch))), kbd_ch, string(rune((kbd_ch))), kbd_cycle)
	} else {
		// On Even cycles, release it.
		kbd_ch = 0
	}
	L("kbdService: irq -> kbd_ch %x %q (kbd_cycle = %d)\n", kbd_ch, string(rune(kbd_ch)), kbd_cycle)
}

func irq() {
	Assert(0 == (ccreg&CC_INHIBIT_IRQ), ccreg)

	interrupt(VECTOR_IRQ)
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

		luser := 0
		if Level == 1 && dpreg != 0 {
			luser = 1
		}
		if Level == 2 && MmuTask != 0 {
			luser = 1
		}
		moduleName := ""
		if luser == 1 {
			pmodul := W0(proc + sym.P_PModul)
			moduleName = Os9String(pmodul + W(pmodul+4))
		}

		if moduleName != "" {
			L("{proc=%x%q} OS9KERNEL%d: %s", pid, moduleName, luser, describe)
		} else {
			L("{proc=%x} OS9KERNEL%d: %s", pid, luser, describe)
		}

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
	//if Cycles >= *FlagTraceAfter {
	DoDumpProcsAndPathsPrime()
	DoDumpSysMap()
	//}

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
	if V['M'] {
		DoDumpAllMemory() // yak
	}
}
