//go:build copico

package emu

import (
	"flag"
	"fmt"
	. "github.com/strickyak/gomar/gu"
	"log"
	"net"
)

var FlagCopicoServer = flag.String("copico_server", "",
	"Lemma server for Copico Connections")

var Copico = struct {
	conn   *net.TCPConn
	rx, tx []byte
	c      byte
}{}

func CopicoInit() {
    if Copico.conn != nil {
        return
    }

	Assert(*FlagCopicoServer != "")

	p := &Copico
	raddy, err := net.ResolveTCPAddr("tcp", *FlagCopicoServer)
	if err != nil {
		log.Panicf("WIZ: cannot ResolveTCPAddr: %v", err)
	}
	tconn, err := net.DialTCP("tcp", nil /*local addy*/, raddy)
	if err != nil {
		log.Panicf("WIZ: cannot ListenUDP: %v", err)
	}
	fmt.Printf("copico/init: DialTcp: success: %q\n", *FlagCopicoServer)
	log.Printf("copico/init: DialTcp: success: %q", *FlagCopicoServer)
	p.conn = tconn
}

func GetCopico(a Word) (z byte) {
	fmt.Printf("GetCopico(%04x) ...\n", a)
	log.Printf("GetCopico(%04x) ...", a)
    CopicoInit()
	p := &Copico
	Assert(p.conn != nil)
	switch a {
	case 0xFF78:
		panic(a)
	case 0xFF79:
		z = 2
	case 0xFF7A:
		n := byte(len(p.tx))
		i := n - p.c
	    fmt.Printf("    n=%d c=%d i=%d z=...\n", n, p.c, i)
	    log.Printf("    n=%d c=%d i=%d z=...\n", n, p.c, i)
		z = p.tx[i]
	    fmt.Printf("    n=%d c=%d i=%d z=%d\n", n, p.c, i, z)
	    log.Printf("    n=%d c=%d i=%d z=%d\n", n, p.c, i, z)
		p.c--

		if p.c == 0 {
			p.tx = nil
		}
	case 0xFF7B:
		panic(a)
	default:
		panic(a)
	}

	fmt.Printf("GetCopico(%04x -> $%02x) [ c=%d rx=%d tx=%d ]\n", a, z, p.c, len(p.rx), len(p.tx))
	log.Printf("GetCopico(%04x -> $%02x) [ c=%d rx=%d tx=%d ]", a, z, p.c, len(p.rx), len(p.tx))
	return
}

func PutCopico(a Word, b byte) {
	fmt.Printf("PutCopico(%04x <- $%02x) ...\n", a, b)
	log.Printf("PutCopico(%04x <- $%02x) ...", a, b)
    CopicoInit()
	p := &Copico
	Assert(p.conn != nil)
	switch a {
	case 0xFF78:
		if 1 <= b && b <= 100 {
			n := b
			p.rx, p.tx = make([]byte, n), nil
			p.c = n
		} else if 101 <= b && b <= 200 {
			n := b - 100
			p.rx, p.tx = nil, make([]byte, n)
			p.c = n

			cc := Value(p.conn.Read(p.tx))
			AssertEQ(cc, len(p.tx))

		} else {
			log.Panicf("PutCopico: bad cmd: %d.", b)
		}
	case 0xFF79:
		panic(a)
	case 0xFF7A:
		panic(a)
	case 0xFF7B:
		AssertLE(1, p.c)
		n := byte(len(p.rx))
		i := n - p.c
	    fmt.Printf("    n=%d c=%d i=%d b=%d\n", n, p.c, i, b)
	    log.Printf("    n=%d c=%d i=%d b=%d\n", n, p.c, i, b)
		p.rx[i] = b
		p.c--

		if p.c == 0 {
			cc := Value(p.conn.Write(p.rx))
			AssertEQ(cc, len(p.rx))
			p.rx = nil
		}

	default:
		panic(a)
	}
	fmt.Printf("PutCopico(%04x <- $%02x) [ c=%d rx=%d tx=%d ]\n", a, b, p.c, len(p.rx), len(p.tx))
	log.Printf("PutCopico(%04x <- $%02x) [ c=%d rx=%d tx=%d ]", a, b, p.c, len(p.rx), len(p.tx))
}
