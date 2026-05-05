//go:build bonobo

package emu

import (
	"flag"
	"fmt"
	. "github.com/strickyak/gomar/gu"
	"log"
	"net"
)

var FlagBonoboServer = flag.String("bonobo_server", "",
	"Lemma server for Bonobo Connections")

var Bonobo = struct {
	conn   *net.TCPConn
	rx, tx []byte
	c      byte
}{}

func InitBonobo() {
	if Bonobo.conn != nil {
		return
	}

	Assert(*FlagBonoboServer != "")

	p := &Bonobo
	raddy, err := net.ResolveTCPAddr("tcp", *FlagBonoboServer)
	if err != nil {
		log.Panicf("WIZ: cannot ResolveTCPAddr: %v", err)
	}
	tconn, err := net.DialTCP("tcp", nil /*local addy*/, raddy)
	if err != nil {
		log.Panicf("WIZ: cannot ListenUDP: %v", err)
	}
	fmt.Printf("bonobo/init: DialTcp: success: %q\n", *FlagBonoboServer)
	log.Printf("bonobo/init: DialTcp: success: %q", *FlagBonoboServer)
	p.conn = tconn
}

func GetBonobo(a Word) (z byte) {
	fmt.Printf("GetBonobo(%04x) ...\n", a)
	log.Printf("GetBonobo(%04x) ...", a)
	BonoboInit()
	p := &Bonobo
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

	fmt.Printf("GetBonobo(%04x -> $%02x) [ c=%d rx=%d tx=%d ]\n", a, z, p.c, len(p.rx), len(p.tx))
	log.Printf("GetBonobo(%04x -> $%02x) [ c=%d rx=%d tx=%d ]", a, z, p.c, len(p.rx), len(p.tx))
	return
}

func PutBonobo(a Word, b byte) {
	fmt.Printf("PutBonobo(%04x <- $%02x) ...\n", a, b)
	log.Printf("PutBonobo(%04x <- $%02x) ...", a, b)
	BonoboInit()
	p := &Bonobo
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
			log.Panicf("PutBonobo: bad cmd: %d.", b)
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
	fmt.Printf("PutBonobo(%04x <- $%02x) [ c=%d rx=%d tx=%d ]\n", a, b, p.c, len(p.rx), len(p.tx))
	log.Printf("PutBonobo(%04x <- $%02x) [ c=%d rx=%d tx=%d ]", a, b, p.c, len(p.rx), len(p.tx))
}
