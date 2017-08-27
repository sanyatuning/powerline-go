// Copyright 2014 Matt Martz <matt@sivel.net>
// Modifications copyright (c) 2013 Anton Chebotaev <anton.chebotaev@gmail.com>
//
// All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package main

import (
	"fmt"
	"github.com/sanyatuning/powerline-go/powerline"
	"os"
	"time"
)

func main() {
	shell := "bash"

	if len(os.Args) > 1 {
		shell = os.Args[1]
	}

	exitCode := "0"
	if len(os.Args) > 2 {
		exitCode = os.Args[2]
	}
	width := "0"
	if len(os.Args) > 3 {
		width = os.Args[3]
	}

	theme := powerline.Dark()
	symbols := powerline.DefaultSymbols()
	p := powerline.NewPowerline(shell, symbols, theme)

	cwd, cwdParts := powerline.GetCurrentWorkingDir()
	gitInfo := powerline.GetGitInformation()
	username := os.Getenv("USER")
	hostname, _ := os.Hostname()

	p.AppendLeft(powerline.UserSegment(theme, username))
	p.AppendLeft(powerline.HostSegment(theme, hostname))
	p.AppendLeft(powerline.PathSegment(cwdParts, theme))
	p.AppendLeft(powerline.LockSegment(cwd, theme, symbols))
	p.AppendRight(powerline.GitSegment(theme, gitInfo))
	p.AppendRight(powerline.TimeSegment(time.Now(), theme))
	p.AppendDown(powerline.BashSegment(theme, username))
	p.AppendDown(powerline.ExitCodeSegment(exitCode, theme))

	fmt.Print(p.PrintAll(width))
}
