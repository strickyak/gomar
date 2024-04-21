//go:build gime

package emu

import (
	"bytes"
	"flag"
	"fmt"

	. "github.com/strickyak/gomar/gu"
)

var FlagShowGIMEScreen = flag.Bool("show_gime_screen", false, "show GIME screens on stdout")

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
	if !*FlagShowGIMEScreen {
		return
	}
	o.Dirty = false

	p := GetCocoDisplayParams()

	if p.Graphics {
		fmt.Printf("=-=-=-=-=-=-=-=-=-= #%d  GRAPHICS %v\n", step, p)
	} else {
		nc := p.AlphaCharsPerRow
		cShift := Cond(p.AlphaHasAttrs, 1, 0)

		fmt.Printf("=-=-=-=-=-=-=-=-=-= #%d  gime %02x %02x \n", step, o.Ports[0x90], o.Ports[0x91])
		for r := 0; r < o.NumRows; r++ {
			var bb bytes.Buffer
			bb.WriteByte('|')
			for c := 0; c < nc; c++ {
				x := mem[o.Addr+(uint(c+r*nc)<<cShift)]
				if x < 32 || 126 < x {
					x = '~'
				}
				bb.WriteByte(x)
			}
			bb.WriteByte('|')
			for c := 0; c < nc; c++ {
				x := mem[o.Addr+(uint(c+r*nc)<<cShift)]
				fmt.Fprintf(&bb, " %02x", x)
			}
			fmt.Printf("%s\n", bb.String())
		}
		fmt.Printf("=-=-=-=-=-=-=-=-=-= #%d\n", step)

	}
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
		if o.Addr <= longAddr && longAddr < o.Addr+8*1024 {
			// for now, just the text screen
			o.Dirty = true
		}
	}
}

func (o *Gime) storePort(addr uint, x byte) {
	L("Gime storePort %04x <- %02x nando", addr, x)
	o.Ports[addr-0xFF00] = x
}

////////////////////////////////////////////////////////
////////////////////////////////////////////////////////
////////////////////////////////////////////////////////

type CocoDisplayParams struct {
	BasicText           bool // 32x16 at 0x400
	Gime                bool // else use VDG
	Graphics            bool // else use Alpha
	AttrsIfAlpha        bool // if every other byte is attrs
	VirtOffsetAddr      int  // Start of data.
	HorzOffsetAddr      int
	VirtScroll          int
	LinesPerField       int
	LinesPerCharRow     int
	Monochrome          bool
	HRES                int
	CRES                int
	HVEN                bool
	GraphicsBytesPerRow int
	GraphicsColorBits   int
	AlphaCharsPerRow    int
	AlphaHasAttrs       bool
	ColorMap            [16]byte
}

func ExplainColor(b byte) string {
	return F("rgb=$%02x=(%x,%x,%x)", b&63,
		((b&0x20)>>4)|((b&0x04)>>2),
		((b&0x10)>>3)|((b&0x02)>>1),
		((b&0x08)>>2)|((b&0x01)>>0))
}

/*
HRES:
	http://users.axess.com/twilight/sock/gime.html
Horizontal resolution using graphics:
000=16 bytes per row
001=20 bytes per row
010=32 bytes per row
011=40 bytes per row
100=64 bytes per row
101=80 bytes per row
110=128 bytes per row
111=160 bytes per row

When using text:
0x0=32 characters per row
0x1=40 characters per row
1x0=64 characters per row
1x1=80 characters per row
*/

var GraphicsBytesPerRowHRES = []int{16, 20, 32, 40, 64, 80, 128, 160}
var AlphaCharsPerRowHRES = []int{32, 40, 32, 40, 64, 80, 64, 80}

/*
CRES:
	http://users.axess.com/twilight/sock/gime.html
Color Resolution using graphics:
00=2 colors (8 pixels per byte)
01=4 colors (4 pixels per byte)
10=16 colors (2 pixels per byte)
11=Undefined (would have been 256 colors)

When using text:
x0=No color attributes
x1=Color attributes enabled
*/

var GraphicsColorBitsCRES = []int{1, 2, 4, 8}
var AlphaHasAttrsCRES = []bool{false, true, false, true}

var FF92Bits = []string{
	"?", "?", "TimerIRQ", "HorzIRQ", "VertIRQ", "SerialIRQ", "KbdIRQ", "CartIRQ"}
var FF93Bits = []string{
	"?", "?", "TimerFIRQ", "HorzFIRQ", "VertFIRQ", "SerialFIRQ", "KbdFIRQ", "CartFIRQ"}

var GimeLinesPerField = []int{192, 200, 210, 225}
var GimeLinesPerCharRow = []int{1, 2, 3, 8, 9, 10, 12, -1}

func GetCocoDisplayParams() *CocoDisplayParams {
	a := PeekB(0xFF98)
	b := PeekB(0xFF99)
	c := PeekB(0xFF9C)
	d := PeekB(0xFF9F)
	z := &CocoDisplayParams{
		BasicText:       *FlagBasicText,
		Gime:            true,
		Graphics:        (a>>7)&1 != 0,
		AttrsIfAlpha:    (a>>6)&1 != 0,
		VirtOffsetAddr:  int(HiLo(PeekB(0xFF9D), PeekB(0xFF9E))) << 3,
		HorzOffsetAddr:  int(d & 127),
		VirtScroll:      int(c & 15),
		LinesPerField:   GimeLinesPerField[(b>>5)&3],
		LinesPerCharRow: GimeLinesPerCharRow[a&7],
		Monochrome:      (a>>4)&1 != 0,
		HRES:            int((b >> 2) & 7),
		CRES:            int(b & 3),
		HVEN:            d>>7 != 0,
	}
	if z.Graphics {
		z.GraphicsBytesPerRow = GraphicsBytesPerRowHRES[z.HRES]
		z.GraphicsColorBits = GraphicsColorBitsCRES[z.CRES]
	} else {
		z.AlphaCharsPerRow = AlphaCharsPerRowHRES[z.HRES]
		z.AlphaHasAttrs = AlphaHasAttrsCRES[z.CRES]
	}
	for i := 0; i < 16; i++ {
		z.ColorMap[i] = PeekB(0xFFB0 + Word(i))
	}
	return z
}

func DumpGimeStatus() {
	for i := Word(0); i < 16; i += 4 {
		L("GIME/palette[%x..%x]: %s %s %s %s", i, i+3,
			ExplainColor(PeekB(0xFFB0+i)),
			ExplainColor(PeekB(0xFFB1+i)),
			ExplainColor(PeekB(0xFFB2+i)),
			ExplainColor(PeekB(0xFFB3+i)))
	}
	L("GIME/CpuSpeed: %x", PeekB(0xFFD9))
	L("GIME/MmuEnable: %v", PeekB(0xFF90)&0x40 != 0)
	L("GIME/MmuTask: %v; clock rate: %v", MmuTask, 0 != (PeekB(0xFF91)&0x40))
	L("GIME/IRQ bits: %s", ExplainBits(PeekB(0xFF92), FF92Bits))
	L("GIME/FIRQ bits: %s", ExplainBits(PeekB(0xFF93), FF93Bits))
	L("GIME/Timer=$%x", HiLo(PeekB(0xFF94), PeekB(0xFF95)))
	b := PeekB(0xFF98)
	L("GIME/GraphicsNotAlpha=%x AttrsIfAlpha=%x Artifacting=%x Monochrome=%x 50Hz=%x LinesPerCharRow=%x=%d.",
		(b>>7)&1,
		(b>>6)&1,
		(b>>5)&1,
		(b>>4)&1,
		(b>>3)&1,
		(b & 7),
		GimeLinesPerCharRow[b&7])
	b = PeekB(0xFF99)
	L("GIME/LinesPerField=%x=%d. HRES=%x CRES=%x",
		(b>>5)&3,
		GimeLinesPerField[(b>>5)&3],
		(b>>2)&7,
		b&3)

	b = PeekB(0xFF9C)
	L("GIME/Virt Scroll (alpha) = %x", b&15)
	L("GIME/VirtOffsetAddr=$%05x",
		uint64(HiLo(PeekB(0xFF9D), PeekB(0xFF9E)))<<3)
	b = PeekB(0xFF9F)
	L("GIME/HVEN=%x HorzOffsetAddr=%x", (b >> 7), b&127)
	L("GIME/GetCocoDisplayParams = %#v", *GetCocoDisplayParams())
}
