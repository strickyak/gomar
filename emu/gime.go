//go:build gime

package emu

import (
	"bytes"
	"fmt"
)

type Gime struct {
	Dirty   bool
	NumCols int
	NumRows int
	V       uint
	Addr    uint
	Len     uint
	P       uint
	R       uint
	M       uint
	Ty      uint
	Ports   [256]byte
	VDG     *VDG
}

func NewGime() Screen {
	n := uint(64 * 25)
	return &Gime{
		Dirty:   true,
		NumCols: 64,
		NumRows: 25,
		Addr:    0,
		Len:     n,
		VDG:     NewVDG(),
	}
}

func (o *Gime) Tick(step int64) {
	// fmt.Printf("=-=-=-=-=-=-=-=-=-= GIME TICK #%d\n", step)
	if (o.Ports[0x90] & 0x80) != 0 {
		o.VDG.Tick(step)
		return
	}
	if !o.Dirty {
		return
	}
	o.Dirty = false

	fmt.Printf("=-=-=-=-=-=-=-=-=-= #%d  gime %02x %02x \n", step, o.Ports[0x90], o.Ports[0x91])
	for r := 0; r < o.NumRows; r++ {
		var bb bytes.Buffer
		bb.WriteByte('|')
		for c := 0; c < o.NumCols; c++ {
			x := mem[o.Addr+uint(c+r*o.NumCols)]
			if x < 32 || 126 < x {
				x = '~'
			}
			bb.WriteByte(x)
		}
		bb.WriteByte('|')
		for c := 0; c < o.NumCols; c++ {
			x := mem[o.Addr+uint(c+r*o.NumCols)]
			fmt.Fprintf(&bb, " %02x", x)
		}
		fmt.Printf("%s\n", bb.String())
	}
	fmt.Printf("=-=-=-=-=-=-=-=-=-= #%d\n", step)
}

func (o *Gime) Poke(addr uint, longAddr uint, x byte) {
	if 0xFF00 <= addr {
		o.storePort(addr, x)
	}
	o.VDG.Poke(addr, longAddr, x)

	if 0xFF80 <= addr && addr < 0xFFC0 {
		o.Dirty = true
		switch addr {
		case 0xFF9D:
			// Start of Video Ram, longAddr, bits 18..11
			o.Addr &^= (255 << 11)  // clear the 8 bits
			o.Addr |= uint(x) << 11 // shift x into the 8 bits
			fmt.Printf("GIME POKE %x %x %x -> %06x\n", addr, longAddr, x, o.Addr)
			o.Dirty = true
		case 0xFF9E:
			// Start of Video Ram, longAddr, bits 10..3
			o.Addr &^= (255 << 3)  // clear the 8 bits
			o.Addr |= uint(x) << 3 // shift x into the 8 bits
			fmt.Printf("GIME POKE %x %x %x -> %06x\n", addr, longAddr, x, o.Addr)
			o.Dirty = true
		}
	} else if (o.Ports[0x90] & 0x80) == 0 {
		// Coco3 modes
		if o.Addr <= longAddr && longAddr < o.Addr+1024 {
			// for now, just the text screen
			o.Dirty = true
		}
	}
}

func (o *Gime) storePort(addr uint, x byte) {
	L("Gime storePort %04x <- %02x nando", addr, x)
	o.Ports[addr-0xFF00] = x
}
