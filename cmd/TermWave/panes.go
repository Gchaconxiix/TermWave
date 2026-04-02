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
	leftTitle := " Radios "
	rightTitle := " Now Playing "
	leftContent := ""
	rightContent := ""
	stationName := "None"

	if m.viewState == "search" {
		leftTitle = " Search Results "
		if m.err != nil {
			leftContent += fmt.Sprintf("Error fetching stations: %v", m.err)
		} else if len(m.stations) == 0 {
			leftContent += "No Stations Found"
		} else {
			for i, s := range m.stations {
				var currentSelection string
				if m.stationCursor == i && m.focused == "stations" {
					currentSelection = listItemStyle.Background(lipgloss.Color("#585858")).Render(s.Name)
				} else {
					currentSelection = listItemStyle.Render(s.Name)
				}
				leftContent += fmt.Sprintf(" %d. %s %s\n", i+1, currentSelection, s.Saved)
			}
		}
	} else if m.viewState == "saved" { //Saved Stations logic
		leftTitle = " Saved Stations "

		if len(m.savedStations) == 0 {
			leftContent += "No stations saved\nSearch for a station in Stations->Add Station"
		} else {
			itemsPerPage := m.getItemsPerPage()

			startIndex := m.savedPage * itemsPerPage
			endIndex := startIndex + itemsPerPage

			if endIndex > len(m.savedStations) {
				endIndex = len(m.savedStations)
			}

			pageItems := m.savedStations[startIndex:endIndex]

			for i, s := range pageItems {
				var currentSelection string

				if m.stationCursor == i && m.focused == "stations" {
					currentSelection = listItemStyle.Background(lipgloss.Color("#585858")).Render(s.Name)
				} else {
					currentSelection = listItemStyle.Render(s.Name)
				}
				leftContent += fmt.Sprintf(" %d. %s %s\n", (startIndex+i)+1, currentSelection, s.Saved)
			}

			totalPages := (len(m.savedStations) + itemsPerPage - 1) / itemsPerPage
			leftContent += fmt.Sprintf("\n\n  < Page %d of %d >", m.savedPage+1, totalPages)
		}
	}

	if m.currentStation.Name != "" {
		stationName = m.currentStation.Name
	}
	//Now to modify the right pane code

	if !m.isPlaying && m.currentStation.Name == "" {
		rightContent = "Please Select a Station... \nSearch for a new station in Station->Add Station"
	} else {
		rightContent += fmt.Sprintf("Station: %s\nTitle: %s\n\n", stationName, m.currentTitle)
	}
	styledLeftTitle := titleBadgeStyle.Render(leftTitle)
	styledRightTitle := titleBadgeStyle.Render(rightTitle)
	finalLeftContent := fmt.Sprintf("%s\n\n%s", styledLeftTitle, leftContent)
	finalRightContent := fmt.Sprintf("%s\n\n%s", styledRightTitle, rightContent)

	leftPane := paneStyle.Width(leftW).Height(paneH).Render(finalLeftContent)
	rightPane := paneStyle.Width(rightW).Height(paneH).Render(finalRightContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
}
