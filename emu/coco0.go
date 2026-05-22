//go:build coco0

package emu

import (
	// . "github.com/strickyak/gomar/gu"

	"flag"
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
	InitHardware13()
}

func ExplainMMU() string             { return "" }
func DoExplainMmuBlock(i int) string { return "" }

func UseRom(addr Word) bool {
	return 0x8000 <= addr && addr < 0xFF00
}
func RomOffset(addr Word) uint {
	return 0x7FFF & uint(addr)
}
func RamOffset(addr Word) uint {
	return uint(addr)
}
func UseExternalRomAssumingRom(addr Word) bool {
	return true
}

// B is fundamental func to get byte.  Hack register access into here.
func B(addr Word) byte {
	var z byte
	if AddressInDeviceSpace(addr) {
		z = GetIOByte(addr)
		L("GetIO %04x -> %02x : %c %c", addr, z, PrettyH(z), PrettyT(z))
		devmem[255&addr] = z
	} else if UseRom(addr) {
		o := RomOffset(addr)
		z = cartRom[o]
	} else {
		z = mem[addr]
	}
	if TraceMem {
		L("\t\t\t\tGetB %04x -> %02x : %c %c", addr, z, PrettyH(z), PrettyT(z))
	}
	return z
}

var WRITE_ROM_FAIL = flag.Bool("write_rom_fail", false, "Abort on writes to ROM")

func PokeB(addr Word, b byte) {
	if enableRom && 0x8000 <= addr && addr < 0xFF00 {
		if *WRITE_ROM_FAIL {
			log.Fatalf("PokeB: ROM MODE inhibits write (%04x <- %02x", addr, b)
		} else {
			L("PokeB: ROM MODE inhibits write (%04x <- %02x", addr, b)
		}
	} else {
		mem[addr] = b
	}
}

func PeekB(addr Word) byte {
	var z byte
	if UseRom(addr) {
		o := RomOffset(addr)
		z = cartRom[o]
	} else {
		// Not ROM so use RAM
		z = mem[addr]
	}
	return z
}

// PutB is fundamental func to set byte.  Hack register access into here.
func PutB(addr Word, x byte) {
	display.Poke(uint(addr), uint(addr), x)

	old := mem[addr]
	if enableRom && 0x8000 <= addr && addr < 0xFF00 {
		L("PutB: ROM MODE inhibits write (%04x <- %02x", addr, x)
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
	log.Printf("Coco1: UNKNOWN PutGimeIOByte address: 0x%04x <- 0x%02x", a, b)
}

func GetGimeIOByte(a Word) byte {
	log.Printf("Coco1: UNKNOWN GetGimeIOByte address: 0x%04x -> 0xFF", a)
	return 0xFF
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
