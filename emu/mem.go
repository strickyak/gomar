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
	panic("DoMemoryDumps")
	// log.Printf("# pre timer interrupt #")
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

/*
Two Megabits and Beyond:

[9:40 AM]t3chD0c: I have a 64K CoCo2 with a socketed 6809, and an available 6309 CPU I could swap in. Is it worth it, or should I save the 6309 for my CoCo3 - once I get the guts to desolder the original CPU off the motherboard 😬
[9:59 AM]Ciaran: anyone know offhand how >2MB works on a coco3?  Richard Goedeken's stress tester seems to set the top 4 bits of $FF9B and i assume what that's doing is providing extra bits whenever there's a write to the bank registers?
[10:00 AM]Ciaran: other sources suggest the lowest 2 (but i guess more where relevant) bits control which 512K bank the vram addresses end up (and that's what xroar does so far up to 2MB)
[10:02 AM]Ciaran: are there any other subtleties than that?  e.g. could imagine only supporting >2MB for one of the tasks to save space
[10:14 AM]Dave Philipsen: I say go ahead and use the 6309 in your CoCo 2. There are still plenty to be found.
[10:16 AM]MDSiegel: I had assumed that without a substition of the GIME anything above 512 was banked
[10:19 AM]Dave Philipsen: There are add ons, I believe, which will allow those extra two unused bits to be used while still using the original GIME.
[10:21 AM]Dave Philipsen: Since the registers are write-only it’s a fairly simple matter to decode the FFA0-F addresses and mirror the GIME to provide the extra two bits for addressing up to 2MB.
[10:54 AM]Ciaran: yeah up to 2MB is straightforward - the GIME doesn't use those top two bits when written, nor assert them when reading.  >2MB means using some other location.  i see $FF9B used, and i wonder if whatever you write there (bits 4 & 5) get used to supplement writes to the more usual bank register addresses
[10:54 AM]Ciaran: probably a nitros-9 dev would know for sure @L. Curtis Boyle ? :)
[10:56 AM]Ciaran: e.g. you write $10 to $FF9B, that doesn't change anything right away, but if you write $FF to $FFA0, it effectively makes the bank used for that address space $1FF
[10:56 AM]Ciaran: i suppose i could just implement it that way, run nitros-9, see whether it reports > 2MB ;)
[11:30 AM]Deek: The OS doesn't use >2MB at all
[11:32 AM]Ciaran: ah right istr reading that - just sets it up as a ram disk?
[11:32 AM]Ciaran: does EOU do that if detected, or is that something you have to do manually?
[11:33 AM]Deek: There's a "NoCan" ramdisk driver that can be used
[11:33 AM]Ciaran: hm, sounds like you need to know what you're doing
[11:33 AM]Deek: I don't know how it keeps from swapping out its own code though
[11:34 AM]Deek: Or if it needs to
[11:34 AM]Ciaran: yeah that's why i assumed the above - that a write to that address is latched and then stored along with writes to the bank registers, not applied straight away
[11:34 AM]Ciaran: couldn't imagine how else it could work
[11:35 AM]Ciaran: unless you stored some trampoliney stuff in $FExx maybe
[11:35 AM]Ciaran: but that sounds horrid
[11:38 AM]Deek: In the NitrOS-9 repo, 3rdparty/drivers/nocan/rammer.asm
[11:42 AM]Ciaran: reading that my first thought is "you do know a tab character takes as much space as a space character, right?"
[11:44 AM]Ciaran: the second is that there are clearly at least two different ways this has been done
[11:45 AM]Ciaran: yeah make that three different ways
[11:46 AM]Ciaran: but it does seem to tally - ff9b or ff80 or ff70 is updated with appropriate bits, then the bank reg in ffax is written, then ff9b or ff80 or ff70 is rewritten with whatever it was before
[11:48 AM]Ciaran: so unless the extended bank register is only affecting bank 0/task 0 (which is the bank the driver seems to play with), seems reasonable to assume it's just a latch that affects subsequent writes
[11:49 AM]Deek: I am working on that
[11:54 AM]Ciaran: heh, wasn't really blaming anyone, just going "jeez that's hard to read"
[11:55 AM]Deek: I know, it's extremely wide, but that's what boisy went with.
[11:57 AM]Deek: It will soon switch to the other extreme, with no extra spacing at all and a script that can be used to prettify it however you like.
[11:57 AM]Ciaran: oh you know what i think we're looking at different things
[11:57 AM]Ciaran: what i'm seeing is single space before the opcode
[11:58 AM]Ciaran: because i checked this out $whoknowshowlong ago.  lemme do a pull...
[11:58 AM]Ciaran: realises again it's not git
[11:58 AM]Ciaran: checks to see what the bloody mercurial equivalents are
[11:58 AM]Deek: Oh, then you're looking at how it was originally written.
[12:00 PM]Ciaran: ok i've just remembered all the recent chatter about putting this in github
[12:00 PM]Ciaran: as an hg pull -u didn't get anything (which is what the ever useful https://wiki.mercurial-scm.org/GitConcepts reminded me of ;)
[12:01 PM]Deek: It will probably wind up (slightly) harder to read than you're seeing now, because the comments will get cuddled up too.
[12:01 PM]Deek: Nocan set 1 0=64Meg Nocan 1=8Meg MESS and Nocan3 2=16Meg Collyer
[12:02 PM]Ciaran: https://github.com/nitros9project/nitros9.git yeah?  not mikey's?
[12:02 PM]Deek: Yes
[12:03 PM]Ciaran: ok no that's much nicer formatting :)
[12:03 PM]Ciaran: i personally stopped tabbing the operands, but that's just personal preference.  still readable.
[12:04 PM]Deek: Once I'm done, setting up the clean and smudge filters and 'renormalizing' will get you 'preferred formatting'
[12:09 PM]Ciaran: while i feel "spaces only, no tabs" to be a perverse inversion of moral correctitude, i'm unlikely to be editing nitros9 source any time soon
[12:14 PM]L. Curtis Boyle: Here are the notes by Robert Gault and Paul Barton for when the NoCan3 was made:
My Notes on NoCan3 & NoCan4 This 8MB interface is just an extension to the already existing 2MB interface.

A 2MB interface is not needed, this board emulates the original 2MB interface.
All 2MB bits work as before, no changes.

For NitrOS-9 users, no changes are required to use the 2MB.
Expand
NoCan3_docs.txt
3 KB
[12:18 PM]Ciaran: excellent, cheers - that seems to confirm it.  top bits of FF9B are latched, used for later writes to FFAx.  lower bits of FF9B are used directly as bank for video.
[12:18 PM]Ciaran: all sounds eminently implementable
[12:19 PM]Ciaran: as for all the other approaches mentioned in that source - nocan64 and "collyer", i think i'll just ignore those
*/

/*
Where exactly does memory get overwritten in the Cartridge space?

546369 #d M c000: 444b 1a50 5f1f 9b30  8d00 5910 8e04 00a6  8484 3fa7 a0a6 802a  f617 180b 1a50 308d   DK.P_..0..Y....&..?' &.*v....P0.
546370 #d M c020: 0042 108e 0500 a684  843f a7a2 a680 2af6  20fe ff7f f8ce 7ff8  1ad0 363f ee8d 000a   .B....&..?'"&.*v ~..xN.x.P6?n...
546371 #d M c040: ff7f fa10 ff7f fcfe  7ff8 39be 7ffa bf7f  fe8e c05c bf7f fa10  ce7f f03b 10fe 7ffc   ..z...|~.x9>.z?.~.@\?.z.N.p;.~.|
546372 #d M c060: 6e9f 7ffe 202d 2d20  5354 5249 434b 5941  4b20 4652 4f42 494f  2050 5245 424f 4f54   n..~ -- STRICKYAK FROBIO PREBOOT
546373 #d M c080: 202d 2da0 ffff ffff  ffff ffff ffff ffff  ffff ffff ffff ffff  ffff ffff ffff ffff    -- ............................
546374 #d M c0a0: ffff ffff ffff ffff  ffff ffff ffff ffff  ffff ffff ffff ffff  ffff ffff ffff ffff   ................................
546375 #d M c0c0: ffff ffff ffff ffff  ffff ffff ffff ffff  ffff ffff ffff ffff  ff7e e29d ffff ffff   .........................~b.....   <==
546376 #d M c0e0: ffff ffff ffff ffff  ffff ffff ffff ffff  ffff ffff ffff ffff  ffff ffff ffff ffff   ................................

546433 #d M c800: 1f02 3440 aee8 1134  10ae 6bbd d6e6 3420  ece8 1533 cb34 40ee  68ae 44bd d6e6 3268   ..4@.h.4..k=Vf4 lh.3K4@nh.D=Vf2h
546434 #d M c820: 2011 ece8 1134 06ae  e811 3410 ae6b bdd6  e632 64ae e4ec 84c3  0028 ed64 ec62 e3e8    .lh.4..h.4..k=Vf2d.dl.C.(mdlbch
546435 #d M c840: 1134 06ae 66bd d79c  c640 ae62 bdd7 0132  6220 04c6 01e7 66e6  6632 6935 e034 6032   .4..f=W.F@.b=W.2b .F.gfff2i5`4`2
546436 #d M c860: 7e10 ae6a 3384 cc00  0020 09ec 68e7 c0ec  e4c3 0001 ede4 10ac  e426 f032 6235 e0ff   ~..j3.L.. .lhg@ldC..md.,d&p2b5`.
546437 #d M c880: ffff ffff ffff ffff  ffff ffff ffff ffff  ffff ffff ffff ffff  ffff ffff ffff ffff   ................................
546438 #d M c8a0: ffff ffff ffff ffff  ffff ffff ffff ffff  ffff ffff 1212 1212  1212 1212 1212 12ff   ................................   <==
546439 #d M c8c0: ffff ffff ffff ffff  ffff ffff ffff ffff  ffff ffff ffff ffff  ffff ffff ffff ffff   ................................
546440 #d M c8e0: ffff ffff ffff ffff  ffff ffff ffff ffff  ffff ffff ffff ffff  ffff ffff ffff ffff   ................................
546441 #d M c900: 3440 8eff 23e6 84c4  fbe7 84ce ff22 e6c4  ca02 e7c4 e684 ca04  e784 e6c4 ca02 e7c4   4@..#f.D{g.N."fDJ.gDf.J.g.fDJ.gD
546442 #d M c920: 8e00 0abd c946 e6c4  c4fd e7c4 8e00 0a35  407e c946 327e 4f5f  8e98 211f 161f 60ed   ...=IFfDD}gD...5@~IF2~O_..!...`m
546443 #d M c940: e4ae e432 6239 3440  3384 bdce 615d 2709  200d 3d3d 3d3d 3d33  5f11 8300 0026 f335   d.d2b94@3.=Na]'. .=====3_....&s5

Clock Speeds:

$ livy
      14.31818 / 8
   _0 = (*livy.Num)
1.7897725
      14.31818 / 16
   _1 = (*livy.Num)
0.89488625

*/
