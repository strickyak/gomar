package emu

import (
	//"github.com/strickyak/gomar/display"
	//"github.com/strickyak/gomar/sym"
	. "github.com/strickyak/gomar/gu"

	//"bufio"
	//"bytes"
	//"flag"
	//"fmt"
	//"io/ioutil"
	"log"
	//"os"
	//"regexp"
	//"sort"
	//"strconv"
	//"strings"
)

var _ = Log // use gu

type Word uint16

// EA is Effective Address, which may be a Word or a special value for a register.
type EA uint32

// 16 bit (0x08 bit clear)
const DRegEA EA = 0x10000000
const XRegEA EA = 0x10000001
const YRegEA EA = 0x10000002
const URegEA EA = 0x10000003
const SRegEA EA = 0x10000004
const PCRegEA EA = 0x10000005

// 8 bit (0x08 bit set)
const ARegEA EA = 0x10000008
const BRegEA EA = 0x10000009
const CCRegEA EA = 0x1000000A
const DPRegEA EA = 0x1000000B

func Hi(a Word) byte {
	return byte(255 & (a >> 8))
}
func Lo(a Word) byte {
	return byte(255 & a)
}
func HiLo(hi, lo byte) Word {
	return (Word(hi) << 8) | Word(lo)
}
func HiMidLo(hi, mid, lo byte) uint {
	return (uint(hi) << 16) | (uint(mid) << 8) | uint(lo)
}

func SignExtend(a byte) Word {
	if (a & 0x80) != 0 {
		return 0xFF80 | Word(a)
	} else {
		return Word(a)
	}
}

// W is fundamental func to get Word.
func W(addr Word) Word {
	hi := B(addr)
	lo := B(addr + 1)
	return HiLo(hi, lo)
}

func PeekW(addr Word) Word {
	hi := PeekB(addr)
	lo := PeekB(addr + 1)
	return HiLo(hi, lo)
}

// PutW is fundamental func to set Word.
func PutW(addr, x Word) {
	PutB(addr, Hi(x))
	PutB(addr+1, Lo(x))
}

func (addr EA) GetB() byte {
	if (addr & 0xFFFF0000) != 0 {
		switch addr {
		case ARegEA:
			return GetAReg()
		case BRegEA:
			return GetBReg()
		case CCRegEA:
			return ccreg
		case DPRegEA:
			return dpreg
		default:
			log.Panicf("bad B_ea EA: 0x%x", addr)
			return 0
		}
	} else {
		x := B(Word(addr))
		TraceByte(addr, x)
		return x
	}
}

func (addr EA) PutB(x byte) {
	if (addr & 0xFFFF0000) != 0 {
		switch addr {
		case ARegEA:
			PutAReg(x)
		case BRegEA:
			PutBReg(x)
		case CCRegEA:
			ccreg = x
		case DPRegEA:
			dpreg = x
		default:
			log.Panicf("bad PutB_ea EA: 0x%x", addr)
		}
	} else {
		TraceByte(addr, x)
		PutB(Word(addr), x)
	}
}

func (addr EA) RegPtrW() *Word {
	switch addr {
	case DRegEA:
		return &dreg
	case XRegEA:
		return &xreg
	case YRegEA:
		return &yreg
	case URegEA:
		return &ureg
	case SRegEA:
		return &sreg
	case PCRegEA:
		return &pcreg
	default:
		log.Panicf("Unknown RegPtr EA: 0x%x", addr)
		return nil
	}
}

func (addr EA) GetW() Word {
	if (addr & 0xFFFF0000) != 0 {
		p := addr.RegPtrW()
		return *p
	} else {
		x := W(Word(addr))
		TraceWord(addr, x)
		return x
	}
}

func (addr EA) PutW(x Word) {
	if (addr & 0xFFFF0000) != 0 {
		p := addr.RegPtrW()
		*p = x
	} else {
		TraceWord(addr, x)
		PutW(Word(addr), x)
	}
}

// For using page 0 for system variables.
func SysMemW(a Word) Word {
	/*
		if a >= 0x2000 {
			log.Panicf("SysMemW: addr too big: %x", a)
		}
	*/
	return HiLo(mem[a], mem[a+1])
}
func SysMemB(a Word) byte {
	/*
		if a >= 0x2000 {
			log.Panicf("SysMemW: addr too big: %x", a)
		}
	*/
	return mem[a]
}

func GetAReg() byte  { return Hi(dreg) }
func GetBReg() byte  { return Lo(dreg) }
func PutAReg(x byte) { dreg = HiLo(x, Lo(dreg)) }
func PutBReg(x byte) { dreg = HiLo(Hi(dreg), x) }

func DoMemoryDumps() {
	log.Printf("# pre timer interrupt #")
	// DoDumpAllMemory()
	// log.Printf("# pre timer interrupt #")
	// DoDumpAllMemoryPhys()
	// log.Printf("# pre timer interrupt #")
}

func B0(addr Word) byte {
	var b byte
	WithMmuTask(0, func() {
		b = B(addr)
	})
	log.Printf("==== kern byte @%x -> %x", addr, b)
	return b
}

func W0(addr Word) Word {
	var w Word
	WithMmuTask(0, func() {
		w = W(addr)
	})
	log.Printf("==== kern word @%x -> %x", addr, w)
	return w
}

func B1(addr Word) byte {
	var b byte
	WithMmuTask(1, func() {
		b = B(addr)
	})
	log.Printf("==== kern byte @%x -> %x", addr, b)
	return b
}

func W1(addr Word) Word {
	var w Word
	WithMmuTask(1, func() {
		w = W(addr)
	})
	log.Printf("==== kern word @%x -> %x", addr, w)
	return w
}

func LoadRom(start Word, m []byte) {
	start = start & 0x7FFF
	size := Word(len(m))
	for i := Word(0); i < size; i++ {
		internalRom[start+i] = m[i]
	}
}

func LoadCart(m []byte) {
	size := Word(len(m))
	offset := Word(0)
	// If 16K or less, goes in second half of 32K cartRom.
	if len(m) <= 0x4000 {
		offset = 0x4000
	}
	for i := Word(0); i < size; i++ {
		cartRom[i+offset] = m[i]
	}
}

func Loadm(loadm []byte) Word {
	size := Word(len(loadm))
	i := Word(0)
	for i < size {
		switch loadm[i] {
		case 0x00:
			n := HiLo(loadm[i+1], loadm[i+2])
			p := HiLo(loadm[i+3], loadm[i+4])
			for j := Word(0); j < n; j++ {
				PokeB(p+j, loadm[i+5+j])
			}
			i += 5 + n
		case 0xFF:
			return HiLo(loadm[i+3], loadm[i+4])
		default:
			log.Panicf("bad clause in loadm: [%x]: %02x", i, loadm[i])
		}
	}
	panic("no end to loadm")
}

func PeekBWithInt(addr int) byte {
	return PeekB(Word(addr))
}
