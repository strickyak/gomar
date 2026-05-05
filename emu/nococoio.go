//go:build !cocoio

package emu

func InitCocoIO()              {}
func PutCocoIO(a Word, b byte) {}
func GetCocoIO(a Word) byte    { return 126 }
