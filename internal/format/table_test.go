package format_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/sorolens/sorolens-cli/internal/format"
)

func TestRenderTableGolden(t *testing.T) {
	format.SetColor(false)
	headers := []string{"Key Hash", "Durability", "Ledgers Left"}
	rows := [][]string{
		{"deadbeef", "persistent", "12000"},
		{"cafebabe", "temporary", "800"},
	}

	var buf bytes.Buffer
	format.RenderTable(&buf, headers, rows)

	goldenPath := "../../testdata/golden/table_ttl.txt"
	if os.Getenv("UPDATE_GOLDEN") == "1" {
		os.WriteFile(goldenPath, buf.Bytes(), 0644)
		t.Logf("updated golden file %s", goldenPath)
		return
	}

	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("golden file missing; run with UPDATE_GOLDEN=1 to create: %v", err)
	}
	if strings.TrimRight(buf.String(), "\n") != strings.TrimRight(string(want), "\n") {
		t.Errorf("table output mismatch\ngot:\n%s\nwant:\n%s", buf.String(), string(want))
	}
}

func TestKeyValueTableGolden(t *testing.T) {
	format.SetColor(false)
	rows := [][2]string{
		{"Alias", "my-contract"},
		{"Network", "mainnet"},
		{"Event Count", "42"},
	}

	var buf bytes.Buffer
	format.KeyValueTable(&buf, rows)

	goldenPath := "../../testdata/golden/table_kv.txt"
	if os.Getenv("UPDATE_GOLDEN") == "1" {
		os.WriteFile(goldenPath, buf.Bytes(), 0644)
		t.Logf("updated golden file %s", goldenPath)
		return
	}

	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("golden file missing; run with UPDATE_GOLDEN=1 to create: %v", err)
	}
	if strings.TrimRight(buf.String(), "\n") != strings.TrimRight(string(want), "\n") {
		t.Errorf("table output mismatch\ngot:\n%s\nwant:\n%s", buf.String(), string(want))
	}
}

func TestTTLLevelFor(t *testing.T) {
	cases := []struct {
		ledgers int64
		want    format.TTLLevel
	}{
		{0, format.TTLDanger},
		{999, format.TTLDanger},
		{1000, format.TTLWarning},
		{9999, format.TTLWarning},
		{10000, format.TTLSafe},
		{50000, format.TTLSafe},
	}
	for _, tc := range cases {
		got := format.TTLLevelFor(tc.ledgers)
		if got != tc.want {
			t.Errorf("TTLLevelFor(%d) = %d, want %d", tc.ledgers, got, tc.want)
		}
	}
}

func TestTruncateHash(t *testing.T) {
	h := "abc123def456789"
	got := format.TruncateHash(h)
	if len(got) > 15 {
		t.Errorf("truncated hash too long: %q", got)
	}
	short := "abc"
	if format.TruncateHash(short) != short {
		t.Error("short hash should not be truncated")
	}
}

func TestColorTTL_NoColor(t *testing.T) {
	format.SetColor(false)
	s := format.ColorTTL("some text", format.TTLSafe)
	if s != "some text" {
		t.Errorf("expected plain text, got %q", s)
	}
}
