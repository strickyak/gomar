// https://sw.kovidgoyal.net/kitty/graphics-protocol/
package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func SendKitty(w io.Writer, control string, payload []byte) {
	var buf bytes.Buffer
	enc := base64.NewEncoder(base64.StdEncoding, &buf)
	enc.Write(payload)
	enc.Close()
	pay64 := buf.Bytes()

	w.Write([]byte{27, '_', 'G'})
	w.Write([]byte(control))
	w.Write([]byte{';'})
	w.Write(pay64)
	w.Write([]byte{27, '\\'})
}

func main() {
	fmt.Printf("\n------\n")
	var z []byte
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if (i & 3) == 0 {
				z = append(z, []byte{128, 128, 128}...)
			} else if (i & 3) == 0 {
				z = append(z, []byte{0, 128, 128}...)
			} else if i == j {
				z = append(z, []byte{0, 200, 0}...)
			} else {
				z = append(z, []byte{50, 50, 50}...)
			}
		}
	}
	SendKitty(os.Stdout, "q=2,f=24,s=10,v=10,a=T,t=d", z)
	fmt.Printf("\n------\n")
}
