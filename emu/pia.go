package emu

import (
	. "github.com/strickyak/gomar/gu"
)

type Pia struct {
	name string
	addr uint

	dirA, dirB         uint
	outA, outB         uint
	controlA, controlB uint

	enableIrqA, enableIrqB bool
	irqA, irqB             bool
}

func (o *Pia) PortADataDirectionSelected() bool { return False(o.controlA & 0x04) }
func (o *Pia) PortBDataDirectionSelected() bool { return False(o.controlB & 0x04) }

func (o *Pia) Ca1IrqEnabled() bool { return True(o.controlA & 0x01) }
func (o *Pia) Cb1IrqEnabled() bool { return True(o.controlB & 0x01) }

func (o *Pia) Ca1IrqFiring() bool { return True(o.controlA & 0x80) }
func (o *Pia) Ca2IrqFiring() bool { return True(o.controlA & 0x40) }
func (o *Pia) Cb1IrqFiring() bool { return True(o.controlB & 0x80) }
func (o *Pia) Cb2IrqFiring() bool { return True(o.controlB & 0x40) }

func (o *Pia) Ca2IsIrqMode() bool { return 0x04 == (o.controlA & 0x14) }

func (o *Pia) Cb2IsIrqMode() bool          { return 0x20 == (o.controlB & 0x20) }
func (o *Pia) Cb2IsOutputStrobeMode() bool { return 0x20 == (o.controlB & 0x30) }
func (o *Pia) Cb2IsOutputBitMode() bool    { return 0x30 == (o.controlB & 0x30) }
