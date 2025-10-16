//go:build coco1 || coco3

package emu

import (
	"bufio"
	"bytes"
	"fmt"
	. "github.com/strickyak/gomar/gu"
	"github.com/strickyak/gomar/listings"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// 'Assembly Language Programming for the CoCo 3 (1987)(Laurence A Tepolt).pdf'
// figure 3-5

var usedRom bool

// ?//var romMode byte
var enableRom bool
var enableTramp bool
var internalRom [0x8000]byte // up to 32K
var cartRom [0x8000]byte     // up to 32K
var portMem [256]byte        // anything written to $FFxx

var Pia0HorzSyncInterruptEnable bool  // 15738 Hz
var Pia0FrameSyncInterruptEnable bool // 60 Hz
var GimeHorzSyncInterruptEnable bool  // 15738 Hz
var GimeVirtSyncInterruptEnable bool  // 60 Hz

var horzCycles int64
var frameCycles int64
var gimeHorzCycles int64
var gimeVirtCycles int64

var horzPending bool
var framePending bool
var gimeHorzPending bool
var gimeVirtPending bool
var nmiPending bool

var sam Sam

var InitialModules []*ModuleFound
var InternalRomSrc *listings.ModSrc
var ExternalRomSrc *listings.ModSrc
var GlobalSrc *listings.ModSrc

/*
Section: .bss (_nekot.o) load at 0008, length 0046
Section: .data.more (_nekot.o) load at 1002, length 0049
Section: .text (_nekot.o) load at 104B, length 0857
Section: .final.kern (_nekot.o) load at 18A2, length 0002
Section: .data.startup (_nekot.o) load at 18A4, length 0018
Section: .text.startup (_nekot.o) load at 18BC, length 0135
Section: .final.startup (_nekot.o) load at 19F1, length 0002
*/

type Section struct {
	name  string
	begin uint
	size  uint
}

var MatchSection = regexp.MustCompile(`^Section: (\S+) [(].*[)] load at (....), length (....)$`)

func LoadMapAndList(mapfile, listfile string) *listings.ModSrc {
	z := &listings.ModSrc{
		Src: make(map[uint]string),
	}

	if mapfile == "" && listfile == "" {
		return z
	}
	if mapfile == "" {
		log.Fatalf("has listfile %q but missing mapfile", listfile)
	}
	if listfile == "" {
		log.Fatalf("has mapfile %q but missing listfile", mapfile)
	}

	sections := make(map[string]Section)
	for _, line := range ReadFileLines(mapfile) {
		if m := MatchSection.FindStringSubmatch(line); m != nil {
			name := m[1]
			begin := uint(Value(strconv.ParseUint(m[2], 16, 16)))
			size := uint(Value(strconv.ParseUint(m[3], 16, 16)))
			sections[name] = Section{name, begin, size}
			log.Printf("SECTION: %q %04x %04x", name, begin, size)
		}
	}

	in := listings.LoadFile(listfile)
	in = Value(in, in.Err)
	newSrcMap := make(map[uint]string)
	for addr, line := range in.Src {
		if area, ok := in.Area[addr]; ok {
			if section, ok := sections[area]; ok {
				newSrcMap[addr+section.begin] = line
			}
		}
	}

	return &listings.ModSrc{
		Src:      newSrcMap,
		Filename: in.Filename,
		Err:      in.Err,
	}
}

func ReadFileLines(filename string) []string {
	var z []string

	fd := Value(os.Open(filename))
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		z = append(z, strings.TrimSuffix(scanner.Text(), "\n"))
	}

	return Value(z, scanner.Err())
}

type ModuleFound struct {
	Addr     uint32
	Len      uint32
	CRC      uint32
	Name     string
	Filename string // If not an OS9 module
}

func (m ModuleFound) Id() string {
	return strings.ToLower(fmt.Sprintf("%s.%04x%06x", m.Name, m.Len, m.CRC))
}

func LoadRomListings() {
	GlobalSrc = LoadMapAndList(*FlagGlobalMap, *FlagGlobalListing)
	if *FlagInternalRomListing != "" {
		InternalRomSrc = listings.LoadFile(*FlagInternalRomListing)
	}
	if *FlagInternalRomDup != "" && InternalRomSrc != nil {
		words := strings.Split(*FlagInternalRomDup, ":")
		//T(words)
		//T(len(words))
		if len(words) == 3 {
			begin := DeHex(words[0])
			end := DeHex(words[1])
			to := DeHex(words[2])
			//T(begin)
			//T(end)
			//T(to)
			tmp := make(map[uint]string)
			for k, v := range InternalRomSrc.Src {
				tmp[k] = v // make a copy
			}
			for k, v := range tmp {
				if begin <= k && k < end {
					//T(k, k + to - begin , v)
					InternalRomSrc.Src[k+to-begin] = v
				}
			}
		}
	}
	if *FlagExternalRomListing != "" {
		ExternalRomSrc = listings.LoadFile(*FlagExternalRomListing)
	}
}

/*
func AddressInTrampSpace(addr Word) bool {
	if BitFixedFExx {
		return (addr&0xFF00) == 0xFE00 || (addr&0xFFF0) == 0xFFF0
	} else {
		return (addr & 0xFFF0) == 0xFFF0
	}
}
*/

/*
func BAD MappedAddressInRomSpace(addr Word, mapped int) bool {
	physPage := uint(mapped) >> 13
	return 0x3C <= physPage && physPage <= 0x3F && !AddressInDeviceSpace(addr)
}
*/

func AddressInRomSpace(addr Word) bool {
	z := 0x8000 <= addr && !AddressInDeviceSpace(addr)
	if z {
		// T("ROM %s %04x ", Cond(UseExternalRomAssumingRom(addr), "X", "N"), addr)
	}
	return z
}

func AddressInDeviceSpace(addr Word) bool {
	return 0xFF00 <= addr && addr < 0xFFF0
}

func GetIOByte(a Word) byte {
	z := GetIOByteI(a)
	L("io GetIOByte %x --> %02x", a, z)
	return z
}
func GetIOByteI(a Word) byte {
	var z byte

	if 0xFF00 <= a && a <= 0xFF40 {
		a &^= 0x003C // Wipe out the don't-care bits of PIAs.
	}

	if 0xFF90 <= a && a < 0xFFC0 {
		return GetGimeIOByte(a)
	}

	switch a {
	case 0xFF78,
		0xFF79,
		0xFF7a,
		0xFF7b:
		return GetCopico(a)

	/* PIA 0 */
	case 0xFF00:
		horzPending = false // Reading Output Register A clears CA1 interrupt.
		z = 255

		if PeekB(0xFF02) == 0xFF {
			// Not strobing keyboard, so answer mouse buttons.
			if MouseDown {
				z = 0xFC // buttons 1 and 2.
			}
		} else {
			// Strobing keyboard.
			if kbd_ch != 0 {
				z = keypress(kbd_probe, kbd_ch)
				Logd("KEYBOARD: %02x %q -> %02x\n", kbd_probe, string(rune(kbd_ch)), z)
			} else {
				Logd("KEYBOARD: %02x      -> %02x\n", kbd_probe, z)
			}
		}

		dac := float64(PeekB(0xFF20)&0xFC) / 256.0
		var mouse float64
		if PeekB(0xFF01)&0x08 == 0 {
			mouse = MouseX // or vice versa
		} else {
			mouse = MouseY // or vice versa
		}
		if mouse <= dac {
			z &= 0x7F
		} else {
			z |= 0x80
		}
		Logd("PIA: Get IO byte $%04x -> $%02x\n", a, z)
		return z
	case 0xFF01:
		return 0
	case 0xFF02:
		framePending = false // Reading Output Register B clears CB1 interrupt.
		return kbd_probe     // Reset IRQ when this is read. TODO: multiple sources of IRQ.
	case 0xFF03:
		return 0x80 // Negative bit set: Yes the PIA caused IRQ.

	/* PIA 1 */
	case 0xFF22:
		Logd("TODO: Get Io byte 0x%04x\n", a)
		return 0

	case 0xFF48: /* STATREG */
		return 0 /* low bit 0 means Ready, other bits are errors or not ready */

	case 0xFF4A /*cocosdc boot*/, 0xFF4B /*floppy*/ : /* Read Data */
		z = 0
		if disk_i < 256 {
			z = disk_stuff[disk_i]
			Logd("fnord %x -> %x\n", disk_i, z)
		} else {
			z = 0
		}
		disk_i++
		if disk_i == 257 {
			Logd("Read SET NMI_PENDING\n")
			nmiPending = true
			z = 0
			disk_i = 0
		}
		return z

	case 0xFF83: /* emudsk */
		return EmudskGetIOByte(a)

	case 0xFF68,
		0xFF69,
		0xFF6a,
		0xFF6b:
		return GetCocoIO(a)

	default:
		Logd("UNKNOWN GetIOByte: 0x%04x\n", a)
		return 0
	}
	panic("notreached")
}

func LogicalSector(sector, side, track byte) int64 {
	log.Printf("LogiclSector (fmt=%d.) sector=%d. side=%d. track=%d.", disk_dd_fmt, sector, side, track)
	switch disk_dd_fmt {
	case 2:
		if side != 0 {
			// ddt
			return int64(disk_sector) - 0 + int64(disk_track)*18
		}
		return int64(disk_sector) - 1 + int64(disk_track)*18
	case 3:
		return int64(disk_sector) - 1 + int64(disk_side)*18 + int64(disk_track)*36
	}
	log.Panicf("bad disk_dd_fmt: %d", disk_dd_fmt)
	panic(0)
}

var FF22Bits = []string{
	"VdgGraphics", "VdgGM2", "VdgGM1/invert", "VdgGM0/shiftToggle",
	"VdgColorSet", "RamSize/Input", "SingleBitSound/Out", "Rs232/Input"}

func ExplainBits(b byte, meanings []string) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "$%02x=", b)
	mask := byte(128)
	for i := 0; i < 8; i++ {
		if b&mask != 0 {
			buf.WriteString(meanings[i])
		}
		if i < 7 {
			buf.WriteByte('|')
		}
		mask >>= 1
	}
	return buf.String()
}

func PutIOByte(a Word, b byte) {
	L("io PutIOByte %x <-- %02x", a, b)
	portMem[a-0xFF00] = b
	PokeB(a, b)

	if 0xFF90 <= a && a < 0xFFC0 {
		PutGimeIOByte(a, b)
		return
	}

	if 0xFF00 <= a && a <= 0xFF40 {
		a &^= 0x003C // Wipe out the don't-care bits of PIA addresses.
	}

	if 0xFFE0 <= a {
		Logd("PutIOByte: WUT GIMEX !? : $%04x", a)
		return
	}

	switch a {
	default:
		log.Panicf("UNKNOWN PutIOByte address: 0x%04x", a)

	case 0xFF78,
		0xFF79,
		0xFF7a,
		0xFF7b:
		PutCopico(a, b)

	// http://tlindner.macmess.org/wp-content/uploads/2006/09/cocopias-R3.pdf
	case 0xFF00, 0xFF1C:
		Logd("PIA0: Put IO byte $%04x <- $%02x\n", a, b)

	case 0xFF01, 0xFF1D:
		Logd("PIA0: Put IO byte $%04x <- $%02x\n", a, b)
		Pia0HorzSyncInterruptEnable = (b & 1) != 0
		return

	case 0xFF02, 0xFF1E:
		Logd("PIA0: Put IO byte $%04x <- $%02x\n", a, b)
		kbd_probe = b
		return

	case 0xFF03, 0xFF1F:
		Logd("PIA0: Put IO byte $%04x <- $%02x\n", a, b)
		Pia0FrameSyncInterruptEnable = (b & 1) != 0

	case 0xFF20,
		0xFF21,
		0xFF23:
		Logd("PIA1: Put IO byte $%04x <- $%02x\n", a, b)
		return

	case 0xFF22:
		Logd("VDG: %s", ExplainBits(b, FF22Bits))
		Logd("PIA1: Put IO byte $%04x <- $%02x\n", a, b)
		return

	case 0xFF40: /* CONTROL */
		{
			disk_control = b
			disk_side = byte(Cond(b&0x40 != 0, 1, 0))
			disk_drive = byte(Cond((b&1 != 0), 1, Cond((b&2 != 0), 2, Cond((b&4 != 0), 3, 0))))

			Logd("CONTROL: disk_command %x (control %x side %x drive %x)\n", disk_command, disk_control, disk_side, disk_drive)
			if b == 0 {
				// log.Panicf("panic: disk_command 0")
				break
			}

			log.Printf("...... Disk Command ($%x) Fnord", disk_command)
			switch disk_command {
			default:
				{
					log.Printf("Unknown Disk Command ($%x) Fnord", disk_command)
				}
			case 0x43:
				{
					log.Printf("Start Command Mode ($43) Fnord")
				}
			case 0xD0:
				{
					log.Printf("Stop any disk command in progress Fnord")
				}
			case 0x80:
				{
					prev_disk_command = disk_command
					disk_offset = 256 * LogicalSector(disk_sector, disk_side, disk_track)
					if disk_drive != 1 {
						log.Panicf("ERROR: R: Drive %d not supported\n", disk_drive)
					}
					if disk_fd == nil {
						log.Panicf("ERROR: R: No file for Disk Read Sector\n")
					}

					disk_stuff = zero_disk_stuff
					log.Printf("disk sector seek: offset=%d. -- disk_sector=%d. disk_side=%d. disk_track=%d.", disk_offset, disk_sector, disk_side, disk_track)
					_, err := disk_fd.Seek(disk_offset, 0)
					if err != nil {
						log.Panicf("Bad disk sector seek: offset=%d. err=%v disk_sector=%d. disk_side=%d. disk_track=%d.", disk_offset, err, disk_sector, disk_side, disk_track)
					}
					n, err := disk_fd.Read(disk_stuff[:])
					if err != nil {
						log.Panicf("Bad disk sector read: err=%v", err)
					}
					if n != 256 {
						log.Panicf("Short disk sector read: n=%d", n)
					}

					AssertEQ(n, 256)
					disk_i = 0
					Logd("READ fnord (Track, Sector-1) %d:%d:%d:%d == %d\n", disk_drive, disk_track, disk_side, disk_sector-1, disk_offset>>8)
				}
			case 0xA0:
				{
					prev_disk_command = disk_command
					disk_offset = 256 * LogicalSector(disk_sector, disk_side, disk_track)
					if disk_drive != 1 {
						log.Panicf("ERROR: W: Drive %d not supported\n", disk_drive)
					}
					if disk_fd == nil {
						log.Panicf("ERROR: W: No file for Disk Read Sector\n")
					}
					disk_stuff = zero_disk_stuff
					_, err := disk_fd.Seek(int64(disk_offset), 0)
					if err != nil {
						log.Panicf("Bad disk sector seek: err=%v", err)
					}

					disk_i = 0
					Logd("WRITE fnord (Track, Sector-1) %d:%d:%d:%d == %d\n", disk_drive, disk_track, disk_side, disk_sector-1, disk_offset>>8)
				}
			}
			disk_command = 0
		}
	case 0xFF48:
		{ // CMDREG //
			disk_command = b
			switch b {
			case 0x10:
				{
					disk_track = disk_data
					disk_status = 0
					Logd("Seek : %d\n", disk_data)
				}
			case 0x80:
				{ // Read Sector //
					// We have set disk_command.  Next control write defines disk & side. //

				}
			case 0xD0:
				{
					disk_drive = 0
					disk_side = 0
					disk_track = 0
					disk_sector = 0
					disk_i = 0
					disk_stuff = zero_disk_stuff
					Logd("Reset Disk\n")
				}
			}
		}
	case 0xFF49: /* TRACK */
		disk_track = b
		Logd("Track : %d\n", b)

	case 0xFF4A: /* SECTOR */
		disk_sector = b
		Logd("Sector-1 : %d\n", b-1)

	case 0xFF4B:
		{ /* DATA */
			if (prev_disk_command & 0xF0) != 0xA0 {
				disk_i = 0
				disk_data = b
			} // else
			if true {
				if disk_i < 256 {
					Logd("fnord %x %x <- %x\n", prev_disk_command, disk_i, b)
					disk_stuff[disk_i] = b
					///++disk_i;
				}
			}
			if (prev_disk_command & 0xF0) == 0xA0 {
				if disk_i < 256 {
					disk_i++
				}
				// TODO -- fix writing.
				if disk_i >= 256 {
					Logd("Write SET NMI_PENDING\n")
					nmiPending = true
					disk_i = 0

					// TODO -- fix writing.
					n, err := disk_fd.Write(disk_stuff[:])
					if err != nil {
						log.Panicf("Error in disk_fd.Write: %v", err)
					}
					if n != 256 {
						log.Panicf("Error in disk_fd.Write: Short n=%d", n)
					}
					Logd("DID_WRITE fnord (Track, Sector-1) %d:%d:%d:%d == %d\n", disk_drive, disk_track, disk_side, disk_sector-1, disk_offset>>8)
				}
			}

		}

	case 0xFF42:
		Logd("Write to $FF42")

	case 0xFF77:
		L("WTF: Write to $FF77")
	case 0xFF7F:
		Logd("Write to $FF7F")

	case 0xFFE1:
		Logd("Write to $FFE1")
	case 0xFFE2:
		Logd("Write to $FFE2")
	case 0xFFE3:
		Logd("Write to $FFE3")
	case 0xFFE8:
		Logd("Write to $FFE8")
	case 0xFF51:
		Logd("Write to $FF51")
	case 0xFF56:
		Logd("tinyide: head = %02x", b)

		/* VDG */
	case 0xFFC0:
		Logd("VDG PutByte OK: %x <- %x\n", a, b)
		sam.Vx &^= 1
		Logd("VDG sam.Vx <- $%x", sam.Vx)
	case 0xFFC1:
		Logd("VDG PutByte OK: %x <- %x\n", a, b)
		sam.Vx |= 1
		Logd("VDG sam.Vx <- $%x", sam.Vx)
	case 0xFFC2:
		Logd("VDG PutByte OK: %x <- %x\n", a, b)
		sam.Vx &^= 2
		Logd("VDG sam.Vx <- $%x", sam.Vx)
	case 0xFFC3:
		Logd("VDG PutByte OK: %x <- %x\n", a, b)
		sam.Vx |= 2
		Logd("VDG sam.Vx <- $%x", sam.Vx)
	case 0xFFC4:
		Logd("VDG PutByte OK: %x <- %x\n", a, b)
		sam.Vx &^= 4
		Logd("VDG sam.Vx <- $%x", sam.Vx)
	case 0xFFC5:
		Logd("VDG PutByte OK: %x <- %x\n", a, b)
		sam.Vx |= 4
		Logd("VDG sam.Vx <- $%x", sam.Vx)

	case 0xFFC6:
		Logd("VDG PutByte OK: %x <- %x\n", a, b)
		sam.Fx &^= 1
		Logd("VDG sam.Fx <- $%x", sam.Fx)
	case 0xFFC7:
		Logd("VDG PutByte OK: %x <- %x\n", a, b)
		sam.Fx |= 1
		Logd("VDG sam.Fx <- $%x", sam.Fx)
	case 0xFFC8:
		Logd("VDG PutByte OK: %x <- %x\n", a, b)
		sam.Fx &^= 2
		Logd("VDG sam.Fx <- $%x", sam.Fx)
	case 0xFFC9:
		Logd("VDG PutByte OK: %x <- %x\n", a, b)
		sam.Fx |= 2
		Logd("VDG sam.Fx <- $%x", sam.Fx)
	case 0xFFCA:
		Logd("VDG PutByte OK: %x <- %x\n", a, b)
		sam.Fx &^= 4
		Logd("VDG sam.Fx <- $%x", sam.Fx)
	case 0xFFCB:
		Logd("VDG PutByte OK: %x <- %x\n", a, b)
		sam.Fx |= 4
		Logd("VDG sam.Fx <- $%x", sam.Fx)
	case 0xFFCC:
		Logd("VDG PutByte OK: %x <- %x\n", a, b)
		sam.Fx &^= 8
		Logd("VDG sam.Fx <- $%x", sam.Fx)
	case 0xFFCD:
		Logd("VDG PutByte OK: %x <- %x\n", a, b)
		sam.Fx |= 8
		Logd("VDG sam.Fx <- $%x", sam.Fx)
	case 0xFFCE:
		Logd("VDG PutByte OK: %x <- %x\n", a, b)
		sam.Fx &^= 16
		Logd("VDG sam.Fx <- $%x", sam.Fx)
	case 0xFFCF:
		Logd("VDG PutByte OK: %x <- %x\n", a, b)
		sam.Fx |= 16
		Logd("VDG sam.Fx <- $%x", sam.Fx)
	case 0xFFD0:
		Logd("VDG PutByte OK: %x <- %x\n", a, b)
		sam.Fx &^= 32
		Logd("VDG sam.Fx <- $%x", sam.Fx)
	case 0xFFD1:
		Logd("VDG PutByte OK: %x <- %x\n", a, b)
		sam.Fx |= 32
		Logd("VDG sam.Fx <- $%x", sam.Fx)
	case 0xFFD2:
		Logd("VDG PutByte OK: %x <- %x\n", a, b)
		sam.Fx &^= 64
		Logd("VDG sam.Fx <- $%x", sam.Fx)
	case 0xFFD3:
		Logd("VDG PutByte OK: %x <- %x\n", a, b)
		sam.Fx |= 64
		Logd("VDG sam.Fx <- $%x", sam.Fx)

	case 0xFFD4:
		sam.P1RamSwap = 0
		Logd("VDG sam.P1RamSwap <- $%x", sam.P1RamSwap)
	case 0xFFD5:
		sam.P1RamSwap = 1
		Logd("VDG sam.P1RamSwap <- $%x", sam.P1RamSwap)

	case 0xFFD6:
		sam.Rx &^= 1
		Logd("VDG sam.Rx <- $%x", sam.Rx)
	case 0xFFD7:
		sam.Rx |= 1
		Logd("VDG sam.Rx <- $%x", sam.Rx)
	case 0xFFD8:
		sam.Rx &^= 2
		Logd("VDG sam.Rx <- $%x", sam.Rx)
	case 0xFFD9:
		sam.Rx |= 2
		Logd("VDG sam.Rx <- $%x", sam.Rx)

	case 0xFFDA:
		sam.Mx &^= 1
		Logd("VDG sam.Mx <- $%x", sam.Mx)
	case 0xFFDB:
		sam.Mx |= 1
		Logd("VDG sam.Mx <- $%x", sam.Mx)
	case 0xFFDC:
		sam.Mx &^= 2
		Logd("VDG sam.Mx <- $%x", sam.Mx)
	case 0xFFDD:
		sam.Mx |= 2
		Logd("VDG sam.Mx <- $%x", sam.Mx)

	case 0xFFDE:
		sam.TyAllRam = false
		enableRom = true
		Logd("VDG TyAllRam <- FALSE")
	case 0xFFDF:
		sam.TyAllRam = true
		enableRom = false
		Logd("VDG TyAllRam <- TRUE")

	case 0xFF80,
		0xFF81,
		0xFF82,
		0xFF83,
		0xFF84,
		0xFF85,
		0xFF86:
		EmudskPutIOByte(a, b)

	case 0xFF68,
		0xFF69,
		0xFF6a,
		0xFF6b:
		PutCocoIO(a, b)
	}
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

func DoDumpSamBits() {
	Logd("VDG/SAM BITS: F=%x M=%x R=%x V=%x TyAllRam=%x P1RamSwap=%x",
		sam.Fx, sam.Mx, sam.Rx, sam.Vx, sam.TyAllRam, sam.P1RamSwap)
}

func DoDumpAllMemory() {
	if !V['m'] {
		return
	}
	DoDumpSamBits()
	DumpGimeStatus()
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
