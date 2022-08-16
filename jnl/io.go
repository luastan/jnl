package jnl

import (
	"fmt"
	"io"
	"sync"
)

/*

Process stdin/stdout/stderr management

*/

func OutErrRoutine(wg *sync.WaitGroup, reader io.ReadCloser, w io.Writer) {
	defer wg.Done()

	for {
		tmp := make([]byte, 1024)
		_, err := reader.Read(tmp)
		if err != nil {
			break
		}
		_, _ = fmt.Fprint(w, string(tmp))
	}
}
