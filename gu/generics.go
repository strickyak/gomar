package gu // Go Utilities for gomar.

import (
	"errors"
	"fmt"
	"log"
	"runtime/debug"
)

func Cond[T any](a bool, b T, c T) T {
	if a {
		return b
	}
	return c
}

func Min[T Number](b T, c T) T {
	if b < c {
		return b
	}
	return c
}

func Max[T Number](b T, c T) T {
	if b > c {
		return b
	}
	return c
}

func Value[T any](value T, err error) T {
	Check(err)
	return value
}

func Errorf(f string, args ...any) error {
	return errors.New(fmt.Sprintf(f, args...))
}

type Number interface {
	~int8 | ~int16 | ~int32 | ~uint8 | ~uint16 | ~uint32 | ~int | ~uint | ~int64 | ~uint64 | ~uintptr
}

func Assert(b bool, args ...any) {
	if !b {
		s := "Assert Fails"
		for _, x := range args {
			s += fmt.Sprintf(" ; %v", x)
		}
		s += "\n[[[[[[\n" + string(debug.Stack()) + "\n]]]]]]\n"
		log.Panic(s)
	}
}

func AssertEQ[N Number](a, b N, args ...any) {
	if a != b {
		s := fmt.Sprintf("AssertEQ Fails: (%v .EQ. %v)", a, b)
		for _, x := range args {
			s += fmt.Sprintf(" ; %v", x)
		}
		s += "\n[[[[[[\n" + string(debug.Stack()) + "\n]]]]]]\n"
		log.Panic(s)
	}
}

func AssertNE[N Number](a, b N, args ...any) {
	if a == b {
		s := fmt.Sprintf("AssertNE Fails: (%v .NE. %v)", a, b)
		for _, x := range args {
			s += fmt.Sprintf(" ; %v", x)
		}
		s += "\n[[[[[[\n" + string(debug.Stack()) + "\n]]]]]]\n"
		log.Panic(s)
	}
}

func AssertLT[N Number](a, b N, args ...any) {
	if a >= b {
		s := fmt.Sprintf("AssertLT Fails: (%v .LT. %v)", a, b)
		for _, x := range args {
			s += fmt.Sprintf(" ; %v", x)
		}
		s += "\n[[[[[[\n" + string(debug.Stack()) + "\n]]]]]]\n"
		log.Panic(s)
	}
}

func AssertLE[N Number](a, b N, args ...any) {
	if a > b {
		s := fmt.Sprintf("AssertLE Fails: (%v .LE. %v)", a, b)
		for _, x := range args {
			s += fmt.Sprintf(" ; %v", x)
		}
		s += "\n[[[[[[\n" + string(debug.Stack()) + "\n]]]]]]\n"
		log.Panic(s)
	}
}

func AssertGT[N Number](a, b N, args ...any) {
	if a <= b {
		s := fmt.Sprintf("AssertGT Fails: (%v .GT. %v)", a, b)
		for _, x := range args {
			s += fmt.Sprintf(" ; %v", x)
		}
		s += "\n[[[[[[\n" + string(debug.Stack()) + "\n]]]]]]\n"
		log.Panic(s)
	}
}

func AssertGE[N Number](a, b N, args ...any) {
	if a < b {
		s := fmt.Sprintf("AssertGE Fails: (%v .GE. %v)", a, b)
		for _, x := range args {
			s += fmt.Sprintf(" ; %v", x)
		}
		s += "\n[[[[[[\n" + string(debug.Stack()) + "\n]]]]]]\n"
		log.Panic(s)
	}
}

func Check(err error, args ...any) {
	if err != nil {
		s := fmt.Sprintf("Check Fails: %v", err)
		for _, x := range args {
			s += fmt.Sprintf(" ; %v", x)
		}
		s += "\n[[[[[[\n" + string(debug.Stack()) + "\n]]]]]]\n"
		log.Panic(s)
	}
}
