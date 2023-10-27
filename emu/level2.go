//go:build level2

package emu

import (
	"bytes"
	"log"

	. "github.com/strickyak/gomar/gu"
	"github.com/strickyak/gomar/sym"
)

const Level = 2

// While booting OS9 Level2, the screen seems to be doubleByte
// at 07c000 to 07d000.  Second line begins at 07c0a0,
// that is 160 bytes from start, or 80 doubleBytes.
// 4096 div 160 is 25.6 lines.

const P_Path = sym.P_Path // vs P_PATH in level 1

func VerboseValidateModuleSyscall() string {
	mapping := GetMapping(dreg)
	hdr := PeekWWithMapping(xreg, mapping)
	mod := "-"
	if hdr == 0x87CD {
		nameOffset := PeekWWithMapping(xreg+4, mapping)
		mod = Os9StringWithMapping(xreg+nameOffset, mapping)
	}
	p := F("addr=%04x=%q map=%x", xreg, mod, mapping)

	{
		temp := V['p'] // TODO
		V['p'] = true
		DoDumpAllMemoryPhys()
		V['p'] = temp
	}
	return p
}
func DoDumpSysMap() {
	L("SMAP")
	begin := SysMemW(sym.D_SysMem)
	end := begin + 256
	for i := begin; i < end; i += 16 {
		var bb bytes.Buffer
	J:
		for j := Word(0); j < 32; j++ {
			if j == 15 {
				bb.WriteByte(' ')
			}
			bit := byte(0x80)
			for k := byte(0); k < 8; k++ {
				x := SysMemB(i + j)
				if (x & bit) != 0 {
					bb.WriteByte('8' - k)
					continue J
				}
				bit >>= 1
			}
			bb.WriteByte('.')
		}
		L("SMAP: %02x: %s", i-begin, bb.String())
	}
}

func DoDumpPageZero() {
	saved_mmut := MmuTask
	MmuTask = 0
	saved_map00 := MmuMap[0][0]
	MmuMap[0][0] = 0
	defer func() {
		MmuTask = saved_mmut
		MmuMap[0][0] = saved_map00
	}()
	////////////////////////////

	L("PageZero:\n")

	/* some Level1:
	L("PageZero: FreeBitMap=%x:%x MemoryLimit=%x ModDir=%x RomBase=%x\n",
		W(sym.D_FMBM), W(sym.D_FMBM+2), W(sym.D_MLIM), W(sym.D_ModDir), W(sym.D_Init))
	*/
	L("  D_SWI3=%x D_SWI2=%x FIRQ=%x IRQ=%x SWI=%x NMI=%x SvcIRQ=%x Poll=%x\n",
		W(sym.D_SWI3), W(sym.D_SWI2), W(sym.D_FIRQ), W(sym.D_IRQ), W(sym.D_SWI), W(sym.D_NMI), W(sym.D_SvcIRQ), W(sym.D_Poll))
	/* some Level1:
	L("  BTLO=%x BTHI=%x  IO Free Mem Lo=%x Hi=%x D_DevTbl=%x D_PolTbl=%x D_PthDBT=%x D_Proc=%x\n",
		W(sym.D_BTLO), W(sym.D_BTHI), W(sym.D_IOML), W(sym.D_IOMH), W(sym.D_DevTbl), W(sym.D_PolTbl), W(sym.D_PthDBT), W(sym.D_Proc))
	*/
	L("  D_Slice=%x D_TSlice=%x\n",
		W(sym.D_Slice), W(sym.D_TSlice))

	var buf bytes.Buffer
	Z(&buf, " D.Tasks=%04x", PeekW(sym.D_Tasks))
	Z(&buf, " D.TmpDAT=%04x", PeekW(sym.D_TmpDAT))
	Z(&buf, " D.Init=%04x", PeekW(sym.D_Init))
	Z(&buf, " D.Poll=%04x", PeekW(sym.D_Poll))
	Z(&buf, " D.Tick=%02x", PeekB(sym.D_Tick))
	Z(&buf, " D.Slice=%02x", PeekB(sym.D_Slice))
	Z(&buf, " D.TSlice=%02x", PeekB(sym.D_TSlice))
	Z(&buf, " D.Boot=%02x", PeekB(sym.D_Boot))
	Z(&buf, " D.MotOn=%02x", PeekB(sym.D_MotOn))
	Z(&buf, " D.ErrCod=%02x", PeekB(sym.D_ErrCod))
	Z(&buf, " D.Daywk=%02x", PeekB(sym.D_Daywk))
	Z(&buf, " D.TkCnt=%02x", PeekB(sym.D_TkCnt))
	Z(&buf, " D.BtPtr=%04x", PeekW(sym.D_BtPtr))
	Z(&buf, " D.BtSz=%04x", PeekW(sym.D_BtSz))
	L("%s", buf.String())
	buf.Reset()

	Z(&buf, " D.CRC=%02x", PeekB(sym.D_CRC))
	Z(&buf, " D.Tenths=%02x", PeekB(sym.D_Tenths))
	Z(&buf, " D.Task1N=%02x", PeekB(sym.D_Task1N))
	Z(&buf, " D.Quick=%02x", PeekB(sym.D_Quick))
	Z(&buf, " D.QIRQ=%02x", PeekB(sym.D_QIRQ))
	Z(&buf, " D.BlkMap=%04x,%04x", PeekW(sym.D_BlkMap), PeekW(sym.D_BlkMap+2))
	Z(&buf, " D.ModDir=%04x,%04x", PeekW(sym.D_ModDir), PeekW(sym.D_ModDir+2))
	Z(&buf, " D.PrcDBT=%04x", PeekW(sym.D_PrcDBT))
	Z(&buf, " D.SysPrc=%04x", PeekW(sym.D_SysPrc))
	Z(&buf, " D.SysDAT=%04x", PeekW(sym.D_SysDAT))
	// Z(&buf, " D.Mem=%04x", PeekW(sym.D_Mem))
	Z(&buf, " D.Proc=%04x", PeekW(sym.D_Proc))
	Z(&buf, " D.AProcQ=%04x", PeekW(sym.D_AProcQ))
	Z(&buf, " D.WProcQ=%04x", PeekW(sym.D_WProcQ))
	Z(&buf, " D.SProcQ=%04x", PeekW(sym.D_SProcQ))
	L("%s", buf.String())
	buf.Reset()

	Z(&buf, " D.ModEnd=%04x", PeekW(sym.D_ModEnd))
	Z(&buf, " D.ModDAT=%04x", PeekW(sym.D_ModDAT))
	Z(&buf, " D.CldRes=%04x", PeekW(sym.D_CldRes))
	Z(&buf, " D.BtBug=%04x%02x", PeekW(sym.D_BtBug), PeekB(sym.D_BtBug+2))
	Z(&buf, " D.Pipe=%04x", PeekW(sym.D_Pipe))

	Z(&buf, " D.QCnt=%02x", PeekB(sym.D_QCnt))
	Z(&buf, " D.DevTbl=%04x", PeekW(sym.D_DevTbl))
	Z(&buf, " D.PolTbl=%04x", PeekW(sym.D_PolTbl))
	Z(&buf, " D.PthDBT=%04x", PeekW(sym.D_PthDBT))
	Z(&buf, " D.DMAReq=%02x", PeekB(sym.D_DMAReq))
	L("%s", buf.String())
	buf.Reset()
}

func DoDumpProcDesc(a Word, queue string, followQ bool) {
	if true {

		// PrettyDumpHex64(a, 0x100)

		tmp := MmuTask
		MmuTask = 0
		defer func() {
			MmuTask = tmp
		}()

		currency := ""
		if W(sym.D_Proc) == a {
			currency = " CURRENT "
		}
		// L("a=%04x", a)
		// switch Level {
		// case 1, 2:
		// {
		begin := PeekW(a + sym.P_PModul)
		name_str := "?"
		mod_str := "?"
		if begin != 0 {
			m := GetMappingTask0(a + sym.P_DATImg)
			modPhys := MapAddrWithMapping(begin, m)
			modPhysPlus4 := PeekWPhys(modPhys + 4)
			if modPhysPlus4 > 0 {
				name := begin + modPhysPlus4
				name_str = Os9StringWithMapping(name, m)
				mod_str = F("%q @%04x %v", name_str, begin, m)
			}
		}
		L("Process %x %s %s @%x: id=%x pid=%x sid=%x cid=%x module=%s", B(a+sym.P_PID), queue, currency, a, B(a+sym.P_ID), B(a+sym.P_PID), B(a+sym.P_SID), B(a+sym.P_CID), mod_str)

		L("   sp=%x task=%x PagCnt=%x User=%x Pri=%x Age=%x State=%x",
			W(a+sym.P_SP), B(a+sym.P_Task), B(a+sym.P_PagCnt), W(a+sym.P_User), B(a+sym.P_Prior), B(a+sym.P_Age), B(a+sym.P_State))

		L("   Queue=%x IOQP=%x IOQN=%x Signal=%x SigVec=%x SigDat=%x",
			W(a+sym.P_Queue), B(a+sym.P_IOQP), B(a+sym.P_IOQN), B(a+sym.P_Signal), B(a+sym.P_SigVec), B(a+sym.P_SigDat))
		L("   DIO %x %x %x  %x %x %x  PATH %x %x %x %x  %x %x %x %x  %x %x %x %x  %x %x %x %x",
			W(a+sym.P_DIO), W(a+sym.P_DIO+2), W(a+sym.P_DIO+4),
			W(a+sym.P_DIO+6), W(a+sym.P_DIO+8), W(a+sym.P_DIO+10),
			B(a+sym.P_Path+0), B(a+sym.P_Path+1), B(a+sym.P_Path+2), B(a+sym.P_Path+3),
			B(a+sym.P_Path+4), B(a+sym.P_Path+5), B(a+sym.P_Path+6), B(a+sym.P_Path+7),
			B(a+sym.P_Path+8), B(a+sym.P_Path+9), B(a+sym.P_Path+10), B(a+sym.P_Path+11),
			B(a+sym.P_Path+12), B(a+sym.P_Path+13), B(a+sym.P_Path+14), B(a+sym.P_Path+15))

		if paranoid {
			if B(a+sym.P_ID) > 10 {
				panic("P_ID")
			}
			if B(a+sym.P_PID) > 10 {
				panic("P_PID")
			}
			if B(a+sym.P_SID) > 10 {
				panic("P_SID")
			}
			if B(a+sym.P_CID) > 10 {
				panic("P_CID")
			}
			if W(a+sym.P_User) > 10 {
				panic("P_User")
			}
		}

		if followQ && W(a+sym.P_Queue) != 0 && queue != "Current" {
			DoDumpProcDesc(W(a+sym.P_Queue), queue, followQ)
		}

		// }
		// }
	}
}
func MemoryModuleOf(addr Word) (name string, offset Word) {
	// TODO -- cache current regions.

	//if enableRom {
	//return "(rom)", addr
	//}
	// TODO: speed up with caching.
	if addr >= 0xFF00 {
		log.Panicf("PC in IO page: $%x", addr)
	}
	if addr >= 0xFE00 {
		return "(tramp)", addr
	}
	if addr < 0x0100 {
		return "(zero)", addr
	}

	addrPhys := MapAddr(addr, true)
	addr32 := uint32(addrPhys)

	// First scan for initial modules.
	for _, m := range InitialModules {
		if addr32 >= m.Addr && addr32 < (m.Addr+m.Len) {
			return m.Id(), Word(addr32 - m.Addr)
		}
	}

	modDirStart := SysMemW(sym.D_ModDir)
	modDirLimit := SysMemW(sym.D_ModEnd)
	if modDirStart == 0 || modDirLimit == 0 {
		return "==", addr
	}
	for i := modDirStart; i < modDirLimit; i += 8 {
		datPtr := SysMemW(i + 0)
		if datPtr == 0 {
			continue
		}
		links := SysMemW(i + 6)
		if links == 0 { // ddt Mon May 29 12:59:19 PM PDT 2023
			// continue
		}
		begin := SysMemW(i + 4)
		//unused// usedBytes := SysMemW(i + 2)

		m := GetMapping(datPtr)
		magic := PeekWWithMapping(begin, m)
		if magic != 0x87CD {
			return "====", addr
		}
		// log.Printf("DDT: TRY i=%x begin=%x %q .....", i, begin, ModuleId(begin, m))

		// Module offset 2 is module size.
		remaining := int(PeekWWithMapping(begin+2, m))
		// Module offset 4 is offset to name string.
		//unused// namePtr := begin + PeekWWithMapping(begin+4, m)
		// log.Printf("DDT: len=%x remaining=%x trying=%q", usedBytes, remaining, Os9StringWithMapping(namePtr, m))

		//-------------
		// beginP := MapAddrWithMapping(begin, m)

		region := begin
		offset := Word(0) // offset into module.
		for remaining > 0 {
			// If module crosses paged blocks, it has more than one region.
			regionP := MapAddrWithMapping(region, m)
			endOfRegionBlockP := 1 + (regionP | 0x1FFF)
			regionSize := remaining
			if int(regionSize) > endOfRegionBlockP-regionP {
				// A smaller region of the module.
				regionSize = endOfRegionBlockP - regionP
			}

			// log.Printf("DDT: try regionP=%x (phys=%x) regionEnds=%x remain=%x", regionP, addrPhys, regionP+int(regionSize), remaining)
			if regionP <= addrPhys && addrPhys < regionP+int(regionSize) {
				if links == 0 {
					// return "unlinkedMod", addr
					// log.Panicf("in unlinked module: i=%x addr=%x", i, addr)
				}
				id := ModuleId(begin, m)
				delta := offset + Word(int(addrPhys)-regionP)
				// log.Printf("DDT: [links=%x] FOUND %q+%x", links, id, delta)
				return id, delta
			}
			remaining -= regionSize
			regionP += regionSize
			region += Word(regionSize)
			offset += Word(regionSize)
			// log.Printf("DDT: advanced remaining=%x regionSize=%x", remaining, regionSize)
		}
	}
	// log.Printf("DDT: NOT FOUND")
	return "", 0 // No module found for the addr.
}
func MemoryModules() {
	WithKernelTask(func() {

		L("[all mem]")
		DumpAllMemory()
		L("[page zero]")
		DumpPageZero()
		L("[processes]")
		DumpProcesses()
		L("[all path descs]")
		DumpAllPathDescs()
		L("[block zero]")
		DoDumpBlockZero()
		L("\n#MemoryModules(")

		var buf bytes.Buffer
		Z(&buf, "MOD name begin:end(len/blocklen) [addr:dat,blocklen,begin,links] dat\n")

		modDirStart := SysMemW(sym.D_ModDir)
		modDirLimit := SysMemW(sym.D_ModEnd)
		for i := modDirStart; i < modDirLimit; i += 8 {
			datPtr := SysMemW(i + 0)
			if datPtr == 0 {
				continue
			}
			usedBytes := SysMemW(i + 2)
			begin := SysMemW(i + 4)
			links := SysMemW(i + 6)

			m := GetMapping(datPtr)
			end := begin + PeekWWithMapping(begin+2, m)
			name := begin + PeekWWithMapping(begin+4, m)

			Z(&buf, "MOD %s %x:%x(%x/%x) [%x:%x,%x,%x,%x] %v\n", Os9StringWithMapping(name, m), begin, end, end-begin, usedBytes, i, datPtr, usedBytes, begin, links, m)
		}
		L("%s", buf.String())
		L("#MemoryModules)")
	})
}
func PrettyDumpHex64(addr Word, size uint) {
	// L(";")
	const PERLINE = 64
	var p Word
	for k := uint(0); k < size; k += PERLINE {
		p = addr + Word(k)
		var i Word
		for i = 0; i < PERLINE; i += 2 {
			if PeekW(p+i) != 0 {
				break
			}
		}
		if i == PERLINE {
			continue // don't print all zeros row.
		}
		var buf bytes.Buffer
		Z(&buf, "%04x:", p)
		for q := Word(0); q < PERLINE; q += 2 {
			if q&7 == 0 {
				Z(&buf, " ")
			}
			if q&15 == 0 {
				Z(&buf, " ")
			}
			w := PeekW(p + q)
			if w == 0 {
				Z(&buf, "---- ")
			} else {
				Z(&buf, "%04x ", w)
			}
		}
		for q := Word(0); q < PERLINE; q += 1 {
			x := PeekB(p + q)
			if ' ' <= x && x <= '~' {
				Z(&buf, "%c", x)
			} else {
				Z(&buf, ".")
			}
			if (q & 7) == 7 {
				Z(&buf, "|")
			}
		}
		L("%s", buf.String())
	}
	// L(";")
}

func DoDumpProcsAndPaths() {
	WithMmuTask(0, DoDumpProcsAndPathsPrime)
}
func DoDumpProcsAndPathsPrime() {
	if V['p'] {
		Logp("ProcsAndPaths:\n")

		DoDumpPageZero()

		Logp("ProcQ:Active")
		active := PeekW(sym.D_AProcQ)
		if active != 0 {
			DoDumpProcDesc(active, "active", true)
		}

		Logp("ProcQ:Waiting")
		waiting := PeekW(sym.D_WProcQ)
		if waiting != 0 {
			DoDumpProcDesc(waiting, "waiting", true)
		}

		Logp("ProcQ:Sleeping")
		sleeping := PeekW(sym.D_SProcQ)
		if sleeping != 0 {
			DoDumpProcDesc(sleeping, "sleeping", true)
		}

		DoDumpPaths()
		DoDumpDevices()
	}
}

func DoDumpPaths() {
	dbt := PeekW(sym.D_PthDBT)
	Logp("DoDumpPaths=%x", dbt)
	if dbt == 0 {
		return
	}
	PrettyDumpHex64(dbt, 64)
	for e := Word(0); e < 64; e++ {
		ext := PeekB(dbt + e)
		// Logp("e=%x ext=%x", e, ext)
		if ext != 0 {
			p := Word(ext) << 8
			for j := Word(0); j < 4; j++ {
				i := e*4 + j
				if i == 0 {
					continue
				} // Skip directory slot.
				pth := p + j*64
				if PeekB(pth) == byte(i) {
					// Logp("e=%x j=%x PATH=%x pth=%x", e, j, i, pth)
					DoDumpPath(i, pth)
					PrettyDumpHex64(pth, 64)
				}
			}
		}
	}
	Logp(";;")
}

func DoDumpPath(i Word, pth Word) {
	/*
	                  ORG       0
	   PD.PD          RMB       1                   Path Number
	   PD.MOD         RMB       1                   Mode (Read/Write/Update)
	   PD.CNT         RMB       1                   Number of Open Images
	   PD.DEV         RMB       2                   Device Table Entry Address
	   PD.CPR         RMB       1                   Current Process
	   PD.RGS         RMB       2                   Caller's Register Stack
	   PD.BUF         RMB       2                   Buffer Address
	   PD.FST         RMB       32-.                File Manager's Storage
	   PD.OPT         EQU       .                   PD GetSts(0) Options
	   PD.DTP         RMB       1                   Device Type
	                  RMB       64-.                Path options
	   PDSIZE         EQU       .
	*/
	Logp("PATH[%d] at $%x", i, pth)
	var mod_name string
	dte := PeekW(pth + sym.PD_DEV)
	if dte != 0 {
		descriptorMod := PeekW(dte + 4)
		if descriptorMod != 0 {
			mod_name = ModuleName(descriptorMod)
		}
	}
	Logp("  mode %x count %x dev %x %q cur_proc %x regs %x buf %x type %x",
		PeekB(pth+sym.PD_MOD),
		PeekB(pth+sym.PD_CNT),
		dte,
		mod_name,
		PeekB(pth+sym.PD_CPR),
		PeekW(pth+sym.PD_RGS),
		PeekW(pth+sym.PD_BUF),
		PeekB(pth+sym.PD_DTP))
}

func DoDumpDevices() {
	/*
	    983 *********************
	    984 * Device Table Format
	    985 *
	    986                ORG       0
	    987 V$DRIV         RMB       2                   Device Driver module
	    988 V$STAT         RMB       2                   Device Driver Static storage
	    989 V$DESC         RMB       2                   Device Descriptor module
	    990 V$FMGR         RMB       2                   File Manager module
	    991 V$USRS         RMB       1                   use count
	    992                IFGT      Level-1
	    993 V$DRIVEX       RMB       2                   Device Driver execution address
	    994 V$FMGREX       RMB       2                   File Manager execution address
	    995                ENDC
	    996 DEVSIZ         EQU       .
	    997
	    998 *******************************
	    999 * Device Static Storage Offsets
	   1000 *
	   1001                ORG       0
	   1002 V.PAGE         RMB       1                   Port Extended Address
	   1003 V.PORT         RMB       2                   Device 'Base' Port Address
	   1004 V.LPRC         RMB       1                   Last Active Process ID
	   1005 V.BUSY         RMB       1                   Active Process ID (0=UnBusy)
	   1006 V.WAKE         RMB       1                   Active PD if Driver MUST Wake-up
	   1007 V.USER         EQU       .                   Driver Allocation Origin
	*/

	devTable := PeekW(sym.D_DevTbl)
	Logp("Device Table at $%x", devTable)
	if devTable == 0 {
		return
	}
	/*----*/
	init := PeekW(sym.D_Init)
	devCount := PeekB(init + 13) //sym.DevCnt
	Logp("   ... %q Init at %x, devCount=%x [ init=%q sys=%q std=%q boot=%q os=%q install=%q level=%x %x.%x.%x +%x +%x ]",
		ModuleName(init), init, devCount,
		Os9String(init+PeekW(init+14)),
		Os9String(init+PeekW(init+16)),
		Os9String(init+PeekW(init+18)),
		Os9String(init+PeekW(init+20)),
		Os9String(init+PeekW(init+0x1D)),
		Os9String(init+PeekW(init+0x1F)),
		PeekB(init+0x17),
		PeekB(init+0x18),
		PeekB(init+0x19),
		PeekB(init+0x1A),
		PeekB(init+0x1B),
		PeekB(init+0x1C),
	)
	/*----*/
	const devSize = 13 // sym.DEVSIZ

	for i := byte(0); i < devCount; i++ {
		p := devTable + Word(i)*Word(devSize)
		driverMod := PeekW(p + 0)
		staticStorage := PeekW(p + 2)
		descriptorMod := PeekW(p + 4)
		managerMod := PeekW(p + 6)
		count := PeekB(p + 8)
		if descriptorMod != 0 {
			Logp("   [%x] p=%x %q=desc=%x %q=driv=%x %q=mgr=%x store=%x count=%x", i, p, ModuleName(descriptorMod), descriptorMod, ModuleName(driverMod), driverMod, ModuleName(managerMod), managerMod, staticStorage, count)
			if staticStorage != 0 {
				PrettyDumpHex64(staticStorage, 0x100)
			}
		}
	}
	Logp(";;")

	/*
		  5005                ldb       #DEVSIZ             Size of each device table entry
		5006                ldx       <D.Init             Get ptr to INIT module
		5007                lda       DevCnt,x            Get # of entries allowed in device table
		5008                ldx       <D.DevTbl           Get start of device table
		5009                mul                           Calculate offset to end of device table
		5010                leay      d,x                 Point Y to end of Device table
		5011                ldb       #DEVSIZ             Get device table entry size again
		5012 DevLoop        ldu       V$DRIV,x            Get driver ptr for device we are checking
		5013                ifne      H6309
		5014                cmpr      u,w                 Same as original window?
		5015                else
		5016                cmpu      >GrfMem+gr00B5
		5017                endc
		5018                bne       NextEnt             No, skip to next entry
		5019                ldu       V$STAT,x            Get static mem ptr for CC3/TC9IO device
		5020                lda       V.WinType,u         Is this a Windint/Grfint window?
		5021                bne       NextEnt             No, VDGINT so skip
		5022                lda       V.InfVld,u          Is this static mem properly initialized?
		5023                beq       NextEnt             No, skip
	*/
}
