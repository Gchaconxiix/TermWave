package main

import (
	"fmt"

	lipgloss "charm.land/lipgloss/v2"
)

func (m model) drawPanes() string {
	availW := m.width - 6
	leftW := int(float64(availW) * 0.6)
	rightW := availW - leftW
	paneH := m.height - 5
	leftTitle := "Radios"
	leftContent := fmt.Sprintf("%s\n\n", leftTitle)
	rightContent := "Now Playing\n\n"
	stationName := "None"

	if m.viewState == "search" {
		leftContent = "Stations (Search Results)\n\n"
		if m.err != nil {
			leftContent += fmt.Sprintf("Error fetching stations: %v", m.err)
		} else if len(m.stations) == 0 {
			leftContent += "No Stations Found"
		} else {
			for i, s := range m.stations {
				cursor := "  "
				if m.stationCursor == i && m.focused == "stations" {
					cursor = "> "
				}
				leftContent += fmt.Sprintf("%s%d. %s %s\n", cursor, i+1, s.Name, s.Saved)
			}
		}
	} else if m.viewState == "saved" { //Saved Stations logic
		leftContent = "Stations\n\n"

		if len(m.savedStations) == 0 {
			leftContent += "No stations saved\nSearch for a station in Stations->Add Station"
		} else {
			itemsPerPage := 25

			startIndex := m.savedPage * itemsPerPage
			endIndex := startIndex + itemsPerPage

			if endIndex > len(m.savedStations) {
				endIndex = len(m.savedStations)
			}

			pageItems := m.savedStations[startIndex:endIndex]

			for i, s := range pageItems {
				cursor := "  "

				if m.stationCursor == i && m.focused == "stations" {
					cursor = "> "
				}
				leftContent += fmt.Sprintf("%s%d. %s %s\n", cursor, (startIndex+i)+1, s.Name, s.Saved)
			}

			totalPages := (len(m.savedStations) + itemsPerPage - 1) / itemsPerPage
			leftContent += fmt.Sprintf("\n\n  --- Page %d of %d ---", m.savedPage+1, totalPages)
		}
	}

	if m.currentStation.Name != "" {
		stationName = m.currentStation.Name
	}
	//Now to modify the right pane code

	rightContent += fmt.Sprintf("Station: %s\nTitle: %s\n\n", stationName, m.currentTitle)
	if !m.isPlaying && m.currentStation.Name == "" {
		rightContent = "Please Select a Station... \nSearch for a new station in Station->Add Station"
	}

	leftPane := paneStyle.Width(leftW).Height(paneH).Render(leftContent)
	rightPane := paneStyle.Width(rightW).Height(paneH).Render(rightContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
}
