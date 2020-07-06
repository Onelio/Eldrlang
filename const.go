package main

import (
	"os"
	"os/exec"
	"runtime"
)

const LOGO = `
        ░░
      ██░░
    ████
    ████░░██
    ██▒▒██████       _______ __     __            __                    
    ██▒▒██▒▒██      |    ___|  |.--|  |.----.    |  |.---.-.-----.-----. ©
  ░░██▒▒▒▒▒▒▓▓      |    ___|  ||  _  ||   _|    |  ||  _  |     |  _  |
  ████▒▒░░▒▒██      |_______|__||_____||__|      |__||___._|__|__|___  |
  ██▒▒░░░░▒▒████                                                 |_____|
  ██▒▒░░░░░░▒▒██
  ▓▓▒▒▒▒░░▒▒▒▒██
  ▒▒▓▓▒▒▒▒▒▒▓▓▒▒
    ░░██████░░`

func cleanConsole() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
