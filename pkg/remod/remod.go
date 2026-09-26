// Package remod provides core utilities used by the gear CLI commands.
package remod

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ErrNoModule is returned by ParseGoMod when the go.mod file does not
// contain a `module` directive.
var ErrNoModule = errors.New("go.mod does not contain a module directive")

// ParseGoMod reads a go.mod file located at the given root and returns the
// module path declared inside it. If no `module` directive is found,
// ErrNoModule is returned.
func ParseGoMod(root string) (string, error) {
	f, err := os.Open(filepath.Join(root, "go.mod"))
	if err != nil {
		return "", fmt.Errorf("open go.mod: %w", err)
	}
	defer f.Close()

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
// stdout and stderr are routed to stdout and stderr in real time.
//
// Either writer may be nil; nil writers are treated as io.Discard.
func RunGoMod(stdout, stderr io.Writer, root string, args ...string) error {
	if stdout == nil {
		stdout = io.Discard
	}
	if stderr == nil {
		stderr = io.Discard
	}

	cmd := exec.Command("go", append([]string{"mod"}, args...)...)
	cmd.Dir = root
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
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

	if err := RunGoMod(stdout, stderr, root, "init", module); err != nil {
		return fmt.Errorf("go mod init failed: %w", err)
	}

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

	if err := Remod(root, stdout, stderr, module); err != nil {
		return err
	}

	fmt.Fprintf(stdout, "module %s rebuilt in %s\n", module, root)
	return nil
}
