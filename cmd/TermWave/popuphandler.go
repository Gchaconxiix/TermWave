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
		popupContent = fmt.Sprintf("About:\nAuthor: %s\nBug tester: %s\nVersion: %s %s", author, contrib, appName, appVer)
	case "theme":
		popupContent = m.renderThemePopup()
	case "info":
		popupContent = m.renderStationInfo()
	case "documentation":
		popupContent = m.renderDocPopup()
	case "license":
		popupContent = "Licensed under MIT\nCopyright (c) 2026 Gabriel Chacon\nSee LICENSE file for Terms of Use"
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
		docLine("+/-", "Move Station Up/Down List"),
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

func (m model) renderThemePopup() string {
	title := docHeaderStyle.Render("Theme Settings")

	instruction := "Use Hex Values (#FFFFFF) or ANSI\n"

	inputFields := lipgloss.JoinVertical(lipgloss.Left,
		fmt.Sprintf("%s\n", m.borderColor.View()),
		fmt.Sprintf("%s\n", m.toolbarColor.View()),
		fmt.Sprintf("%s\n", m.panelColor.View()))

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		instruction,
		inputFields,
		"\n[Enter] to set, [Esc] to cancel")

	return content
}

func (m model) renderManualEntry() string {
	title := docHeaderStyle.Render("Input a Station")

	inputFields := lipgloss.JoinVertical(lipgloss.Left,
		fmt.Sprintf("%s\n", m.manualName.View()),
		fmt.Sprintf("%s\n", m.manualLink.View()),
		fmt.Sprintf("%s\n", m.manualHome.View()),
		fmt.Sprintf("%s\n", m.manualImage.View()),
		fmt.Sprintf("%s\n", m.manualTags.View()),
		fmt.Sprintf("%s\n", m.manualCountry.View()),
		fmt.Sprintf("%s\n", m.manualState.View()),
		fmt.Sprintf("%s\n", m.manualCodec.View()))

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		inputFields,
		"\nPress [Enter] to Add Station")

	return content
}

func (m model) renderStationInfo() string {
	title := docHeaderStyle.Render("Station Info")
	temp := m.currentStation
	bitrate, err := temp.Bitrate.Int64()
	if err != nil {
		bitrate = 0
	}

	infoLines := lipgloss.JoinVertical(lipgloss.Left,
		fmt.Sprintf("Name:     %s\n", temp.Name),
		fmt.Sprintf("Tags:     %s\n", temp.Tags),
		fmt.Sprintf("Homepage: %s\n", temp.Home),
		fmt.Sprintf("URL:      %s\n", temp.URL),
		fmt.Sprintf("Image:    %s\n", temp.Image),
		fmt.Sprintf("Country:  %s\n", temp.Country),
		fmt.Sprintf("State:    %s\n", temp.State),
		fmt.Sprintf("Codec:    %s %dk\n", temp.Codec, bitrate),
		fmt.Sprintf("UUID:     %s\n", temp.UUID))

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		infoLines,
		"\nPress [Esc] to Exit")

	return content
}
