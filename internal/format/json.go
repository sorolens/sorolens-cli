package format

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// PrintJSON marshals v as indented JSON to w (defaults to os.Stdout).
func PrintJSON(w io.Writer, v interface{}) error {
	if w == nil {
		w = os.Stdout
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// Errorf writes a formatted error line to stderr and returns a non-nil error.
func Errorf(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintln(os.Stderr, "error:", msg)
	return fmt.Errorf("%s", msg)
}
