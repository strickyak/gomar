package emu

import (
// . "github.com/strickyak/gomar/gu"
)

var pia0 = &Pia{
	name: "Pia0",
	addr: 0xFF00,
}
var pia1 = &Pia{
	name: "Pia1",
	addr: 0xFF00,
}

const cocoVsyncPeriod = 900000 / 60 // assume 0.9 MHz, 60 Hz

var cocoVsyncCounter int = cocoVsyncPeriod

func CocoVsyncTick() {
	cocoVsyncCounter -= 1
	for cocoVsyncCounter <= 0 {
		cocoVsyncCounter += cocoVsyncPeriod
		pia0.controlA |= 0x80
	}
}

func CocoVsyncIrqEnabled() bool {
	return pia0.Ca1IrqEnabled()
}
func CocoVsyncIrqFiring() bool {
	return pia0.Ca1IrqFiring()
}
func CocoVsyncIrqEffective() bool {
	return CocoVsyncIrqEnabled() && CocoVsyncIrqFiring()
}
