package main

import (
	"fmt"

	lipgloss "charm.land/lipgloss/v2"
)

func (m model) drawToolbar(height int) string {
	var buttons []string

	for i, title := range m.menuTitles {
		activeStyle := buttonStyle

		if i == m.activeMenuIndex && m.focused == "toolbar" {
			if m.menuOpen {
				activeStyle = activeStyle.Background(lipgloss.Color("62"))
			} else {
				activeStyle = activeStyle.Background(lipgloss.Color("240"))
			}
		} else {
			activeStyle = activeStyle.Background(lipgloss.Color("0"))
		}
		buttons = append(buttons, activeStyle.Render(fmt.Sprintf(" %s ", title)))
	}
	toolbarContent := lipgloss.JoinHorizontal(lipgloss.Top, buttons...)

	playbackIcon := "⏹"
	if m.isPlaying && m.currentStation.Name != "" {
		playbackIcon = "▶︎"
	}
	leftWidth := lipgloss.Width(toolbarContent)
	remainingWidth := (m.width - 6) - leftWidth
	playbackButton := playButtonStyle.Width(remainingWidth).Render(playbackIcon)
	finalContent := lipgloss.JoinHorizontal(lipgloss.Top, toolbarContent, playbackButton)

	toolbar := toolbarStyle.Width(m.width - 2).
		Width(m.width - 2).
		Height(height).
		Render(finalContent)

	return toolbar
}
