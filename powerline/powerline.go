package powerline

import (
	"fmt"
	"bytes"
	"strings"
	"strconv"
	"regexp"
	"unicode/utf8"
)

type Symbols struct {
	Lock               string
	Network            string
	Separator          string
	SeparatorThin      string
	SeparatorRight     string
	SeparatorThinRight string
	Ellipsis           string
	NewLine            string
}

func DefaultSymbols() Symbols {
	return Symbols{
		Lock:               "\uE0A2",
		Network:            "\uE0A2",
		Separator:          "\uE0B0",
		SeparatorThin:      "\uE0B1",
		SeparatorRight:     "\uE0B2",
		SeparatorThinRight: "\uE0B3",
		Ellipsis:           "\u2026",
		NewLine:            "\n",
	}
}
func TestSymbols() Symbols {
	return Symbols{
		Lock:               "L",
		Network:            "N",
		Separator:          "->",
		SeparatorThin:      ">",
		SeparatorRight:     "<-",
		SeparatorThinRight: "<",
		Ellipsis:           "...",
		NewLine:            "\n",
	}
}

type Powerline struct {
	ShTemplate    string // still not quite get it
	ColorTemplate string // how to output color
	ShellBg       string
	Reset         string
	Symbols       Symbols
	SegmentsLeft  []Segment
	SegmentsRight []Segment
	SegmentsDown  []Segment
}

func (p *Powerline) color(prefix string, code string) string {
	return fmt.Sprintf(
		p.ShTemplate,
		fmt.Sprintf(p.ColorTemplate, prefix, code),
	)
}

func (p *Powerline) fgColor(code string) string {
	return p.color("38", code)
}

func (p *Powerline) bgColor(code string) string {
	return p.color("48", code)
}

func (p *Powerline) AppendLeft(s Segment) {
	p.SegmentsLeft = append(p.SegmentsLeft, s)
}

func (p *Powerline) AppendRight(s Segment) {
	p.SegmentsRight = append(p.SegmentsRight, s)
}

func (p *Powerline) AppendDown(s Segment) {
	p.SegmentsDown = append(p.SegmentsDown, s)
}

func (p *Powerline) PrintAll(width string) string {
	var buffer bytes.Buffer
	var re = regexp.MustCompile(`\\\[\\e.*?\\]`)
	left := p.PrintSegments(p.SegmentsLeft, false)
	right := p.PrintSegments(p.SegmentsRight, true)

	s1 := re.ReplaceAllString(left, ``)
	s2 := re.ReplaceAllString(right, ``)

	i, err := strconv.Atoi(width)
	count := 5
	if err == nil {
		count = i - utf8.RuneCountInString(s1) - utf8.RuneCountInString(s2)
	}
	if count < 0 {
		count = 5
	}
	buffer.WriteString(left)
	buffer.WriteString(strings.Repeat(" ", count))
	buffer.WriteString(right)
	buffer.WriteString(p.Symbols.NewLine)
	buffer.WriteString(p.PrintSegments(p.SegmentsDown, false))
	buffer.WriteString(" ")
	return buffer.String()
}

func (p *Powerline) PrintSegments(segments []Segment, right bool) string {
	if len(segments) == 0 {
		return ""
	}

	var buffer bytes.Buffer

	for i, cur := range segments {
		next := getNext(segments, i)
		buffer.WriteString(p.PrintSegment(cur, next, right))
	}

	buffer.WriteString(p.Reset)

	return buffer.String()
}

func (p *Powerline) PrintSegment(segment Segment, next *Segment, right bool) string {
	if segment.values == nil {
		return ""
	}
	var buffer bytes.Buffer
	var nextBg string
	if next == nil {
		// if it is the last one, switch to shell bg
		nextBg = p.Reset
	} else {
		nextBg = p.bgColor(next.Bg)
	}

	// print parts with correct foregrounds
	for j, segPart := range segment.values {
		if right {
			if j == 0 {
				buffer.WriteString(p.fgColor(segment.Bg))
				buffer.WriteString(p.Symbols.SeparatorRight)
			} else {
				// while not last part
				buffer.WriteString(p.fgColor(segment.sepFg))
				buffer.WriteString(p.Symbols.SeparatorThinRight)
			}

			buffer.WriteString(p.fgColor(segment.Fg))
			buffer.WriteString(p.bgColor(segment.Bg))
			buffer.WriteString(fmt.Sprintf(" %s ", segPart))
		} else {
			buffer.WriteString(p.bgColor(segment.Bg))
			buffer.WriteString(p.fgColor(segment.Fg))
			buffer.WriteString(fmt.Sprintf(" %s ", segPart))
			if (j + 1) == len(segment.values) {
				// last part switches background to next
				buffer.WriteString(nextBg)
				buffer.WriteString(p.fgColor(segment.Bg))
				buffer.WriteString(p.Symbols.Separator)
			} else {
				// while not last part
				buffer.WriteString(p.fgColor(segment.sepFg))
				buffer.WriteString(p.Symbols.SeparatorThin)
			}
		}
	}

	return buffer.String()
}

func NewPowerline(shell string, sym Symbols, t Theme) Powerline {
	var p Powerline
	if shell == "test" {
		p = Powerline{
			ShTemplate:    "%s",
			ColorTemplate: "%.s%.s",
			Reset:         ".R",
		}
	} else if shell == "zsh" {
		p = Powerline{
			ShTemplate:    "%s",
			ColorTemplate: "%%{[%s;5;%sm%%}",
			Reset:         "%{$reset_color%}",
		}
	} else {
		p = Powerline{
			ShTemplate:    "\\[\\e%s\\]",
			ColorTemplate: "[%s;5;%sm",
			Reset:         "\\[\\e[0m\\]",
		}
	}
	p.ShellBg = t.ShellBg
	p.Symbols = sym

	return p
}

func getNext(segments []Segment, i int) *Segment {
	i++
	for i < len(segments) {
		if segments[i].values != nil {
			return &segments[i]
		}
		i++
	}
	return nil
}