package main

import (
	"fmt"
	"os/exec"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

type stationsLoadedMsg []Station
type errMsg struct{ err error }
type imageLoadMsg string
type titleUpdateMsg string
type stationClickMsg string

var titleChannel = make(chan string, 10) //buffer size

type model struct {
	width           int
	height          int
	activeMenuIndex int
	menuOpen        bool
	menuCursors     []int
	menuTitles      []string
	menuItems       [][]string
	stations        []Station
	savedStations   []Station //Page number, each will hold 8
	stationCursor   int
	viewState       string //Tells me if I am in search or saved sations mode
	savedPage       int    //Page #s
	focused         string //The window that the cursor should be on
	err             error
	currentStation  Station
	currentTitle    string
	searchInput     textinput.Model
	isPlaying       bool
	currImgData     string
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Search Station by Name..."
	ti.CharLimit = 50
	ti.SetWidth(30)

	//Loading saved stations
	loadedStations, err := loadStations()
	if err != nil || loadedStations == nil {
		loadedStations = []Station{}
	}

	return model{
		stations:        []Station{},
		savedStations:   loadedStations,
		stationCursor:   0,
		focused:         "stations", //This tells me which part is in focus
		viewState:       "saved",    //For the left pane
		isPlaying:       false,
		savedPage:       0,
		searchInput:     ti,
		activeMenuIndex: 0,
		menuOpen:        false,
		menuCursors:     []int{0, 0, 0},
		menuTitles:      []string{"Stations", "Settings", "Help"},
		menuItems: [][]string{
			{"Add Station", "Saved Stations", "Quit"},
			{"Audio Settings", "Theme Settings", "Preferences"},
			{"About", "Documentation", "License"},
		},
	}
}

func fetchStations(query string) tea.Cmd {
	return func() tea.Msg {
		stations, err := StationSearch(query)
		if err != nil {
			return errMsg{err}
		}
		return stationsLoadedMsg(stations)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) showImage() tea.Cmd {
	return func() tea.Msg {
		imageUrl, err := DownloadImage(m.currentStation.Image)
		if err != nil {
			return nil
		}
		cmd := exec.Command("chafa", "-f", "sixels", "--scale", "1", "--size", "40x20", imageUrl)
		currentImage, err := cmd.CombinedOutput()
		if err != nil {
			errorMsg := fmt.Sprintf("Chafa failed: %v\nReason: %s\nLink: %s", err, string(currentImage), m.currentStation.Image)
			return imageLoadMsg(errorMsg)
		}
		return imageLoadMsg(string(currentImage))
	}
}

func waitForTitle(sub chan string) tea.Cmd {
	return func() tea.Msg {
		newTitle := <-sub
		return titleUpdateMsg(newTitle)
	}
}

func sendClickBack(UUID string) tea.Cmd {
	return func() tea.Msg {
		err := RegisterStationClick(UUID)
		if err != nil {
			return nil
		}
		return nil
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case stationsLoadedMsg:
		//Making a loop so I can check if the station is already saved
		for i, newStation := range msg {
			for _, savedStation := range m.savedStations {
				if newStation.URL == savedStation.URL {
					msg[i].Saved = "*"
					break
				}
			}
		}
		m.stations = msg
		m.viewState = "search"
		m.stationCursor = 0
		return m, nil

	case errMsg:
		m.err = msg.err
		return m, nil

	case imageLoadMsg:
		m.currImgData = string(msg)
		return m, nil

	case titleUpdateMsg:
		m.currentTitle = string(msg)
		return m, waitForTitle(titleChannel)

	case stationClickMsg:
		return m, nil

	case tea.KeyPressMsg:
		currModel, cmd := m.keyHandler(msg)
		return currModel, cmd

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("Initializing")
	}
	toolbarHeight := 1
	contentPanes := m.drawPanes()
	toolbar := m.drawToolbar(toolbarHeight)

	border := baseBorderStyle.
		Width(m.width - 2).
		Height(m.height - toolbarHeight - 4).
		Render(contentPanes)

	ui := lipgloss.JoinVertical(lipgloss.Left, toolbar, border)
	bgLayer := lipgloss.NewLayer(ui).X(0).Y(0).Z(0)
	layers := []*lipgloss.Layer{bgLayer}

	if m.menuOpen {
		currentItems := m.menuItems[m.activeMenuIndex]
		currentCursor := m.menuCursors[m.activeMenuIndex]

		menuText := fmt.Sprintf("%s\n\n", m.menuTitles[m.activeMenuIndex])
		for i, item := range currentItems {
			cursor := " "
			if currentCursor == i {
				cursor = ">"
			}
			menuText += fmt.Sprintf(" %s %s\n", cursor, item)
		}

		menu := menuStyle.Background(lipgloss.Color("0")).Render(menuText)

		xOffset := 2 + m.activeMenuIndex*14
		menuLayer := lipgloss.NewLayer(menu).X(xOffset).Y(3).Z(1)

		layers = append(layers, menuLayer)
	}

	switch m.focused {
	case "search":
		popupContent := fmt.Sprintf("Search Station by Name:\n\n%s", m.searchInput.View())

		searchPopup := popup.Render(popupContent)
		popupWidth := lipgloss.Width(searchPopup)
		popupHeight := lipgloss.Height(searchPopup)
		x := (m.width / 2) - (popupWidth / 2) - 4
		y := (m.height / 2) - (popupHeight / 2)

		searchLayer := lipgloss.NewLayer(searchPopup).X(x).Y(y).Z(2)
		layers = append(layers, searchLayer)
	case "about":
		popupContent := fmt.Sprintf("About:\nAuthor: Gabriel Chacon\nBug tester: Wolfie574\nVersion: TermWave 0.2")

		aboutPopup := popup.Render(popupContent)
		popupWidth := lipgloss.Width(aboutPopup)
		popupHeight := lipgloss.Height(aboutPopup)
		x := (m.width / 2) - (popupWidth / 2) - 4
		y := (m.height / 2) - (popupHeight / 2)

		aboutLayer := lipgloss.NewLayer(aboutPopup).X(x).Y(y).Z(2)
		layers = append(layers, aboutLayer)

	}
	compositor := lipgloss.NewCompositor(layers...)
	finalUI := compositor.Render()

	//Image needs to be drawn AFTER the rendering for BubbleTea
	if m.currImgData != "" {
		availW := m.width - 6
		leftW := int(float64(availW) * 0.6)
		col := leftW + 6
		row := 12

		cursorMove := fmt.Sprintf("\x1b[%d;%dH", row, col) //Have to put the exact escape sequence in the right place
		hideCursor := "\x1b[?25l"

		finalUI += cursorMove + m.currImgData + hideCursor
	}

	v := tea.NewView(finalUI)
	v.AltScreen = true
	return v
}
