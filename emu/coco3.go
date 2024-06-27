//go:build coco3
// +build coco3

package emu

import (
	. "github.com/strickyak/gomar/gu"
	"github.com/strickyak/gomar/sym"

	"bytes"
	"flag"
	"fmt"
	"log"
	"strings"
)

var FuzixModelFlag = flag.Bool("fuzix", false, "special for fuxiz")

const TraceMem = false // TODO: restore this some day.

var Init0 byte
var Init1 byte
var MmuEnable bool
var MmuTask byte
var MmuMap [2][8]byte

var BitCoCo12Compat bool
var BitFixedFExx bool
var BitMC0, BitMC1 bool // Rom Mode: low bits at FF90

var videoEpoch int64

var DisabledMmuMap = []byte{0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f}

func UseExternalRomAssumingRom(addr Word) bool {
	var z bool
	if BitMC1 {
		if 0xFFF0 <= addr {
			z = false
		} else if BitMC0 {
			z = true
		} else {
			z = false
		}
	} else {
		z = 0xC000 <= addr && addr < 0xFFF0
	}
	//T("ext?addr,mc1,mc0", addr, BitMC1, BitMC0, "-->", z)
	return z
}
func InternalRomOffset(addr Word) uint {
	if 0xFFF0 <= addr {
		return 0x3FFF & uint(addr)
	}
	if BitMC1 {
		return 0x7FFF & uint(addr)
	} else {
		return 0x3FFF & uint(addr)
	}
}
func ExternalRomOffset(addr Word) uint {
	if BitMC1 {
		return 0x7FFF & uint(addr)
	} else {
		return 0x4000 + 0x3FFF&uint(addr)
	}
}

////////////////////////////////////////

// Coco3Contract ensures the contract between Coco3's disk booting mechanism
// and the OS/9 Level2 kernel, documented at
// nitros9/level2/modules/kernel/ccbkrn.txt
func InitHardware() {
	display = NewGime()
	if usedRom {
		Coco3ContractRaw()
	} else {
		Coco3ContractForDos()
	}
	PutB(0xFF90, 0x80) // Start in VDG mode
}
func InitializeMemoryMap() {
	for task := 0; task < 2; task++ {
		for block, phys := range DisabledMmuMap {
			MmuMap[task][block] = phys
		}
	}
}

func Coco3ContractRaw() {
	if *FuzixModelFlag {
		DisabledMmuMap = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
	} else {
		DisabledMmuMap = []byte{0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f}
	}
	InitializeMemoryMap()
}

func Coco3ContractForDos() {
	if *FuzixModelFlag {
		DisabledMmuMap = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
	} else {
		DisabledMmuMap = []byte{0x00, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f}
	}
	InitializeMemoryMap()

	// Initialize physical block 3b to spaces, except 0x0008 at the beginning.
	const block3b = 0x3b * 0x2000
	mem[block3b+0] = 0x00
	mem[block3b+1] = 0x08
	for i := 2; i < 0x2000; i++ {
		mem[block3b+i] = ' '
	}

	/*   starting at 0xff90:
	6c      init0
	00      init1
	00      irq enable
	00      firq enable
	0900    timer register
	0000    unused
	0320    screen settings
	0000    ????
	00      ????
	ec01    physical video address (block 3b offset 0x0008 )
	00      horizontal offset / scroll

	A mirror of these bytes will appear at 0x0090-0x009f in the DP
	*/
	for i, b := range []byte{0x6c, 0, 0, 0, 9, 0, 0, 0, 3, 0x20, 0, 0, 0, 0x3c, 1, 0} {
		PutIOByte(Word(0xFF90+i), b)
		// DONT // mem[0x90+i] = b // Probably don't need to set the mirror, but doing it anyway.
	}
	if *FuzixModelFlag {
		PutIOByte(0xFF91, 0x01) // BOOT sets task 1
	}
}

type Mapping [8]Word

func GetMapping(addr Word) Mapping {
	// Mappings are in SysMem (block 0).
	return Mapping{
		// TODO: drop the "0x3F &".
		0x3F & SysMemW(addr),
		0x3F & SysMemW(addr+2),
		0x3F & SysMemW(addr+4),
		0x3F & SysMemW(addr+6),
		0x3F & SysMemW(addr+8),
		0x3F & SysMemW(addr+10),
		0x3F & SysMemW(addr+12),
		0x3F & SysMemW(addr+14),
	}
}

func WithMmuTask(task byte, fn func()) {
	tmp := MmuTask
	MmuTask = task
	defer func() {
		MmuTask = tmp
	}()
	fn()
}

func GetBytesTask0(addr Word, n Word) []byte {
	// Use Task 0 for the mapping.
	tmp := MmuTask
	MmuTask = 0
	defer func() {
		MmuTask = tmp
	}()

	var bb []byte
	for i:=Word(0); i<n; i++ {
		bb = append(bb, PeekB(addr+i))
	}
	return bb
}
func GetMappingTask0(addr Word) Mapping {
	// Use Task 0 for the mapping.
	tmp := MmuTask
	MmuTask = 0
	defer func() {
		MmuTask = tmp
	}()

	return Mapping{
		// TODO: drop the "0x3F &".
		0x3F & PeekW(addr),
		0x3F & PeekW(addr+2),
		0x3F & PeekW(addr+4),
		0x3F & PeekW(addr+6),
		0x3F & PeekW(addr+8),
		0x3F & PeekW(addr+10),
		0x3F & PeekW(addr+12),
		0x3F & PeekW(addr+14),
	}
}
func TaskNumberToMapping(task byte) Mapping {
	dope := PeekW(0x00A1) // D.TskIPt
	dat := PeekW(dope + 2*Word(task))
	var m Mapping
	for i := Word(0); i < 8; i++ {
		m[i] = PeekW(dat + 2*i)
	}
	return m
}
func PeekBWithTask(addr Word, task byte) byte {
	m := TaskNumberToMapping(task)
	return PeekBWithMapping(addr, m)
}
func PeekWWithTask(addr Word, task byte) Word {
	m := TaskNumberToMapping(task)
	return PeekWWithMapping(addr, m)
}
func PeekBWithMapping(addr Word, m Mapping) byte {
	logBlock := (addr >> 13) & 7
	physBlock := m[logBlock]
	ptr := int(addr&0x1FFF) | (int(physBlock) << 13)
	return mem[ptr]
}
func PeekWWithMapping(addr Word, m Mapping) Word {
	hi := PeekBWithMapping(addr, m)
	lo := PeekBWithMapping(addr+1, m)
	return (Word(hi) << 8) | Word(lo)
}

/*
func Os9StringWithMapping(addr Word, m Mapping) string {
	var buf bytes.Buffer
	for {
		var b byte = PeekBWithMapping(addr, m)
		var ch byte = 0x7F & b
		if '!' <= ch && ch <= '~' {
			buf.WriteByte(ch)
		} else {
			break
		}
		if (b & 128) != 0 {
			break
		}
		addr++
	}
	return buf.String()
}
*/

func ExplainMMU() string {
	romText := ""
	if sam.TyAllRam {
		romText += "AllRam, "
	} else {
		if BitFixedFExx {
			romText += "stuckFExx, "
		}
		if BitMC1 {
			if BitMC0 {
				romText += "ROM=32kExt, "
			} else {
				romText += "ROM=32kInt, "
			}
		} else {
			romText += "ROM=16kInt+16kExt, "
		}
	}
	if sam.P1RamSwap != 0 {
		romText += "RAM=Swapped32k32k, "
	}
	return F("mmu:%d task:%d [[ %02x %02x %02x %02x  %02x %02x %02x %02x || %02x %02x %02x %02x  %02x %02x %02x %02x ]] Init=%02x,%02x  %s  (%02x %02x %02x %02x  %02x %02x %02x %02x) (%02x %02x %02x %02x  %02x %02x %02x %02x) %02x %02x %02x %02x %02x %02x / %02x %02x %02x %02x  %02x %02x %02x %02x / %02x%02x%02x",
		Cond(MmuEnable, 1, 0),
		MmuTask&1,
		MmuMap[0][0],
		MmuMap[0][1],
		MmuMap[0][2],
		MmuMap[0][3],
		MmuMap[0][4],
		MmuMap[0][5],
		MmuMap[0][6],
		MmuMap[0][7],
		MmuMap[1][0],
		MmuMap[1][1],
		MmuMap[1][2],
		MmuMap[1][3],
		MmuMap[1][4],
		MmuMap[1][5],
		MmuMap[1][6],
		MmuMap[1][7],
		Init0,
		Init1,
		romText,

		portMem[0x90],
		portMem[0x91],
		portMem[0x92],
		portMem[0x93],
		portMem[0x94],
		portMem[0x95],
		portMem[0x96],
		portMem[0x97],

		portMem[0x98],
		portMem[0x99],
		portMem[0x9a],
		portMem[0x9b],
		portMem[0x9c],
		portMem[0x9d],
		portMem[0x9e],
		portMem[0x9f],

		PeekB(0x00EA),
		PeekB(0x00EB),
		PeekB(0x00EC),
		PeekB(0x00ED),
		PeekB(0x00EE),
		PeekB(0x00EF),

		PeekB(0x00F0),
		PeekB(0x00F1),
		PeekB(0x00F2),
		PeekB(0x00F3),

		PeekB(0x00F4),
		PeekB(0x00F5),
		PeekB(0x00F6),
		PeekB(0x00F7),

		PeekB(0xD938),
		PeekB(0xD939),
		PeekB(0xD93A),
	)
}

func MapAddrWithMapping(logical Word, m Mapping) int {
	slot := 7 & (logical >> 13)
	low := int(logical & 0x1FFF)
	physicalPage := m[slot]
	return (int(physicalPage) << 13) | low
}

func MapAddr(logical Word, quiet bool) int {
	slot := byte(logical >> 13)
	low := int(logical & 0x1FFF)
	var physicalPage byte

	if BitFixedFExx && logical >= 0xFE00 {
		physicalPage = 0x3F
	} else if logical >= 0xFF00 {
		physicalPage = 0x3F
	} else if MmuEnable {
		physicalPage = MmuMap[MmuTask][slot]
	} else {
		physicalPage = DisabledMmuMap[slot]
	}

	z := (int(physicalPage) << 13) | low

	if !quiet && TraceMem {
		L("\t\t\t\t\t\t MapAddr: %04x -> %06x ... task=%x  slot=%x  page=%x", logical, z, MmuTask, slot, physicalPage)
	}
	return z
}

// B is fundamental func to get byte.  Hack register access into here.
func B(addr Word) byte {
	var z byte
	mapped := MapAddr(addr, false)

	if AddressInDeviceSpace(addr) {
		z = GetIOByte(addr)
		Logd("GetIO (%06x) %04x -> %02x : %c %c", mapped, addr, z, PrettyH(z), PrettyT(z))
		mem[mapped] = z
	} else {
		z = PeekB(addr)
	}
	if TraceMem {
		L("\t\t\t\tGetB (%06x) %04x -> %02x : %c %c", mapped, addr, z, PrettyH(z), PrettyT(z))
	}
	return z
}

func PeekB(addr Word) byte {
	var z byte
	//T(Fmt("$%04x", addr))

	if !sam.TyAllRam && AddressInRomSpace(addr) {
		if UseExternalRomAssumingRom(addr) {
			//T()
			o := ExternalRomOffset(addr)
			z = cartRom[o]
		} else {
			//T()
			o := InternalRomOffset(addr)
			z = internalRom[o]
		}
	} else {
		//T()
		// Not ROM so use RAM
		mapped := MapAddr(addr, true)
		z = mem[mapped]
		//T("PeekB", Hex(addr), Hex(mapped), Hex(z))
	}
	//T(z)
	return z
}

func PokeB(addr Word, x byte) {
	if !sam.TyAllRam && AddressInRomSpace(addr) {
		if addr < 0xFF00 {
			log.Panicf("PokeB: write to ROM addr $%04x <- $%02x", addr, x)
		}
	}

	mapped := MapAddr(addr, true)
	mem[mapped] = x
}

// PutB is fundamental func to set byte.  Hack register access into here.
func PutB(addr Word, x byte) {
	if !sam.TyAllRam && AddressInRomSpace(addr) {
		if addr < 0xFF00 {
			log.Panicf("PutB: write to ROM addr $%04x <- $%02x", addr, x)
		}
	}

	mapped := MapAddr(addr, false)
	display.Poke(uint(addr), uint(mapped), x)

	old := mem[mapped]
	mem[mapped] = x
	if TraceMem {
		Logd("\t\t\t\tPutB (%06x) %04x <- %02x (was %02x)", mapped, addr, x, old)
	}
	if addr >= 0xfff0 { // XXX
		L("\t\t\t\tPutB VECTOR (%06x) %04x <- %02x (was %02x)", mapped, addr, x, old)
	}
	if AddressInDeviceSpace(addr) {
		PutIOByte(addr, x)
		Logd("PutIO (%06x) %04x <- %02x (was %02x)", mapped, addr, x, old)
	}
}

func PeekWPhys(addr int) Word {
	if addr+1 > len(mem) {
		panic(addr)
		// return 0
	}
	return Word(mem[addr])<<8 | Word(mem[addr+1])
}

//////// DUMP

func DoDumpAllMemoryPhys() {
	if !V['p'] {
		return
	}
	var i, j int
	var buf bytes.Buffer
	L("\n#DumpAllMemoryPhys(\n")
	n := len(mem)
	for i = 0; i < n; i += 32 {
		if i&0x1FFF == 0 {
			L("P [%02x] %06x:", i>>13, i)
		}
		// Look ahead for something interesting on this line.
		something := false
		for j = 0; j < 32; j++ {
			x := mem[i+j]
			// if x != 0 && x != ' ' //
			if x != 0 {
				something = true
				break
			}
		}

		if !something {
			continue
		}

		buf.Reset()
		Z(&buf, "P %06x: ", i)
		for j = 0; j < 32; j += 8 {
			Z(&buf,
				"%02x%02x %02x%02x %02x%02x %02x%02x  ",
				mem[i+j+0], mem[i+j+1], mem[i+j+2], mem[i+j+3],
				mem[i+j+4], mem[i+j+5], mem[i+j+6], mem[i+j+7])
		}
		buf.WriteRune(' ')
		for j = 0; j < 32; j++ {
			ch := 0x7F & mem[i+j]
			var r rune = '.'
			if ' ' <= ch && ch <= '~' {
				r = rune(ch)
			}
			buf.WriteRune(r)
		}
		L("%s\n", buf.String())
	}
	L("#DumpAllMemoryPhys)\n")
}

func DoExplainMmuBlock(i int) {
	blk := (i >> 13) & 0x3F
	blkPhys := MmuMap[MmuTask][blk]
	L("[%x -> %02x] %06x", blk, blkPhys, MapAddr(Word(i), true))
}

func DoDumpBlockZero() {
	PrettyDumpHex64(0, 0xFF00)
}

func DoDumpPathDesc(a Word) {
	PrettyDumpHex64(a, 0x40)
	if 0 == B(a+sym.PD_PD) {
		return
	}
	pd_pd := B(a + sym.PD_PD)
	if pd_pd > 32 {
		// Doesn't seem likely > 32
		L("???????? PathDesc %x @%x: mode=%x count=%x entry=%x\n", pd_pd, a, B(a+sym.PD_MOD), B(a+sym.PD_CNT), W(a+sym.PD_DEV))
		return
	}

	L("PathDesc %x @%x: mode=%x count=%x entry=%x\n", pd_pd, a, B(a+sym.PD_MOD), B(a+sym.PD_CNT), W(a+sym.PD_DEV))
	L("   curr_process=%x regs=%x buf=%x  dev_type=%x\n",
		B(a+sym.PD_CPR), W(a+sym.PD_RGS), W(a+sym.PD_BUF), B(a+sym.PD_DTP))

	// the Device Table Entry:
	dev := W(a + sym.PD_DEV)
	var buf bytes.Buffer
	Z(&buf, "   dev: @%x driver_mod=%x=%s ",
		dev, W(dev+sym.V_DRIV), ModuleName(W(dev+sym.V_DRIV)))
	Z(&buf, "driver_static_store=%x descriptor_mod=%x=%s ",
		W(dev+sym.V_STAT), W(dev+sym.V_DESC), ModuleName(W(dev+sym.V_DESC)))
	Z(&buf, "file_man=%x=%s use=%d\n",
		W(dev+sym.V_FMGR), ModuleName(W(dev+sym.V_FMGR)), B(dev+sym.V_USRS))
	L("%s", buf.String())

	if false && paranoid {
		if B(a+sym.PD_PD) > 10 {
			panic("PD_PD")
		}
		if B(a+sym.PD_CNT) > 20 {
			panic("PD_CNT")
		}
		if B(a+sym.PD_CPR) > 10 {
			panic("PD_CPR")
		}
	}
}

func DoDumpAllPathDescs() {
	if true || Level == 1 {
		p := W(sym.D_PthDBT)
		if 0 == p {
			L("DoDumpAllPathDescs: D_PthDPT is zero.")
			return
		}
		AssertEQ(p&255, 0, p)
		PrettyDumpHex64(p, 64)

		for i := Word(0); i < 64; i++ {
			q := Word(B(p + i)) << 8
			if q != 0 {
				L("PathDesc[%x]: %x", i, q)

				for j := Word(0); j < 4; j++ {
					k := i*4 + j
					if k == 0 {
						continue
					} // There is no path desc 0 (it's the table of allocs).
					n := B(q + j*64)
					if n == 0 {
						continue
					}
					AssertEQ(Word(n), k)
					L("........[%x]: %x", j, k)
					DoDumpPathDesc(q + j*64)
				}

			}
		}
	}
}

func DoDumpProcesses() {
	saved_mmut := MmuTask
	MmuTask = 0
	saved_map00 := MmuMap[0][0]
	MmuMap[0][0] = 0
	defer func() {
		MmuTask = saved_mmut
		MmuMap[0][0] = saved_map00
	}()
	///////////////////////////////////
	p := W(sym.D_PrcDBT)
	AssertNE(p, 0)
	AssertEQ(p&255, 0, p)
	PrettyDumpHex64(p, 64)

	for i := 0; i < 64; i++ {
		pg := B(p + Word(i))
		if pg == 0 {
			break
		}
		DoDumpProcDesc(Word(pg)<<8, F("TABLE_%d", i), false)
	}

	///////////////////////////////////

	if W(sym.D_Proc) != 0 {
		DoDumpProcDesc(W(sym.D_Proc), "Current", false)
	}
	if W(sym.D_AProcQ) != 0 {
		// L("D_AProcQ: Active:")
		DoDumpProcDesc(W(sym.D_AProcQ), "ActiveQ", true)
	}
	if W(sym.D_WProcQ) != 0 {
		// L("D_WProcQ: Wait:")
		DoDumpProcDesc(W(sym.D_WProcQ), "WaitQ", true)
	}
	if W(sym.D_SProcQ) != 0 {
		// L("D_SProcQ: Sleep")
		DoDumpProcDesc(W(sym.D_SProcQ), "SleepQ", true)
	}
}

func GetGimeIOByte(a Word) byte {
	if 0xFFA0 <= a && a <= 0xFFBF {
		log.Printf("GetGimeIOByte: readable range: $%04x -> $%04x.")
		return portMem[a&0xFF]
	}
	switch a {
	case 0xFF92: /* GIME IRQ */
		Logd("GIME -- Read FF92 (IRQ)")
		var z byte

		if gimeVirtPending {
			z |= 0x08
		}
		gimeVirtPending = false

		if gimeHorzPending {
			z |= 0x10
		}
		gimeHorzPending = false

		return z
		//case 0xFF93: /* GIME FIRQ */
		//Logd("GIME -- Read FF93 (FIRQ) NOT IMP")
	}
	log.Panicf("GetGimeIOByte %04x", a)
	return 0 // NOTREACHED
}

func PutGimeIOByte(a Word, b byte) {
	L("GIME %x <= %02x", a, b)
	PokeB(a, b)

	switch a {
	default:
		log.Panicf("UNKNOWN PutGimeIOByte address: 0x%04x", a)

	case 0xFFB0,
		0xFFB1,
		0xFFB2,
		0xFFB3,
		0xFFB4,
		0xFFB5,
		0xFFB6,
		0xFFB7,
		0xFFB8,
		0xFFB9,
		0xFFBA,
		0xFFBB,
		0xFFBC,
		0xFFBD,
		0xFFBE,
		0xFFBF:
		L("GIME\t\t$%x: palette[$%x] <- %s", a, a&15, ExplainColor(b))

	case 0xFFD9:
		L("GIME\t\t$%x: Cpu Speed <- %02x", a, b)

	case 0xFF90:
		Init0 = b
		MmuEnable = 0 != (b & 0x40)
		BitFixedFExx = 0 != (b & 0x08)
		BitMC1 = 0 != (b & 0x02)
		BitMC0 = 0 != (b & 0x01)
		L("GIME MmuEnable <- %v; MC=%d", MmuEnable, (b & 3))

	case 0xFF91:
		Init1 = b
		MmuTask = b & 0x01
		L("GIME MmuTask <- %v; clock rate <- %v", MmuTask, 0 != (b&0x40))

	case 0xFF92:
		L("GIME\t\tIRQ bits: %s", ExplainBits(b, FF92Bits))
		// 0x08: Vertical IRQ.  0x01: Cartridge.
		if (b & 0x08) != 0 {
			GimeVirtSyncInterruptEnable = true
		} else {
			GimeVirtSyncInterruptEnable = false
		}
		new_b := b &^ 0x08 // zap Virt bit
		new_b &^= 0x01     // also zap Cart bit
		if new_b != 0 {
			log.Panicf("GIME IRQ Enable for unsupported emulated bits: %04x %02x", a, b)
		}

	case 0xFF93:
		L("GIME\t\tFIRQ bits: %s", ExplainBits(b, FF93Bits))
		if b != 0 {
			log.Panicf("GIME FIRQ Enable for unsupported emulated bits: %04x %02x", a, b)
		}

	case 0xFF94:
		L("GIME\t\tTimer=$%x Start!", HiLo(PeekB(0xFF94), PeekB(0xFF95)))
	case 0xFF95:
		L("GIME\t\tTimer=$%x", HiLo(PeekB(0xFF94), PeekB(0xFF95)))
	case 0xFF96:
		L("GIME\t\treserved")
	case 0xFF97:
		L("GIME\t\treserved")
	case 0xFF98:
		L("GIME\t\tGraphicsNotAlpha=%x AttrsIfAlpha=%x Artifacting=%x Monochrome=%x 50Hz=%x LinesPerCharRow=%x=%d.",
			(b>>7)&1,
			(b>>6)&1,
			(b>>5)&1,
			(b>>4)&1,
			(b>>3)&1,
			(b & 7),
			GimeLinesPerCharRow[b&7])
	case 0xFF99:
		L("GIME\t\tLinesPerField=%x=%d. HRES=%x CRES=%x",
			(b>>5)&3,
			GimeLinesPerField[(b>>5)&3],
			(b>>2)&7,
			b&3)

	case 0xFF9A:
		L("GIME\t\tBorder: %s", ExplainColor(b))
	case 0xFF9B:
		L("GIME\t\t512K bank selector: %02x", b)
	case 0xFF9C:
		L("GIME\t\tVirt Scroll (alpha) = %x", b&15)
	case 0xFF9D,
		0xFF9E:
		L("GIME\t\tVirtOffsetAddr=$%05x",
			uint64(HiLo(PeekB(0xFF9D), PeekB(0xFF9E)))<<3)
	case 0xFF9F:
		L("GIME\t\tHVEN=%x HorzOffsetAddr=%x", (b >> 7), b&127)

	case 0xFFA0,
		0xFFA1,
		0xFFA2,
		0xFFA3,
		0xFFA4,
		0xFFA5,
		0xFFA6,
		0xFFA7,
		0xFFA8,
		0xFFA9,
		0xFFAA,
		0xFFAB,
		0xFFAC,
		0xFFAD,
		0xFFAE,
		0xFFAF:
		{
			task := byte((a >> 3) & 1)
			slot := byte(a & 7)
			was := MmuMap[task][slot]
			MmuMap[task][slot] = b & 0x3F
			L("GIME MmuMap[%d][%d] <- %02x  (was %02x)", task, slot, b, was)
			// if task == 0 && slot == 7 && b != 0x3F {
			// panic("bad MmuMap[0][7]")
			// }
			// yak ddt TODO
			// MmuMap[0][7] = 0x3F // Never change slot 7.
			// MmuMap[1][7] = 0x3F // Never change slot 7.
		}

	}
}
func ModuleId(begin Word, m Mapping) string {
	namePtr := begin + PeekWWithMapping(begin+4, m)
	modname := strings.ToLower(Os9StringWithMapping(namePtr, m))
	sz := PeekWWithMapping(begin+2, m)
	crc1 := PeekBWithMapping(begin+sz-3, m)
	crc2 := PeekBWithMapping(begin+sz-2, m)
	crc3 := PeekBWithMapping(begin+sz-1, m)
	return fmt.Sprintf("%s.%04x%02x%02x%02x", modname, sz, crc1, crc2, crc3)
}

func WithKernelTask(fn func()) {
	saved_mmut := MmuTask
	MmuTask = 0
	saved_map00 := MmuMap[0][0]
	MmuMap[0][0] = 0
	defer func() {
		MmuTask = saved_mmut
		MmuMap[0][0] = saved_map00
	}()

	fn()
}

func IsTermPath(path byte) bool {
	isTerm := false
	kpath := path
	task := MmuTask & 1
	WithMmuTask(0, func() {
		proc := PeekW(sym.D_Proc)
		procID := PeekB(proc + sym.P_ID)
		if task == 1 {
			// User mode: translate path to kernel path.
			kpath = PeekB(proc + P_Path + Word(path))
		}
		pathDBT := PeekW(sym.D_PthDBT)
		// fmt.Printf(" [dbt:%x] ", pathDBT)

		for i := Word(0); i < 8; i++ {
			// fmt.Printf(" [%x:%x] ", i, PeekW(pathDBT+2*i))
		}

		var pdPage Word
		if kpath > 3 {
			// no // pdPage = PeekW(pathDBT + 2*(Word(kpath)>>2))  // Use indirect DBT page.
			// wait -- shouldnt we just peek the high byte?
			pdPage = Word(PeekB(pathDBT+(Word(kpath)>>2))) << 8 // Use indirect DBT page.
		} else {
			pdPage = pathDBT // Use main DBT page.
		}
		if pdPage != 0 { // this should always be true.
			pd := pdPage + 64*(Word(kpath)&3)
			dev := PeekW(pd + sym.PD_DEV)
			desc := PeekW(dev + sym.V_DESC)
			name := ModuleName(desc)
			_ = procID
			// fmt.Printf("<<< #%d %x.t%x.p%x/kpath=%x/dbt=%x/page=%x/pd=%x/dev=%x/desc=%x/name=%s>>>\n", Cycles, procID, task, path, kpath, pathDBT, pdPage, pd, dev, desc, name) // yak
			if (pdPage & 255) != 0 {
				// CoreDump(fmt.Sprintf("/tmp/core#%d", Cycles))
			}
			isTerm = (name == *FlagTerm)
		}
	})
	return isTerm
}

func InitializeVectors() {
	panic("dont InitializeVectors")
	PutW(0xFFF2, 0xFEEE) // SWI3
	PutW(0xFFF4, 0xFEF1) // SWI2
	PutW(0xFFFA, 0xFEFA) // SWI
	PutW(0xFFFC, 0xFEFD) // NMI
	PutW(0xFFF8, 0xFEF7) // IRQ
	PutW(0xFFF6, 0xFEF4) // FIRQ
}
