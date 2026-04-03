package main

import (
	"os"

	lipgloss "charm.land/lipgloss/v2"
)

var (
	hasDarkBG       = lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	lightDark       = lipgloss.LightDark(hasDarkBG)
	greenCheck      = lightDark(lipgloss.Color("#43BF6D"), lipgloss.Color("#73F59F"))
	purpleBubbletea = lipgloss.Color("#5F5FD7")
	baseColor       = lipgloss.Color("#585858")
	textColor       = lipgloss.Color("#ffffd7")

	baseBorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(baseColor).
			Align(lipgloss.Center)

	toolbarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(baseColor).
			Align(lipgloss.Top)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFD7")).
			Background(lipgloss.Color("#000000")).
			Padding(0, 1).
			Bold(true)

	menuStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(purpleBubbletea).
			BorderBackground(lipgloss.Color("0")).
			Padding(0, 1)
		//.Align(lipgloss.Center)

	paneStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(baseColor).
			Padding(1)

	popup = baseBorderStyle.
		BorderForeground(purpleBubbletea). // Make it pop with a purple border
		Padding(1, 2)

	playButtonStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Align(lipgloss.Right)

	footerStyle = lipgloss.NewStyle().
			Foreground(baseColor).
			MarginLeft(2)

	checkMark = lipgloss.NewStyle().
			SetString("✓").
			Foreground(greenCheck).
			PaddingRight(1).String()

	titleBadgeStyle = lipgloss.NewStyle().
			Background(baseColor).
			Foreground(textColor).
			Bold(true).
			Padding(0, 1)
	listItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))
	docHeaderStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Background(purpleBubbletea).
			Bold(true).
			Padding(0, 1).
			MarginBottom(1)

	docKeyStyle = lipgloss.NewStyle().
			Foreground(purpleBubbletea).
			Bold(true)
)

//Helpers down here

func updateBorderColor(newColor string) {
	if newColor == "" {
		return
	}
	color := lipgloss.Color(newColor)
	baseBorderStyle = baseBorderStyle.BorderForeground(color)
}

func updateToolbarColor(newColor string) {
	if newColor == "" {
		return
	}
	color := lipgloss.Color(newColor)
	toolbarStyle = toolbarStyle.BorderForeground(color)
}

func updatePanelColor(newColor string) {
	if newColor == "" {
		return
	}
	color := lipgloss.Color(newColor)
	paneStyle = paneStyle.BorderForeground(color)
}
