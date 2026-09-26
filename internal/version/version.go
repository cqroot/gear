package version

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

// Version info variables. Overridden at build time via -ldflags.
var (
	version   = "dev"
	commit    = "none"
	date      = "unknown"
	builtWith = fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
)

// labelColor styles the labels in the version block.
var labelColor = color.New(color.FgCyan)

// Info contains version and build information.
type Info struct {
	Version   string
	Commit    string
	Date      string
	BuiltWith string
}

// Get returns the current version information.
func Get() Info {
	return Info{
		Version:   version,
		Commit:    commit,
		Date:      date,
		BuiltWith: builtWith,
	}
}

// String returns a formatted version info string.
func (i Info) String() string {
	var sb strings.Builder
	sb.WriteString("\n  ")
	sb.WriteString(labelColor.Sprint("• Version:      "))
	sb.WriteString(i.Version)

	sb.WriteString("\n  ")
	sb.WriteString(labelColor.Sprint("• Commit:       "))
	sb.WriteString(i.Commit)

	sb.WriteString("\n  ")
	sb.WriteString(labelColor.Sprint("• Built at:     "))
	sb.WriteString(i.Date)

	sb.WriteString("\n  ")
	sb.WriteString(labelColor.Sprint("• Built with:   "))
	sb.WriteString(i.BuiltWith)
	return sb.String()
}
