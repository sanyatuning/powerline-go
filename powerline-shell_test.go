package main

import (
	"github.com/sanyatuning/powerline-go/powerline"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

var shell = "test"

func Test_root(t *testing.T) {
	cwd := "/home/username/dirname"
	cwdParts := strings.Split(cwd, "/")

	theme := powerline.Dark()
	symbols := powerline.TestSymbols()
	p := powerline.NewPowerline(shell, symbols, theme)
	testNow, err := time.Parse(time.RFC3339Nano, "2013-06-05T14:10:43.678Z")
	if err != nil {
		t.Error(err)
	}
	gitInfo := powerline.GitInfo{
		Branch: "master",
		Staged: false,
	}
	width := "80"
	username := "root"

	p.AppendLeft(powerline.UserSegment(theme, username))
	p.AppendLeft(powerline.HostSegment(theme, "hostname"))
	p.AppendLeft(powerline.PathSegment(cwdParts, theme))
	p.AppendLeft(powerline.LockSegment(cwd, theme, symbols))
	p.AppendRight(powerline.GitSegment(theme, gitInfo))
	p.AppendRight(powerline.TimeSegment(testNow, theme))
	p.AppendDown(powerline.BashSegment(theme, username))
	p.AppendDown(powerline.ExitCodeSegment("0", theme))

	rootSegments := p.PrintAll(width)
	want := " root -> hostname -> /home/username/dirname -> L .R->.R     <-  master <- Wed 5 14:10:43 .R\n # .R->.R "

	assert.Equal(t, want, rootSegments)
}

func Test_user(t *testing.T) {
	cwd := "~/dirname"
	cwdParts := strings.Split(cwd, "/")

	theme := powerline.Dark()
	symbols := powerline.TestSymbols()
	p := powerline.NewPowerline(shell, symbols, theme)
	testNow, err := time.Parse(time.RFC3339Nano, "2013-06-05T14:10:43.678Z")
	if err != nil {
		t.Error(err)
	}
	gitInfo := powerline.GitInfo{
		Branch: "master",
		Staged: false,
		Tag: "2.0.2",
	}

	width := "80"
	username := "user"

	p.AppendLeft(powerline.UserSegment(theme, username))
	p.AppendLeft(powerline.HostSegment(theme, "hostname"))
	p.AppendLeft(powerline.PathSegment(cwdParts, theme))
	p.AppendLeft(powerline.LockSegment(cwd, theme, symbols))
	p.AppendRight(powerline.GitSegment(theme, gitInfo))
	p.AppendRight(powerline.TimeSegment(testNow, theme))
	p.AppendDown(powerline.BashSegment(theme, username))
	p.AppendDown(powerline.ExitCodeSegment("0", theme))

	rootSegments := p.PrintAll(width)
	want := " user -> hostname -> ~/dirname -> L .R->.R     <-  master \"2.0.2\" <- Wed 5 14:10:43 .R\n $ .R->.R "

	assert.Equal(t, want, rootSegments)
}
