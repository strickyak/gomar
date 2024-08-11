//go:build !copico

package emu

func GetCopico(a Word) (z byte) {
    panic("copico disabled (add --tag=copico)")
}

func PutCopico(a Word, b byte) {
    panic("copico disabled (add --tag=copico)")
}
