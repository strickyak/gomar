package emu

import (
	. "github.com/strickyak/gomar/gu"
	"github.com/strickyak/gomar/sym"

	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
)

var _ = Log

func DecodeOs9Error(b byte) string {
	s := "???"
	switch b {
	case 0x0A:
		s = "E$UnkSym :Unknown symbol"
		break
	case 0x0B:
		s = "E$ExcVrb :Excessive verbage"
		break
	case 0x0C:
		s = "E$IllStC :Illegal statement construction"
		break
	case 0x0D:
		s = "E$ICOvf  :I-code overflow"
		break
	case 0x0E:
		s = "E$IChRef :Illegal channel reference"
		break
	case 0x0F:
		s = "E$IllMod :Illegal mode"
		break
	case 0x10:
		s = "E$IllNum :Illegal number"
		break
	case 0x11:
		s = "E$IllPrf :Illegal prefix"
		break
	case 0x12:
		s = "E$IllOpd :Illegal operand"
		break
	case 0x13:
		s = "E$IllOpr :Illegal operator"
		break
	case 0x14:
		s = "E$IllRFN :Illegal record field name"
		break
	case 0x15:
		s = "E$IllDim :Illegal dimension"
		break
	case 0x16:
		s = "E$IllLit :Illegal literal"
		break
	case 0x17:
		s = "E$IllRet :Illegal relational"
		break
	case 0x18:
		s = "E$IllSfx :Illegal type suffix"
		break
	case 0x19:
		s = "E$DimLrg :Dimension too large"
		break
	case 0x1A:
		s = "E$LinLrg :Line number too large"
		break
	case 0x1B:
		s = "E$NoAssg :Missing assignment statement"
		break
	case 0x1C:
		s = "E$NoPath :Missing path number"
		break
	case 0x1D:
		s = "E$NoComa :Missing comma"
		break
	case 0x1E:
		s = "E$NoDim  :Missing dimension"
		break
	case 0x1F:
		s = "E$NoDO   :Missing DO statement"
		break
	case 0x20:
		s = "E$MFull  :Memory full"
		break
	case 0x21:
		s = "E$NoGoto :Missing GOTO"
		break
	case 0x22:
		s = "E$NoLPar :Missing left parenthesis"
		break
	case 0x23:
		s = "E$NoLRef :Missing line reference"
		break
	case 0x24:
		s = "E$NoOprd :Missing operand"
		break
	case 0x25:
		s = "E$NoRPar :Missing right parenthesis"
		break
	case 0x26:
		s = "E$NoTHEN :Missing THEN statement"
		break
	case 0x27:
		s = "E$NoTO   :Missing TO statement"
		break
	case 0x28:
		s = "E$NoVRef :Missing variable reference"
		break
	case 0x29:
		s = "E$EndQou :Missing end quote"
		break
	case 0x2A:
		s = "E$SubLrg :Too many subscripts"
		break
	case 0x2B:
		s = "E$UnkPrc :Unknown procedure"
		break
	case 0x2C:
		s = "E$MulPrc :Multiply defined procedure"
		break
	case 0x2D:
		s = "E$DivZer :Divice by zero"
		break
	case 0x2E:
		s = "E$TypMis :Operand type mismatch"
		break
	case 0x2F:
		s = "E$StrOvf :String stack overflow"
		break
	case 0x30:
		s = "E$NoRout :Unimplemented routine"
		break
	case 0x31:
		s = "E$UndVar :Undefined variable"
		break
	case 0x32:
		s = "E$FltOvf :Floating Overflow"
		break
	case 0x33:
		s = "E$LnComp :Line with compiler error"
		break
	case 0x34:
		s = "E$ValRng :Value out of range for destination"
		break
	case 0x35:
		s = "E$SubOvf :Subroutine stack overflow"
		break
	case 0x36:
		s = "E$SubUnd :Subroutine stack underflow"
		break
	case 0x37:
		s = "E$SubRng :Subscript out of range"
		break
	case 0x38:
		s = "E$ParmEr :Paraemter error"
		break
	case 0x39:
		s = "E$SysOvf :System stack overflow"
		break
	case 0x3A:
		s = "E$IOMism :I/O type mismatch"
		break
	case 0x3B:
		s = "E$IONum  :I/O numeric input format bad"
		break
	case 0x3C:
		s = "E$IOConv :I/O conversion: number out of range"
		break
	case 0x3D:
		s = "E$IllInp :Illegal input format"
		break
	case 0x3E:
		s = "E$IOFRpt :I/O format repeat error"
		break
	case 0x3F:
		s = "E$IOFSyn :I/O format syntax error"
		break
	case 0x40:
		s = "E$IllPNm :Illegal path number"
		break
	case 0x41:
		s = "E$WrSub  :Wrong number of subscripts"
		break
	case 0x42:
		s = "E$NonRcO :Non-record type operand"
		break
	case 0x43:
		s = "E$IllA   :Illegal argument"
		break
	case 0x44:
		s = "E$IllCnt :Illegal control structure"
		break
	case 0x45:
		s = "E$UnmCnt :Unmatched control structure"
		break
	case 0x46:
		s = "E$IllFOR :Illegal FOR variable"
		break
	case 0x47:
		s = "E$IllExp :Illegal expression type"
		break
	case 0x48:
		s = "E$IllDec :Illegal declarative statement"
		break
	case 0x49:
		s = "E$ArrOvf :Array size overflow"
		break
	case 0x4A:
		s = "E$UndLin :Undefined line number"
		break
	case 0x4B:
		s = "E$MltLin :Multiply defined line number"
		break
	case 0x4C:
		s = "E$MltVar :Multiply defined variable"
		break
	case 0x4D:
		s = "E$IllIVr :Illegal input variable"
		break
	case 0x4E:
		s = "E$SeekRg :Seek out of range"
		break
	case 0x4F:
		s = "E$NoData :Missing data statement"
		break
	case 0xB7:
		s = "E$IWTyp  :Illegal window type"
		break
	case 0xB8:
		s = "E$WADef  :Window already defined"
		break
	case 0xB9:
		s = "E$NFont  :Font not found"
		break
	case 0xBA:
		s = "E$StkOvf :Stack overflow"
		break
	case 0xBB:
		s = "E$IllArg :Illegal argument"
		break
	case 0xBD:
		s = "E$ICoord :Illegal coordinates"
		break
	case 0xBE:
		s = "E$Bug    :Bug (should never be returned)"
		break
	case 0xBF:
		s = "E$BufSiz :Buffer size is too small"
		break
	case 0xC0:
		s = "E$IllCmd :Illegal command"
		break
	case 0xC1:
		s = "E$TblFul :Screen or window table is full"
		break
	case 0xC2:
		s = "E$BadBuf :Bad/Undefined buffer number"
		break
	case 0xC3:
		s = "E$IWDef  :Illegal window definition"
		break
	case 0xC4:
		s = "E$WUndef :Window undefined"
		break
	case 0xC5:
		s = "E$Up     :Up arrow pressed on SCF I$ReadLn with PD.UP enabled"
		break
	case 0xC6:
		s = "E$Dn     :Down arrow pressed on SCF I$ReadLn with PD.DOWN enabled"
		break
	case 0xC8:
		s = "E$PthFul :Path Table full"
		break
	case 0xC9:
		s = "E$BPNum  :Bad Path Number"
		break
	case 0xCA:
		s = "E$Poll   :Polling Table Full"
		break
	case 0xCB:
		s = "E$BMode  :Bad Mode"
		break
	case 0xCC:
		s = "E$DevOvf :Device Table Overflow"
		break
	case 0xCD:
		s = "E$BMID   :Bad Module ID"
		break
	case 0xCE:
		s = "E$DirFul :Module Directory Full"
		break
	case 0xCF:
		s = "E$MemFul :Process Memory Full"
		break
	case 0xD0:
		s = "E$UnkSvc :Unknown Service Code"
		break
	case 0xD1:
		s = "E$ModBsy :Module Busy"
		break
	case 0xD2:
		s = "E$BPAddr :Bad Page Address"
		break
	case 0xD3:
		s = "E$EOF    :End of File"
		break
	case 0xD5:
		s = "E$NES    :Non-Existing Segment"
		break
	case 0xD6:
		s = "E$FNA    :File Not Accesible"
		break
	case 0xD7:
		s = "E$BPNam  :Bad Path Name"
		break
	case 0xD8:
		s = "E$PNNF   :Path Name Not Found"
		break
	case 0xD9:
		s = "E$SLF    :Segment List Full"
		break
	case 0xDA:
		s = "E$CEF    :Creating Existing File"
		break
	case 0xDB:
		s = "E$IBA    :Illegal Block Address"
		break
	case 0xDC:
		s = "E$HangUp :Carrier Detect Lost"
		break
	case 0xDD:
		s = "E$MNF    :Module Not Found"
		break
	case 0xDF:
		s = "E$DelSP  :Deleting Stack Pointer memory"
		break
	case 0xE0:
		s = "E$IPrcID :Illegal Process ID"
		break
	case 0xE2:
		s = "E$NoChld :No Children"
		break
	case 0xE3:
		s = "E$ISWI   :Illegal SWI code"
		break
	case 0xE4:
		s = "E$PrcAbt :Process Aborted"
		break
	case 0xE5:
		s = "E$PrcFul :Process Table Full"
		break
	case 0xE6:
		s = "E$IForkP :Illegal Fork Parameter"
		break
	case 0xE7:
		s = "E$KwnMod :Known Module"
		break
	case 0xE8:
		s = "E$BMCRC  :Bad Module CRC"
		break
	case 0xE9:
		s = "E$USigP  :Unprocessed Signal Pending"
		break
	case 0xEA:
		s = "E$NEMod  :Non Existing Module"
		break
	case 0xEB:
		s = "E$BNam   :Bad Name"
		break
	case 0xEC:
		s = "E$BMHP   :(bad module header parity)"
		break
	case 0xED:
		s = "E$NoRAM  :No (System) RAM Available"
		break
	case 0xEE:
		s = "E$DNE    :Directory not empty"
		break
	case 0xEF:
		s = "E$NoTask :No available Task number"
		break
	case 0xF0:
		s = "E$Unit   :Illegal Unit (drive)"
		break
	case 0xF1:
		s = "E$Sect   :Bad Sector number"
		break
	case 0xF2:
		s = "E$WP     :Write Protect"
		break
	case 0xF3:
		s = "E$CRC    :Bad Check Sum"
		break
	case 0xF4:
		s = "E$Read   :Read Error"
		break
	case 0xF5:
		s = "E$Write  :Write Error"
		break
	case 0xF6:
		s = "E$NotRdy :Device Not Ready"
		break
	case 0xF7:
		s = "E$Seek   :Seek Error"
		break
	case 0xF8:
		s = "E$Full   :Media Full"
		break
	case 0xF9:
		s = "E$BTyp   :Bad Type (incompatable) media"
		break
	case 0xFA:
		s = "E$DevBsy :Device Busy"
		break
	case 0xFB:
		s = "E$DIDC   :Disk ID Change"
		break
	case 0xFC:
		s = "E$Lock   :Record is busy (locked out)"
		break
	case 0xFD:
		s = "E$Share  :Non-sharable file busy"
		break
	case 0xFE:
		s = "E$DeadLk :I/O Deadlock error"
		break
	}
	return s
}

func DecodeOs9GetStat(b byte) string {
	s := "???"
	switch b {
	case 0x00:
		s = "SS.Opt    : Read/Write PD Options"
		break
	case 0x01:
		s = "SS.Ready  : Check for Device Ready"
		break
	case 0x02:
		s = "SS.Size   : Read/Write File Size"
		break
	case 0x03:
		s = "SS.Reset  : Device Restore"
		break
	case 0x04:
		s = "SS.WTrk   : Device Write Track"
		break
	case 0x05:
		s = "SS.Pos    : Get File Current Position"
		break
	case 0x06:
		s = "SS.EOF    : Test for End of File"
		break
	case 0x07:
		s = "SS.Link   : Link to Status routines"
		break
	case 0x08:
		s = "SS.ULink  : Unlink Status routines"
		break
	case 0x09:
		s = "SS.Feed   : Issue form feed"
		break
	case 0x0A:
		s = "SS.Frz    : Freeze DD. information"
		break
	case 0x0B:
		s = "SS.SPT    : Set DD.TKS to given value"
		break
	case 0x0C:
		s = "SS.SQD    : Sequence down hard disk"
		break
	case 0x0D:
		s = "SS.DCmd   : Send direct command to disk"
		break
	case 0x0E:
		s = "SS.DevNm  : Return Device name (32-bytes at [X])"
		break
	case 0x0F:
		s = "SS.FD     : Return File Descriptor (Y-bytes at [X])"
		break
	case 0x10:
		s = "SS.Ticks  : Set Lockout honor duration"
		break
	case 0x11:
		s = "SS.Lock   : Lock/Release record"
		break
	case 0x12:
		s = "SS.DStat  : Return Display Status (CoCo)"
		break
	case 0x13:
		s = "SS.Joy    : Return Joystick Value (CoCo)"
		break
	case 0x14:
		s = "SS.BlkRd  : Block Read"
		break
	case 0x15:
		s = "SS.BlkWr  : Block Write"
		break
	case 0x16:
		s = "SS.Reten  : Retension cycle"
		break
	case 0x17:
		s = "SS.WFM    : Write File Mark"
		break
	case 0x18:
		s = "SS.RFM    : Read past File Mark"
		break
	case 0x19:
		s = "SS.ELog   : Read Error Log"
		break
	case 0x1A:
		s = "SS.SSig   : Send signal on data ready"
		break
	case 0x1B:
		s = "SS.Relea  : Release device"
		break
	case 0x1C:
		s = "SS.AlfaS  : Return Alfa Display Status (CoCo, SCF/GetStat)"
		break
	case 0x1D:
		s = "SS.Break  : Send break signal out acia"
		break
	case 0x1E:
		s = "SS.RsBit  : Reserve bitmap sector (do not allocate in) LSB(X)=sct#"
		break
	case 0x20:
		s = "SS.DirEnt : Reserve bitmap sector (do not allocate in) LSB(X)=sct#"
		break
	case 0x24:
		s = "SS.SetMF  : Reserve $24 for Gimix G68 (Flex compatability?)"
		break
	case 0x25:
		s = "SS.Cursr  : Cursor information for COCO"
		break
	case 0x26:
		s = "SS.ScSiz  : Return screen size for COCO"
		break
	case 0x27:
		s = "SS.KySns  : Getstat/SetStat for COCO keyboard"
		break
	case 0x28:
		s = "SS.ComSt  : Getstat/SetStat for Baud/Parity"
		break
	case 0x29:
		s = "SS.Open   : SetStat to tell driver a path was opened"
		break
	case 0x2A:
		s = "SS.Close  : SetStat to tell driver a path was closed"
		break
	case 0x2B:
		s = "SS.HngUp  : SetStat to tell driver to hangup phone"
		break
	case 0x2C:
		s = "SS.FSig   : New signal for temp locked files"
		break
	}
	return s
}

func Os9StringN(addr Word, n Word) string {
	var buf bytes.Buffer
	for i := Word(0); i < n; i++ {
		var ch byte = 0x7F & PeekB(addr+i)
		if '!' <= ch && ch <= '~' {
			buf.WriteByte(ch)
		} else {
			Z(&buf, "{%d}", PeekB(addr+i))
		}
	}
	return buf.String()
}

func StringSomeBytesWithMapping(addr Word, mapping Mapping) string {
	var buf bytes.Buffer
	for i := 0; i < 16; i++ {
		var b byte = PeekBWithMapping(addr, mapping)
		fmt.Fprintf(&buf, "%2x", b)
	}
	return buf.String()
}

func SomeBytesWithMapping(addr Word, mapping Mapping) []byte {
	var bb []byte
	for i := Word(0); i < 16; i++ {
		var b byte = PeekBWithMapping(addr+i, mapping)
		bb = append(bb, b)
	}
	return bb
}

func Os9StringWithMapping(addr Word, mapping Mapping) string {
	var buf bytes.Buffer
	for {
		var b byte = PeekBWithMapping(addr, mapping)
		var ch byte = 0x7F & b
		if '!' <= ch && ch <= '~' {
			buf.WriteByte(ch)
		} else {
			break
		}
		if (b & 128) != 0 {
			break
		}
		addr++
	}
	return buf.String()
}

func Os9String(addr Word) string {
	var buf bytes.Buffer
	for {
		var b byte = PeekB(addr)
		var ch byte = 0x7F & b
		if '!' <= ch && ch <= '~' {
			buf.WriteByte(ch)
		} else {
			break
		}
		if (b & 128) != 0 {
			break
		}
		addr++
	}
	return buf.String()
}

func Os9StringPhys(addr int) string {
	var buf bytes.Buffer
	for {
		var b byte = mem[addr]
		var ch byte = 0x7F & b
		if '!' <= ch && ch <= '~' {
			buf.WriteByte(ch)
		} else {
			break
		}
		if (b & 128) != 0 {
			break
		}
		addr++
	}
	return buf.String()
}

func PrintableStringThruEOS(a Word, max Word) string {
	result := ""

	// It turns out that WritLn needs to get the buffer from User Task.
	// I didn't realize that, at first.
	WithMmuTask(1, func() {

		var buf bytes.Buffer
		//debug// Z(&buf, " [ %04x/%04x:%x: ", a, max, MmuTask) // yak
		//debug// for i := Word(0); i < max; i++ {
		//debug// Z(&buf, "%02x~", PeekB(a + i)) // yak
		//debug// }

		for i := Word(0); i < max; i++ {
			ch := PeekB(a + i)
			if 32 <= ch && ch < 127 {
				buf.WriteByte(ch)
			} else if ch == '\n' || ch == '\r' {
				buf.WriteByte('\n')
			} else if ch == 0 {
				break
			} else {
				Z(&buf, "{%d}", ch)
			}
			if ch == '\r' {
				break
			}
		}
		//debug// Z(&buf, " ] ") // yak
		result = buf.String()

	})
	return result
}

func StrungMemory(a Word, max Word) string {
	var buf bytes.Buffer
	for i := Word(0); i < yreg && i < max; i++ {
		buf.WriteByte(PeekB(a + i))
	}
	return buf.String()
}

func PrintableMemory(a Word, max Word) string {
	result := ""

	WithMmuTask(1, func() {

		var buf bytes.Buffer
		for i := Word(0); i < yreg && i < max; i++ {
			ch := PeekB(a + i)
			if 32 <= ch && ch < 127 {
				buf.WriteByte(ch)
			} else if ch == '\n' || ch == '\r' {
				buf.WriteByte('\n')
			} else {
				fmt.Fprintf(&buf, "{%d}", ch)
			}
		}
		result = buf.String()

	})
	return result
}

func ModuleName(module_loc Word) string {
	name_loc := module_loc + PeekW(module_loc+4)
	return Os9String(name_loc)
}

func Regs() string {
	var buf bytes.Buffer
	Z(&buf, "a=%02x b=%02x x=%04x:%04x y=%04x:%04x u=%04x:%04x s=%04x:%04x,%04x cc=%s dp=%02x #%d",
		GetAReg(), GetBReg(), xreg, PeekW(xreg), yreg, PeekW(yreg), ureg, PeekW(ureg), sreg, PeekW(sreg), PeekW(sreg+2), ccbits(ccreg), dpreg, Cycles)
	return buf.String()
}

var Expectations []string

func InitExpectations() {
	if *FlagExpectFile != "" && Expectations == nil {
		bb := Value(os.ReadFile(*FlagExpectFile))
		Expectations = strings.Split(string(bb), "\n")
	}
	if *FlagExpect != "" && Expectations == nil {
		Expectations = strings.Split(*FlagExpect, ";")
		fmt.Printf("\n===@=== SET Expectations: %q\n", *FlagExpect)
		log.Printf("\n===@=== SET Expectations: %q\n", *FlagExpect)
	}
}

func CheckExpectation(got string) {
	// Skip out if no expectations were defined.
	if len(Expectations) == 0 {
		return
	}

	// Skip empty expectations
	for len(Expectations) > 0 && len(Expectations[0]) == 0 {
		Expectations = Expectations[1:]
	}

	// Process one expectation, if possible, if valid.
	if len(Expectations) > 0 {
		if strings.Contains(got, Expectations[0]) {
			fmt.Printf("\n===@=== GOT Expectation: %q\n", Expectations[0])
			log.Printf("\n===@=== GOT Expectation: %q\n", Expectations[0])
			Expectations = Expectations[1:]

		}
	}

	// Skip empty expectations
	for len(Expectations) > 0 && len(Expectations[0]) == 0 {
		Expectations = Expectations[1:]
	}

	// Exit(0) if all expectations are met.
	if len(Expectations) == 0 {
		fmt.Printf("\n===@=== SUCCESS -- FINISHED Expectations.\n")
		log.Printf("\n===@=== SUCCESS -- FINISHED Expectations.\n")
		os.Exit(0)
	}
}

// Returns a string and whether this operation typically returns to caller.
func DecodeOs9Opcode(b byte) (string, bool) {
	MemoryModules()
	DoDumpPaths()
	s, p := "", ""
	returns := true
	switch b {
	case 0x00:
		s = "F$Link   : Link to Module"
		p = F("type/lang=%02x module/file='%s'", GetAReg(), Os9String(xreg))

	case 0x01:
		s = "F$Load   : Load Module from File"
		p = F("type/lang=%02x filename='%s'", GetAReg(), Os9String(xreg))

	case 0x02:
		s = "F$UnLink : Unlink Module"
		p = F("u=%04x magic=%04x module='%s'", ureg, PeekW(ureg), ModuleName(ureg))

	case 0x03:
		s = "F$Fork   : Start New Process"
		p = F("Module/file='%s' param=%q lang/type=%x pages=%x", Os9String(xreg), Os9StringN(ureg, yreg), GetAReg(), GetBReg())

	case 0x04:
		s = "F$Wait   : Wait for Child Process to Die"

	case 0x05:
		s = "F$Chain  : Chain Process to New Module"
		p = F("Module/file='%s' param=%q lang/type=%x pages=%x", Os9String(xreg), Os9StringN(ureg, yreg), GetAReg(), GetBReg())
		returns = false

	case 0x06:
		s = "F$Exit   : Terminate Process"
		p = F("status=%x", GetBReg())
		returns = false

	case 0x07:
		s = "F$Mem    : Set Memory Size"
		p = F("desired_size=%x", dreg)

	case 0x08:
		s = "F$Send   : Send Signal to Process"
		p = F("pid=%02x signal=%02x", GetAReg(), GetBReg())

	case 0x09:
		s = "F$Icpt   : Set Signal Intercept"
		p = F("routine=%04x storage=%04x", xreg, ureg)

	case 0x0A:
		s = "F$Sleep  : Suspend Process with Sleep"
		p = F("ticks=%04x", xreg)

	case 0x0B:
		s = "F$SSpd   : Suspend Process with SSpd (unused?)"

	case 0x0C:
		s = "F$ID     : Return Process ID"

	case 0x0D:
		s = "F$SPrior : Set Process Priority"
		p = F("pid=%02x priority=%02x", GetAReg(), GetBReg())

	case 0x0E:
		s = "F$SSWI   : Set Software Interrupt"
		p = F("code=%02x addr=%04x", GetAReg(), xreg)

	case 0x0F:
		s = "F$PErr   : Print Error"

	case 0x10:
		s = "F$PrsNam : Parse Pathlist Name"
		p = F("path='%s'", Os9String(xreg))
	case 0x11:
		s = "F$CmpNam : Compare Two Names"
		p = F("first=%q second=%q", Os9StringN(xreg, Word(GetBReg())), Os9String(yreg))

	case 0x12:
		s = "F$SchBit : Search Bit Map"
		p = F("bitmap=%04x end=%04x first=%x count=%x", xreg, ureg, dreg, yreg)

	case 0x13:
		s = "F$AllBit : Allocate in Bit Map"
		p = F("bitmap=%04x first=%x count=%x", xreg, dreg, yreg)

	case 0x14:
		s = "F$DelBit : Deallocate in Bit Map"
		p = F("bitmap=%04x first=%x count=%x", xreg, dreg, yreg)

	case 0x15:
		s = "F$Time   : Get Current Time"
		p = F("buf=%x", xreg)

	case 0x16:
		s = "F$STime  : Set Current Time"
		p = F("y%d m%d d%d h%d m%d s%d", PeekB(xreg+0), PeekB(xreg+1), PeekB(xreg+2), PeekB(xreg+3), PeekB(xreg+4), PeekB(xreg+5))

	case 0x17:
		s = "F$CRC    : Generate CRC ($1"
		p = F("addr=%04x len=%04x buf=%04x", xreg, yreg, ureg)

	// NitrOS9:

	case 0x27:
		s = "F$VIRQ   : Install/Delete Virtual IRQ"

	case 0x28:
		s = "F$SRqMem : System Memory Request"
		p = F("size=%x", dreg)

	case 0x29:
		s = "F$SRtMem : System Memory Return"
		p = F("size=%x start=%x", dreg, ureg)

	case 0x2A:
		s = "F$IRQ    : Enter IRQ Polling Table"

	case 0x2B:
		s = "F$IOQu   : Enter I/O Queue"
		p = F("pid=%02x", GetAReg())

	case 0x2C:
		s = "F$AProc  : Enter Active Process Queue"
		p = F("proc=%x", xreg)

	case 0x2D:
		s = "F$NProc  : Start Next Process"
		returns = false

	case 0x2E:
		s = "F$VModul : Validate Module"
		p = VerboseValidateModuleSyscall()

	case 0x2F:
		s = "F$Find64 : Find Process/Path Descriptor"
		p = F("base=%04x id=%x", xreg, GetAReg())

	case 0x30:
		s = "F$All64  : Allocate Process/Path Descriptor"
		p = F("table=%x", xreg)

	case 0x31:
		s = "F$Ret64  : Return Process/Path Descriptor"
		p = F("block_num=%x address=%x", GetAReg(), xreg)

	case 0x32:
		s = "F$SSvc   : Service Request Table Initialization"
		p = F("table=%x", yreg)

	case 0x33:
		s = "F$IODel  : Delete I/O Module"
		p = F("module=%x", xreg)

		// Level 2:
	case 0x34:
		s = "F$SLink  : System Link"
		mapping := GetMappingTask0(yreg)
		p = F("%q type %x name@ %x dat@ %x", Os9StringWithMapping(xreg, mapping), GetAReg(), xreg, yreg)

	case 0x38:
		s = "F$Move   : Move data (low bound first)"
		p = F("srcTask=%x destTask=%x srcPtr=%04x destPtr=%04x size=%04x", GetAReg(), GetBReg(), xreg, ureg, yreg)

	case 0x39:
		s = "F$AllRAM : Allocate RAM blocks"
		p = F("numBlocks=%x", GetBReg())

	case 0x3A:
		s = "F$AllImg : Allocate Image RAM blocks"
		p = F("beginBlock=%x numBlocks=%x processDesc=%04x", GetAReg(), GetBReg(), xreg)

	case 0x3B:
		s = "F$DelImg : Deallocate Image RAM blocks"
		p = F("beginBlock=%x numBlocks=%x processDesc=%04x", GetAReg(), GetBReg(), xreg)

	case 0x3F:
		s = "F$AllTsk : Allocate process Task number"
		p = F("processDesc=%04x", xreg)

	case 0x40:
		s = "F$DelTsk : Deallocate Task Number"
		p = F("proc_desc=%x", xreg)

	case 0x44:
		s = "F$DATLog : Convert DAT block/offset to Logical Addr"
		p = F("DatImageOffset=%x blockOffset=%x", GetBReg(), xreg)

	case 0x4B:
		s = "F$AllPrc : Allocate Process descriptor"

	case 0x4E:
		s = "F$FModul   : Find Module Directory Entry"
		mapping := GetMappingTask0(yreg)
		p = F("%q type %x name@ %x dat@ %x", Os9StringWithMapping(xreg, mapping), GetAReg(), xreg, yreg)

	case 0x4F:
		s = "F$MapBlk   : Map specific block"
		p = F("beginningBlock=%x numBlocks=%x", xreg, GetBReg())

	case 0x50:
		s = "F$ClrBlk : Clear specific Block"
		p = F("numBlocks=%x firstBlock=%x", GetBReg(), ureg)

	case 0x51:
		s = "F$DelRam : Deallocate RAM blocks"
		p = F("numBlocks=%x firstBlock=%x", GetBReg(), xreg)

	// IOMan:

	case 0x80:
		s = "I$Attach : Attach I/O Device"
		p = F("%04x='%s'", xreg, Os9String(xreg))

	case 0x81:
		s = "I$Detach : Detach I/O Device"
		p = F("%04x", ureg)

	case 0x82:
		s = "I$Dup    : Duplicate Path"
		p = F("$%x", GetAReg())

	case 0x83:
		s = "I$Create : Create New File"
		p = F("%04x='%s'", xreg, Os9String(xreg))

	case 0x84:
		s = "I$Open   : Open Existing File"
		p = F("%04x='%s'", xreg, Os9String(xreg))

	case 0x85:
		s = "I$MakDir : Make Directory File"
		p = F("%04x='%s'", xreg, Os9String(xreg))

	case 0x86:
		s = "I$ChgDir : Change Default Directory"
		p = F("%04x='%s'", xreg, Os9String(xreg))

	case 0x87:
		s = "I$Delete : Delete File"
		p = F("%04x='%s'", xreg, Os9String(xreg))

	case 0x88:
		s = "I$Seek   : Change Current Position"
		p = F("path=%x pos=%04x%04x", GetAReg(), xreg, ureg)

	case 0x89:
		s = "I$Read   : Read Data"
		p = F("path=%x buf=%04x size=%x", GetAReg(), xreg, yreg)

	case 0x8A:
		s = "I$Write  : Write Data"
		path := GetAReg()
		if IsTermPath(path) {
			p = PrintableMemory(xreg, yreg)
			if *FlagQuotedTerminal {
				fmt.Printf("[%q]", StrungMemory(xreg, yreg))
			} else if *FlagBracketTerminal {
				fmt.Printf("[%s]", StrungMemory(xreg, yreg))
			} else {
				fmt.Printf("%s", p)
			}
			CheckExpectation(p)
		}

	case 0x8B:
		s = "I$ReadLn : Read Line of ASCII Data"

	case 0x8C:
		s = "I$WritLn : Write Line of ASCII Data"
		{
			path := GetAReg()
			if IsTermPath(path) {
				str := PrintableStringThruEOS(xreg, yreg)
				if *FlagQuotedTerminal {
					fmt.Printf("%q ", str)
				} else if *FlagBracketTerminal {
					fmt.Printf("{%s}", str)
				} else {
					fmt.Printf("%s", str)
				}
				CheckExpectation(str)
			}
		}

	case 0x8D:
		s = "I$GetStt : Get Path Status"
		p = F("path=%x %x==%s", GetAReg(), GetBReg(), DecodeOs9GetStat(GetBReg()))

	case 0x8E:
		s = "I$SetStt : Set Path Status"
		p = F("path=%x %s", GetAReg(), DecodeOs9GetStat(GetBReg()))

	case 0x8F:
		s = "I$Close  : Close Path"
		p = F("path=%x", GetAReg())

	case 0x90:
		s = "I$DeletX : Delete from current exec dir"

	}
	if true || s == "" {
		s, _ = sym.SysCallNames[b]
	}
	return F("OS9$%02x <%s> {%s} #%d", b, s, p, Cycles), returns
}
