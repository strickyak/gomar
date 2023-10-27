//go:build f256flat

package emu

import (
	// "github.com/strickyak/gomar/display"
	// . "github.com/strickyak/gomar/gu"

	"flag"
	"io/ioutil"
	"log"
)

var FlagF256FlatBooter = flag.String("f256flat-booter", "", "The OS9Boot file for f256 level 1")

func DoDumpProcsAndPathsPrime() {
	log.Printf("TODO: DoDumpProcsAndPathsPrime")
}

var MmuTask byte // but not used in coco1.

const TraceMem = false // TODO: restore this some day.

func EmitHardware() {}
func InitHardware() {
	if *FlagF256FlatBooter != "" {
		bb, err := ioutil.ReadFile(*FlagF256FlatBooter)
		if err != nil {
			log.Fatalf("Cannot read -f256flat-booter %q: %v", *FlagF256FlatBooter, err)
		}
		for i, b := range bb {
			dest := i + 0x1000
			mem[dest] = b
		}
	}
	pcreg = W(0xFFFE)
}

func ExplainMMU() string             { return "" }
func DoExplainMmuBlock(i int) string { return "" }

func FireTimerInterrupt() {
	irqs_pending |= IRQ_PENDING
	Waiting = false
}

// B is fundamental func to get byte.  Hack register access into here.
func B(addr Word) byte {
	var z byte
	if AddressInDeviceSpace(addr) {
		z = GetIOByte(addr)
		L("GetIO %04x -> %02x : %c %c", addr, z, H(z), T(z))
		mem[addr] = z
	} else {
		z = mem[addr]
	}
	if TraceMem {
		L("\t\t\t\tGetB %04x -> %02x : %c %c", addr, z, H(z), T(z))
	}
	return z
}

func PokeB(addr Word, b byte) {
	if enableRom && 0x8000 <= addr && addr < 0xFF00 {
		L("ROM MODE inhibits write")
	} else {
		mem[addr] = b
	}
}

func PeekB(addr Word) byte {
	return mem[addr]
}

// PutB is fundamental func to set byte.  Hack register access into here.
func PutB(addr Word, x byte) {
	old := mem[addr]
	if enableRom && 0x8000 <= addr && addr < 0xFF00 {
		L("ROM MODE inhibits write")
	} else {
		mem[addr] = x

		if TraceMem {
			L("\t\t\t\tPutB %04x <- %02x (was %02x)", addr, x, old)
		}
		if AddressInDeviceSpace(addr) {
			PutIOByte(addr, x)
			L("PutIO %04x <- %02x (was %02x)", addr, x, old)
		}
	}
}

func WithMmuTask(task byte, fn func()) {
	fn()
}

func PutGimeIOByte(a Word, b byte) {
	// not used on coco1.
	// TODO -- but, coco13 has this line:
	// if 0xFF90 <= a && a < 0xFFC0 { PutGimeIOByte(a, b) }

	log.Printf("TODO UNKNOWN PutGimeIOByte address: 0x%04x <- 0x%02x", a, b)
	// log.Panicf("UNKNOWN PutGimeIOByte address: 0x%04x <- 0x%02x", a, b)
}

// TODO -- assume True for now.
func IsTermPath(path byte) bool {
	return true
}

// coco1 has no tasks, so ignore task.
func PeekWWithTask(addr Word, task byte) Word {
	return PeekW(addr)
}

// coco1 has no tasks, so ignore task.
func PeekBWithTask(addr Word, task byte) byte {
	return PeekB(addr)
}
