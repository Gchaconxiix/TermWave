package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

var (
	appVer  string = "0.4"
	author  string = "Gabriel Chacon"
	contrib string = "Wolfie574"
)

func main() {
	p := tea.NewProgram(initialModel())
	tea.RequestCapability("RGB")
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
