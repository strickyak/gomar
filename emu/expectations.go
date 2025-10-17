package emu

import (
	. "github.com/strickyak/gomar/gu"

	"fmt"
	"log"
	"os"
	"strings"
)

var Expectations []string

func InitExpectations() {
	if *FlagExpectFile != "" && Expectations == nil {
		bb := Value(os.ReadFile(*FlagExpectFile))
		Expectations = strings.Split(string(bb), "\n")
	}
	if *FlagExpect != "" && Expectations == nil {
		Expectations = strings.Split(*FlagExpect, ";")
		fmt.Printf("\n===@=== SET Expectations: %q\n", *FlagExpect)
		log.Printf("\n===@=== SET Expectations: %q\n", *FlagExpect)
	}
	log.Printf("EXPECTATIONS LEN: %d", len(Expectations))
	log.Printf("EXPECTATIONS: %#v", Expectations)
}

func CheckExpectation(got string) {
	// Skip out if no expectations were defined.
	if len(Expectations) == 0 {
		return
	}

	// Skip empty expectations
	for len(Expectations) > 0 && len(Expectations[0]) == 0 {
		Expectations = Expectations[1:]
	}

	// Process one expectation, if possible, if valid.
	if len(Expectations) > 0 {
		if strings.Contains(got, Expectations[0]) {
			fmt.Printf("\n===@=== GOT Expectation: %q\n", Expectations[0])
			log.Printf("\n===@=== GOT Expectation: %q\n", Expectations[0])
			Expectations = Expectations[1:]

		}
	}

	// Skip empty expectations
	for len(Expectations) > 0 && len(Expectations[0]) == 0 {
		Expectations = Expectations[1:]
	}

	// Exit(0) if all expectations are met.
	if len(Expectations) == 0 {
		fmt.Printf("\n===@=== SUCCESS -- FINISHED Expectations.\n")
		log.Printf("\n===@=== SUCCESS -- FINISHED Expectations.\n")
		os.Exit(0)
	}
}
