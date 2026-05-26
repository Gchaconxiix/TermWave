package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

var (
	appName   string = "TermWave"
	appVer    string = "0.6"
	author    string = "Gabriel Chacon"
	contrib   string = "Wolfie574"
	userAgent string = fmt.Sprintf("%s/%s", appName, appVer)
)

func main() {
	p := tea.NewProgram(initialModel())
	tea.RequestCapability("RGB")
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
