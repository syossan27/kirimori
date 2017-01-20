package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func cmdInit(c *cli.Context) error {
	if fileExists(setting_file_path) {
		fatal("Error: Setting file exist.")
	}

	var vimrc_file_path string
	fmt.Println("Type your .vimrc path. (default: ~/.vimrc)")
	fmt.Print("> ")
	fmt.Scanln(&vimrc_file_path)
	if vimrc_file_path == "" {
		vimrc_file_path = "~/.vimrc"
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

	writer := bufio.NewWriter(file)
	writer.Write(createSettingFileText(vimrc_file_path, manager_type))
	writer.Flush()

	success("Success: Create setting file.")
	return nil
}
