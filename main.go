package main

import (
	"bd-install/utils"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

// Default Channel
var releaseChannel = "stable"

func main() {
	// Get Target Exe
	var targetExe = utils.GetExeName(releaseChannel)

	// Kill Discord if it's running
	var exe = utils.GetProcessExe(targetExe)
	if len(exe) > 0 {
		if err := utils.KillProcess(targetExe); err != nil {
			log.Fatal("Could not kill Discord")
			return
		}
	}

	// Make BD directories
	if err := os.MkdirAll(utils.Data, 0755); err != nil {
		log.Fatal("Could not create BetterDiscord folder")
	}

	if err := os.MkdirAll(utils.Plugins, 0755); err != nil {
		log.Fatal("Could not create plugins folder")
	}

	if err := os.MkdirAll(utils.Themes, 0755); err != nil {
		log.Fatal("Could not create theme folder")
	}

	// Download asar into the BD folder
	var downloadUrl = "https://github.com/BetterDiscord/BetterDiscord/releases/latest/download/betterdiscord.asar"
	var asarPath = path.Join(utils.Data, "betterdiscord.asar")
	var err = utils.DownloadFile(downloadUrl, asarPath)
	if err != nil {
		log.Fatal("Could not download asar")
	}

	// Inject shim loader
	var corePath = utils.DiscordPath(releaseChannel)

	var indString = `require("` + asarPath + `");`
	indString = strings.ReplaceAll(indString, `\`, "/")
	indString = indString + "\nmodule.exports = require(\"./core.asar\");"

	if err := os.WriteFile(path.Join(corePath, "index.js"), []byte(indString), 0755); err != nil {
		log.Fatal("Could not write index.js in discord_desktop_core!")
	}

	// Launch Discord if we killed it
	if len(exe) > 0 {
		var cmd = exec.Command(exe)
		var err = cmd.Start()
		if err != nil {
			fmt.Println(err)
		}
	}
}
