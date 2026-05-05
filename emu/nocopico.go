//go:build !bonobo

package emu

func InitBonobo() {}

func GetBonobo(a Word) (z byte) {
	panic("bonobo disabled (add --tag=bonobo)")
}

func PutBonobo(a Word, b byte) {
	panic("bonobo disabled (add --tag=bonobo)")
}
