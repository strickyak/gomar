//go:build f256flat

package emu

import (
	"bytes"
	"fmt"
	// "github.com/strickyak/gomar/display"
	. "github.com/strickyak/gomar/gu"
	"log"
	"strings"
)

// 'Assembly Language Programming for the CoCo 3 (1987)(Laurence A Tepolt).pdf'
// figure 3-5

var usedRom bool
var romMode byte
var enableRom bool
var enableTramp bool
var internalRom [0x8000]byte // up to 32K
var cartRom [0x8000]byte     // up to 32K

var sam display.Sam

var InitialModules []*ModuleFound

type ModuleFound struct {
	Addr uint32
	Len  uint32
	CRC  uint32
	Name string
}

func (m ModuleFound) Id() string {
	return strings.ToLower(fmt.Sprintf("%s.%04x%06x", m.Name, m.Len, m.CRC))
}

func AddressInDeviceSpace(addr Word) bool {
	return (addr&0xFE00) == 0xFE00 && (addr&0xFFF0) != 0xFFF0
}

func GetIOByte(a Word) byte {
	z := GetIOByteI(a)
	L("io GetIOByte %x --> %02x", a, z)
	return z
}
func GetIOByteI(a Word) byte {
	return mem[a]
}

func PutIOByte(a Word, b byte) {
	L("io PutIOByte %x <-- %02x", a, b)
	PokeB(a, b)
}

func DumpHexLines(label string, bb []byte) {
	for i := 0; i < len(bb); i += 32 {
		DumpHexLine(F("%s$%04x", label, i), bb[i:i+32])
	}
}

func DumpHexLine(label string, bb []byte) {
	var buf bytes.Buffer
	buf.WriteString(label)
	for i, b := range bb {
		if i&1 == 0 {
			buf.WriteByte(' ')
		}
		fmt.Fprintf(&buf, "%02x", b)
	}
	buf.WriteRune(' ')
	for _, b := range bb {
		c := b & 127
		if ' ' <= c && c <= '~' {
			buf.WriteByte(c)
		} else {
			buf.WriteByte('.')
		}
	}
	log.Print(buf.String())
}

func DoDumpAllMemory() {
	if !V['m'] {
		return
	}
	Logd("ExplainMMU: %s", ExplainMMU())

	JustDoDumpAllMemory()
}

func JustDoDumpAllMemory() {
	if !V['d'] {
		return
	}

	var i, j int
	var buf bytes.Buffer
	Logd("\n#DumpAllMemory(\n")
	for i = 0; i < 0x10000; i += 32 {
		if (i & 0x1FFF) == 0 {
			// For coco3
			DoExplainMmuBlock(i)
		}
		// Look ahead for something interesting on this line.
		something := false
		for j = 0; j < 32; j++ {
			x := PeekB(Word(i + j))
			if x != 0 && x != ' ' {
				something = true
				break
			}
		}

		if !something {
			continue
		}

		buf.Reset()
		Z(&buf, "M %04x: ", i)
		for j = 0; j < 32; j += 8 {
			Z(&buf,
				"%02x%02x %02x%02x %02x%02x %02x%02x  ",
				PeekB(Word(i+j+0)), PeekB(Word(i+j+1)), PeekB(Word(i+j+2)), PeekB(Word(i+j+3)),
				PeekB(Word(i+j+4)), PeekB(Word(i+j+5)), PeekB(Word(i+j+6)), PeekB(Word(i+j+7)))
		}
		buf.WriteRune(' ')
		for j = 0; j < 32; j++ {
			ch := 0x7F & PeekB(Word(i+j))
			var r rune = '.'
			if ' ' <= ch && ch <= '~' {
				r = rune(ch)
			}
			buf.WriteRune(r)
		}
		Logd("%s\n", buf.String())
	}
	Logd("#DumpAllMemory)\n")
}

func ScanRamForOs9Modules() []*ModuleFound {
	var z []*ModuleFound
	for i := 256; i < len(mem)-256; i++ {
		if mem[i] == 0x87 && mem[i+1] == 0xCD {
			parity := byte(255)
			for j := 0; j < 9; j++ {
				parity ^= mem[i+j]
			}
			if parity == 0 {
				sz := int(HiLo(mem[i+2], mem[i+3]))
				nameAddr := i + int(HiLo(mem[i+4], mem[i+5]))
				got := uint32(HiMidLo(mem[i+sz-3], mem[i+sz-2], mem[i+sz-1]))
				crc := 0xFFFFFF ^ Os9CRC(mem[i:i+sz])
				if got == crc {
					log.Printf("SCAN (at $%x sz $%x) %q %06x %06x", i, sz, Os9StringPhys(nameAddr), mem[i+sz-3:i+sz], 0xFFFFFF^Os9CRC(mem[i:i+sz]))
					z = append(z, &ModuleFound{
						Addr: uint32(i),
						Len:  uint32(sz),
						CRC:  crc,
						Name: Os9StringPhys(nameAddr),
					})
				} else {
					log.Printf("SCAN BAD CRC (@%04x) %06x %06x", i, got, crc)

				}
			} else {
				log.Printf("SCAN BAD PARITY (@%04x) %02x", i, parity)
			}
		}
	}
	return z
}

func Os9CRC(a []byte) uint32 {
	var crc uint32 = 0xFFFFFF
	for k := 0; k < len(a)-3; k++ {
		crc ^= uint32(a[k]) << 16
		for i := 0; i < 8; i++ {
			crc <<= 1
			if (crc & 0x1000000) != 0 {
				crc ^= 0x800063
			}
		}
	}
	return crc & 0xffffff
}
