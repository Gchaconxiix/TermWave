package main

import (
	"fmt"

	lipgloss "charm.land/lipgloss/v2"
)

func (m model) renderPopup() *lipgloss.Layer {
	var x int
	var y int
	var popupContent string
	switch m.focused {
	case "search":
		popupContent = fmt.Sprintf("Search Station by Name:\n\n%s", m.searchInput.View())
	case "about":
		popupContent = "About:\nAuthor: Gabriel Chacon\nBug tester: Wolfie574\nVersion: TermWave 0.2"
	case "theme":
		popupContent = fmt.Sprintf("Theme Settings\n\nBorder Color:  %s\nToolbar Color: %s\nPanel Color:   %s", m.borderColor.View(), m.toolbarColor.View(), m.panelColor.View())
	case "error":
		popupContent = "Error: Invalid Entry"
	}

	popupWindow := popup.Render(popupContent)
	popupWidth := lipgloss.Width(popupWindow)
	popupHeight := lipgloss.Height(popupWindow)
	x = (m.width / 2) - (popupWidth / 2) - 4
	y = (m.height / 2) - (popupHeight / 2)

	popupLayer := lipgloss.NewLayer(popupWindow).X(x).Y(y).Z(2)
	return popupLayer
}
