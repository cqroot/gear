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

// Package remod provides core utilities used by the gear CLI commands.
package remod

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// ErrNoModule is returned by ParseGoMod when the go.mod file does not
// contain a `module` directive.
var ErrNoModule = errors.New("go.mod does not contain a module directive")

// bannerColor styles gear's own banner lines.
var bannerColor = color.New(color.FgCyan)

// goModColor styles the raw output forwarded from the go tool. Bright black
// renders as gray on both light and dark terminals.
var goModColor = color.New(color.FgHiBlack)

// stepIndent and goModIndent indent nested output: step banners sit two
// spaces in and the go tool's own output six, beneath the top-level banner.
const (
	stepIndent  = "  "
	goModIndent = "      "
)

// colorWriter writes each complete line to dst wrapped in style and prefixed
// with indent. Call flush to emit a trailing partial line.
type colorWriter struct {
	dst    io.Writer
	indent string
	style  *color.Color
	buf    bytes.Buffer
}

// writeLine writes one logical line to dst via style, indenting non-blank
// lines.
func (w *colorWriter) writeLine(line string) error {
	if w.indent != "" && line != "\n" && line != "" {
		line = w.indent + line
	}
	_, err := w.style.Fprint(w.dst, line)
	return err
}

// Write buffers p and emits every complete line it now contains.
func (w *colorWriter) Write(p []byte) (int, error) {
	w.buf.Write(p)
	for {
		line, err := w.buf.ReadString('\n')
		if err != nil {
			// No newline yet: keep the partial line for the next call.
			w.buf.WriteString(line)
			break
		}
		if writeErr := w.writeLine(line); writeErr != nil {
			return len(p), writeErr
		}
	}
	return len(p), nil
}

// flush emits any buffered partial line that never got a trailing newline.
func (w *colorWriter) flush() error {
	if w.buf.Len() == 0 {
		return nil
	}
	line := w.buf.String()
	w.buf.Reset()
	return w.writeLine(line)
}

// ParseGoMod reads a go.mod file located at the given root and returns the
// module path declared inside it. If no `module` directive is found,
// ErrNoModule is returned.
func ParseGoMod(root string) (string, error) {
	f, err := os.Open(filepath.Join(root, "go.mod"))
	if err != nil {
		return "", fmt.Errorf("open go.mod: %w", err)
	}
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 2 && fields[0] == "module" {
			return fields[1], nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("read go.mod: %w", err)
	}
	return "", ErrNoModule
}

// RemoveFiles deletes the listed relative paths from root. Missing files
// are silently ignored so that the operation is idempotent.
func RemoveFiles(root string, names ...string) error {
	for _, name := range names {
		if err := os.Remove(filepath.Join(root, name)); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("remove %s: %w", name, err)
		}
	}
	return nil
}

// RunGoMod executes `go mod <args...>` inside root. The go command's
// stdout and stderr are routed to stdout and stderr in real time, with
// each line indented and rendered in gray.
//
// Either writer may be nil; nil writers are treated as io.Discard.
func RunGoMod(stdout, stderr io.Writer, root string, args ...string) error {
	if stdout == nil {
		stdout = io.Discard
	}
	if stderr == nil {
		stderr = io.Discard
	}

	out := &colorWriter{dst: stdout, indent: goModIndent, style: goModColor}
	errOut := &colorWriter{dst: stderr, indent: goModIndent, style: goModColor}

	cmd := exec.Command("go", append([]string{"mod"}, args...)...)
	cmd.Dir = root
	cmd.Stdout = out
	cmd.Stderr = errOut
	runErr := cmd.Run()
	// Flush partial lines even when the go command fails so no output is
	// lost before the caller reports the error.
	if err := out.flush(); err != nil && runErr == nil {
		runErr = err
	}
	if err := errOut.flush(); err != nil && runErr == nil {
		runErr = err
	}
	if runErr != nil {
		return runErr
	}
	return nil
}

// Remod rebuilds the Go module at root without going through ParseGoMod.
// module must be non-empty; the function returns an error before touching
// any files otherwise.
func Remod(root string, stdout, stderr io.Writer, module string) error {
	if root == "" {
		return errors.New("remod: empty root directory")
	}
	// Validate the module input before touching any files so a bad call
	// cannot silently wipe an existing go.mod / go.sum.
	if module == "" {
		return errors.New("remod: module must be non-empty")
	}

	if err := RemoveFiles(root, "go.mod", "go.sum"); err != nil {
		return err
	}

	_, _ = bannerColor.Fprint(stdout, stepIndent+"• ")
	_, _ = fmt.Fprintln(stdout, "gear remod: initialising module")

	if err := RunGoMod(stdout, stderr, root, "init", module); err != nil {
		return fmt.Errorf("go mod init failed: %w", err)
	}

	_, _ = bannerColor.Fprint(stdout, stepIndent+"• ")
	_, _ = fmt.Fprintln(stdout, "gear remod: resolving dependencies")

	if err := RunGoMod(stdout, stderr, root, "tidy"); err != nil {
		return fmt.Errorf("go mod tidy failed: %w", err)
	}
	return nil
}

// Run rebuilds the Go module at root.
//
// If root is empty, the current working directory is used. Either writer
// may be nil; nil writers are treated as io.Discard. Errors are wrapped
// so callers can match ErrNoModule with errors.Is.
func Run(root string, stdout, stderr io.Writer) error {
	if stdout == nil {
		stdout = io.Discard
	}
	if stderr == nil {
		stderr = io.Discard
	}

	if root == "" {
		var err error
		root, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("cannot determine working directory: %w", err)
		}
	}

	module, err := ParseGoMod(root)
	if err != nil {
		return fmt.Errorf("parse %s/go.mod: %w", root, err)
	}

	_, _ = bannerColor.Fprintf(stdout, "• ")
	_, _ = fmt.Fprintf(stdout, "gear remod: rebuilding %s\n", module)

	if err := Remod(root, stdout, stderr, module); err != nil {
		return err
	}

	_, _ = bannerColor.Fprintf(stdout, "• ")
	_, _ = fmt.Fprintf(stdout, "gear remod: rebuilt %s\n", module)
	return nil
}
