//go:build noos

package emu

import ()

const Level = -1

const P_Path = ""

func VerboseValidateModuleSyscall() string { return "" }
func DoDumpSysMap() {
}

func MemoryModuleOf(addr Word) (string, Word) {
	return "", addr
}

func ScanModDir() {
}

func DoDumpProcDesc(a Word, queue string, followQ bool) {
}

func MemoryModules() {
}

func DoDumpAllMemoryPhys() {}
func DoDumpPageZero()      {}
func DoDumpProcesses()     {}
func DoDumpAllPathDescs()  {}
func DumpGimeStatus()      {}

func MapAddr(logical Word, quiet bool) int {
	return int(logical)
}

func DoDumpPaths() {
}

func DecodeOs9Level2Opcode(b byte) (s string, p string, returns bool) {
	panic(0)
}
func Swi2ForOs9()                         {}
func Os9HypervisorCall(syscall byte) bool { return false }

func PreRTI() (stack int, describe string) { return }
func PostRTI(stack int, describe string)   {}

func ScanRamForOs9Modules() []*ModuleFound {
	return nil
}
