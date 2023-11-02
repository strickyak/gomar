//go:build trace
// +build trace

package emu

import (
	// "github.com/strickyak/gomar/sym"
	"github.com/strickyak/gomar/listings"

	"bytes"
	"fmt"
	"log"
	"strings"
)

var beenThere [0x10000]byte

/* max. bytes of instruction code per trace line */
const kMaximumBytesPerOpcode = 4
const kNoEffAddr = 0xFFFFFFFF

var effAddr EA = kNoEffAddr
var effByte int = -1
var effWord int = -1

func TraceByte(addr EA, x byte) {
	effAddr = addr
	effByte = int(x)
}

func TraceWord(addr EA, x Word) {
	effAddr = addr
	effWord = int(x)
}

var dis_length Word // do we need this?

func Dis_len(n Word) {
	dis_length = n
}
func Dis_len_incr(n Word) {
	dis_length += n
}

func Trace() {
	var buf bytes.Buffer
	wh := where(pcreg_prev)
	// oldnew would be improved with Memory Block.

	newOp := PeekB(pcreg_prev)
	oldnew := 'N'
	if beenThere[pcreg_prev] == newOp {
		oldnew = 'o'
	} else {
		beenThere[pcreg_prev] = newOp
	}

	Z(&buf, "%s%c %04x:", wh, oldnew, pcreg_prev)

	var ilen int
	if dis_length != 0 {
		ilen = int(dis_length)
	} else {
		ilen = int(pcreg - pcreg_prev)
		if ilen < 0 {
			ilen = -ilen
		}
	}
	for i := Word(0); i < kMaximumBytesPerOpcode; i++ {
		if int(i) < ilen {
			Z(&buf, "%02x", B(pcreg_prev+i)) // two hex chars
		} else {
			Z(&buf, "  ") // two spaces
		}
	}

	module, offset := MemoryModuleOf(pcreg_prev)

	text := ""
	if module != "" {
		moduleLower := strings.ToLower(module)
		text = listings.Lookup(moduleLower, uint(offset), func() {
			// TODO -- why??? // *FlagTraceAfter = 1
		})
	}

	eff := ""
	if effAddr < 0x10000 {
		if effByte != -1 {
			eff = fmt.Sprintf(" %04x:%02x", effAddr, byte(effByte))
		} else if effWord != -1 {
			eff = fmt.Sprintf(" %04x:%04x", effAddr, Word(effWord))
		}
	}

	Z(&buf, " {%-5s %-17s}  ", dinst.String(), dops.String())
	log.Printf("%s%s {{%s}} %s", buf.String(), Regs(), text, eff)
	log.Printf("")
	dis_length = 0

	if pcreg < pcreg_prev || pcreg > pcreg_prev+4 {
		log.Printf("")
		log.Printf("    %s", ExplainMMU())
		log.Printf("")
	}

	wh = strings.Trim(wh, " ")
	for _, w := range Watches {
		if wh == w.Where {
			var val Word
			switch w.Register {
			case "d":
				val = dreg
			}
			log.Printf("@WATCH@ %s == %04x == %q", w.Where, val, w.Message)
		}
	}
	effAddr = kNoEffAddr
	effByte = -1
	effWord = -1
}

func Finish() {
	L("Finish:")
	L("Cycles: #%d", Cycles)
	L("")
	DoDumpAllMemory()
	L("")
	DoDumpAllMemoryPhys()
	L("")
	L("Cycles: #%d", Cycles)
}

func where(addr Word) string {
	//if Level == 2 {
	name, offset := MemoryModuleOf(addr)
	if name != "" {
		return F("%q+%04x ", name, offset)
	} else {
		return "\"\" "
	}
	//}

	/*
		// if Level == 1 ...
		// TODO -- did this ever work for Level 1?
		var buf bytes.Buffer

		start := W(0x26)
		limit := W(0x28)

		for i := start; i < limit; i += 4 {
			mod := W(i)
			if mod != 0 {
				size := W(mod + 2)
				if mod < addr && addr < mod+size {
					cp := mod + W(mod+4)
					for {
						b := B(cp)
						ch := 127 & b
						if '!' <= ch && ch <= '~' {
							buf.WriteByte(ch)
						}
						if (b & 128) != 0 {
							Z(&buf, ",%04x ", addr-mod)
							return buf.String()
						}
						cp++
					}
				}
			}
		}
		return "? "
	*/
}

func Dis_inst(inst string, reg string, cyclecount int64) {
	dinst.Reset()
	dops.Reset()
	dinst.WriteString(inst)
	dinst.WriteString(reg)
	Cycles += cyclecount
}

func Dis_inst_cat(inst string, cyclecount int64) {
	dinst.WriteString(inst)
	Cycles += cyclecount
}

func Dis_ops(part1 string, part2 string, cyclecount int64) {
	dops.WriteString(part1)
	dops.WriteString(part2)
	Cycles += cyclecount
}

var reg_for_da_reg = []string{"d", "x", "y", "u", "s", "pc", "?", "?", "a", "b", "cc", "dp", "?", "?", "?", "?"}

func Dis_reg(b byte) {
	dops.WriteString(reg_for_da_reg[(b>>4)&0xf])
	dops.WriteString(",")
	dops.WriteString(reg_for_da_reg[b&0xf])
}

func DumpAllMemory()    { DoDumpAllMemory() }
func DumpPageZero()     { DoDumpPageZero() }
func DumpProcesses()    { DoDumpProcesses() }
func DumpAllPathDescs() { DoDumpAllPathDescs() }

//func LogIO(f string, args ...interface{}) {
//	L(f, args...)
//}

// Call this before each instruction until it returns false.
func EarlyAction() bool {
	// OS9 boots with PC in the first half of memory space.
	// When it jumps into the higher half, it jumps into modules.
	if pcreg > 0x8000 {
		DumpAllMemory()
		InitialModules = ScanRamForOs9Modules()
		return false
	}
	return true
}
