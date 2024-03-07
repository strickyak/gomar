package gu // Go Utilities for gomar.

import (
	"fmt"
	"log"
	"strconv"
)

func Fmt(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

func QFmt(format string, args ...any) string {
	return fmt.Sprintf("%q", fmt.Sprintf(format, args...))
}

func Log(format string, args ...any) {
	log.Printf(format, args...)
}

func Str(x any) string {
	return fmt.Sprintf("%v", x)
}

func Repr(x any) string {
	return fmt.Sprintf("%#v", x)
}

func QStr(x any) string {
	return fmt.Sprintf("%q", fmt.Sprintf("%v", x))
}

func QRepr(x any) string {
	return fmt.Sprintf("%q", fmt.Sprintf("%#v", x))
}

func DeHex(s string) uint {
	x, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		log.Panicf("Cannot convert hex $q: %v", s, err)
	}
	return uint(x)
}
