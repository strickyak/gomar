//go:build vdg

package emu

import (
	"github.com/strickyak/gomar/kitty"

	"bytes"
	"flag"
	"fmt"
	"os"
)

var FlagShowVDGScreen = flag.Bool("show_vdg_screen", false, "show VDG screens on stdout")
var FlagSemiGraphicsNotDirty = flag.Bool("semi-graphics-not-dirty", true, "don't let semigraphics (blinking cursor!) dirty it")
var FlagMag = flag.Int("mag", 8, "graphics pixel magnification")

type VDG struct {
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
}

func NewVDG() *VDG {
	n := uint(32 * 16)
	return &VDG{
		Dirty:   true,
		NumCols: 32,
		NumRows: 16,
		Addr:    0,
		Len:     n,
	}
}

func (o *VDG) DrawText() {
	for r := 0; r < o.NumRows; r++ {
		var bb bytes.Buffer
		bb.WriteByte('|')
		for c := 0; c < o.NumCols; c++ {
			x := PeekB(Word(o.Addr) + Word(c+r*o.NumCols))
			if 128 <= x {
				if (x & 0x0F) == 0x00 {
					// Use space if no spots are set
					x = ' '
				} else if (x & 0x0F) == 0x0F {
					// // Use letters [A-H] if all 4 spots are set
					// x = 'A' + (0x0F & (x >> 4)) - 8
					// Use octothorpe if all 4 spots are set
					x = '#'
				} else {
					// Use '+' if some but not all spots are set
					x = '+'
				}
			} else {
				x = 63 & x
				if x < 32 {
					x += 64
				}
			}
			bb.WriteByte(x)
		}
		bb.WriteByte('|')
		for c := 0; c < o.NumCols; c++ {
			x := PeekB(Word(o.Addr) + Word(c+r*o.NumCols))
			fmt.Fprintf(&bb, " %02x", x)
		}
		fmt.Printf("%s\n", bb.String())
	}
}
func (o *VDG) DrawPMode1() {
	const NUM_ROWS = 96
	const NUM_COLS = 128
	const PIXELS_PER_BYTE = 4
	const BITS_PER_PIXEL = 8 / PIXELS_PER_BYTE

	lookup := [4]struct{ R, G, B byte }{
		{50, 250, 50},
		{50, 50, 50},
		{250, 250, 50},
		{50, 50, 250},
	}

	var payload []byte
	for r := 0; r < NUM_ROWS; r++ {
		var bb []byte
		for c := 0; c < NUM_COLS/PIXELS_PER_BYTE; c++ {
			one_byte := PeekB(Word(o.Addr) + Word(c+r*NUM_COLS/PIXELS_PER_BYTE))
			for pix := PIXELS_PER_BYTE - 1; pix >= 0; pix-- {
				color := ((one_byte >> (pix * BITS_PER_PIXEL)) & (PIXELS_PER_BYTE - 1))
				red := lookup[color].R
				green := lookup[color].G
				blue := lookup[color].B
				for mag := 0; mag < *FlagMag; mag++ {
					bb = append(bb, red, green, blue)
				}
			}
		}
		for mag := 0; mag < *FlagMag; mag++ {
			payload = append(payload, bb...)
		}
	}
	kitty.Draw(os.Stdout, uint(*FlagMag*NUM_COLS), uint(*FlagMag*NUM_ROWS), payload)
}
func (o *VDG) Tick(step int64) {
	if !o.Dirty {
		return
	}
	if !*FlagShowVDGScreen {
		return
	}
	o.Dirty = false

	fmt.Printf("V=-=-=-=-=-=-=-=-=-= #%d (v:%d &%d p:%d r:%d m:%d ty:%d)\n", step, o.V, o.Addr, o.P, o.R, o.M, o.Ty)
	switch o.V {
	case 0:
		o.DrawText()
	case 4:
		o.DrawPMode1()
	default:
		fmt.Printf("[ video mode v=%d not handled yet ]\n", o.V)
	}

	fmt.Printf("V=-=-=-=-=-=-=-=-=-= #%d\n", step)
}

func (o *VDG) Poke(addr uint, longAddr uint, x byte) {
	if 0xFF00 <= addr {
		o.storePort(addr, x)
	}
	if 0xFFC0 <= addr && addr < 0xFFE0 {
		o.changeBit(addr)
	} else if o.Addr <= addr && addr < o.Addr+o.Len {
		if *FlagSemiGraphicsNotDirty && (x&128) != 0 && (x&15) == 15 {
			// dont set dirty bit
		} else {
			o.Dirty = true
		}
	}
}

func (o *VDG) storePort(addr uint, x byte) {
	L("vdg storePort %04x <- %02x nando", addr, x)
	o.Ports[addr-0xFF00] = x
}

func (o *VDG) changeBit(addr uint) {
	o.Dirty = true
	switch addr - 0xFFC0 {
	case 0x00:
		o.V &^= 1
	case 0x01:
		o.V |= 1
	case 0x02:
		o.V &^= 2
	case 0x03:
		o.V |= 2
	case 0x04:
		o.V &^= 4
	case 0x05:
		o.V |= 4
	case 0x06:
		o.Addr &^= 0x0200
	case 0x07:
		o.Addr |= 0x0200
	case 0x08:
		o.Addr &^= 0x0400
	case 0x09:
		o.Addr |= 0x0400
	case 0x0A:
		o.Addr &^= 0x0800
	case 0x0B:
		o.Addr |= 0x0800
	case 0x0C:
		o.Addr &^= 0x1000
	case 0x0D:
		o.Addr |= 0x1000
	case 0x0E:
		o.Addr &^= 0x2000
	case 0x0F:
		o.Addr |= 0x2000
	case 0x10:
		o.Addr &^= 0x4000
	case 0x11:
		o.Addr |= 0x4000
	case 0x12:
		o.Addr &^= 0x8000
	case 0x13:
		o.Addr |= 0x8000
	case 0x14:
		o.P &^= 1
	case 0x15:
		o.P |= 1
	case 0x16:
		o.R &^= 1
	case 0x17:
		o.R |= 1
	case 0x18:
		o.R &^= 2
	case 0x19:
		o.R |= 2
	case 0x1A:
		o.M &^= 1
	case 0x1B:
		o.M |= 1
	case 0x1C:
		o.M &^= 2
	case 0x1D:
		o.M |= 2
	case 0x1E:
		o.Ty &^= 1
	case 0x1F:
		o.Ty |= 1
	}
	enableRom = (o.Ty & 1) == 0
}
