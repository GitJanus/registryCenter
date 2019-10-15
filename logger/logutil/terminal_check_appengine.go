// +build appengine

package logutil

import (
	"io"
)

func checkIfTerminal(w io.Writer) bool {
	return true
}
