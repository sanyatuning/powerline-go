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

type Theme struct {
	ShellBg string
	Home ColorPair
	Path ColorTriplet
	Git
	Lock ColorPair
	Error ColorPair
}

func Dark() Theme {
	return Theme{
		ShellBg: "3",
		Home: ColorPair{Bg: "237", Fg: "7"},
		Path: ColorTriplet{Bg: "25", Fg: "254"},
		Git: Git{
			Clean: ColorPair{Bg: "148", Fg: "0"},
			Dirty: ColorPair{Bg: "161", Fg: "15"},
		},
		Lock: ColorPair{Bg: "3", Fg: "7"},
		Error: ColorPair{Bg: "1", Fg: "7"},
	}
}

func SolarizedDark() Theme {
	return Theme{
		ShellBg: "0",
		Home: ColorPair{Bg: "10", Fg: "0"},
		Path: ColorTriplet{Bg: "8", Fg: "12", SepFg: "0"},
		Git: Git{
			Clean: ColorPair{Bg: "14", Fg: "0"},
			Dirty: ColorPair{Bg: "2", Fg: "0"},
		},
		Lock: ColorPair{Bg: "4", Fg: "7"},
		Error: ColorPair{Bg: "1", Fg: "7"},
	}
}
