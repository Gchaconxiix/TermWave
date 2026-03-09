package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type Station struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Tags    string `json:"tags"`
	Country string `json:"country"`
	Image   string `json:"favicon"`
	Saved   string
}

func StationSearch(searchTerm string) ([]Station, error) {
	safeQuery := url.QueryEscape(searchTerm)
	endPoint := fmt.Sprintf("http://de1.api.radio-browser.info/json/stations/search?name=%s&limit=20", safeQuery)

	resp, err := http.Get(endPoint)
	if err != nil {
		return nil, fmt.Errorf("Request failed: %w", err)
	}

	defer resp.Body.Close() //Defer: Run before function is done

	var stations []Station
	err = json.NewDecoder(resp.Body).Decode(&stations)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse JSON: %w", err)
	}
	return stations, nil
}

func DownloadImage(imageUrl string) (string, error) {
	if imageUrl == "" {
		return "", fmt.Errorf("URL empty")
	}

	resp, err := http.Get(imageUrl)
	if err != nil {
		return "", fmt.Errorf("Failed to retrieve image: %w", err)
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
