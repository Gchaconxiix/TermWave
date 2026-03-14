package main

import tea "charm.land/bubbletea/v2"

func (m model) keyHandler(msg string) (tea.Model, tea.Cmd) {
	s := msg
	if m.focused == "search" { //This is for the search window. Note to me: I should make this less confusing later
		switch s {
		case "esc":
			m.focused = "stations"
			m.searchInput.Blur()
			return m, nil
		case "enter":
			query := m.searchInput.Value()
			m.focused = "stations"
			m.searchInput.Blur()
			return m, fetchStations(query) //Search for the list

		}
		//This is to set up the field for typing. Otherwise it will do weird stuff
		var cmd tea.Cmd
		m.searchInput, cmd = m.searchInput.Update(msg)
		return m, cmd

	} else if m.focused == "about" {
		if s == "enter" || s == "esc" || s == "backspace" {
			m.focused = "stations"
			return m, nil
		}
	}
	//Keeping my "Supposed to always work" keys here
	if s == "ctrl+c" {
		StopStream()
		return m, tea.Quit
	}
	if s == "space" && m.isPlaying {
		StopStream()
		m.isPlaying = false
	} else if s == "space" && !m.isPlaying && m.currentStation.Name != "" {
		_ = PlayStream(m.currentStation.URL)
		m.isPlaying = true
		return m, waitForTitle(titleChannel)
	}

	if s == "tab" && !m.menuOpen {
		if m.focused == "stations" {
			m.focused = "toolbar"
		} else {
			m.focused = "stations"
		}
		return m, nil
	}
	if m.menuOpen {
		switch s {
		case "up", "k":
			if m.menuCursors[m.activeMenuIndex] > 0 {
				m.menuCursors[m.activeMenuIndex]--
			} else if m.menuCursors[m.activeMenuIndex] == 0 {
				m.menuOpen = false
			}
		case "down", "j":
			maxIndex := len(m.menuItems[m.activeMenuIndex]) - 1
			if m.menuCursors[m.activeMenuIndex] < maxIndex {
				m.menuCursors[m.activeMenuIndex]++
			}
		case "esc", "q", "left", "right":
			m.menuOpen = false
		case "enter":
			selectedItem := m.menuItems[m.activeMenuIndex][m.menuCursors[m.activeMenuIndex]]

			switch selectedItem {
			case "Add Station":
				m.focused = "search"
				m.searchInput.Focus()
				m.searchInput.SetValue("")
				m.menuOpen = false

			case "Saved Stations":
				m.focused = "stations"
				m.viewState = "saved"
				m.menuOpen = false
				m.stationCursor = 0

			case "About":
				m.focused = "About"
				m.menuOpen = false

			case "Quit":
				StopStream()
				return m, tea.Quit

			default:
				m.menuOpen = false
			}

		}

	} else {
		// Menu is closed
		if m.focused == "toolbar" {
			switch s {
			case "right", "l":
				m.activeMenuIndex = (m.activeMenuIndex + 1) % len(m.menuTitles)
			case "left", "h":
				m.activeMenuIndex = ((m.activeMenuIndex - 1) + len(m.menuTitles)) % len(m.menuTitles)
			case "enter", "down", "s":
				m.menuOpen = true
			case "q", "esc":
				StopStream()
				return m, tea.Quit
			}

		} else if m.focused == "stations" {
			//Radio list focused
			var currentListLen int
			if m.viewState == "search" {
				currentListLen = len(m.stations)
			} else {
				startIndex := m.savedPage * 25
				endIndex := startIndex + 25
				if endIndex > len(m.savedStations) {
					endIndex = len(m.savedStations)
				}
				currentListLen = endIndex - startIndex
			}
			switch s {
			case "up", "k":
				if m.stationCursor > 0 {
					m.stationCursor--
				}
			case "down", "j":
				if m.stationCursor < currentListLen-1 {
					m.stationCursor++
				}
			case "right":
				itemsPerPage := 25
				totalPages := (len(m.savedStations) + itemsPerPage - 1) / itemsPerPage
				if m.savedPage != (totalPages - 1) {
					m.savedPage += 1
				}
				startIndex := m.savedPage * itemsPerPage
				itemsOnThisPage := len(m.savedStations) - startIndex
				if m.stationCursor >= itemsOnThisPage {
					m.stationCursor = itemsOnThisPage - 1
				}
			case "left":
				if m.savedPage != 0 {
					m.savedPage -= 1
				}
			case "s":
				if m.viewState == "search" && len(m.stations) > 0 && m.stations[m.stationCursor].Saved != "*" {
					m.savedStations = append(m.savedStations, m.stations[m.stationCursor])
					m.stations[m.stationCursor].Saved = "*"
					_ = saveStations(m.savedStations)
				}
			case "x", "delete":
				if m.viewState == "saved" && len(m.savedStations) > 0 {
					startIndex := m.savedPage * 25
					actualIndex := startIndex + m.stationCursor

					m.savedStations = append(m.savedStations[:actualIndex], m.savedStations[actualIndex+1:]...)
					_ = saveStations(m.savedStations)

					itemsLeftOnPage := len(m.savedStations) - startIndex
					if m.stationCursor > 0 && m.stationCursor >= itemsLeftOnPage {
						m.stationCursor--
					}

					if startIndex >= len(m.savedStations) && m.savedPage > 0 {
						m.savedPage--
						m.stationCursor = 7
					}
				}
			case "enter":
				m.isPlaying = true
				if m.viewState == "search" && len(m.stations) > 0 {
					m.currentStation = m.stations[m.stationCursor]
					_ = PlayStream(m.currentStation.URL)
				} else if m.viewState == "saved" && currentListLen > 0 {
					startIndex := m.savedPage * 25
					m.currentStation = m.savedStations[startIndex+m.stationCursor]
					_ = PlayStream(m.currentStation.URL)
				}
				m.currentTitle = ""
				m.currImgData = "" //Reminding myself that the old photo needs to be cleared
				return m, tea.Batch(m.showImage(), waitForTitle(titleChannel), sendClickBack(m.currentStation.UUID))

			case "backspace":
				if m.viewState == "search" {
					m.stationCursor = 0
					m.viewState = "saved"
				}

			case "q", "esc":
				if m.viewState == "search" {
					m.stationCursor = 0
					m.viewState = "saved"
				} else {
					StopStream()
					return m, tea.Quit
				}
			}
		}
	}

	return m, tea.RequestWindowSize
}

func (m model) keyReleaseHandler(msg string) (tea.Model, tea.Cmd) {
	s := msg
	switch s {

	}
	return m, tea.RequestWindowSize
}
