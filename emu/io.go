package emu

type ByteGetter func(a Word) byte
type BytePutter func(a Word, b byte)

const Const64K = 1 << 16

var IOGetters [Const64K]ByteGetter
var IOPutters [Const64K]BytePutter
