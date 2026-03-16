package main

import (
	"fmt"
	"os/exec"
	"strings"
	"bufio"
)

var currentPlayer *exec.Cmd

func PlayStream(streamURL string) error {
	StopStream()

	currentPlayer = exec.Command("mpv", "--no-video", streamURL)

	//set up stdout pipe for metadata
	stdout, err := currentPlayer.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Failed to get stdout pipe: %w", err)
	}

	err = currentPlayer.Start()
	if err != nil {
		return fmt.Errorf("MPV failed to start: %w", err)
	}

	//Making a background routine to always check for new lines from mpv
	go func() {
		scanner := bufio.NewScanner(stdout)

		for scanner.Scan() {
			line := scanner.Text()

			if strings.Contains(line, "icy-title:") {
				parts := strings.SplitN(line, "icy-title:", 2)
				if len(parts) == 2 {
					cleanTitle := strings.TrimSpace(parts[1])
					select {
					case titleChannel <- cleanTitle:
					default:
					}
				}
			}
		}
	}()

	return nil
}

func StopStream() {
	if currentPlayer != nil && currentPlayer.Process != nil {
		currentPlayer.Process.Kill()
		currentPlayer.Process.Wait()
		currentPlayer = nil
	}
}

