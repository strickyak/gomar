/*

Gomar is an emulator for the Radio Shack Color Computer
(type 1, 2, and 3) and the Motorola 6809 CPU.

The focus of this emulator is on systems programming and debugging,
not on exact hardware emulation and gaming.

Copyright (C) 2019-2023 Henry Strickland (github.com/strickyak)

Gomar is based on older code "sbc09.c" with the following notices:
   """
      created 1994 by L.C. Benschop.
      copyleft (c) 1994-2014 by the sbc09 team, see AUTHORS for more details.
      license: GNU General Public License version 2, see LICENSE for more details.
   """
That code was coverted to Go Language by Henry Strickland in 2019,
and has been enhanced ever since.

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 675 Mass Ave, Cambridge, MA 02139, USA.

*/
package main

import (
	"github.com/strickyak/gomar/emu"

	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

var FlagCpuProfile = flag.String("cpu-profile", "", "write cpu profile to file")
var FlagMemProfile = flag.String("mem-profile", "", "write memory profile to file")
var FlagTTL = flag.Duration("ttl", parseDurationOrDie("5m"), "max duration to live, or 0 for unlimited")
var FlagLogFile = flag.String("log-file", "_log", "file to write logging messages to")
var FlagSplash = flag.Bool("splash", true, "Enables printing an initial notice")

func parseDurationOrDie(s string) time.Duration {
	value, err := time.ParseDuration(s)
	if err != nil {
		log.Fatalf("cannot ParseDuration %q: %v", s, err)
	}
	return value
}

func main() {
	log.SetFlags(0)
	flag.Parse()

	if *FlagSplash {
		os.Stderr.Write(([]byte(SPLASH))[1:]) // skip initial newline.
	}

	logger, err := os.Create(*FlagLogFile)
	if err != nil {
		log.Fatalf("Cannot create log file %q: %v", *FlagLogFile, err)
	}
	log.SetOutput(logger)

	if *FlagSplash {
		fmt.Fprintf(os.Stderr, "Verbose logging is going to file %q\n\n", *FlagLogFile)
	}

	if *FlagTTL != 0 {
		go func() {
			time.Sleep(*FlagTTL)
			log.Fatal("gomar: TTL Expired: %v", *FlagTTL)
		}()
	}

	if *FlagCpuProfile != "" {
		f, err := os.Create(*FlagCpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	emu.Main()

	if *FlagMemProfile != "" {
		f, err := os.Create(*FlagMemProfile)
		if err != nil {
			log.Fatalf("could not create memory profile: %v", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatalf("could not write memory profile: %v", err)
		}
	}
}

const SPLASH = `
Gomar 6809 Emulator, Copyright (C) 2019-2023 Henry Strickland.
This is free software, and you are welcome to redistribute it
under the terms of the GNU General Public License version 2.
`

/*
  HINT
  go run --tags=level2,coco3,hyper  gomar.go -boot .o/drive/boot2coco3 -disk .o/drive/disk2 -v=abcdefghijklmnopqrstuvwxyz
*/
