package main

import (
	"github.com/sanyatuning/powerline-go/powerline"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

var shell = "test"

func Symbols() powerline.Symbols {
	return powerline.Symbols{
		Branch:         "ß",
		CommitsAhead:   "^",
		CommitsBehind:  "v",
		Ellipsis:       "...",
		Lock:           "L",
		NewLine:        "\n",
		GitDiff:        "▲",
		Separator:      "->",
		SeparatorRight: "<-",
	}
}

func Test_root(t *testing.T) {
	username := "root"
	cwd := "/home/username/dirname"
	cwdParts := strings.Split(cwd, "/")

	theme := powerline.Dark()
	symbols := Symbols()
	p := powerline.NewPowerline(shell, symbols, theme)
	builder := powerline.NewSegmentBuilder(theme, symbols)

	testNow, _ := time.Parse(time.RFC3339Nano, "2013-06-05T14:10:43.678Z")
	width := "80"

	p.AppendLeft(builder.UserSegment(username))
	p.AppendLeft(builder.HostSegment("hostname"))
	p.AppendLeft(builder.PathSegment(cwdParts))
	p.AppendLeft(builder.LockSegment(cwd))
	p.AppendRight(builder.TimeSegment(testNow))
	p.AppendDown(builder.BashSegment(username))
	p.AppendDown(builder.ExitCodeSegment("0"))

	rootSegments := p.PrintAll(width)
	want := " root -> hostname -> /home/username/dirname -> L .R->.R     <- Wed 5 14:10:43 .R\n # .R->.R "

	assert.Equal(t, want, rootSegments)
}

func Test_user(t *testing.T) {
	cwd := "~/dirname"
	cwdParts := strings.Split(cwd, "/")

	theme := powerline.Dark()
	symbols := Symbols()
	p := powerline.NewPowerline(shell, symbols, theme)
	builder := powerline.NewSegmentBuilder(theme, symbols)

	testNow, _ := time.Parse(time.RFC3339Nano, "2013-06-05T14:10:43.678Z")

	width := "80"
	username := "user"

	p.AppendLeft(builder.UserSegment(username))
	p.AppendLeft(builder.HostSegment("hostname"))
	p.AppendLeft(builder.PathSegment(cwdParts))
	p.AppendLeft(builder.LockSegment(cwd))
	p.AppendRight(builder.TimeSegment(testNow))
	p.AppendDown(builder.BashSegment(username))
	p.AppendDown(builder.ExitCodeSegment("0"))

	rootSegments := p.PrintAll(width)
	want := " user -> hostname -> ~/dirname -> L .R->.R                  <- Wed 5 14:10:43 .R\n $ .R->.R "

	assert.Equal(t, want, rootSegments)
}

func Test_git(t *testing.T) {
	cwd := "~/dirname"
	cwdParts := strings.Split(cwd, "/")

	theme := powerline.Dark()
	symbols := Symbols()
	p := powerline.NewPowerline(shell, symbols, theme)
	builder := powerline.NewSegmentBuilder(theme, symbols)

	gitInfo := powerline.GitInfo{
		Branch:        "master",
		CommitsBehind: 12,
		CommitsAhead:  3,
		Staged:        true,
		Tag:           "2.0.2",
	}

	width := "80"
	username := "user"

	p.AppendLeft(builder.UserSegment(username))
	p.AppendLeft(builder.HostSegment("hostname"))
	p.AppendLeft(builder.PathSegment(cwdParts))
	p.AppendLeft(builder.GitSegment(gitInfo))
	p.AppendDown(builder.BashSegment(username))
	p.AppendDown(builder.ExitCodeSegment("0"))

	rootSegments := p.PrintAll(width)
	want := " user -> hostname -> ~/dirname -> ß master \"2.0.2\" 12v 3^ ▲ .R->.R\n $ .R->.R "

	assert.Equal(t, want, rootSegments)
}
