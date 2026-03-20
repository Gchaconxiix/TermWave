package main

import (
	lipgloss "charm.land/lipgloss/v2"
)

var (
	baseColor string = "240"

	baseBorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(baseColor)).
			Align(lipgloss.Center)

	toolbarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(baseColor)).
			Align(lipgloss.Top)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("0")).
			Padding(0, 1).
			Bold(true)

	menuStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("62")).
			BorderBackground(lipgloss.Color("0")).
			Padding(0, 1)
		//.Align(lipgloss.Center)

	paneStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color(baseColor)).
			Padding(1)

	popup                                         = baseBorderStyle.
		BorderForeground(lipgloss.Color("62")). // Make it pop with a purple border
		Padding(1, 2)

	playButtonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Align(lipgloss.Right)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginLeft(2)
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
