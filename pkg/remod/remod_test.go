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

package remod_test

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cqroot/gear/pkg/remod"
)

const sampleGoMod = `module github.com/cqroot/gear

go 1.27.1

require github.com/spf13/cobra v1.10.2
`

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func TestParseGoMod(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "go.mod", sampleGoMod)

	mod, err := remod.ParseGoMod(dir)
	if err != nil {
		t.Fatalf("ParseGoMod: %v", err)
	}
	if mod != "github.com/cqroot/gear" {
		t.Fatalf("unexpected module path: %q", mod)
	}
}

func TestParseGoMod_missing(t *testing.T) {
	dir := t.TempDir()
	if _, err := remod.ParseGoMod(dir); err == nil {
		t.Fatal("expected error for missing go.mod, got nil")
	}
}

func TestParseGoMod_noModuleDirective(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "go.mod", "// just a comment\n")
	_, err := remod.ParseGoMod(dir)
	if err == nil || !errors.Is(err, remod.ErrNoModule) {
		t.Fatalf("expected ErrNoModule, got %v", err)
	}
}

func TestRemoveFiles_idempotent(t *testing.T) {
	dir := t.TempDir()
	if err := remod.RemoveFiles(dir, "go.mod", "go.sum"); err != nil {
		t.Fatalf("RemoveFiles: %v", err)
	}

	target := filepath.Join(dir, "go.mod")
	writeFile(t, dir, "go.mod", sampleGoMod)
	if err := remod.RemoveFiles(dir, "go.mod", "go.sum"); err != nil {
		t.Fatalf("RemoveFiles: %v", err)
	}
	if _, err := os.Stat(target); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected go.mod removed: %v", err)
	}
}

func TestRun_writesSummaryAndReinitializesModule(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "go.mod", sampleGoMod)
	writeFile(t, dir, "go.sum", "placeholder\n")
	// A real .go file forces `go mod tidy` to actually resolve the
	// declared dependency, which makes the go tool emit diagnostics on
	// stdout ("go: finding module for ...", "go: found ...") so the
	// test can assert that those bytes were forwarded through Run.
	writeFile(t, dir, "main.go", "package main\n")

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	if err := remod.Run(context.Background(), dir, stdout, stderr); err != nil {
		t.Fatalf("Run: %v", err)
	}

	// go.mod must have been recreated.
	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err != nil {
		t.Fatalf("expected regenerated go.mod: %v", err)
	}
	if !strings.Contains(stdout.String(), "• gear remod: rebuilt github.com/cqroot/gear") {
		t.Fatalf("expected success banner on stdout, got %q", stdout.String())
	}
	// The go tool must have produced output that flowed through to one
	// of the two streams. We don't pin stdout vs stderr (go mixes
	// them), only that nothing was silently dropped.
	combined := stdout.String() + stderr.String()
	if !strings.Contains(combined, "go:") {
		t.Fatalf("expected go tool output on stdout/stderr, got stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
}

func TestRun_missingModuleDirective(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "go.mod", "// nothing here\n")

	err := remod.Run(context.Background(), dir, nil, nil)
	if err == nil || !errors.Is(err, remod.ErrNoModule) {
		t.Fatalf("expected ErrNoModule, got %v", err)
	}
}

// Regression test: an invalid module input must make Remod fail before
// it deletes the existing go.mod / go.sum.
func TestRemod_rejectsEmptyModuleWithoutTouchingFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "go.mod", sampleGoMod)

	if err := remod.Remod(context.Background(), dir, nil, nil, ""); err == nil {
		t.Fatal("expected error for empty module, got nil")
	}

	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err != nil {
		t.Fatalf("existing go.mod was deleted: %v", err)
	}
}

// Regression test: a failed rebuild must restore the original go.mod / go.sum
// instead of leaving the project without them.
func TestRemod_restoresModuleFilesOnFailure(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "go.mod", sampleGoMod)
	writeFile(t, dir, "go.sum", "placeholder\n")

	// An invalid module path makes `go mod init` fail after the original
	// files have already been removed.
	err := remod.Remod(context.Background(), dir, &bytes.Buffer{}, &bytes.Buffer{}, "not a valid module")
	if err == nil {
		t.Fatal("expected error for invalid module, got nil")
	}

	assertFileContent(t, filepath.Join(dir, "go.mod"), sampleGoMod)
	assertFileContent(t, filepath.Join(dir, "go.sum"), "placeholder\n")
}

func assertFileContent(t *testing.T, path, want string) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if string(got) != want {
		t.Fatalf("%s = %q, want %q", path, got, want)
	}
}
