//go:build !canvas

package emu

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec" // Command stty
	"time"

	. "github.com/strickyak/gomar/gu"
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
var flagInkeyFile = flag.String("inkey_file", "", "Filename from which to (first) inject keystrokes")
var flagInkeySleep = flag.Duration("inkey_sleep", 100*time.Millisecond, "how long to sleep before each keystroke")

func InputRoutine(keystrokes chan<- byte) {
	keystrokes2 := make(chan byte, 1)
	go InputRoutine2(keystrokes2)

	for k := range keystrokes2 {
		time.Sleep(*flagInkeySleep)
		keystrokes <- k
	}
}

func InputRoutine2(keystrokes chan<- byte) {
	// Start with the --inkey flag.
	s := *flagInkey
	// But if the --inkey_file flag is set, insert that in front.
	if *flagInkeyFile != "" {
		bb, err := ioutil.ReadFile(*flagInkeyFile)
		if err != nil {
			log.Fatalf("Cannot read --inkey_file %q: %v", *flagInkeyFile, err)
		}
		s = string(bb) + s
	}
	// Use all the flag chars.
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

	// Do not read stdin if -n flag (like the -n flag to ssh).
	if *flagN {
		return
	}
	// Do read stdin.
	if true {
		// e := exec.Command("stty", "-cbreak").Run()
		cmd := &exec.Cmd{
			Path:   "/usr/bin/stty",
			Args:   []string{"stty", "cbreak"},
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}
		e := cmd.Run()
		// e := exec.Command("sh", "-c", "stty -cbreak").Run()
		Check(e)
		log.Printf("ZXC === stty -cbreak ===")
		char1 := make([]byte, 1)
		for {
			nn, err := os.Stdin.Read(char1)
			if err != nil {
				log.Panicf("os.Stdin.Read (* stty -cbreak *): %v", char1)
			}
			if nn != 1 {
				log.Panicf("os.Stdin.Read (* stty -cbreak *): short, want 1, got %d", nn)
			}
			c1 := char1[0]
			log.Printf("ZXC GOT_KEYSTROKE %d.", c1)
			keystrokes <- Cond(c1 == 10, '\r', c1)
		}
	} else {
		in := bufio.NewScanner(os.Stdin)
		for in.Scan() {
			for _, r := range in.Text() {
				keystrokes <- byte(r)
			}
			keystrokes <- '\r'
		}
	}
	close(keystrokes)
	log.Fatal("Stdin ended")
}

func PublishVideoText() {}
