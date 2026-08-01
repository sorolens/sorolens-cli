package format

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/charmbracelet/lipgloss"
)

// TTLLevel describes how urgently a storage entry is expiring.
type TTLLevel int

const (
	TTLSafe    TTLLevel = iota // > 10000 ledgers
	TTLWarning                 // 1000-9999 ledgers
	TTLDanger                  // < 1000 ledgers or expired
)

var (
	styleHeader = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	styleGreen  = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	styleYellow = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	styleRed    = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

// RenderTable writes a tab-aligned table to w (defaults to os.Stdout).
func RenderTable(w io.Writer, headers []string, rows [][]string) {
	if w == nil {
		w = os.Stdout
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	header := strings.Join(headers, "\t")
	if ColorEnabled() {
		fmt.Fprintln(tw, styleHeader.Render(header))
	} else {
		fmt.Fprintln(tw, header)
	}

	sep := make([]string, len(headers))
	for i, h := range headers {
		sep[i] = strings.Repeat("-", len(h))
	}
	fmt.Fprintln(tw, strings.Join(sep, "\t"))

	for _, row := range rows {
		fmt.Fprintln(tw, strings.Join(row, "\t"))
	}
	tw.Flush()
}

// ColorTTL applies a color style to a string based on TTL level.
// Returns plain text when color is disabled.
func ColorTTL(s string, level TTLLevel) string {
	if !ColorEnabled() {
		return s
	}
	switch level {
	case TTLSafe:
		return styleGreen.Render(s)
	case TTLWarning:
		return styleYellow.Render(s)
	default:
		return styleRed.Render(s)
	}
}

// TTLLevelFor classifies ledgers-left into a TTL level.
func TTLLevelFor(ledgersLeft int64) TTLLevel {
	switch {
	case ledgersLeft < 1000:
		return TTLDanger
	case ledgersLeft < 10000:
		return TTLWarning
	default:
		return TTLSafe
	}
}

// TruncateHash shortens a transaction hash to 12 chars + "...".
func TruncateHash(h string) string {
	if len(h) <= 15 {
		return h
	}
	return h[:12] + "..."
}

// KeyValueTable renders a two-column key/value table to w.
func KeyValueTable(w io.Writer, rows [][2]string) {
	tableRows := make([][]string, len(rows))
	for i, r := range rows {
		tableRows[i] = []string{r[0], r[1]}
	}
	RenderTable(w, []string{"Field", "Value"}, tableRows)
}

// SectionHeader renders a bold section header line.
func SectionHeader(w io.Writer, title string) {
	if w == nil {
		w = os.Stdout
	}
	if ColorEnabled() {
		fmt.Fprintln(w, styleHeader.Render(title))
	} else {
		fmt.Fprintln(w, strings.ToUpper(title))
	}
}
