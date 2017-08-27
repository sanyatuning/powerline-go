package powerline

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Segment struct {
	Bg    string
	Fg    string
	sepFg string
	value string
}

type GitInfo struct {
	Branch        string
	CommitsAhead  int
	CommitsBehind int
	Staged        bool
	Tag           string
}

func isWritableDir(dir string) bool {
	tmpPath := path.Join(dir, ".powerline-write-test")
	_, err := os.Create(tmpPath)
	if err != nil {
		return false
	}
	os.Remove(tmpPath)
	return true
}

func LockSegment(cwd string, t Theme, s Symbols) Segment {
	if isWritableDir(cwd) {
		return Segment{value: ""}
	} else {
		return Segment{
			Bg:    t.Lock.Bg,
			Fg:    t.Lock.Fg,
			value: s.Lock,
		}
	}
}

func GetCurrentWorkingDir() (string, []string) {
	dir, err := filepath.Abs(".")
	if err != nil {
		log.Fatal(err)
	}
	userDir := strings.Replace(dir, os.Getenv("HOME"), "~", 1)
	userDir = strings.TrimSuffix(userDir, "/")
	parts := strings.Split(userDir, "/")
	return dir, parts
}

func UserSegment(t Theme, username string) Segment {
	c := t.User
	if username == "root" {
		c = t.Root
	}
	return Segment{
		Bg:    c.Bg,
		Fg:    c.Fg,
		value: username,
	}
}

func HostSegment(t Theme, hostname string) Segment {
	c := t.Host.Other
	if m, _ := regexp.MatchString("-desktop$", hostname); m {
		c = t.Host.Desktop
	}
	return Segment{
		Bg:    c.Bg,
		Fg:    c.Fg,
		value: hostname,
	}
}

func PathSegment(cwdParts []string, t Theme) Segment {
	var c = t.Path
	if cwdParts[0] == "~" {
		c = t.Home
	}

	return Segment{
		Bg:    c.Bg,
		Fg:    c.Fg,
		value: strings.Join(cwdParts, "/"),
	}
}

func GetGitInformation() GitInfo {
	var status string
	var staged bool
	stdout, _ := exec.Command("git", "status", "--ignore-submodules").Output()
	reBranch := regexp.MustCompile(`^(HEAD detached at|HEAD detached from|On branch) (\S+)`)
	matchBranch := reBranch.FindStringSubmatch(string(stdout))
	if len(matchBranch) > 0 {
		if matchBranch[2] == "detached" {
			status = matchBranch[2]
		} else {
			status = matchBranch[2]
		}
	}

	reStatus := regexp.MustCompile(`Your branch is (ahead|behind).*?([0-9]+) comm`)
	matchStatus := reStatus.FindStringSubmatch(string(stdout))
	if len(matchStatus) > 0 {
		status = fmt.Sprintf("%s %s", status, matchStatus[2])
		if matchStatus[1] == "behind" {
			status = fmt.Sprintf("%s\u21E3", status)
		} else if matchStatus[1] == "ahead" {
			status = fmt.Sprintf("%s\u21E1", status)
		}
	} else {
		reStatus := regexp.MustCompile(`Your branch and.*\n.*(\d+) and (\d+) diff`)
		matchStatus := reStatus.FindStringSubmatch(string(stdout))
		if len(matchStatus) > 0 {
			status = fmt.Sprintf(
				"%s %s\u21E3 %s\u21E1",
				status,
				matchStatus[1],
				matchStatus[2],
			)
		}
	}

	staged = !strings.Contains(string(stdout), "nothing to commit")
	if strings.Contains(string(stdout), "Untracked files") {
		status = fmt.Sprintf("%s +", status)
	}

	tag, _ := exec.Command("git", "describe", "--tags", "--exact").Output()

	return GitInfo{
		Branch:        status,
		CommitsAhead:  0,
		CommitsBehind: 0,
		Staged:        staged,
		Tag:           strings.TrimSpace(string(tag)),
	}
}

func GitSegment(t Theme, gitInfo GitInfo) Segment {

	gitStatus := gitInfo.Branch
	if gitStatus != "" {
		var bg = t.Git.Clean.Bg
		var fg = t.Git.Clean.Fg
		gitStatus = " " + gitStatus
		if gitInfo.Staged {
			bg = t.Git.Dirty.Bg
			fg = t.Git.Dirty.Fg
			gitStatus += " ▲"
			//gitStatus += " ▲↑↓"
		}
		if gitInfo.Tag != "" {
			gitStatus += " \"" + gitInfo.Tag + "\""
		}
		return Segment{
			Bg:    bg,
			Fg:    fg,
			value: gitStatus,
		}
	} else {
		return Segment{value: ""}
	}
}

func ExitCodeSegment(code string, t Theme) Segment {
	i, err := strconv.Atoi(code)
	if err != nil || i == 0 {
		return Segment{value: ""}
	} else {
		return Segment{
			Bg:    t.Error.Bg,
			Fg:    t.Error.Fg,
			value: code,
		}
	}
}

func BashSegment(t Theme, username string) Segment {
	v := "$"
	if username == "root" {
		v = "#"
	}
	return Segment{
		Bg:    t.Path.Bg,
		Fg:    t.Path.Fg,
		value: v,
	}
}

func TimeSegment(time time.Time, t Theme) Segment {
	return Segment{
		Bg:    t.User.Bg,
		Fg:    t.User.Fg,
		value: time.Format("Mon 2 15:04:05"),
	}
}
