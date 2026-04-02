package main

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	colPink   = lipgloss.Color("213")
	colDim    = lipgloss.Color("246")
	colBody   = lipgloss.Color("255")
	colBorder = lipgloss.Color("240")
	colTag    = lipgloss.Color("245")
	colHint   = lipgloss.Color("244")
)

var (
	stylePink     = lipgloss.NewStyle().Foreground(colPink)
	styleDim      = lipgloss.NewStyle().Foreground(colDim)
	styleBody     = lipgloss.NewStyle().Foreground(colBody)
	styleBorder   = lipgloss.NewStyle().Foreground(colBorder)
	styleHint     = lipgloss.NewStyle().Foreground(colHint)
	styleTag      = lipgloss.NewStyle().Foreground(colTag)
	styleTitleBar = lipgloss.NewStyle().Foreground(colDim)
	stylePortrait = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	styleName     = lipgloss.NewStyle().Foreground(colPink).Bold(true)

	stylePanel = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(colBorder).
			Padding(1, 2)

	styleNavActive = lipgloss.NewStyle().
			Foreground(colPink).
			Bold(true)

	styleNavInactive = lipgloss.NewStyle().
				Foreground(colDim)
)

const (
	breakpointSmall  = 95
	breakpointMedium = 110
)

var portraitLines = []string{
	"$XXXX$XXXXxx+xXXXX$$$$$$$$$$X$$$$$X$$$XxXX$$$X+x++X$Xx+Xxxxxx",
	"XX$XX$$XXXXXXXXXX$$$$$$$$$$$$$$$$$$$$$X+X$$$$$$X++X;+x;XX+x+X",
	"XX$$XX$$XX$$XXXX$&$$$$$$$$$$$$$$$$XXXXX+X$$$$$$$xxx+;xx;;++.+",
	"XX$$$$X$$$$$XX$$&$$$$$&&&&$$$$Xx+;;;;++xXX$$$$$$X;;++X:..xx..",
	"XX$Xx$$$$$$$$X$&&&&&&&&&&$$$Xx++;;;::;;;++xX&&$$$XXXXxxx;++xx",
	"$XX$+x$+$X$$$$&&&&&&&&&&$&$Xx++++;;;;;;;;;+xX&$$$$XX$Xxxxx+;X",
	"XXXXx;x+X+$$$&&&&&&&&&&&&&X++;;;;;;;;;;;;;;+x$&&$&$XX+$$xx+;:",
	"+XXXx++xX+X$&&&&&&&&&&&&&Xxx++;;;;;;;;;;;;;;+$&&$&$x+xXXXXXxx",
	"++XXX++xxxX&&&&&&&&&&&&$XxxXXX$X+;;;;;:;:::;+$&&&&$+xX$$XXX$$",
	"++$$Xxxx+x$&&&&&&&&&&&Xxx$$&&XXXXx++;;;;;;;;+X&&&&$x$$$XXx$$$",
	"+x+Xx+xxxx&&&&&&&&&&&X+;+xxxx+XXXx++;+xXXXXXX&&&&&$X$$+X$$$$X",
	";x+X+xx+x$&&&&&$&&&&$+;:::;;+++++;::;x$&$XXx$&&&&&&$+$XXXXxXx",
	"xxx+xxx+X&&&&&&&&&&$x+;;::::;;;;::.:;+xXx+XX&&&&&&$Xx;x$X+XXX",
	"Xxxxxxx$&&&&&&&$&&&Xx++;;;;;;+;;::::;;;++++X&&&&&&&xxxxx+xxxx",
	"XXxXxxX&&&&&&&&&&&$xxxx+++++++::.::::::::;+&&&&&&&&XxxxxxXx++",
	"xxxxXX$&&&&&&&&&&&$x++++x++;+xXx;;::;;:::;$&&&&&&&&XXxXxx+;;x",
	"xXXXX$&&&&&&&&&&&&$Xxx++xX$Xx+xxx+;;+x;;+$&&&&&&&&&Xxxx+;xx++",
	"xX$$$&&&&&&&&&&&&&$Xxxx+;++xxXXxxx+;+++X+&&&&&&&&&&$xx+xxx+xx",
	"XXX$&&&&&&&&&&&&&&$XXxx+;;+xxxxx++++;+&&&&&&&&&&&&&$xXx+:xxXx",
	"X$$$&&&&&&&&&&&&&$$XXXXx+;;;;;;;;;;x&&&&&&&&&&&&&&&&$Xx+X+++x",
	"$xXX&&&&&&&&&&&&&$$XXXxxxx++;;;;;X&&&&&&&&&&&&&&&&&&$$xXxx+xX",
	"::;++X$$$$&&&&$&&$$XXxxxxxxxxxxXX&&&&&&&&&&&&&&&&&&&$XxXX+Xxx",
	";;::xXX$$&&&&$$&&$$XXXxxxxxxxXXX$&$&&&&&$&&&&&&&&&&&&x+xx$xxx",
	";.;x+X$$$&&&XX&&&$$XXXXxxxxXXX$$&$$$$&&$X$XX$&&&&&&&&$x...xXx",
	" :xXXX$$&&$$x$&&$$XXXXXXXXX$XX$&$$$$&&&$$++x$&&&XxXX$$$$X:::x",
	" ;+xXX$$&&$xX&&&$XXXXXXXXX$XX$&X$$$$&&$x;+X&&&X+;::..;X$&&+::",
	":+xxX$$$$$x+$&&&XXxXXxXXXXXxX&XX$$&&&$;::X&$X+;::......X$XXx;",
	"x$$$$$$$$X+X&&$XXXxxxxXXXx+x$xX$$&&&$;;x$Xx;:::........:$$XxX",
	"X$$$$$$&X+x&&&+::xxx+xxxx;;XxX$$$&&$+x$+;:::..........::X$$XX",
	"$&&&&&&&++$&&X;:...     ::+xXX$$&&&xXx;:::...........:::X$$$$",
	"$&&&&&&+;X&&&;:::..   ..:;xx$$$$&&$x+;::..........::::;+&$X$$",
	"&&&&&&X;x$&&X.:::...  .::+XX$$$$&&X+::..........:::::;+&&&$$X",
	"&&&&&$;+$&$&;...........:xX$$$&&&&;:..........:::::;;;+&&$$$&",
	"&&&&$;+$&$$$;:::..:::.:.+X$$$$&&X:.........:::::::;;;++&&&&$&",
	"&&&$;+$$$X$X.::::::::;..X$$$&&$:.........:::::::;;;;;:X$$$$&$",
}

var nameArtLarge = []string{
	"(_)                          ",
	" _ _ __ ___  _ __ ___ _ __   ",
	"| | '_ ` _ \\| '__/ _ \\ '_ \\  ",
	"| | | | | | | | |  __/ | | | ",
	"|_|_| |_| |_|_|  \\___|_| |_| ",
}

var nameArtMedium = []string{
	" _",
	"(_)_ __ ___  _ __ ___ _ __",
	"| | '_ ` _ \\| '__/ _ \\ '_ \\",
	"| | | | | | | | |  __/ | | |",
	"|_|_| |_| |_|_|  \\___|_| |_|",
}

var nameArtSmall = []string{
	"imren",
}

type entry struct {
	key, val string
	tags     []string
}

type section struct {
	title   string
	entries []entry
}

var sections = map[string]section{
	"projects": {
		"/ projects",
		[]entry{
			{"HomeLab", "Virtual home lab for malware analysis, threat intelligence and penetration testing.", []string{"virtualbox", "kali", "ubuntu"}},
			{"RapidDeploymentKit", "Streamlined network deployment.", []string{"hardware", "python", "switchconfig"}},
			{"WyncoServices", "Real-time project status and financial dashboard.", []string{"typescript", "react", "dataviz"}},
		},
	},
	"links": {
		"/ links",
		[]entry{
			{"github", "github.com/imrenmore", nil},
			{"linkedin", "linkedin.com/in/imrenmore", nil},
			{"resume", "imrenmore.com", nil},
		},
	},
	"contact": {
		"/ contact",
		[]entry{
			{"email", "imrenkmore@gmail.com", nil},
			{"availability", "Open to opportunities.", nil},
			{"response time", "Within 24 hours.", nil},
		},
	},
}

var navOrder = []string{"projects", "links", "contact"}

type tickMsg struct{}

func typeCmd() tea.Cmd {
	return tea.Tick(14*time.Millisecond, func(time.Time) tea.Msg { return tickMsg{} })
}

type model struct {
	width, height int
	navIdx        int
	typedBio      string
	bioFull       string
	typing        bool
	typeTick      int
}

func newModel() model {
	return model{
		bioFull: "CS & Mathematics graduate building secure,\nhuman-centered technology.\n\nCybersecurity + design — exploring how digital systems shape identity & UX.\n\nExplore the directories below:",
		typing:  true,
	}
}

func (m model) Init() tea.Cmd {
	return typeCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tickMsg:
		if m.typing {
			r := []rune(m.bioFull)
			if m.typeTick < len(r) {
				m.typeTick++
				m.typedBio = string(r[:m.typeTick])
				return m, typeCmd()
			}
			m.typing = false
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab", "right", "l":
			m.navIdx = (m.navIdx + 1) % len(navOrder)
			return m, tea.ClearScreen
		case "shift+tab", "left", "h":
			m.navIdx = (m.navIdx - 1 + len(navOrder)) % len(navOrder)
			return m, tea.ClearScreen
		case "1":
			m.navIdx = 0
			return m, tea.ClearScreen
		case "2":
			m.navIdx = 1
			return m, tea.ClearScreen
		case "3":
			m.navIdx = 2
			return m, tea.ClearScreen
		}
	}

	return m, nil
}

func clamp(n, min, max int) int {
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func pickNameArt(w int) []string {
	switch {
	case w >= breakpointMedium:
		return nameArtLarge
	case w >= breakpointSmall:
		return nameArtMedium
	default:
		return nameArtSmall
	}
}

func renderNav(idx int, compact bool) string {
	if compact {
		var out []string
		for i, key := range navOrder {
			label := strings.ToUpper(key[:1]) + key[1:]
			if i == idx {
				out = append(out, styleNavActive.Render("["+label+"]"))
			} else {
				out = append(out, styleNavInactive.Render(label))
			}
		}
		return strings.Join(out, " ")
	}

	cellW := 12
	var out []string
	for i, key := range navOrder {
		label := strings.ToUpper(key[:1]) + key[1:]
		if i == idx {
			out = append(out, styleNavActive.Width(cellW).Render("+ "+label))
		} else {
			out = append(out, styleNavInactive.Width(cellW).Render("  "+label))
		}
	}
	return strings.Join(out, "")
}

func renderPanel(sec section, width int, compact bool) string {
	rows := []string{
		stylePink.Render(strings.ToUpper(sec.title)),
		styleBorder.Render(strings.Repeat("─", width)),
		"",
	}

	if compact {
		for _, e := range sec.entries {
			rows = append(rows, stylePink.Bold(true).Render(e.key))
			rows = append(rows, "  "+styleBody.Width(max(10, width-2)).Render(e.val))

			if len(e.tags) > 0 {
				var tags []string
				for _, t := range e.tags {
					tags = append(tags, styleTag.Render("["+t+"]"))
				}
				rows = append(rows, "  "+strings.Join(tags, " "))
			}

			rows = append(rows, "")
		}
		return strings.Join(rows, "\n")
	}

	keyW := 18
	valW := width - keyW - 3
	if valW < 10 {
		valW = 10
	}

	for _, e := range sec.entries {
		key := stylePink.Width(keyW).Render(e.key)
		val := styleBody.Width(valW).Render(e.val)
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, key, "   ", val))

		if len(e.tags) > 0 {
			indent := strings.Repeat(" ", keyW+3)
			var tags []string
			for _, t := range e.tags {
				tags = append(tags, styleTag.Render("["+t+"]"))
			}
			rows = append(rows, indent+strings.Join(tags, " "))
		}

		rows = append(rows, "")
	}

	return strings.Join(rows, "\n")
}

func scalePortrait(lines []string, maxW, maxH int) []string {
	if len(lines) == 0 {
		return nil
	}

	srcW := lipgloss.Width(lines[0])
	if srcW == 0 {
		return lines
	}

	stepX := 1
	for srcW/stepX > maxW {
		stepX++
	}

	stepY := 1
	for len(lines)/stepY > maxH {
		stepY++
	}

	var out []string
	for i := 0; i < len(lines); i += stepY {
		runes := []rune(lines[i])
		var sampled []rune
		for j := 0; j < len(runes); j += stepX {
			sampled = append(sampled, runes[j])
		}
		out = append(out, string(sampled))
	}
	return out
}

func (m model) View() string {
	if m.width < 50 || m.height < 12 {
		msg := stylePink.Render("Terminal too small — please resize.")
		return lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			PaddingTop(1).
			PaddingLeft(1).
			Render(msg)
	}

	showPortrait := m.width >= breakpointMedium && m.height >= 20
	compact := m.width < breakpointSmall

	titleText := "  imren -- ssh.imren.online"
	quitLabel := "[ q to quit ]"
	gap := m.width - lipgloss.Width(titleText) - lipgloss.Width(quitLabel)
	if gap < 1 {
		gap = 1
	}

	title := styleTitleBar.Render(titleText + strings.Repeat(" ", gap) + quitLabel)
	sep := styleBorder.Render(strings.Repeat("─", m.width))

	var nameLines []string
	for _, l := range pickNameArt(m.width) {
		nameLines = append(nameLines, styleName.Render(l))
	}
	name := strings.Join(nameLines, "\n")

	bio := m.typedBio
	if !m.typing {
		bio = m.bioFull
	}

	bioWidth := 50
	if compact {
		bioWidth = clamp(m.width-8, 24, 42)
	}

	right := lipgloss.JoinVertical(
		lipgloss.Left,
		name,
		"",
		styleBody.Width(bioWidth).Render(bio),
		"",
		renderNav(m.navIdx, compact),
	)

	var top string
	if showPortrait {
		portraitMaxW := clamp(m.width/3, 22, 34)
		portraitMaxH := 24
		if m.height < 28 {
			portraitMaxH = 18
		}

		scaled := scalePortrait(portraitLines, portraitMaxW, portraitMaxH)
		portraitHeight := len(scaled)

		portrait := stylePortrait.Render(strings.Join(scaled, "\n"))

		divLines := make([]string, portraitHeight)
		for i := 0; i < portraitHeight; i++ {
			divLines[i] = "│"
		}
		divider := styleBorder.Render(strings.Join(divLines, "\n"))

		// Vertically center 
		rightBox := lipgloss.NewStyle().
			Height(portraitHeight).
			PaddingLeft(1).
			AlignVertical(lipgloss.Center).
			Render(right)

		top = lipgloss.JoinHorizontal(
			lipgloss.Top,
			portrait,
			divider,
			rightBox,
		)
	} else {
		top = right
	}

	panelOuterW := clamp(m.width-4, 32, 120)
	panelInnerW := panelOuterW - 6
	if panelInnerW < 20 {
		panelInnerW = 20
	}

	panel := stylePanel.Width(panelOuterW).Render(
		renderPanel(sections[navOrder[m.navIdx]], panelInnerW, compact),
	)

	hintText := "  tab/h/l: navigate   1/2/3: jump   q: quit"
	if compact {
		hintText = "  tab: navigate   q: quit"
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		sep,
		"",
		lipgloss.NewStyle().PaddingLeft(1).Render(top),
		"",
		lipgloss.NewStyle().PaddingLeft(1).Render(panel),
		"",
		styleHint.Render(hintText),
	)

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Render(content)
}