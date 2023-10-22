package emu

// See credits.go

import (
	//NODISPLAY// "github.com/strickyak/gomar/display"
	. "github.com/strickyak/gomar/gu"
	"github.com/strickyak/gomar/sym"

	"bufio"
	"bytes"
	"flag"
	"fmt"
	//"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var FlagTerm = flag.String("term", "Term", "name of terminal device")
var FlagLinkerMapFilename = flag.String("map", "", "")
var FlagBootImageFilename = flag.String("boot", "", "")
var FlagLoadmFilename = flag.String("loadm", "", "")
var FlagCartFilename = flag.String("cart", "", "")
var FlagRom8000Filename = flag.String("rom_8000", "", "")
var FlagRomA000Filename = flag.String("rom_a000", "", "")
var FlagKernelFilename = flag.String("kernel", "", "")
var FlagDiskImageFilename = flag.String("disk", "../_disk_", "")
var FlagMaxSteps = flag.Int64("max", 0, "")
var FlagClock = flag.Int64("clock", 5*1000*1000, "")
var FlagBasicText = flag.Bool("basic_text", false, "")
var FlagUserResetVector = flag.Bool("use_reset_vector", false, "")
var FlagQuotedTerminal = flag.Bool("quoted_terminal", false, "quote terminal output for debugging")

var FlagWatch = flag.String("watch", "", "Sequence of module:addr:reg:message,...")
var FlagTriggerPc = flag.Uint64("trigger_pc", 0xC00D, "")
var FlagTriggerOp = flag.Uint64("trigger_op", 0x17, "")
var FlagTraceOnOS9 = flag.String("trigger_os9", "", "")
var RegexpTraceOnOS9 *regexp.Regexp

type Watch struct {
	Where    string
	Register string
	Message  string
}

var Watches []*Watch

func CompileWatches() {
	for _, s := range strings.Split(*FlagWatch, ",") {
		if s != "" {
			v := strings.Split(s, ":")
			if len(v) != 3 {
				log.Fatalf("Watch was %q, split on colon, len was %d, want 3", v, len(v))
			}
			Watches = append(Watches, &Watch{
				Where:    v[0],
				Register: v[1],
				Message:  v[2],
			})
		}
	}
}

const IRQ_FREQ = (500 * 1000)

const paranoid = false // Do paranoid checks.
const hyp = true       // Use hyperviser code.

// F is for FORMAT (i.e. fmt.Sprintf)
func F(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

// L is for LOG (i.e. log.Printf)
func L(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// Z is for Printf to Buffer (i.e. fmt.Fprintf)
func Z(w *bytes.Buffer, format string, args ...interface{}) {
	fmt.Fprintf(w, format, args...)
}

var SymbolLine = regexp.MustCompile(`^Symbol: _(.*) = ([0-9A-F]+)`)

type LinkerRec struct {
	Sym  string
	Addr int
}

type LinkerMapType []*LinkerRec

func (m LinkerMapType) Len() int           { return len(m) }
func (m LinkerMapType) Swap(a, b int)      { m[a], m[b] = m[b], m[a] }
func (m LinkerMapType) Less(a, b int) bool { return m[a].Addr < m[b].Addr }

var LinkerMap LinkerMapType

func ReadLinkerMap() {
	if *FlagLinkerMapFilename == "" {
		return
	}
	fd, err := os.Open(*FlagLinkerMapFilename)
	if err != nil {
		log.Fatalf("cannot open %q: %v", *FlagLinkerMapFilename, err)
	}
	defer fd.Close()
	sc := bufio.NewScanner(fd)
	for sc.Scan() {
		s := sc.Text()
		m := SymbolLine.FindStringSubmatch(s)
		if m != nil {
			sym := m[1]
			hex := m[2]
			addr, err := strconv.ParseUint(hex, 16, 16)
			if err != nil {
				log.Fatalf("cannot ParseUint hex: %q: %v", hex, err)
			}
			rec := &LinkerRec{
				Sym:  sym,
				Addr: int(addr),
			}
			LinkerMap = append(LinkerMap, rec)
		}
	}
	sort.Sort(LinkerMap)
}

func CoreDump(filename string) {
	fd, err := os.Create(filename)
	if err != nil {
		log.Fatalf("cannot create %q: %v", filename, err)
	}
	w := bufio.NewWriter(fd)
	for i := 0; i < 0x10000; i++ {
		w.WriteByte(B(Word(i)))
		// w.WriteByte(EA(i).GetB())
	}
	for i := DRegEA; i <= PCRegEA; i++ {
		word := EA(i).GetW()
		w.WriteByte(byte(word >> 8))
		w.WriteByte(byte(word >> 0))
	}
	w.WriteByte(CCRegEA.GetB())
	w.WriteByte(DPRegEA.GetB())
	w.Flush()
	fd.Close()
}

func FatalCoreDump() {
	const NAME = "/tmp/coredump09"

	ReadLinkerMap()
	CoreDump(NAME)

	fmt.Printf(" ... Wrote %q ... Begin Frame Chain\n", NAME)

	fp := EA(URegEA.GetW())
	codeOffset := (int(fp)/0x2000)*0x2000 + 0x2000
	p := EA(SRegEA.GetW())
	fmt.Printf("S: $%04x  U: $%04x\n", p, fp)
	gap := int(fp) - int(p)
	firstGap := true
	for 0 <= gap && gap <= 64 {
		fmt.Printf("\n@$%04x: ", int(p))
		for p < fp {
			fmt.Printf("%02x ", EA(p).GetB())
			p += 1
		}

		if false && firstGap {
			firstGap = false
		} else if LinkerMap != nil {
			fp2 := EA(fp + 2)
			pc := fp2.GetW()

			found := sort.Search(len(LinkerMap), func(i int) bool {
				return (codeOffset+LinkerMap[i].Addr > int(pc))
			})
			if found > 0 {
				prev := LinkerMap[found-1]
				fmt.Printf("\n ............ pc=$%x is $%x + %q=$%x",
					pc,
					int(pc)-codeOffset+prev.Addr,
					prev.Sym,
					codeOffset+prev.Addr)
			} else {
				fmt.Printf("\n ............ pc=$%x is too low", pc)
			}
		}

		fp = EA(fp.GetW())
		gap = int(fp) - int(p)
	}
	fmt.Printf("\nEnd Frame Chain\n")

	log.Fatalf("EMULATOR CORE DUMPED: %q", NAME)
}

func TfrReg(b byte) EA {
	if 6 == b || b == 7 || b > 11 {
		log.Panicf("Bad TfrReg byte: 0x%x", b)
	}
	return DRegEA + EA(b)
}

//NODISPLAY//var CocodChan chan *display.CocoDisplayParams
//NODISPLAY//var Disp *display.Display

var fdump int
var Cycles int64

var Os9Description = make(map[int]string) // Describes OS9 kernel call at this big stack addr.

/* 6809 registers */
var ccreg, dpreg byte
var xreg, yreg, ureg, sreg, pcreg Word
var dreg Word

var iflag byte /* flag to indicate prebyte $10 or $11 */
var ireg byte  /* Instruction register */
var pcreg_prev Word

var mem [0x40 * 0x2000]byte

var ixregs = []*Word{&xreg, &yreg, &ureg, &sreg}

var idx byte

/* disassembled instruction buffer */
var dinst bytes.Buffer

/* disassembled operand buffer */
var dops bytes.Buffer

var Waiting bool
var irqs_pending byte

var instructionTable [256]func()

//////////////////////////////////////////////////////////////

const NMI_PENDING = CC_ENTIRE /* borrow this bit */
const IRQ_PENDING = CC_INHIBIT_IRQ
const FIRQ_PENDING = CC_INHIBIT_FIRQ

const CC_INHIBIT_IRQ = 0x10
const CC_INHIBIT_FIRQ = 0x40
const CC_ENTIRE = 0x80

const VECTOR_IRQ = 0xFFF8
const VECTOR_FIRQ = 0xFFF6
const VECTOR_NMI = 0xFFFC

// 200 = 0x80 = CLEAR; 033=ESC;  201=F1, 202=F2, 203=BREAK
// 204=up 205=dn 206=left 207=right
const KB_NORMAL = "@ABCDEFGHIJKLMNOPQRSTUVWXYZ\204\205\206\207 0123456789:;,-./\r\200\033\000\000\201\202\000"
const KB_SHIFT = "`abcdefghijklmnopqrstuvwxyz____ 0!\"#$%&'()*+<=>?___..__."
const KB_CTRL = `.................................|.~...^[]..{_}\........`

func keypress(probe byte, ch byte) byte {
	shifted, controlled := false, false
	sense := byte(0)
	probe = ^probe
	for j := uint(0); j < 8; j++ {
		for i := uint(0); i < 7; i++ {
			if KB_NORMAL[i*8+j] == ch {
				if (byte(1<<j) & probe) != 0 {
					sense |= 1 << i
				}
			} else if KB_SHIFT[i*8+j] == ch && ch != '.' {
				if (byte(1<<j) & probe) != 0 {
					sense |= byte(1 << i)
				}
				shifted = true
			} else if KB_CTRL[i*8+j] == ch && ch != '.' {
				if (byte(1<<j) & probe) != 0 {
					sense |= byte(1 << i)
				}
				controlled = true
			}
		}
	}
	if shifted && (probe&0x80) != 0 {
		sense |= 0x40 // Shift key.
	}
	if controlled && (probe&0x10) != 0 {
		sense |= 0x40 // Ctrl key.
	}
	log.Printf("keypress: probe %x char %x sense %x shifted %v", probe, ch, sense, shifted)
	return ^sense
}

var prev_disk_command byte
var disk_command byte
var disk_offset int64
var disk_drive byte
var disk_side byte
var disk_sector byte
var disk_track byte
var disk_status byte
var disk_data byte
var disk_control byte
var disk_fd *os.File
var disk_stuff [256]byte
var zero_disk_stuff [256]byte
var disk_sector_0 [256]byte
var disk_dd_fmt byte // Offset 16.
var disk_i Word

var kbd_ch byte
var kbd_probe byte
var kbd_cycle Word

func MaybeGetChar() byte {
	return 0
}

func inkey(keystrokes <-chan byte) byte {
	select {
	case _ch, _ok := <-keystrokes:
		if _ok {
			if Level == 2 {
				// In Level2, swap case.
				if 'A' <= _ch && _ch <= 'Z' {
					return _ch + 32
				} else if 'a' <= _ch && _ch <= 'z' {
					return _ch - 32
				} else {
					return _ch
				}
			} else {
				return _ch
			}
		} else {
			log.Printf("EXIT: inkey gets end of channel")
			Finish()
			os.Exit(0)
			return 0
		}
	default:
		return 0
	}
}

/*
func printableChar(ch byte) string {
	if ' ' <= ch && ch <= '~' {
		return string(rune(ch))
	} else {
		return F("{%d}", ch)
	}
}
*/

func H(ch byte) byte {
	ch &= 0x7F
	if 32 <= ch && ch <= 126 {
		return ch
	} else {
		return ' '
	}
}
func T(ch byte) byte {
	if ch&128 != 0 && 128+32 <= ch && ch <= 128+126 {
		return '+'
	} else {
		return ' '
	}
}

// Now follow the posbyte addressing modes. //

const (
	AttachModeDev byte = iota
	AttachModeRead
	AttachModeWrite
	AttachModeReadWrite
)

func Os9HypervisorCall(syscall byte) bool {
	handled := false
	L("Hyp::%x", syscall)
	switch Word(syscall) {
	case sym.I_Attach:
		{
			access_mode := GetAReg()
			dev_name := Os9String(xreg)
			L("Hyp I_Attach %q mode %d", dev_name, access_mode)
		}
	case sym.I_ChgDir:
	case sym.I_Close:
	case sym.I_Create:
	case sym.I_Delete:
	case sym.I_DeletX:
	case sym.I_Detach:
		{
			dev_table := ureg
			L("Hyp I_Detach %04x", dev_table)
		}
	case sym.I_Dup:
		L("Hyp I_Dup %d.", GetAReg())
	case sym.I_GetStt:
	case sym.I_MakDir:
	case sym.I_Open:
	case sym.I_Read:
	case sym.I_ReadLn:
	case sym.I_Seek:
	case sym.I_SetStt:
	case sym.I_Write:
	case sym.I_WritLn:
	}
	if handled {
		sreg += 10
		PullWord(&pcreg)
		pcreg++
	}
	return handled
}

const MaxInt64 = 0x7FFFFFFFFFFFFFFF

var PrevBasicText []byte

func EqualBytes(a, b []byte) bool {
	n := len(a)
	if n != len(b) {
		return false
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func ShowBasicText() {
	start := Word(sam.Fx) << 9
	limit := start + 0x200

	Logd("CYCLE: #%d .... at $%04x", Cycles, start)
	for y := start; y < limit; y += 32 {
		text := make([]byte, 32)
		for x := Word(0); x < 32; x++ {
			b := PeekB(x+y) & 63
			if b < 32 {
				b += 64
			}
			text[x] = b
		}
		Logd("TEXT: %q", text)
	}
}
