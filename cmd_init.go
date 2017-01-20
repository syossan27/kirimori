package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
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

	var vimrcFilePath string
	fmt.Println("Type your .vimrc path. (default: ~/" + vimrcName + ")")
	fmt.Print("> ")
	fmt.Scanln(&vimrcFilePath)
	if vimrcFilePath == "" {
		vimrcFilePath = filepath.Join(homePath, vimrcName)
	}

	var managerType string
	fmt.Println("Choose a your vim bundle plugin. (default: 1)")
	fmt.Println("\t1) Vundle")
	fmt.Println("\t2) NeoBundle")
	fmt.Println("\t3) dein.vim")
	fmt.Print("Type number > ")
	fmt.Scanln(&managerType)
	switch managerType {
	case "1":
		managerType = "Vundle"
	case "2":
		managerType = "NeoBundle"
	case "3":
		managerType = "dein.vim"
	default:
		managerType = "Vundle"
	}

	file, err := os.Create(settingFilePath)
	if err != nil {
		fatal("Error: Setting file exist.")
	}
	defer file.Close()

	var conf Config
	conf.ManagerType = managerType
	conf.VimrcPath = vimrcFilePath
	err = toml.NewEncoder(file).Encode(&conf)
	if err == nil {
		success("Success: Create setting file.")
	}
	return err
}
