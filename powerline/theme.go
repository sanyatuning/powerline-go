package powerline

type ColorPair struct {
	Bg string
	Fg string
}

type ColorTriplet struct {
	Bg    string
	Fg    string
	SepFg string
}

type Git struct {
	Clean ColorPair
	Dirty ColorPair
}

type Host struct {
	Desktop    ColorPair
	Other      ColorPair
	Production ColorPair
	Special    ColorPair
}

type Theme struct {
	ShellBg string
	Host    Host
	User    ColorPair
	Root    ColorPair
	Path    ColorPair
	Home    ColorPair
	Git
	Lock  ColorPair
	Error ColorPair
}

func Dark() Theme {
	return Theme{
		ShellBg: "3",
		Root:    ColorPair{Bg: "1", Fg: "7"},
		User:    ColorPair{Bg: "22", Fg: "7"},
		Host: Host{
			Desktop:    ColorPair{Bg: "34", Fg: "0"},
			Other:      ColorPair{Bg: "11", Fg: "0"},
			Production: ColorPair{Bg: "1", Fg: "7"},
			Special:    ColorPair{Bg: "14", Fg: "0"},
		},
		Path: ColorPair{Bg: "0", Fg: "214"},
		Home: ColorPair{Bg: "82", Fg: "0"},
		Git: Git{
			Clean: ColorPair{Bg: "0", Fg: "10"},
			Dirty: ColorPair{Bg: "0", Fg: "202"},
		},
		Lock:  ColorPair{Bg: "1", Fg: "7"},
		Error: ColorPair{Bg: "1", Fg: "7"},
	}
}
