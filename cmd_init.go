package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/Songmu/prompter"
	"github.com/urfave/cli"
)

func cmdInit(c *cli.Context) error {
	if fileExists(settingFilePath) {
		fatal("Error: Setting file exist.")
	}

	var vimrcName string
	if runtime.GOOS == "windows" {
		vimrcName = "_vimrc"
	} else {
		vimrcName = ".vimrc"
	}

	var filename string = prompter.Prompt("Type your .vimrc path.", filepath.Join(homePath, vimrcName))

	var managerType string
	fmt.Println("Choose a your vim bundle plugin.")
	for i, pm := range pluginManagers {
		fmt.Printf("\t%d) %s : %s\n", i+1, pm.Name, pm.URL)
	}
	managerType = "Vundle"
	for {
		fmt.Print("Type number. [1]: ")
		var s string
		if _, err := fmt.Scanln(&s); err != nil {
			if s == "" {
				break
			}
			return err
		}
		s = strings.TrimSpace(s)
		if i, err := strconv.Atoi(s); err == nil {
			if i > 0 && i <= len(pluginManagers) {
				managerType = pluginManagers[i-1].Name
				break
			}
		}
	}

	file, err := os.Create(settingFilePath)
	if err != nil {
		fatal("Error: Setting file exist.")
	}
	defer file.Close()

	var conf Config
	conf.ManagerType = managerType
	conf.VimrcPath = filename
	err = toml.NewEncoder(file).Encode(&conf)
	if err == nil {
		success("Success: Create setting file.")
	}
	return err
}
