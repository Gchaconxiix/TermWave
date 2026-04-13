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
	case "manualEntry":
		popupContent = m.renderManualEntry()
	case "about":
		popupContent = fmt.Sprintf("About:\nAuthor: %s\nBug tester: %s\nVersion: TermWave %s", author, contrib, appVer)
	case "theme":
		popupContent = fmt.Sprintf("Theme Settings\n\nBorder Color:  %s\nToolbar Color: %s\nPanel Color:   %s", m.borderColor.View(), m.toolbarColor.View(), m.panelColor.View())
	case "documentation":
		popupContent = m.renderDocPopup()
	case "error":
		popupContent = "Error: Invalid Entry"
	}

	popupWindow := popup.Render(popupContent)
	popupWidth := lipgloss.Width(popupWindow)
	popupHeight := lipgloss.Height(popupWindow)
	x = (m.width / 2) - (popupWidth / 2)
	y = (m.height / 2) - (popupHeight / 2)

	popupLayer := lipgloss.NewLayer(popupWindow).X(x).Y(y).Z(2)
	return popupLayer
}

func docLine(key, desc string) string {
	formattedKey := fmt.Sprintf(" %-8s ", key)
	return docKeyStyle.Render(formattedKey) + " " + lipgloss.NewStyle().Foreground(lipgloss.Color("#8a8a8a")).Render(desc)
}

func (m model) renderDocPopup() string {
	title := docHeaderStyle.Render(" TermWave Documentation ")

	navHeader := lipgloss.NewStyle().Foreground(lipgloss.Color("#bcbcbc")).Underline(true).Render("Navigation")
	navKeys := lipgloss.JoinVertical(lipgloss.Left,
		docLine("Tab", "Switch active window focus"),
		docLine("↑/k", "Move cursor up"),
		docLine("↓/j", "Move cursor down"),
		docLine("Enter", "Select / Open menu"),
		docLine("Esc", "Close current popup / menu"),
	)

	mediaHeader := lipgloss.NewStyle().Foreground(lipgloss.Color("#bcbcbc")).Underline(true).Render("Station Controls")
	mediaKeys := lipgloss.JoinVertical(lipgloss.Left,
		docLine("s", "Save the currently highlighted station"),
		docLine("x", "Delete a saved station"),
		docLine("Space", "Play / Pause stream"),
	)

	themeHeader := lipgloss.NewStyle().Foreground(lipgloss.Color("#bcbcbc")).Underline(true).Render("Theming")
	themeKeys := lipgloss.JoinVertical(lipgloss.Left,
		docLine("Hex", "Enter valid hex codes (e.g., #FF00FF) in the theme menu"),
		docLine("ANSI", "Enter ANSI numbers (0-255) for standard terminal colors"),
	)

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		navHeader,
		navKeys,
		"",
		mediaHeader,
		mediaKeys,
		"",
		themeHeader,
		themeKeys,
		"\nPress [Esc] to close",
	)
	return content
}

func (m model) renderManualEntry() string {
	title := docHeaderStyle.Render("Input a Station")

	inputFields := lipgloss.JoinVertical(lipgloss.Left,
		fmt.Sprintf("%s\n", m.manualName.View()),
		fmt.Sprintf("%s\n", m.manualLink.View()),
		fmt.Sprintf("%s\n", m.manualTags.View()),
		fmt.Sprintf("%s\n", m.manualCountry.View()),
		fmt.Sprintf("%s\n", m.manualImage.View()))

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		inputFields,
		"\nPress [Enter] to Add Station")

	return content
}
