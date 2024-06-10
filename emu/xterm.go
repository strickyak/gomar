package emu

import (
	"fmt"
	"io"
)

// Thanks: https://tintin.mudhalla.net/info/256color/

var Four2Six = []byte{0, 2, 3, 5}

func FgBgColorXterm(w io.Writer, fgR, fgG, fgB, bgR, bgG, bgB byte) {
	fgR = Four2Six[fgR]
	fgG = Four2Six[fgG]
	fgB = Four2Six[fgB]
	bgR = Four2Six[bgR]
	bgG = Four2Six[bgG]
	bgB = Four2Six[bgB]

	fg := 16 + 36*fgR + 6*fgG + fgB
	bg := 16 + 36*bgR + 6*bgG + bgB

	fmt.Fprintf(w, "\033[38;5;%dm\033[48;5;%dm", fg, bg)
}

func NormalColorXterm(w io.Writer) {
	fmt.Fprintf(w, "\033[0m")
}
