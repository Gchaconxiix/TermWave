package main

import (
	tea "charm.land/bubbletea/v2"
)

func (m model) keyHandler(msg tea.Msg) (tea.Model, tea.Cmd) {
	s := ""
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		s = keyMsg.String()
	}
	switch m.focused { //For the popup windows
	case "search":
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

	case "manualEntry":
		switch s {
		case "esc":
			m.focused = "stations"
			m.manualName.Blur()
			m.manualLink.Blur()
			m.manualTags.Blur()
			m.manualCountry.Blur()
			m.manualImage.Blur()
			return m, nil
		case "enter":
			if m.manualName.Value() == "" || m.manualLink.Value() == "" {
				m.focused = "stations"
				return m, nil
			}
			manualStation := Station{
				UUID:    "",
				Name:    m.manualName.Value(),
				URL:     m.manualLink.Value(),
				Tags:    m.manualTags.Value(),
				Country: m.manualCountry.Value(),
				Image:   m.manualImage.Value(),
				Saved:   "",
			}
			m.savedStations = append(m.savedStations, manualStation)
			_ = saveStations(m.savedStations)

			m.manualName.SetValue("")
			m.manualLink.SetValue("")
			m.manualTags.SetValue("")
			m.manualCountry.SetValue("")
			m.manualImage.SetValue("")

			m.focused = "stations"
			return m, nil
		case "up", "shift+tab":
			m.manualFocusIdx--
			if m.manualFocusIdx < 0 {
				m.manualFocusIdx = 4
			}
		case "down", "tab":
			m.manualFocusIdx++
			if m.manualFocusIdx > 4 {
				m.manualFocusIdx = 0
			}
		}
		m.manualName.Blur()
		m.manualLink.Blur()
		m.manualTags.Blur()
		m.manualCountry.Blur()
		m.manualImage.Blur()
		switch m.manualFocusIdx {
		case 0:
			m.manualName.Focus()
		case 1:
			m.manualLink.Focus()
		case 2:
			m.manualTags.Focus()
		case 3:
			m.manualCountry.Focus()
		case 4:
			m.manualImage.Focus()
		}
		var cmds []tea.Cmd
		var cmd tea.Cmd
		m.manualName, cmd = m.manualName.Update(msg)
		cmds = append(cmds, cmd)
		m.manualLink, cmd = m.manualLink.Update(msg)
		cmds = append(cmds, cmd)
		m.manualTags, cmd = m.manualTags.Update(msg)
		cmds = append(cmds, cmd)
		m.manualCountry, cmd = m.manualCountry.Update(msg)
		cmds = append(cmds, cmd)
		m.manualImage, cmd = m.manualImage.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	case "about", "documentation", "license":
		if s == "enter" || s == "esc" || s == "backspace" {
			m.focused = "stations"
			return m, nil
		}
	case "theme":
		switch s {
		case "esc":
			m.focused = "stations"
			m.borderColor.Blur()
			m.toolbarColor.Blur()
			m.panelColor.Blur()
			return m, nil
		case "up", "shift+tab":
			m.themeFocusIdx--
			if m.themeFocusIdx < 0 {
				m.themeFocusIdx = 2
			}
		case "down", "tab":
			m.themeFocusIdx++
			if m.themeFocusIdx > 2 {
				m.themeFocusIdx = 0
			}
		case "enter":
			updateBorderColor(m.borderColor.Value())
			updateToolbarColor(m.toolbarColor.Value())
			updatePanelColor(m.panelColor.Value())
			m.focused = "stations"
			m.borderColor.Blur()
			m.toolbarColor.Blur()
			m.panelColor.Blur()
			return m, nil
		}
		m.borderColor.Blur()
		m.toolbarColor.Blur()
		m.panelColor.Blur()
		switch m.themeFocusIdx {
		case 0:
			m.borderColor.Focus()
		case 1:
			m.toolbarColor.Focus()
		case 2:
			m.panelColor.Focus()
		}
		var cmds []tea.Cmd
		var cmd tea.Cmd
		m.borderColor, cmd = m.borderColor.Update(msg)
		cmds = append(cmds, cmd)
		m.toolbarColor, cmd = m.toolbarColor.Update(msg)
		cmds = append(cmds, cmd)
		m.panelColor, cmd = m.panelColor.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}
	//Keeping my "Supposed to always work" keys here
	if s == "ctrl+c" {
		StopStream()
		_ = saveStations(m.savedStations)
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
		case "esc", "q":
			m.menuOpen = false
		case "right":
			m.activeMenuIndex = (m.activeMenuIndex + 1) % len(m.menuTitles)
		case "left":
			m.activeMenuIndex = ((m.activeMenuIndex - 1) + len(m.menuTitles)) % len(m.menuTitles)
		case "enter":
			selectedItem := m.menuItems[m.activeMenuIndex][m.menuCursors[m.activeMenuIndex]]

			switch selectedItem {
			case "Search Station":
				m.focused = "search"
				m.searchInput.Focus()
				m.searchInput.SetValue("")
				m.menuOpen = false

			case "Input Station":
				m.focused = "manualEntry"
				m.manualName.Focus()
				m.manualName.SetValue("")
				m.menuOpen = false

			case "Saved Stations":
				m.focused = "stations"
				m.viewState = "saved"
				m.menuOpen = false
				m.stationCursor = 0

			case "Theme Settings":
				m.focused = "theme"
				m.borderColor.Focus()
				m.menuOpen = false

			case "About":
				m.focused = "about"
				m.menuOpen = false

			case "Documentation":
				m.focused = "documentation"
				m.menuOpen = false
			case "License":
				m.focused = "license"
				m.menuOpen = false

			case "Quit":
				StopStream()
				_ = saveStations(m.savedStations)
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
				_ = saveStations(m.savedStations)
				return m, tea.Quit
			}

		} else if m.focused == "stations" {
			//Radio list focused
			var currentListLen int
			if m.viewState == "search" {
				currentListLen = len(m.stations)
			} else {
				itemsPerPage := m.getItemsPerPage()
				startIndex := m.savedPage * itemsPerPage
				endIndex := startIndex + itemsPerPage
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
				itemsPerPage := m.getItemsPerPage()
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
				if m.viewState == "search" && len(m.stations) > 0 && m.stations[m.stationCursor].Saved != checkMark {
					m.savedStations = append(m.savedStations, m.stations[m.stationCursor])
					m.stations[m.stationCursor].Saved = checkMark
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
			case "+":
				if m.viewState == "saved" {
					itemsPerPage := m.getItemsPerPage()
					startIndex := m.savedPage * itemsPerPage
					actualIndex := startIndex + m.stationCursor

					if actualIndex > 0 {
						m.savedStations[actualIndex-1], m.savedStations[actualIndex] = m.savedStations[actualIndex], m.savedStations[actualIndex-1]
						m.stationCursor--

						if m.stationCursor < 0 && m.savedPage > 0 {
							m.savedPage--
							m.stationCursor = 24
						}
					}
				}
			case "-":
				if m.viewState == "saved" {
					itemsPerPage := m.getItemsPerPage()
					startIndex := m.savedPage * itemsPerPage
					actualIndex := startIndex + m.stationCursor

					if actualIndex < len(m.savedStations)-1 {
						m.savedStations[actualIndex+1], m.savedStations[actualIndex] = m.savedStations[actualIndex], m.savedStations[actualIndex+1]
						m.stationCursor++

						if m.stationCursor > 24 {
							m.savedPage++
							m.stationCursor = 0
						}
					}
				}
			case "enter":
				m.isPlaying = true
				if m.viewState == "search" && len(m.stations) > 0 {
					m.currentStation = m.stations[m.stationCursor]
					_ = PlayStream(m.currentStation.URL)
				} else if m.viewState == "saved" && currentListLen > 0 {
					itemsPerPage := m.getItemsPerPage()
					startIndex := m.savedPage * itemsPerPage
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
					_ = saveStations(m.savedStations)
					return m, tea.Quit
				}
			}
		}
	}

	return m, tea.RequestWindowSize
}
