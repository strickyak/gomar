package kitty

// https://sw.kovidgoyal.net/kitty/graphics-protocol/

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
)

func Draw(w io.Writer, width uint, height uint, payload []byte) {
	control := fmt.Sprintf("q=2,f=24,s=%d,v=%d,a=T,t=d", width, height)

	var buf bytes.Buffer
	enc := base64.NewEncoder(base64.StdEncoding, &buf)
	enc.Write(payload)
	enc.Close()
	pay := buf.Bytes()

	const N = 1000
	continued := false
	for len(pay) > N {
		w.Write([]byte{27, '_', 'G'})
		if continued {
			w.Write([]byte("m=1;"))
		} else {
			w.Write([]byte(control))
			w.Write([]byte(",m=1;"))
		}
		w.Write(pay[:N])
		w.Write([]byte{27, '\\'})
		pay = pay[N:]
		continued = true
	}

	w.Write([]byte{27, '_', 'G'})
	if continued {
		w.Write([]byte("m=0;"))
	} else {
		w.Write([]byte(control))
	}
	w.Write(pay)
	w.Write([]byte{27, '\\'})
}
