package emu

import (
	//"github.com/strickyak/gomar/display"
	//"github.com/strickyak/gomar/sym"
	. "github.com/strickyak/gomar/gu"

	//"bufio"
	"bytes"
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

var _ = Log

func ImmByte() byte {
	z := B(pcreg)
	pcreg++
	return z
}
func ImmWord() Word {
	hi := ImmByte()
	lo := ImmByte()
	return HiLo(hi, lo)
}

/* sreg */
func PushByte(b byte) {
	sreg--
	PutB(sreg, b)
}
func PushWord(w Word) {
	PushByte(Lo(w))
	PushByte(Hi(w))
}
func PullByte(bp *byte) {
	*bp = B(sreg)
	sreg++
}
func PullWord(wp *Word) {
	var hi, lo byte
	PullByte(&hi)
	PullByte(&lo)
	*wp = HiLo(hi, lo)
}

/* ureg */
func PushUByte(b byte) {
	ureg--
	PutB(ureg, b)
}
func PushUWord(w Word) {
	PushUByte(Lo(w))
	PushUByte(Hi(w))
}
func PullUByte(bp *byte) {
	*bp = B(ureg)
	ureg++
}
func PullUWord(wp *Word) {
	var hi, lo byte
	PullUByte(&hi)
	PullUByte(&lo)
	*wp = HiLo(hi, lo)
}

func nmi() {
	L("INTERRUPTING with NMI")
	interrupt(VECTOR_NMI)
	irqs_pending &^= NMI_PENDING
}

func illaddr() EA { // illegal addressing mode, defaults to zero //
	log.Panicf("Illegal Addressing Mode")
	panic(0)
}

var dixreg = []string{"x", "y", "u", "s"}

func ainc() EA {
	Dis_ops(",", dixreg[idx], 2)
	Dis_ops("+", "", 0)
	regPtr := ixregs[idx]
	z := *regPtr
	(*regPtr)++
	return EA(z)
}

func ainc2() EA {
	Dis_ops(",", dixreg[idx], 3)
	Dis_ops("++", "", 0)
	regPtr := ixregs[idx]
	z := *regPtr
	(*regPtr) += 2
	return EA(z)
}

func adec() EA {
	Dis_ops(",-", dixreg[idx], 2)
	regPtr := ixregs[idx]
	(*regPtr)--
	return EA(*regPtr)
}

func adec2() EA {
	Dis_ops(",--", dixreg[idx], 3)
	regPtr := ixregs[idx]
	(*regPtr) -= 2
	return EA(*regPtr)
}

func plus0() EA {
	Dis_ops(",", dixreg[idx], 0)
	return EA(*ixregs[idx])
}

func plusa() EA {
	Dis_ops("a,", dixreg[idx], 1)
	return EA((*ixregs[idx]) + SignExtend(GetAReg()))
}

func plusb() EA {
	Dis_ops("b,", dixreg[idx], 1)
	return EA((*ixregs[idx]) + SignExtend(GetBReg()))
}

func plusn() EA {
	off := ""
	b := ImmByte()
	/* negative offsets alway decimal, otherwise hex */
	if (b & 0x80) != 0 {
		off = F("%d,", int(b)-256)
	} else {
		off = F("$%02x,", b)
	}
	Dis_ops(off, dixreg[idx], 1)
	return EA((*ixregs[idx]) + SignExtend(b))
}

func plusnn() EA {
	w := ImmWord()
	off := F("$%04x,", w)
	Dis_ops(off, dixreg[idx], 4)
	return EA(*ixregs[idx] + w)
}

func plusd() EA {
	Dis_ops("d,", dixreg[idx], 4)
	return EA(*ixregs[idx] + dreg)
}

func npcr() EA {
	b := ImmByte()
	off := F("$%04x,pcr", (pcreg+SignExtend(b))&0xffff)
	Dis_ops(off, "", 1)
	return EA(pcreg + SignExtend(b))
}

func nnpcr() EA {
	w := ImmWord()
	off := F("$%04x,pcr", (pcreg+w)&0xffff)
	Dis_ops(off, "", 5)
	return EA(pcreg + w)
}

func direct() EA {
	w := ImmWord()
	off := F("$%04x", w)
	Dis_ops(off, "", 3)
	return EA(w)
}

func zeropage() EA {
	b := ImmByte()
	off := F("$%02x", b)
	Dis_ops(off, "", 2)
	return EA(HiLo(dpreg, b))
}

func immediate() EA {
	off := F("#$%02x", B(pcreg))
	Dis_ops(off, "", 0)
	z := pcreg
	pcreg++
	return EA(z)
}

func immediate2() EA {
	z := pcreg
	off := F("#$%04x", (Word(B(pcreg))<<8)|Word(B(pcreg+1)))
	Dis_ops(off, "", 0)
	pcreg += 2
	return EA(z)
}

var pbtable = []func() EA{
	ainc, ainc2, adec, adec2,
	plus0, plusb, plusa, illaddr,
	plusn, plusnn, illaddr, plusd,
	npcr, nnpcr, illaddr, direct}

func postbyte() EA {
	pb := ImmByte()
	idx = ((pb & 0x60) >> 5)
	if (pb & 0x80) != 0 {
		if (pb & 0x10) != 0 {
			Dis_ops("[", "", 3)
		}
		temp := (pbtable[pb&0x0f])()
		if (pb & 0x10) != 0 {
			temp = EA(temp.GetW())
			Dis_ops("]", "", 0)
		}
		return EA(temp)
	} else {
		temp := Word(pb & 0x1f)
		if (temp & 0x10) != 0 {
			temp |= 0xfff0 /* sign extend */
		}
		var off string
		if (temp & 0x10) != 0 {
			// Use int16 for negative signed number.
			// Sign-extend by or'ing with 0xF0.
			off = F("%d,", int16(0xF0|temp))
		} else {
			off = F("%d,", temp)
		}
		Dis_ops(off, dixreg[idx], 1)
		return EA(*ixregs[idx] + temp)
	}
}

func eaddr0() EA { // effective address for NEG..JMP //
	switch (ireg & 0x70) >> 4 {
	case 0:
		return zeropage()
	case 1, 2, 3: //canthappen//
		log.Panicf("UNKNOWN eaddr0: %02x\n", ireg)
		return 0
	case 4:
		Dis_inst_cat("a", -2)
		return ARegEA
	case 5:
		Dis_inst_cat("b", -2)
		return BRegEA
	case 6:
		Dis_inst_cat("", 2)
		return postbyte()
	case 7:
		return direct()
	}
	panic("notreached")
}

func eaddr8() EA { // effective address for 8-bits ops. //
	switch (ireg & 0x30) >> 4 {
	case 0:
		return immediate()
	case 1:
		return zeropage()
	case 2:
		Dis_inst_cat("", 2)
		return postbyte()
	case 3:
		return direct()
	}
	panic("notreached")
}

func eaddr16() EA { // effective address for 16-bits ops. //
	switch (ireg & 0x30) >> 4 {
	case 0:
		Dis_inst_cat("", -1)
		return immediate2()
	case 1:
		Dis_inst_cat("", -1)
		return zeropage()
	case 2:
		Dis_inst_cat("", 1)
		return postbyte()
	case 3:
		Dis_inst_cat("", -1)
		return direct()
	}
	panic("notreached")
}

func ill() {
	log.Panicf("Illegal Opcode: 0x%x", ireg)
}

// macros to set status flags //
func SEC() { ccreg |= 0x01 }
func CLC() { ccreg &= 0xfe }
func SEZ() { ccreg |= 0x04 }
func CLZ() { ccreg &= 0xfb }
func SEN() { ccreg |= 0x08 }
func CLN() { ccreg &= 0xf7 }
func SEV() { ccreg |= 0x02 }
func CLV() { ccreg &= 0xfd }
func SEH() { ccreg |= 0x20 }
func CLH() { ccreg &= 0xdf }

// set N and Z flags depending on 8 or 16 bit result //
func SETNZ8(b byte) {
	if b != 0 {
		CLZ()
	} else {
		SEZ()
	}
	if (b & 0x80) != 0 {
		SEN()
	} else {
		CLN()
	}
}
func SETNZ16(b Word) {
	if b != 0 {
		CLZ()
	} else {
		SEZ()
	}
	if (b & 0x8000) != 0 {
		SEN()
	} else {
		CLN()
	}
}

func SETSTATUS(a byte, b byte, res Word) {
	if ((a ^ b ^ byte(res)) & 0x10) != 0 {
		SEH()
	} else {
		CLH()
	}
	if ((a ^ b ^ byte(res) ^ byte(res>>1)) & 0x80) != 0 {
		SEV()
	} else {
		CLV()
	}
	if (res & 0x100) != 0 {
		SEC()
	} else {
		CLC()
	}
	SETNZ8(byte(res))
}

func CondB(b bool, x, y byte) byte {
	if b {
		return x
	} else {
		return y
	}
}
func CondW(b bool, x, y Word) Word {
	if b {
		return x
	} else {
		return y
	}
}
func CondI(b bool, x, y int) int {
	if b {
		return x
	} else {
		return y
	}
}
func CondS(b bool, x, y string) string {
	if b {
		return x
	} else {
		return y
	}
}

func AOrB(aIfZero byte) EA {
	if aIfZero == 0 {
		return ARegEA
	} else {
		return BRegEA
	}
}

func add() {
	var aop, bop, res Word
	Dis_inst("add", CondS(0 != (ireg&0x40), "b", "a"), 2)
	accum := AOrB(ireg & 0x40)
	aop = Word(accum.GetB())
	bop = Word(eaddr8().GetB())
	res = (aop) + (bop)
	SETSTATUS(byte(aop), byte(bop), res)
	accum.PutB(byte(res))
}

func sbc() {
	var aop, bop, res Word
	Dis_inst("sbc", CondS(0 != (ireg&0x40), "b", "a"), 2)
	accum := AOrB(ireg & 0x40)
	aop = Word(accum.GetB())
	bop = Word(eaddr8().GetB())
	res = aop - bop - Word(ccreg&0x01)
	SETSTATUS(byte(aop), byte(bop), res)
	accum.PutB(byte(res))
}

func sub() {
	var aop, bop, res Word
	Dis_inst("sub", CondS(0 != (ireg&0x40), "b", "a"), 2)
	accum := AOrB(ireg & 0x40)
	aop = Word(accum.GetB())
	bop = Word(eaddr8().GetB())
	res = aop - bop
	SETSTATUS(byte(aop), byte(bop), res)
	accum.PutB(byte(res))
}

func adc() {
	var aop, bop, res Word
	Dis_inst("adc", CondS(0 != (ireg&0x40), "b", "a"), 2)
	accum := AOrB(ireg & 0x40)
	aop = Word(accum.GetB())
	bop = Word(eaddr8().GetB())
	res = aop + bop + Word(ccreg&0x01)
	SETSTATUS(byte(aop), byte(bop), res)
	accum.PutB(byte(res))
}

func cmp() {
	var aop, bop, res Word
	Dis_inst("cmp", CondS(0 != (ireg&0x40), "b", "a"), 2)
	accum := AOrB(ireg & 0x40)
	aop = Word(accum.GetB())
	bop = Word(eaddr8().GetB())
	res = aop - bop
	SETSTATUS(byte(aop), byte(bop), res)
}

func and() {
	var aop, bop, res byte
	Dis_inst("and", CondS(0 != (ireg&0x40), "b", "a"), 2)
	accum := AOrB(ireg & 0x40)
	aop = (accum.GetB())
	bop = (eaddr8().GetB())
	res = aop & bop
	SETNZ8(res)
	CLV()
	accum.PutB(res)
}
func or() {
	var aop, bop, res byte
	Dis_inst("or", CondS(0 != (ireg&0x40), "b", "a"), 2)
	accum := AOrB(ireg & 0x40)
	aop = (accum.GetB())
	bop = (eaddr8().GetB())
	res = aop | bop
	SETNZ8(res)
	CLV()
	accum.PutB(res)
}
func eor() {
	var aop, bop, res byte
	Dis_inst("eor", CondS(0 != (ireg&0x40), "b", "a"), 2)
	accum := AOrB(ireg & 0x40)
	aop = (accum.GetB())
	bop = (eaddr8().GetB())
	res = aop ^ bop
	SETNZ8(res)
	CLV()
	accum.PutB(res)
}
func bit() {
	var aop, bop, res byte
	Dis_inst("bit", CondS(0 != (ireg&0x40), "b", "a"), 2)
	accum := AOrB(ireg & 0x40)
	aop = (accum.GetB())
	bop = (eaddr8().GetB())
	res = aop & bop
	SETNZ8(res)
	CLV()
}

func ld() {
	Dis_inst("ld", CondS(0 != (ireg&0x40), "b", "a"), 2)
	accum := AOrB(ireg & 0x40)
	res := eaddr8().GetB()
	SETNZ8(res)
	CLV()
	accum.PutB(res)
}

func st() {
	Dis_inst("st", CondS(0 != (ireg&0x40), "b", "a"), 2)
	accum := AOrB(ireg & 0x40)
	res := accum.GetB()
	eaddr8().PutB(res)
	SETNZ8(res)
	CLV()
}

func jsr() {
	Dis_inst("jsr", "", 5)
	Dis_len(-pcreg)
	w := eaddr8()
	Dis_len_incr(pcreg + 1)
	PushWord(pcreg)
	pcreg = Word(w)
}

func bsr() {
	b := ImmByte()
	Dis_inst("bsr", "", 7)
	Dis_len(2)
	PushWord(pcreg)
	pcreg += SignExtend(b)
	off := F("$%04x", pcreg&0xffff)
	Dis_ops(off, "", 0)
}

func neg() {
	var a, r Word

	{
		t := W(pcreg)
		if t == 0 {
			log.Panicf("Executing 0000 instruction at pcreg=%04x", pcreg-1)
			// log.Printf("Warning: Executing 0000 instruction at pcreg=%04x", pcreg-1)
		}
	}

	a = 0
	Dis_inst("neg", "", 4)
	ea := eaddr0()
	a = Word(ea.GetB())
	r = -a
	SETSTATUS(0, byte(a), r)
	ea.PutB(byte(r))
}

func com() {
	Dis_inst("com", "", 4)
	ea := eaddr0()
	r := ^(ea.GetB())
	SETNZ8(r)
	SEC()
	CLV()
	ea.PutB(r)
}

func lsr() {
	Dis_inst("lsr", "", 4)
	ea := eaddr0()
	r := ea.GetB()
	if (r & 0x01) != 0 {
		SEC()
	} else {
		CLC()
	}
	if (r & 0x10) != 0 {
		SEH()
	} else {
		CLH()
	}
	r >>= 1
	SETNZ8(r)
	ea.PutB(r)
}

func ror() {
	c := (ccreg & 0x01) << 7
	Dis_inst("ror", "", 4)
	ea := eaddr0()
	r := ea.GetB()
	if (r & 0x01) != 0 {
		SEC()
	} else {
		CLC()
	}
	r = (r >> 1) + c
	SETNZ8(r)
	ea.PutB(r)
}

func asr() {
	Dis_inst("asr", "", 4)
	ea := eaddr0()
	r := ea.GetB()
	if (r & 0x01) != 0 {
		SEC()
	} else {
		CLC()
	}
	if (r & 0x10) != 0 {
		SEH()
	} else {
		CLH()
	}
	r >>= 1
	if (r & 0x40) != 0 {
		r |= 0x80
	}
	SETNZ8(r)
	ea.PutB(r)
}

func asl() {
	var a, r Word

	Dis_inst("asl", "", 4)
	ea := eaddr0()
	a = Word(ea.GetB())
	r = a << 1
	SETSTATUS(byte(a), byte(a), r)
	ea.PutB(byte(r))
}

func rol() {
	c := (ccreg & 0x01)
	Dis_inst("rol", "", 4)
	ea := eaddr0()
	r := ea.GetB()
	if (r & 0x80) != 0 {
		SEC()
	} else {
		CLC()
	}
	if ((r & 0x80) ^ ((r << 1) & 0x80)) != 0 {
		SEV()
	} else {
		CLV()
	}
	r = (r << 1) + c
	SETNZ8(r)
	ea.PutB(r)
}

func inc() {
	Dis_inst("inc", "", 4)
	ea := eaddr0()
	r := ea.GetB()
	r++
	if r == 0x80 {
		SEV()
	} else {
		CLV()
	}
	SETNZ8(r)
	ea.PutB(r)
}

func dec() {
	Dis_inst("dec", "", 4)
	ea := eaddr0()
	r := ea.GetB()
	r--
	if r == 0x7f {
		SEV()
	} else {
		CLV()
	}
	SETNZ8(r)
	ea.PutB(r)
}

func tst() {
	Dis_inst("tst", "", 4)
	ea := eaddr0()
	r := ea.GetB()
	SETNZ8(r)
	CLV()
}

func jmp() {
	Dis_len(-pcreg)
	Dis_inst("jmp", "", 1)
	ea := eaddr0()
	Dis_len_incr(pcreg + 1)
	pcreg = Word(ea)
}

func clr() {
	Dis_inst("clr", "", 4)
	ea := eaddr0()
	ea.PutB(0)
	CLN()
	CLV()
	SEZ()
	CLC()
}

func flag0() {
	/*
		if iflag != 0 { // in case flag already set by previous flag instr don't recurse //
			pcreg--
			return
		}
	*/
	iflag = 1
	ireg = B(pcreg)
	pcreg++
	Dis_inst("", "", 1)
	(instructionTable[ireg])()
	iflag = 0
}

func flag1() {
	/*
		if iflag != 0 { // in case flag already set by previous flag instr don't recurse //
			pcreg--
			return
		}
	*/
	iflag = 2
	ireg = B(pcreg)
	pcreg++
	Dis_inst("", "", 1)
	(instructionTable[ireg])()
	iflag = 0
}

func nop() {
	Dis_inst("nop", "", 2)
}

func sync_inst() {
	L("sync_inst")
	Waiting = true
}

func cwai() {
	b := B(pcreg) // Immediate operand //
	ccreg &= b
	pcreg++

	L("Waiting, cwai #$%02x.", b)
	Waiting = true

	Dis_inst("cwai", "", 20)
	off := F("#$%02x", b)
	Dis_ops(off, "", 0)
}

func lbra() {
	w := ImmWord()
	pcreg += w
	Dis_len(3)
	Dis_inst("lbra", "", 5)
	off := F("$%04x", pcreg&0xffff)
	Dis_ops(off, "", 0)
}

func lbsr() {
	Dis_len(3)
	Dis_inst("lbsr", "", 9)
	w := ImmWord()
	PushWord(pcreg)
	pcreg += w
	off := F("$%04x", pcreg)
	Dis_ops(off, "", 0)
}

func daa() {
	var a Word
	Dis_inst("daa", "", 2)
	a = Word(GetAReg())
	if (ccreg & 0x20) != 0 {
		a += 6
	}
	if (a & 0x0f) > 9 {
		a += 6
	}
	if (ccreg & 0x01) != 0 {
		a += 0x60
	}
	if (a & 0xf0) > 0x90 {
		a += 0x60
	}
	if (a & 0x100) != 0 {
		SEC()
	}
	PutAReg(byte(a))
}

func orcc() {
	b := ImmByte()
	off := F("#$%02x", b)
	Dis_inst("orcc", "", 3)
	Dis_ops(off, "", 0)
	ccreg |= b
}

func andcc() {
	b := ImmByte()
	off := F("#$%02x", b)
	Dis_inst("andcc", "", 3)
	Dis_ops(off, "", 0)
	ccreg &= b
}

func mul() {
	w := Word(GetAReg()) * Word(GetBReg())
	Dis_inst("mul", "", 11)
	if (w) != 0 {
		CLZ()
	} else {
		SEZ()
	}
	if (w & 0x80) != 0 {
		SEC()
	} else {
		CLC()
	}
	dreg = (w)
}

func sex() {
	Dis_inst("sex", "", 2)
	w := SignExtend(GetBReg())
	SETNZ16(w)
	dreg = (w)
}

func abx() {
	Dis_inst("abx", "", 3)
	xreg += Word(GetBReg())
}

func rts() {
	Dis_inst("rts", "", 5)
	Dis_len(1)
	PullWord(&pcreg)

	if *FlagBasicText {
		ShowBasicText()
	}
}

func tfr() {
	Dis_inst("tfr", "", 7)
	b := ImmByte()
	Dis_reg(b)
	src := TfrReg(15 & (b >> 4))
	dst := TfrReg(15 & b)
	if (src & 8) != (dst & 8) {
		log.Panicf("tfr with inconsistent sizes; src=%d dst=%d", src, dst)
	}
	if (src & 8) == 0 {
		// 16 bit
		dst.PutW(src.GetW())
	} else {
		// 8 bit
		dst.PutB(src.GetB())
	}
}

func exg() {
	Dis_inst("exg", "", 8)
	b := ImmByte()
	Dis_reg(b)
	r1 := TfrReg(15 & (b >> 4))
	r2 := TfrReg(15 & b)
	if (b & 0x80) == 0 {
		// 16 bit
		t1, t2 := r1.GetW(), r2.GetW()
		r1.PutW(t2)
		r2.PutW(t1)
	} else {
		// 8 bit
		t1, t2 := r1.GetB(), r2.GetB()
		r1.PutB(t2)
		r2.PutB(t1)
	}
}

func br(f bool) {
	var dest Word

	if 0 == iflag {
		b := ImmByte()
		dest = pcreg + SignExtend(b)
		if f {
			pcreg += SignExtend(b)
		}
		Dis_len(2)
	} else {
		w := ImmWord()
		dest = pcreg + w
		if f {
			pcreg += w
		}
		Dis_len(3)
	}
	off := F("$%04x", dest&0xffff)
	Dis_ops(off, "", 0)
}

func NXORV() bool {
	return ((ccreg & 0x08) ^ (ccreg & 0x02)) != 0
}
func IFLAG() bool {
	return iflag != 0
}

func bra() {
	if iflag == 0 && B(pcreg) == 0xFE {
		// An infinite loop with the bra statement to itself
		// usually indicates an assertion failed.

		// DoDumpAllMemoryPhys()
		DumpAllMemory()
		log.Panicf("Panic: SELF-BRANCH at pc=$%04x", pcreg-1)
	}
	Dis_inst(CondS(IFLAG(), "l", ""), "bra", int64(CondI(IFLAG(), 5, 3)))
	br(true)
}

func brn() {
	Dis_inst(CondS(IFLAG(), "l", ""), "brn", int64(CondI(IFLAG(), 5, 3)))
	br(false)

	// The magic sequence "NOP ; BRN #offset" (i.e. $12 $21 offset)
	// is the new way to call the hyperviser.
	prevInst := B(pcreg - 3) // What came before the BRN?
	hyperOp := B(pcreg - 1)  // What is the immediate argument to BRN?
	if prevInst == 0x12 /*NOP*/ {
		if hyperOp == 0xFF {
			// Signature Op to set the B register to 'G' if Gomar is running.
			//
			// So your total sequence to test for GOMAR could be
			//     CLRA    ; optional, for 16-bit return value in D.
			//     CLRB
			//     NOP
			//     FCB $21 ; brn....
			//     FCB $FF ;    magic signature value.
			//
			// On GOMAR, that acts like
			//     CLRA    ; optional, for 16-bit return value in D.
			//     LDB #$47  ; 'G'
			//
			// This must be recognized by Gomar, even if other HyperOps are disabled!

			L("NewHyperOp GOMAR SIGNATURE: 'G' -> B")
			PutBReg('G')

		} else {
			L("NewHyperOp %d.", hyperOp)
			HyperOp(hyperOp)
		}
	}
}

func bhi() {
	Dis_inst(CondS(IFLAG(), "l", ""), "bhi", int64(CondI(IFLAG(), 5, 3)))
	br(0 == (ccreg & 0x05))
}

func bls() {
	Dis_inst(CondS(IFLAG(), "l", ""), "bls", int64(CondI(IFLAG(), 5, 3)))
	br(0 != ccreg&0x05)
}

func bcc() {
	Dis_inst(CondS(IFLAG(), "l", ""), "bcc", int64(CondI(IFLAG(), 5, 3)))
	br(0 == (ccreg & 0x01))
}

func bcs() {
	Dis_inst(CondS(IFLAG(), "l", ""), "bcs", int64(CondI(IFLAG(), 5, 3)))
	br(0 != ccreg&0x01)
}

func bne() {
	Dis_inst(CondS(IFLAG(), "l", ""), "bne", int64(CondI(IFLAG(), 5, 3)))
	br(0 == (ccreg & 0x04))
}

func beq() {
	Dis_inst(CondS(IFLAG(), "l", ""), "beq", int64(CondI(IFLAG(), 5, 3)))
	br(0 != ccreg&0x04)
}

func bvc() {
	Dis_inst(CondS(IFLAG(), "l", ""), "bvc", int64(CondI(IFLAG(), 5, 3)))
	br(0 == (ccreg & 0x02))
}

func bvs() {
	Dis_inst(CondS(IFLAG(), "l", ""), "bvs", int64(CondI(IFLAG(), 5, 3)))
	br(0 != ccreg&0x02)
}

func bpl() {
	Dis_inst(CondS(IFLAG(), "l", ""), "bpl", int64(CondI(IFLAG(), 5, 3)))
	br(0 == (ccreg & 0x08))
}

func bmi() {
	Dis_inst(CondS(IFLAG(), "l", ""), "bmi", int64(CondI(IFLAG(), 5, 3)))
	br(0 != ccreg&0x08)
}

func bge() {
	Dis_inst(CondS(IFLAG(), "l", ""), "bge", int64(CondI(IFLAG(), 5, 3)))
	br(!NXORV())
}

func blt() {
	Dis_inst(CondS(IFLAG(), "l", ""), "blt", int64(CondI(IFLAG(), 5, 3)))
	br(NXORV())
}

func bgt() {
	Dis_inst(CondS(IFLAG(), "l", ""), "bgt", int64(CondI(IFLAG(), 5, 3)))
	br(!(NXORV() || 0 != ccreg&0x04))
}

func ble() {
	Dis_inst(CondS(IFLAG(), "l", ""), "ble", int64(CondI(IFLAG(), 5, 3)))
	br(NXORV() || 0 != ccreg&0x04)
}

func leax() {
	Dis_inst("leax", "", 4)
	w := Word(postbyte())
	if w != 0 {
		CLZ()
	} else {
		SEZ()
	}
	xreg = w
}

func leay() {
	Dis_inst("leay", "", 4)
	w := Word(postbyte())
	if w != 0 {
		CLZ()
	} else {
		SEZ()
	}
	yreg = w
}

func leau() {
	Dis_inst("leau", "", 4)
	ureg = Word(postbyte())
}

func leas() {
	Dis_inst("leas", "", 4)
	sreg = Word(postbyte())
}

var reg_for_bit_count = []string{"pc", "u", "y", "x", "dp", "b", "a", "cc"}

func bit_count(b byte) int {
	var mask byte = 0x80
	count := 0
	for i := 0; i <= 7; i++ {
		if (b & mask) != 0 {
			count++
			Dis_ops(CondS(count > 1, ",", ""),
				reg_for_bit_count[i],
				1+int64(CondI(i < 4, 1, 0)))
		}
		mask >>= 1
	}
	return count
}

func pshs() {
	b := ImmByte()
	Dis_inst("pshs", "", 5)
	bit_count(b)
	if (b & 0x80) != 0 {
		PushWord(pcreg)
	}
	if (b & 0x40) != 0 {
		PushWord(ureg)
	}
	if (b & 0x20) != 0 {
		PushWord(yreg)
	}
	if (b & 0x10) != 0 {
		PushWord(xreg)
	}
	if (b & 0x08) != 0 {
		PushByte(dpreg)
	}
	if (b & 0x04) != 0 {
		PushByte(GetBReg())
	}
	if (b & 0x02) != 0 {
		PushByte(GetAReg())
	}
	if (b & 0x01) != 0 {
		PushByte(ccreg)
	}
}

func puls() {
	b := ImmByte()
	Dis_inst("puls", "", 5)
	Dis_len(2)
	bit_count(b)
	if (b & 0x01) != 0 {
		PullByte(&ccreg)
	}
	if (b & 0x02) != 0 {
		var t byte
		PullByte(&t)
		PutAReg(t)
	}
	if (b & 0x04) != 0 {
		var t byte
		PullByte(&t)
		PutBReg(t)
	}
	if (b & 0x08) != 0 {
		PullByte(&dpreg)
	}
	if (b & 0x10) != 0 {
		PullWord(&xreg)
	}
	if (b & 0x20) != 0 {
		PullWord(&yreg)
	}
	if (b & 0x40) != 0 {
		PullWord(&ureg)
	}
	if (b & 0x80) != 0 {
		PullWord(&pcreg)
	}

	if *FlagBasicText {
		ShowBasicText()
	}
}

func pshu() {
	b := ImmByte()
	Dis_inst("pshu", "", 5)
	bit_count(b)
	if (b & 0x80) != 0 {
		PushUWord(pcreg)
	}
	if (b & 0x40) != 0 {
		PushUWord(sreg)
	}
	if (b & 0x20) != 0 {
		PushUWord(yreg)
	}
	if (b & 0x10) != 0 {
		PushUWord(xreg)
	}
	if (b & 0x08) != 0 {
		PushUByte(dpreg)
	}
	if (b & 0x04) != 0 {
		PushUByte(GetBReg())
	}
	if (b & 0x02) != 0 {
		PushUByte(GetAReg())
	}
	if (b & 0x01) != 0 {
		PushUByte(ccreg)
	}
}

func pulu() {
	b := ImmByte()
	Dis_inst("pulu", "", 5)
	Dis_len(2)
	bit_count(b)
	if (b & 0x01) != 0 {
		PullUByte(&ccreg)
	}
	if (b & 0x02) != 0 {
		var t byte
		PullUByte(&t)
		PutAReg(t)
	}
	if (b & 0x04) != 0 {
		var t byte
		PullUByte(&t)
		PutBReg(t)
	}
	if (b & 0x08) != 0 {
		PullUByte(&dpreg)
	}
	if (b & 0x10) != 0 {
		PullUWord(&xreg)
	}
	if (b & 0x20) != 0 {
		PullUWord(&yreg)
	}
	if (b & 0x40) != 0 {
		PullUWord(&sreg)
	}
	if (b & 0x80) != 0 {
		PullUWord(&pcreg)
	}
}

func SETSTATUSD(a, b, res uint32) {
	if (res & 0x10000) != 0 {
		SEC()
	} else {
		CLC()
	}
	if (((res >> 1) ^ a ^ b ^ res) & 0x8000) != 0 {
		SEV()
	} else {
		CLV()
	}
	SETNZ16(Word(res))
}

func addd() {
	var aop, bop, res uint32
	Dis_inst("addd", "", 5)
	aop = uint32(dreg)
	ea := eaddr16()
	bop = uint32(ea.GetW())
	res = aop + bop
	SETSTATUSD(aop, bop, res)
	dreg = Word(res)
}

func subd() {
	var aop, bop, res uint32
	if iflag != 0 {
		Dis_inst("cmpd", "", 5)
	} else {
		Dis_inst("subd", "", 5)
	}
	if iflag == 2 {
		aop = uint32(ureg)
		Dis_inst("cmpu", "", 5)
	} else {
		aop = uint32(dreg)
	}
	ea := eaddr16()
	bop = uint32(ea.GetW())
	res = aop - bop
	SETSTATUSD(aop, bop, res)
	if iflag == 0 {
		dreg = Word(res)
	}
}

func cmpx() {
	var aop, bop, res uint32
	switch iflag {
	case 0:
		Dis_inst("cmpx", "", 5)
		aop = uint32(xreg)
	case 1:
		Dis_inst("cmpy", "", 5)
		aop = uint32(yreg)
	case 2:
		Dis_inst("cmps", "", 5)
		aop = uint32(sreg)
	}
	ea := eaddr16()
	bop = uint32(ea.GetW())
	res = aop - bop
	SETSTATUSD(aop, bop, res)
}

func ldd() {
	Dis_inst("ldd", "", 4)
	ea := eaddr16()
	w := ea.GetW()
	SETNZ16(w)
	dreg = w
}

func ldx() {
	if iflag != 0 {
		Dis_inst("ldy", "", 4)
	} else {
		Dis_inst("ldx", "", 4)
	}
	ea := eaddr16()
	w := ea.GetW()
	SETNZ16(w)
	if iflag == 0 {
		xreg = w
	} else {
		yreg = w
	}
}

func ldu() {
	if iflag != 0 {
		Dis_inst("lds", "", 4)
	} else {
		Dis_inst("ldu", "", 4)
	}
	ea := eaddr16()
	w := ea.GetW()
	SETNZ16(w)
	if iflag == 0 {
		ureg = w
	} else {
		sreg = w
	}
}

func std() {
	Dis_inst("std", "", 4)
	ea := eaddr16()
	w := dreg
	SETNZ16(w)
	ea.PutW(w)
}

func stx() {
	if iflag != 0 {
		Dis_inst("sty", "", 4)
	} else {
		Dis_inst("stx", "", 4)
	}
	ea := eaddr16()
	var w Word
	if iflag == 0 {
		w = xreg
	} else {
		w = yreg
	}
	SETNZ16(w)
	ea.PutW(w)
}

func stu() {
	if iflag == 0 {
		Dis_inst("stu", "", 4)
	} else {
		Dis_inst("sts", "", 4)
	}
	ea := eaddr16()
	var w Word
	if iflag == 0 {
		w = ureg
	} else {
		w = sreg
	}
	SETNZ16(w)
	ea.PutW(w)
}

func ccbits(b byte) string {
	var buf bytes.Buffer
	big := "EFHINZVC"    // bits that are set.
	little := "efhinzvc" // bits that are clear.
	i := 0
	for bm := byte(0x80); bm > 0; bm >>= 1 {
		if b&bm != 0 {
			buf.WriteByte(big[i])
		} else {
			buf.WriteByte(little[i])
		}
		i++
	}

	return buf.String()
}

func init() {
	instructionTable = [256]func(){
		neg, ill, ill, com, lsr, ill, ror, asr,
		asl, rol, dec, ill, inc, tst, jmp, clr,
		flag0, flag1, nop, sync_inst, ill, ill, lbra, lbsr,
		ill, daa, orcc, ill, andcc, sex, exg, tfr,
		bra, brn, bhi, bls, bcc, bcs, bne, beq,
		bvc, bvs, bpl, bmi, bge, blt, bgt, ble,
		leax, leay, leas, leau, pshs, puls, pshu, pulu,
		ill, rts, abx, rti, cwai, mul, ill, swi,
		neg, ill, ill, com, lsr, ill, ror, asr,
		asl, rol, dec, ill, inc, tst, ill, clr,
		neg, ill, ill, com, lsr, ill, ror, asr,
		asl, rol, dec, ill, inc, tst, ill, clr,
		neg, ill, ill, com, lsr, ill, ror, asr,
		asl, rol, dec, ill, inc, tst, jmp, clr,
		neg, ill, ill, com, lsr, ill, ror, asr,
		asl, rol, dec, ill, inc, tst, jmp, clr,
		sub, cmp, sbc, subd, and, bit, ld, st,
		eor, adc, or, add, cmpx, bsr, ldx, stx,
		sub, cmp, sbc, subd, and, bit, ld, st,
		eor, adc, or, add, cmpx, jsr, ldx, stx,
		sub, cmp, sbc, subd, and, bit, ld, st,
		eor, adc, or, add, cmpx, jsr, ldx, stx,
		sub, cmp, sbc, subd, and, bit, ld, st,
		eor, adc, or, add, cmpx, jsr, ldx, stx,
		sub, cmp, sbc, addd, and, bit, ld, st,
		eor, adc, or, add, ldd, std, ldu, stu,
		sub, cmp, sbc, addd, and, bit, ld, st,
		eor, adc, or, add, ldd, std, ldu, stu,
		sub, cmp, sbc, addd, and, bit, ld, st,
		eor, adc, or, add, ldd, std, ldu, stu,
		sub, cmp, sbc, addd, and, bit, ld, st,
		eor, adc, or, add, ldd, std, ldu, stu,
	}
}
