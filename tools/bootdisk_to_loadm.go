//go:build main

/*
Read an OS9 boot disk image and extract the boot track
from track 35.  Write out initial RAM image for gomar,
as raw memory bytes, starting at address $0000, and going up.
*/
package main

/*
	DD.FMT DISK FORMAT:  offset $10:

	BIT B0 - SIDE
	0 = SINGLE SIDED
	1 = DOUBLE SIDED

	BIT B1 - DENSITY
	0 = SINGLE DENSITY
	1 = DOUBLE DENSITY

	BIT B2 - TRACK DENSITY
	0 = SINGLE (48 TPI)
	1 = DOUBLE (96 TPI)
*/

import (
	"flag"
	"io"
	"log"
	"os"
)

const BOOT_SECTOR = 1224
const BOOT_SECTOR_VHD = 612

func main() {
	flag.Parse()
	var sector0 [256]byte
	_, err := io.ReadFull(os.Stdin, sector0[:])
	if err != nil {
		log.Fatalf("cannot read sector0: %v", err)
	}

	var lsn int64
	format := sector0[0x10]
	switch format {
	case 2:
		lsn = 612
	case 3:
		lsn = 1224
	default:
		log.Fatalf("unknown format byte: $%02x", format)
	}

	var track [18 * 256]byte
	_, err = os.Stdin.Seek(lsn*256, 0)
	if err != nil {
		log.Fatalf("cannot Seek bootTrack: %v", err)
	}
	_, err = io.ReadFull(os.Stdin, track[:])
	if err != nil {
		log.Fatalf("cannot Read bootTrack: %v", err)
	}

	if track[0] != 'O' || track[1] != 'S' {
		log.Fatalf("bad Magic Numbers in bootTrack: $%02x $%02x", track[0], track[1])
	}

	var loadm []byte
	loadm = append(loadm, 0, 0x12, 0x00, 0x26, 0x00)
	loadm = append(loadm, track[:]...)
	loadm = append(loadm, 0xFF, 0, 0, 0x26, 0x02)

	n, err := os.Stdout.Write(loadm)
	if err != nil {
		log.Fatalf("cannot write Stdout: %v", err)
	}
	if n != len(loadm) {
		log.Fatalf("short write to Stdout")
	}
}
