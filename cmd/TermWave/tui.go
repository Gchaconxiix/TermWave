package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/dolmen-go/kittyimg"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type stationsLoadedMsg []Station
type errMsg struct{ err error }
type imageLoadMsg string
type titleUpdateMsg string
type stationClickMsg string

var titleChannel = make(chan string, 10)

type model struct {
	width           int
	height          int
	activeMenuIndex int
	menuOpen        bool
	menuCursors     []int
	menuTitles      []string
	menuItems       [][]string
	stations        []Station
	savedStations   []Station
	stationCursor   int
	viewState       string
	savedPage       int
	focused         string
	err             error
	currentStation  Station
	currentTitle    string
	searchInput     textinput.Model
	borderColor     textinput.Model
	toolbarColor    textinput.Model
	panelColor      textinput.Model
	manualName      textinput.Model
	manualLink      textinput.Model
	manualHome      textinput.Model
	manualTags      textinput.Model
	manualCountry   textinput.Model
	manualState     textinput.Model
	manualCodec     textinput.Model
	manualImage     textinput.Model
	manualFocusIdx  int
	themeFocusIdx   int
	isPlaying       bool
	currImgData     string
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Search Station by Name..."
	ti.CharLimit = 50
	ti.SetWidth(30)
	bi := textinput.New()
	bi.Placeholder = "Border"
	bi.CharLimit = 7
	bi.SetWidth(10)
	tbi := textinput.New()
	tbi.Placeholder = "Toolbar"
	tbi.CharLimit = 7
	tbi.SetWidth(10)
	pi := textinput.New()
	pi.Placeholder = "Panel"
	pi.CharLimit = 7
	pi.SetWidth(10)

	mn := textinput.New()
	mn.Placeholder = "Name of Station..."
	mn.SetWidth(50)

	ml := textinput.New()
	ml.Placeholder = "Station/Youtube link..."
	ml.SetWidth(50)

	mh := textinput.New()
	mh.Placeholder = "Homepage of Station (Optional)..."
	mh.SetWidth(50)

	mt := textinput.New()
	mt.Placeholder = "Station Tags (Optional)..."
	mt.SetWidth(50)

	mc := textinput.New()
	mc.Placeholder = "Country of Origin (Optional)..."
	mc.SetWidth(50)

	ms := textinput.New()
	ms.Placeholder = "State of Origin (Optional)..."
	ms.SetWidth(50)

	mCodec := textinput.New()
	mCodec.Placeholder = "Codec (Optional)..."
	mCodec.SetWidth(50)

	mi := textinput.New()
	mi.Placeholder = "Link/Path to Image (Optional)..."
	mi.SetWidth(50)

	loadedStations, err := loadStations()
	if err != nil || loadedStations == nil {
		loadedStations = []Station{}
	}

	return model{
		stations:        []Station{},
		savedStations:   loadedStations,
		stationCursor:   0,
		focused:         "stations",
		viewState:       "saved",
		isPlaying:       false,
		savedPage:       0,
		searchInput:     ti,
		borderColor:     bi,
		toolbarColor:    tbi,
		panelColor:      pi,
		manualName:      mn,
		manualLink:      ml,
		manualHome:      mh,
		manualTags:      mt,
		manualCountry:   mc,
		manualState:     ms,
		manualCodec:     mCodec,
		manualImage:     mi,
		manualFocusIdx:  0,
		themeFocusIdx:   0,
		activeMenuIndex: 0,
		menuOpen:        false,
		menuCursors:     []int{0, 0, 0},
		menuTitles:      []string{"Stations", "Settings", "Help"},
		menuItems: [][]string{
			{"Search Station", "Input Station", "Saved Stations", "Quit"},
			{"Theme Settings", "Station Info"},
			{"About", "Documentation", "License"},
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case stationsLoadedMsg:
		for i, newStation := range msg {
			for _, savedStation := range m.savedStations {
				if newStation.URL == savedStation.URL {
					msg[i].Saved = checkMark
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

		// Margin safety check for dynamic item lists
		if m.viewState == "saved" {
			itemsPerPage := m.getItemsPerPage()
			if len(m.savedStations) > 0 {
				totalPages := (len(m.savedStations) + itemsPerPage - 1) / itemsPerPage
				if m.savedPage >= totalPages {
					m.savedPage = totalPages - 1
				}
			} else {
				m.savedPage = 0
			}

			startIndex := m.savedPage * itemsPerPage
			itemsOnPage := len(m.savedStations) - startIndex

			if itemsOnPage > itemsPerPage {
				itemsOnPage = itemsPerPage
			} else if itemsOnPage < 0 {
				itemsOnPage = 0
			}
			if m.stationCursor >= itemsOnPage && itemsOnPage > 0 {
				m.stationCursor = itemsOnPage - 1
			}
		}
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
			var currentSelection string
			if currentCursor == i {
				currentSelection = listItemStyle.Background(lipgloss.Color("#585858")).Render(item)
			} else {
				currentSelection = listItemStyle.Render(item)
			}
			menuText += fmt.Sprintf(" %s\n", currentSelection)
		}

		menu := menuStyle.Background(lipgloss.Color("0")).Render(menuText)

		xOffset := 2 + m.activeMenuIndex*14
		menuLayer := lipgloss.NewLayer(menu).X(xOffset).Y(3).Z(1)

		layers = append(layers, menuLayer)
	}

	if m.focused != "stations" && m.focused != "toolbar" {
		layers = append(layers, m.renderPopup())
	}

	compositor := lipgloss.NewCompositor(layers...)
	finalUI := compositor.Render()

	// Image rendered after BubbleTea compositing
	if m.currImgData != "" && (m.focused == "stations" || m.focused == "toolbar") {
		availW := m.width - 6
		leftW := int(float64(availW) * 0.6)
		col := leftW + 6
		row := 12

		cursorMove := fmt.Sprintf("\x1b[%d;%dH", row, col)
		hideCursor := "\x1b[?25l"

		imgData := m.currImgData
		if isKittyTerminal() {
			imgData = injectKittyPixelSize(m.currImgData, 320)
		}

		finalUI += cursorMove + imgData + hideCursor
	}

	v := tea.NewView(finalUI)
	v.AltScreen = true
	return v
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

func (m model) showImage() tea.Cmd {
	if isKittyTerminal() {
		return m.showImageKitty()
	}
	return func() tea.Msg {
		imageUrl, err := DownloadImage(m.currentStation.Image)
		if err != nil {
			return nil
		}
		cmd := exec.Command("chafa", "--scale", "1", "--size", "40x20", imageUrl)
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

func (m model) getItemsPerPage() int {
	items := m.height - 14
	if items < 1 {
		return 1
	}
	return items
}

func isKittyTerminal() bool {
	return os.Getenv("KITTY_WINDOW_ID") != "" || os.Getenv("TERM") == "xterm-kitty"
}

func injectKittyPixelSize(kittyData string, pixelWidth int) string {
	cellWidth := 8
	cols := max(pixelWidth / cellWidth, 1)

	sizeParam := fmt.Sprintf(",c=%d", cols)
	
	insertMarker := ",t=d"
	if idx := bytes.Index([]byte(kittyData), []byte(insertMarker)); idx != -1 {
		return kittyData[:idx] + sizeParam + kittyData[idx:]
	}
	if idx := bytes.IndexByte([]byte(kittyData), ';'); idx != -1 {
		return kittyData[:idx] + sizeParam + kittyData[idx:]
	}
	return kittyData
}

func (m model) showImageKitty() tea.Cmd {
	return func() tea.Msg {
		imagePath, err := DownloadImage(m.currentStation.Image)
		if err != nil {
			return nil
		}

		file, err := os.Open(imagePath)
		if err != nil {
			return nil
		}
		defer file.Close()

		var buf bytes.Buffer
		err = kittyimg.Transcode(&buf, file)
		if err != nil {
			return nil
		}

		return imageLoadMsg(buf.String())
	}
}
