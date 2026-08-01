package format

import (
	"os"

	"github.com/mattn/go-isatty"
)

var colorEnabled = func() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	return isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
}()

// SetColor overrides the auto-detected color setting.
func SetColor(enabled bool) {
	colorEnabled = enabled
}

// ColorEnabled returns whether color is currently enabled.
func ColorEnabled() bool {
	return colorEnabled
}
