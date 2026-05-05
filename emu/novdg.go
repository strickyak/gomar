//go:build !vdg

package emu

type VDG struct{}

func NewVDG() *VDG {
	return &VDG{}
}
func (o *VDG) DrawText()                             {}
func (o *VDG) DrawPMode1()                           {}
func (o *VDG) Tick(step int64)                       {}
func (o *VDG) Poke(addr uint, longAddr uint, x byte) {}
