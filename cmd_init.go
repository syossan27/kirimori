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
	if fileExists(setting_file_path) {
		fatal("Error: Setting file exist.")
	}

	var vimrc_name string
	if runtime.GOOS == "windows" {
		vimrc_name = "_vimrc"
	} else {
		vimrc_name = ".vimrc"
	}

	var vimrc_file_path string
	fmt.Println("Type your .vimrc path. (default: ~/" + vimrc_name + ")")
	fmt.Print("> ")
	fmt.Scanln(&vimrc_file_path)
	if vimrc_file_path == "" {
		vimrc_file_path = filepath.Join(home_path, vimrc_name)
	}

	var manager_type string
	fmt.Println("Choose a your vim bundle plugin. (default: 1)")
	fmt.Println("\t1) Vundle")
	fmt.Println("\t2) NeoBundle")
	fmt.Println("\t3) dein.vim")
	fmt.Print("Type number > ")
	fmt.Scanln(&manager_type)
	switch manager_type {
	case "1":
		manager_type = "Vundle"
	case "2":
		manager_type = "NeoBundle"
	case "3":
		manager_type = "dein.vim"
	default:
		manager_type = "Vundle"
	}

	file, err := os.Create(setting_file_path)
	if err != nil {
		fatal("Error: Setting file exist.")
	}
	defer file.Close()

	var conf Config
	conf.ManagerType = manager_type
	conf.VimrcPath = vimrc_file_path
	err = toml.NewEncoder(file).Encode(&conf)
	if err == nil {
		success("Success: Create setting file.")
	}
	return err
}
