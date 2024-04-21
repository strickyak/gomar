//go:build !canvas

package emu

import (
	"bufio"
	"flag"
	"io/ioutil"
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

var flagN = flag.Bool("n", false, "Disable (third) reading keystrokes from stdin")
var flagInkey = flag.String("inkey", "", "(second) Inject keystrokes")
var flagInkeyFile = flag.String("inkey_file", INKEY1, "Filename from which to (first) inject keystrokes")

func InputRoutine(keystrokes chan<- byte) {
	s := *flagInkey
	if *flagInkeyFile != "" {
		bb, err := ioutil.ReadFile(*flagInkeyFile)
		if err != nil {
			log.Fatalf("Cannot read --inkey_file %q: %v", *flagInkeyFile, err)
		}
		s = string(bb) + s
	}
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
