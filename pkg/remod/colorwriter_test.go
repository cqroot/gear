// Copyright (C) 2026 Keith Chu <cqroot@outlook.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package remod

import (
	"bytes"
	"testing"

	"github.com/fatih/color"
)

// withColor forces color output on for the duration of the test so the
// escape sequences can be asserted regardless of the test environment.
func withColor(t *testing.T, enabled bool) {
	t.Helper()
	old := color.NoColor
	color.NoColor = !enabled
	t.Cleanup(func() { color.NoColor = old })
}

func TestColorWriter_colorizesCompleteLines(t *testing.T) {
	withColor(t, true)

	var dst bytes.Buffer
	w := &colorWriter{dst: &dst, style: goModColor}

	if _, err := w.Write([]byte("first\nsecond")); err != nil {
		t.Fatalf("Write: %v", err)
	}
	// The incomplete "second" line must stay buffered until flush.
	if got, want := dst.String(), "\x1b[90mfirst\n\x1b[0m"; got != want {
		t.Fatalf("after Write: got %q, want %q", got, want)
	}

	if err := w.flush(); err != nil {
		t.Fatalf("flush: %v", err)
	}
	if got, want := dst.String(), "\x1b[90mfirst\n\x1b[0m\x1b[90msecond\x1b[0m"; got != want {
		t.Fatalf("after flush: got %q, want %q", got, want)
	}
}

func TestColorWriter_buffersUntilNewline(t *testing.T) {
	withColor(t, true)

	var dst bytes.Buffer
	w := &colorWriter{dst: &dst, style: goModColor}

	for _, chunk := range []string{"go: ", "downloading ", "modules"} {
		if _, err := w.Write([]byte(chunk)); err != nil {
			t.Fatalf("Write(%q): %v", chunk, err)
		}
		if dst.Len() != 0 {
			t.Fatalf("expected no output before newline, got %q", dst.String())
		}
	}

	if _, err := w.Write([]byte("\n")); err != nil {
		t.Fatalf("Write newline: %v", err)
	}
	if got, want := dst.String(), "\x1b[90mgo: downloading modules\n\x1b[0m"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestColorWriter_indentsNonBlankLines(t *testing.T) {
	withColor(t, true)

	var dst bytes.Buffer
	w := &colorWriter{dst: &dst, indent: "    ", style: goModColor}

	if _, err := w.Write([]byte("one\n\ntwo")); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if err := w.flush(); err != nil {
		t.Fatalf("flush: %v", err)
	}

	// "one" and "two" are indented; the blank line stays empty.
	want := "\x1b[90m    one\n\x1b[0m\x1b[90m\n\x1b[0m\x1b[90m    two\x1b[0m"
	if got := dst.String(); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestColorWriter_passthroughWhenColorDisabled(t *testing.T) {
	withColor(t, false)

	var dst bytes.Buffer
	w := &colorWriter{dst: &dst, style: goModColor}

	if _, err := w.Write([]byte("plain\n")); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if err := w.flush(); err != nil {
		t.Fatalf("flush: %v", err)
	}
	if got, want := dst.String(), "plain\n"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
