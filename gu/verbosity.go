package gu // Go Utilities for gomar.

import (
	"bytes"
	"flag"
	"log"
	"runtime"
	"path"
	"fmt"
)

var (
	V [128]bool

	FlagInitialVerbosity = flag.String("v", "", "Initial verbosity chars") // Initial Verbosity
)

func T(args ... any) {
	var bb bytes.Buffer
	fmt.Fprintf(&bb, "#")
	for i := 5; i > 0; i-- {
		_, filename, lineno, ok := runtime.Caller(i)
		if ok {
			fmt.Fprintf(&bb, " %s:%d", path.Base(filename), lineno)
		}
	}
	for _, arg := range args {
		fmt.Fprintf(&bb, " :: %v", arg)
	}
	Log("%s", bb.String())
}

func InitVerbosity() {
	SetVerbosityBits(*FlagInitialVerbosity)
}

func SetVerbosityBits(s string) {
	for _, r := range s {
		if int(r) >= len(V) {
			log.Panicf("Verbosity rune %d too large for Verbosity Array", r)
		}
		if r == 'a' {
			for i := 'a'; i <= 'z'; i++ {
				V[i] = true
			}
			continue
		}
		V[r] = true
	}
}

func Loga(f string, a ...any) {
	if V['a'] {
		Log("#a "+f, a...)
	}
}
func Logb(f string, a ...any) {
	if V['b'] {
		Log("#b "+f, a...)
	}
}
func Logc(f string, a ...any) {
	if V['c'] {
		Log("#c "+f, a...)
	}
}
func Logd(f string, a ...any) {
	if V['d'] {
		Log("#d "+f, a...)
	}
}
func Loge(f string, a ...any) {
	if V['e'] {
		Log("#e "+f, a...)
	}
}
func Logf(f string, a ...any) {
	if V['f'] {
		Log("#f "+f, a...)
	}
}
func Logg(f string, a ...any) {
	if V['g'] {
		Log("#g "+f, a...)
	}
}
func Logh(f string, a ...any) {
	if V['h'] {
		Log("#h "+f, a...)
	}
}
func Logi(f string, a ...any) {
	if V['i'] {
		Log("#i "+f, a...)
	}
}
func Logj(f string, a ...any) {
	if V['j'] {
		Log("#j "+f, a...)
	}
}
func Logk(f string, a ...any) {
	if V['k'] {
		Log("#k "+f, a...)
	}
}
func Logl(f string, a ...any) {
	if V['l'] {
		Log("#l "+f, a...)
	}
}
func Logm(f string, a ...any) {
	if V['m'] {
		Log("#m "+f, a...)
	}
}
func Logn(f string, a ...any) {
	if V['n'] {
		Log("#n "+f, a...)
	}
}
func Logo(f string, a ...any) {
	if V['o'] {
		Log("#o "+f, a...)
	}
}
func Logp(f string, a ...any) {
	if V['p'] {
		Log("#p "+f, a...)
	}
}
func Logq(f string, a ...any) {
	if V['q'] {
		Log("#q "+f, a...)
	}
}
func Logr(f string, a ...any) {
	if V['r'] {
		Log("#r "+f, a...)
	}
}
func Logs(f string, a ...any) {
	if V['s'] {
		Log("#s "+f, a...)
	}
}
func Logt(f string, a ...any) {
	if V['t'] {
		Log("#t "+f, a...)
	}
}
func Logu(f string, a ...any) {
	if V['u'] {
		Log("#u "+f, a...)
	}
}
func Logv(f string, a ...any) {
	if V['v'] {
		Log("#v "+f, a...)
	}
}
func Logw(f string, a ...any) {
	if V['w'] {
		Log("#w "+f, a...)
	}
}
func Logx(f string, a ...any) {
	if V['x'] {
		Log("#x "+f, a...)
	}
}
func Logy(f string, a ...any) {
	if V['y'] {
		Log("#y "+f, a...)
	}
}
func Logz(f string, a ...any) {
	if V['z'] {
		Log("#z "+f, a...)
	}
}

func LogA(f string, a ...any) {
	if V['A'] {
		Log("#A "+f, a...)
	}
}
func LogB(f string, a ...any) {
	if V['B'] {
		Log("#B "+f, a...)
	}
}
func LogC(f string, a ...any) {
	if V['C'] {
		Log("#C "+f, a...)
	}
}
func LogD(f string, a ...any) {
	if V['D'] {
		Log("#D "+f, a...)
	}
}
func LogE(f string, a ...any) {
	if V['E'] {
		Log("#E "+f, a...)
	}
}
func LogF(f string, a ...any) {
	if V['F'] {
		Log("#F "+f, a...)
	}
}
func LogG(f string, a ...any) {
	if V['G'] {
		Log("#G "+f, a...)
	}
}
func LogH(f string, a ...any) {
	if V['H'] {
		Log("#H "+f, a...)
	}
}
func LogI(f string, a ...any) {
	if V['I'] {
		Log("#I "+f, a...)
	}
}
func LogJ(f string, a ...any) {
	if V['J'] {
		Log("#J "+f, a...)
	}
}
func LogK(f string, a ...any) {
	if V['K'] {
		Log("#K "+f, a...)
	}
}
func LogL(f string, a ...any) {
	if V['L'] {
		Log("#L "+f, a...)
	}
}
func LogM(f string, a ...any) {
	if V['M'] {
		Log("#M "+f, a...)
	}
}
func LogN(f string, a ...any) {
	if V['N'] {
		Log("#N "+f, a...)
	}
}
func LogO(f string, a ...any) {
	if V['O'] {
		Log("#O "+f, a...)
	}
}
func LogP(f string, a ...any) {
	if V['P'] {
		Log("#P "+f, a...)
	}
}
func LogQ(f string, a ...any) {
	if V['Q'] {
		Log("#Q "+f, a...)
	}
}
func LogR(f string, a ...any) {
	if V['R'] {
		Log("#R "+f, a...)
	}
}
func LogS(f string, a ...any) {
	if V['S'] {
		Log("#S "+f, a...)
	}
}
func LogT(f string, a ...any) {
	if V['T'] {
		Log("#T "+f, a...)
	}
}
func LogU(f string, a ...any) {
	if V['U'] {
		Log("#U "+f, a...)
	}
}
func LogV(f string, a ...any) {
	if V['V'] {
		Log("#V "+f, a...)
	}
}
func LogW(f string, a ...any) {
	if V['W'] {
		Log("#W "+f, a...)
	}
}
func LogX(f string, a ...any) {
	if V['X'] {
		Log("#X "+f, a...)
	}
}
func LogY(f string, a ...any) {
	if V['Y'] {
		Log("#Y "+f, a...)
	}
}
func LogZ(f string, a ...any) {
	if V['Z'] {
		Log("#Z "+f, a...)
	}
}
