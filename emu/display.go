package emu

import (
	"flag"
	"sync"
)

var FONT = flag.String("font", "/home/strick/go/src/golang.org/x/image/font/gofont/ttfs/Go-Mono.ttf", ".ttf font file")
var SIZE = flag.Float64("fontsize", 25, "font size")

// Emulator hard window size.
const WIDTH = 1280 + 20
const HEIGHT = 800 + 20

// Global vars describing mouse state.
var MouseX, MouseY float64 // 0 to 1
var MouseDown bool
var MouseMutex sync.Mutex

var display Screen

/*
type Display struct {
	Dirty   bool
	Mem     []byte
	Rows    [][]byte
	NumRows int
	NumCols int
	Cocod   <-chan *CocoDisplayParams
	Inkey   chan<- byte
	Sam     *Sam
	PeekB   func(addr int) byte
	x, y    int
	ctrl    bool
}
*/

type Sam struct {
	Fx        byte
	Mx        byte // Memory size for SAM.
	Rx        byte // Clock speed for SAM.
	Vx        byte
	TyAllRam  bool // TY bit
	P1RamSwap byte // P1 bit
}

type Screen interface {
	Tick(step int64)
	Poke(addr uint, longAddr uint, x byte)
}
