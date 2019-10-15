// +build js nacl plan9

package logutil

import (
	"io"
)

func checkIfTerminal(w io.Writer) bool {
	return false
}
