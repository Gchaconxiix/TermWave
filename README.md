
# TermWave
### A Terminal-based internet radio player for Linux

<img width="1258" height="586" alt="Screenshot_20260423_084434" src="https://github.com/user-attachments/assets/152428a3-ea57-4160-b258-e8c2d20bf24d" />

Termwave is a TUI that organizes internet radio stations and plays them back.
Program was made using Go and the BubbleTea V2 framework. It is essentially a frontend for MPV and api.radio-browser.info 
I wanted something that I can run on both my PC and on a very low-power CyberDeck for traveling/camping, so I decided to make this.

⚠️ Termwave is currently still being developed. While I will be actively adding/fixing features, things may break and some things may be removed until I can fix them.
⚠️ MPV is required for audio playback! Chafa is required for displaying station images (Unless using Kitty)

### Currently Working/Recent Changes:
    * Kitty Terminal Support (Finally!)
    * Search for internet radio stations (Using api.radio-browser)
    * Functionality to Save and Remove Stations has been added
    * Added artwork through Chafa
    * Added Documentation Menu
    * Added Theme Settings Menu
    * Added Station Sorting using +/- keys
    * Can now change color of each outline in Theme Settings (ANSI or Hex)
    * Added Dynamic List functionaliy. Stations per page will now change based on the current window size
    * Changed Chafa to automatically select best format based on terminal (Better compatibility)
    * Added https Manual Station addition (Also works with Youtube links with yt-dlp)
### Roadmap/Things to fix:
    * Make the UI look better
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
I have had good artwork results running this on Konsole, Kity, mlterm and Xterm with vt340 emulation enabled. The program should run on any Terminal Emulator though.

Thank you for trying it out! I hope to make it better in the near future!

