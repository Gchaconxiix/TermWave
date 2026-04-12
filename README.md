# TermWave
### A Terminal-based internet radio player for Linux

Termwave is a TUI that searches for internet radio stations and plays them back. It is my first program made in Go, and the first using Bubbletea, so it is still in early-stage development. 
Program was made using Go and the BubbleTea V2 framework. It is essentially a frontend for MPV and api.radio-browser.info 
I wanted something that I can run on both my PC and on a CyberDeck for traveling/camping.

⚠️ Termwave is currently still being developed. While I will be actively adding/fixing features, things may break and some things may be removed until I can fix them.
⚠️ MPV is required for audio playback! Chafa is required for displaying station images.

### Currently Working/Recent Changes:
    * Search for internet radio stations (Using api.radio-browser)
    * Radio Playback using mpv
    * Functionality to Save and Remove Stations has been added
    * Icy-Titles work!
    * Added artwork through Chafa
    * Added Documentation Menu
    * Started Theme Settings Menu
    * Added station sorting
    * Can now change color of each outline in Theme Settings (ANSI or Hex)
    * Added Dynamic List functionaliy. Stations per page will now change based on the current window size
    * Changed Chafa to automatically select best format based on terminal (Better compatibility)
    * Added Manual Station addition (Also works with Youtube links with yt-dlp)
### Roadmap/Things to fix:
    * Make the UI look better
    * Kitty images still do not work properly. (If you have any ideas please let me know!)
    * Implement the rest of the menu options before an actual release
Will be making releases hopefully in the near-future

### Usage:

```
$ git clone https://github.com/MeatballSteakTips/TermWave.git
$ cd TermWave
$ go build ./cmd/TermWave
$ ./TermWave
```
### Controls: 
    * Use Tab to switch focus between Toolbar and Stations Pane.
    * S saves a Station
    * X removes a station
    * Use + and - Keys to move saved stations up and down the list. Saves on Exit
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
I have had good artwork results running this on Konsole, mlterm and Xterm with vt340 emulation enabled. The program should run on any Terminal Emulator though.
⚠️ Note: Artwork will NOT appear on Kitty terminal due to Sixels being unsupported

Thank you for trying it out! I hope to make it better in the near future!

