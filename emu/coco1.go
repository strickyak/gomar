//go:build coco1
// +build coco1

package emu

import (
	// . "github.com/strickyak/gomar/gu"

	"log"
)

func DoDumpProcsAndPathsPrime() {
	log.Printf("TODO: DoDumpProcsAndPathsPrime")
}

var MmuTask byte // but not used in coco1.

const TraceMem = false // TODO: restore this some day.

func EmitHardware() {}
func InitHardware() {
	display = NewVDG()
}

func ExplainMMU() string             { return "" }
func DoExplainMmuBlock(i int) string { return "" }

func UseExternalRomAssumingRom(addr Word) bool {
	return 0xC000 <= addr && addr < 0xFFF0
}
func InternalRomOffset(addr Word) uint {
	// In coco1 & coco2, there is only 16K internal,
	// in the lower half of 32K.
	return 0x3FFF & uint(addr)
}
func ExternalRomOffset(addr Word) uint {
	// In coco1 & coco2, there is only 16K external,
	// but it is in the upper half of 32K.
	return 0x4000 + 0x3FFF&uint(addr)
}
func RamOffset(addr Word) uint {
	return SimpleRamOffset(addr)
}

// Coco1,2 can swap ram pages
func SimpleRamOffset(addr Word) uint {
	if sam.P1RamSwap != 0 {
		return 0x8000 ^ uint(addr)
	} else {
		return uint(addr)
	}
}

// B is fundamental func to get byte.  Hack register access into here.
func B(addr Word) byte {
	var z byte
	if AddressInDeviceSpace(addr) {
		z = GetIOByte(addr)
		L("GetIO %04x -> %02x : %c %c", addr, z, H(z), T(z))
		devmem[255&addr] = z
	} else {
		z = mem[RamOffset(addr)]
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

func GetCocoDisplayParams() *CocoDisplayParams {
	z := &CocoDisplayParams{
		BasicText:       *FlagBasicText,
		Gime:            false,
		Graphics:        false,
		AttrsIfAlpha:    false,
		VirtOffsetAddr:  0x8000, // TODO
		HorzOffsetAddr:  0x80,   // TODO
		VirtScroll:      0x0F,   // TODO
		LinesPerField:   8,      // TODO
		LinesPerCharRow: 8,      // TODO
		Monochrome:      true,
		HRES:            0,     // TODO
		CRES:            0,     // TODO
		HVEN:            false, // TODO
	}
	for i := 0; i < 16; i++ {
		z.ColorMap[i] = byte(i) // TODO
	}
	return z
}

// TODO

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

func InitializeVectors() {
	PutW(0xFFF2, 0x0100) // SWI3
	PutW(0xFFF4, 0x0103) // SWI2
	PutW(0xFFFA, 0x0106) // SWI
	PutW(0xFFFC, 0x0109) // NMI
	PutW(0xFFF8, 0x010C) // IRQ
	PutW(0xFFF6, 0x010F) // FIRQ
}
