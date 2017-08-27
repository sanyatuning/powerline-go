package powerline

import (
	"log"
	"golang.org/x/sys/unix"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type SegmentBuilder struct {
	theme   Theme
	symbols Symbols
}

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
	NoGit         bool
	Staged        bool
	Tag           string
	Untracked     bool
}

func NewSegmentBuilder(theme Theme, symbols Symbols) SegmentBuilder {
	return SegmentBuilder{
		theme:   theme,
		symbols: symbols,
	}
}

func (s *SegmentBuilder) LockSegment(cwd string) Segment {
	if isWritableDir(cwd) {
		return Segment{value: ""}
	} else {
		return Segment{
			Bg:    s.theme.Lock.Bg,
			Fg:    s.theme.Lock.Fg,
			value: s.symbols.Lock,
		}
	}
}

func GetCurrentWorkingDir(username string) (string, []string) {
	dir, err := filepath.Abs(".")
	if err != nil {
		log.Fatal(err)
	}
	userDir := strings.Replace(dir, "/home/"+username, "~", 1)
	userDir = strings.TrimSuffix(userDir, "/")
	parts := strings.Split(userDir, "/")
	return dir, parts
}

func (s *SegmentBuilder) UserSegment(username string) Segment {
	c := s.theme.User
	if username == "root" {
		c = s.theme.Root
	}
	return Segment{
		Bg:    c.Bg,
		Fg:    c.Fg,
		value: username,
	}
}

func (s *SegmentBuilder) HostSegment(hostname string) Segment {
	c := getHostColor(s.theme.Host, hostname)

	return Segment{
		Bg:    c.Bg,
		Fg:    c.Fg,
		value: hostname,
	}
}

func (s *SegmentBuilder) PathSegment(cwdParts []string) Segment {
	var c = s.theme.Path
	if cwdParts[0] == "~" {
		c = s.theme.Home
	}

	return Segment{
		Bg:    c.Bg,
		Fg:    c.Fg,
		value: strings.Join(cwdParts, "/"),
	}
}

func GetGitInformation() GitInfo {
	var r = GitInfo{
		Branch:        "",
		CommitsAhead:  0,
		CommitsBehind: 0,
		Staged:        false,
		Tag:           "",
	}
	stdout, err := exec.Command("git", "status", "--ignore-submodules").Output()
	if err != nil {
		r.NoGit = true
		return r
	}
	reBranch := regexp.MustCompile(`^(HEAD detached at|HEAD detached from|On branch) (\S+)`)
	matchBranch := reBranch.FindStringSubmatch(string(stdout))
	if len(matchBranch) > 0 && matchBranch[1] == "On branch" {
		r.Branch = matchBranch[2]
	}

	reStatus := regexp.MustCompile(`Your branch is (ahead|behind).*?([0-9]+) comm`)
	matchStatus := reStatus.FindStringSubmatch(string(stdout))
	if len(matchStatus) > 0 {
		if matchStatus[1] == "behind" {
			r.CommitsBehind, _ = strconv.Atoi(matchStatus[2])
		} else if matchStatus[1] == "ahead" {
			r.CommitsAhead, _ = strconv.Atoi(matchStatus[2])
		}
	} else {
		reStatus := regexp.MustCompile(`Your branch and.*\n.*(\d+) and (\d+) diff`)
		matchStatus := reStatus.FindStringSubmatch(string(stdout))
		if len(matchStatus) > 0 {
			r.CommitsBehind, _ = strconv.Atoi(matchStatus[2])
			r.CommitsAhead, _ = strconv.Atoi(matchStatus[1])
		}
	}

	r.Staged = !strings.Contains(string(stdout), "nothing to commit")
	r.Untracked = strings.Contains(string(stdout), "Untracked files")

	output, _ := exec.Command("git", "describe", "--tags", "--exact").Output()
	r.Tag = strings.TrimSpace(string(output))

	return r
}

func (s *SegmentBuilder) GitSegment(gitInfo GitInfo) Segment {
	if gitInfo.NoGit {
		return Segment{value: ""}
	}

	bg := s.theme.Git.Clean.Bg
	fg := s.theme.Git.Clean.Fg
	gitStatus := s.symbols.Branch + " " + gitInfo.Branch
	if gitInfo.Branch == "" {
		gitStatus = s.symbols.Branch + " no branch!"
	}
	if gitInfo.Tag != "" {
		gitStatus += " \"" + gitInfo.Tag + "\""
	}
	if gitInfo.CommitsBehind > 0 {
		gitStatus += " " + strconv.Itoa(gitInfo.CommitsBehind) + s.symbols.CommitsBehind
	}
	if gitInfo.CommitsAhead > 0 {
		gitStatus += " " + strconv.Itoa(gitInfo.CommitsAhead) + s.symbols.CommitsAhead
	}
	if gitInfo.Staged {
		bg = s.theme.Git.Dirty.Bg
		fg = s.theme.Git.Dirty.Fg
		gitStatus += " " + s.symbols.GitDiff
	}
	return Segment{
		Bg:    bg,
		Fg:    fg,
		value: gitStatus,
	}
}

func (s *SegmentBuilder) ExitCodeSegment(code string) Segment {
	i, err := strconv.Atoi(code)
	if err != nil || i == 0 {
		return Segment{value: ""}
	}

	return Segment{
		Bg:    s.theme.Error.Bg,
		Fg:    s.theme.Error.Fg,
		value: code,
	}
}

func (s *SegmentBuilder) BashSegment(username string) Segment {
	v := "$"
	if username == "root" {
		v = "#"
	}
	return Segment{
		Bg:    s.theme.Path.Bg,
		Fg:    s.theme.Path.Fg,
		value: v,
	}
}

func (s *SegmentBuilder) TimeSegment(time time.Time) Segment {
	return Segment{
		Bg:    s.theme.User.Bg,
		Fg:    s.theme.User.Fg,
		value: time.Format("Mon 2 15:04:05"),
	}
}

func isWritableDir(dir string) bool {
	return unix.Access(dir, unix.W_OK) == nil
}

func getHostColor(h Host, hostname string) ColorPair {
	if match("-desktop$", hostname) {
		return h.Desktop
	}
	if match("^(p|vcc|[a-z]{2}[1-9])-", hostname) {
		return h.Production
	}
	if match("^syslog-", hostname) {
		return h.Special
	}

	return h.Other
}

func match(pattern string, hostname string) bool {
	m, _ := regexp.MatchString(pattern, hostname)
	return m
}
