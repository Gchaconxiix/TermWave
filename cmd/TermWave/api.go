package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Station struct {
	UUID    string      `json:"stationuuid"`
	Name    string      `json:"name"`
	URL     string      `json:"url"`
	Home    string      `json:"homepage"`
	Tags    string      `json:"tags"`
	Country string      `json:"country"`
	State   string      `json:"state"`
	Codec   string      `json:"codec"`
	Bitrate json.Number `json:"bitrate"`
	Image   string      `json:"favicon"`
	Saved   string
}

func StationSearch(searchTerm string) ([]Station, error) {
	safeQuery := url.QueryEscape(searchTerm)
	endPoint := fmt.Sprintf("http://de1.api.radio-browser.info/json/stations/search?name=%s&limit=25", safeQuery)

	req, err := http.NewRequest("GET", endPoint, nil)
	if err != nil {
		return nil, fmt.Errorf("Request failed: %w", err)
	}
	req.Header.Set("User-Agent", "TermWave/0.1")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request Failed: %w", err)
	}
	defer resp.Body.Close() //Defer: Run before function is done

	var stations []Station
	err = json.NewDecoder(resp.Body).Decode(&stations)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse JSON: %w", err)
	}

	//Limit the size of the Name to 50
	for i, value := range stations {
		runes := []rune(value.Name)
		if len(runes) > 50 {
			truncated := runes[0:49]
			stations[i].Name = string(truncated)
		}
	}
	return stations, nil
}

func DownloadImage(imageUrl string) (string, error) {
	if imageUrl == "" {
		return "", fmt.Errorf("URL empty")
	}

	imageUrl = strings.TrimPrefix(imageUrl, "file://")
	if strings.HasPrefix(imageUrl, "http://") || strings.HasPrefix(imageUrl, "https://") {
		req, err := http.NewRequest("GET", imageUrl, nil)
		if err != nil {
			return "", fmt.Errorf("Failed to retrieve image: %w", err)
		}
		req.Header.Set("User-Agent", "TermWave/0.1")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return "", fmt.Errorf("Request Failed: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("HTTP error: %s", resp.Status)
		}

		tempPath := filepath.Join(os.TempDir(), "termwaveStationArt.jpg")
		file, err := os.Create(tempPath)
		if err != nil {
			return "", fmt.Errorf("Failed to create temp file: %w", err)
		}
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return "", fmt.Errorf("Failed to save Image: %w", err)
		}

		return tempPath, nil
	}

	//If not HTTP, then I will assume it is a local path
	//Check for tilde
	if strings.HasPrefix(imageUrl, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("Failed to find Home Dir: %w", err)
		}
		imageUrl = filepath.Join(homeDir, imageUrl[1:])
	}

	info, err := os.Stat(imageUrl)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("Local Image file does not exist: %s", imageUrl)
	}
	if info.IsDir() {
		return "", fmt.Errorf("Path is a Directory: %s", imageUrl)
	}

	return imageUrl, nil
}

func RegisterStationClick(UUID string) error {
	if UUID == "" {
		return fmt.Errorf("No Station UUID provided")
	}

	endPoint := fmt.Sprintf("http://de1.api.radio-browser.info/json/url/%s", UUID)

	req, err := http.NewRequest("POST", endPoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "TermWave/0.1")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP Error: %s", resp.Status)
	}
	return nil
}
