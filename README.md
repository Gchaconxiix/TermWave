# TermWave
### A Terminal-based internet radio player for Linux

Termwave is a TUI that searches for internet radio stations and plays them back. It is my first program made in Go, and the first using Bubbletea, so it is still in early-stage development. 

⚠️ Termwave is currently still being developed. While I will be actively adding/fixing features, things may break and some things may be removed until I can fix them.

Program was made using Go and the BubbleTea V2 framework. ⚠️ MPV is required for audio playback! Chafa is required for displaying station images.

### Currently Working/Recent Changes:
    * Search for internet radio stations (Using api.radio-browser)
    * Radio Playback using mpv
    * Functionality to Save and Remove Stations has been added
    * Icy-Titles work!
    * Added artwork via Sixels through Chafa
    * Fixed Menus. Each Menu holds 25 Stations
### Roadmap/Things to fix:
    * Make the UI look better
    * Implement Kitty images instead of Sixels, and center images better
    * Implement the rest of the menu options before an actual release
Will be making releases hopefully in the near-future

### Usage:

```
$ git clone https://github.com/MeatballSteakTips/TermWave.git
$ cd TermWave
$ go build ./...
$ ./TermWave
```
### Controls: 
    * Use Tab to switch focus between Toolbar and Stations Pane.
    * S saves a Station
    * X removes a station
    * Enter for everything else

### Requirements:

Be sure you have MPV and chafa installed!

Debian/Ubuntu based:
```
$ sudo apt install mpv
$ sudo apt install chafa
```

Arch
```
$ sudo pacman -S mpv
$ sudo pacman -S chafa
```
Note: Artwork will NOT appear on Kitty terminal due to Sixels being unsupported

Thank you for trying it out! I hope to make it better in the near future!

