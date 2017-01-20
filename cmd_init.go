package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func cmdInit(c *cli.Context) error {
	if fileExists(setting_file_path) {
		fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Setting file exist.\n")
		os.Exit(ExitCodeError)
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
		fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Setting file exist.\n")
		os.Exit(ExitCodeError)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.Write(createSettingFileText(vimrc_file_path, manager_type))
	writer.Flush()

	fmt.Printf("\x1b[32m%s\x1b[0m", "Success: Create setting file.\n")
	return nil
}
