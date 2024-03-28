//go:build !canvas

package emu

import (
	"bufio"
	"flag"
	"log"
	"os"
)

const INKEY1 = `
~~~~~~~~~~~~~~~~~~~~~~~92
~~~~~~~~~~~~~~~~~~~~~~~@
~~~~~~~~~~~~~~~~~~~~~~~10 rem bilbo frodo
~~~~~dir
~~~~~save "a"
~~~~~dir
~~~~~save "b"
~~~~~dir
~~~~~save "c"
~~~~~dir
~~~~~save "d"
~~~~~dir
~~~~~save "e"
~~~~~dir
~~~~~save "f"
~~~~~dir
~~~~~save "g"
~~~~~dir
~~~~~save "h"
~~~~~dir
~~~~~save "i"
~~~~~dir
~~~~~save "j"
~~~~~dir
~~~~~save "k"
~~~~~dir
~~~~~save "l"
~~~~~dir
~~~~~save "m"
~~~~~dir
~~~~~save "o"
~~~~~dir
~~~~~save "p"
~~~~~dir
~~~~~save "q"
~~~~~dir
~~~~~~~~save "r"
~~~~~~~~dir
~~~~~~~~save "s"
~~~~~~~~dir
~~~~~~~~save "t"
~~~~~~~~dir
~~~~~~~~save "u"
~~~~~~~~dir
~~~~~~~~save "v"
~~~~~~~~dir
~~~~~~~~save "w"
~~~~~~~~dir
~~~~~~~~save "x"
~~~~~~~~dir
~~~~~~~~save "y"
~~~~~~~~dir
~~~~~~~~save "z"
~~~~~~~~dir
`

const INKEY2 = `
~~~~~~~~~~~~~~~~~~~~~~~92
~~~~~~~~~~~~~~~~~~~~~~~@
~~~~~~~~~~~~~~~~~~~~~~~10 rem bilbo frodo
~~~~~dir
~~~~~dir
~~~~~dir
~~~~~dir
~~~~~dir
~~~~~dir
~~~~~dir
~~~~~dir
~~~~~dir
~~~~~dir
~~~~~dir
~~~~~dir
~~~~~dir
~~~~~dir
~~~~~dir
~~~~~dir
~~~~~dir
~~~~~dir
~~~~~dir
~~~~~dir
~~~~~dir
~~~~~dir
~~~~~dir
~~~~~dir
~~~~~dir
~~~~~dir
`

var flagN = flag.Bool("n", false, "Disable reading keystrokes from stdin")
var flagInkey = flag.String("inkey", INKEY1, "Inject keystrokes")

func InputRoutine(keystrokes chan<- byte) {
	s := *flagInkey
	for s != "" {
		ch := s[0]
		if ch == '\n' {
			ch = '\r'
		}
		if ch == '~' {
			ch = 0
		}
		keystrokes <- ch
		s = s[1:]
	}

	if *flagN {
		return
	}
	in := bufio.NewScanner(os.Stdin)
	for in.Scan() {
		for _, r := range in.Text() {
			keystrokes <- byte(r)
		}
		keystrokes <- '\r'
	}
	close(keystrokes)
	log.Fatal("Stdin ended")
}

func PublishVideoText() {}
