package tui

import "github.com/charmbracelet/lipgloss"

var (
	accent   = lipgloss.AdaptiveColor{Light: "#d97757", Dark: "#d97757"}
	muted    = lipgloss.AdaptiveColor{Light: "#555049", Dark: "#c0c4cc"}
	faint    = lipgloss.AdaptiveColor{Light: "#82796f", Dark: "#858b96"}
	success  = lipgloss.AdaptiveColor{Light: "#5d9b66", Dark: "#74b87a"}
	warn     = lipgloss.AdaptiveColor{Light: "#b68120", Dark: "#d9a441"}
	errColor = lipgloss.AdaptiveColor{Light: "#b94b4d", Dark: "#e0696a"}

	appStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(faint)

	titleStyle = lipgloss.NewStyle().
			Foreground(accent).
			Bold(true).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(faint).
			Padding(0, 1)

	// Input box: only top + bottom borders, no sides (reasonix style).
	inputBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true, false, true, false).
			BorderForeground(accent).
			PaddingLeft(1)

	statusBlockStyle = lipgloss.NewStyle().
				Foreground(faint)

	workingStyle = lipgloss.NewStyle().
			Foreground(faint).
			Italic(true)

	transcriptAreaStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(faint).
				Padding(0, 1)

	loadingStyle = lipgloss.NewStyle().
			Foreground(warn).
			Italic(true)
)
