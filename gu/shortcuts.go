package gu // Go Utilities for gomar.

import (
	"fmt"
	"log"
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
