package main

import (
	"testing"
	"./powerline"
	"strings"
	"github.com/stretchr/testify/assert"
	"time"
)

var shell = "test"

func Test_root(t *testing.T) {
	cwd := "/home/username/dirname"
	cwdParts := strings.Split(cwd, "/")

	theme := powerline.SolarizedDark()
	symbols := powerline.TestSymbols()
	p := powerline.NewPowerline(shell, symbols, theme)
	testNow, err := time.Parse(time.RFC3339Nano, "2013-06-05T14:10:43.678Z")
	if err != nil {
		t.Error(err)
	}
	gitStatus := "master"
	gitStaged := false
	width := "80"

	p.AppendLeft(powerline.HomeSegment(cwdParts, theme))
	p.AppendLeft(powerline.PathSegment(cwdParts, theme, symbols))
	p.AppendLeft(powerline.LockSegment(cwd, theme, symbols))
	p.AppendRight(powerline.GitSegment(theme, gitStatus, gitStaged))
	p.AppendRight(powerline.TimeSegment(testNow, theme))
	p.AppendDown(powerline.BashSegment(theme))
	p.AppendDown(powerline.ExitCodeSegment("0", theme))

	rootSegments := p.PrintAll(width)
	want := " / > home > username > dirname -> L .R->.R        <- master <- Wed 5 14:10:43 .R\n $ .R->.R "

	assert.Equal(t, want, rootSegments)
}
