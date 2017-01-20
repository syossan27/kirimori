package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"
)

func cmdRemove(c *cli.Context) error {
	// 設定ファイルの読み込み
	plugin_name := c.Args().First()
	var conf Config
	if _, err := toml.DecodeFile(setting_file_path, &conf); err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Can't read setting file.\n")
		os.Exit(ExitCodeError)
	}
	conf.VimrcPath = strings.Replace(conf.VimrcPath, "~", home_path, 1)
	// .vimrcのパスにファイルが存在するかどうか判定
	if fileExists(conf.VimrcPath) {
		vimrc_file, err := os.OpenFile(conf.VimrcPath, os.O_RDONLY, 0644)
		if err != nil {
			fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Can't open .vimrc file.\n")
			os.Exit(ExitCodeError)
		}
		defer vimrc_file.Close()
		// true: プラグインマネージャーの種類を取得し、case文でそれぞれ処理
		switch conf.ManagerType {
		case "Vundle":
			line := scanRemoveLineForVundle(vimrc_file, plugin_name)

			_, err := vimrc_file.Seek(0, 0)
			if err != nil {
				fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail change file offset.\n")
				os.Exit(ExitCodeError)
			}

			vimrc_content, err := createRemovePluginContentForVundle(vimrc_file, plugin_name, line)
			if err != nil {
				fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Can't read .vimrc file.\n")
				os.Exit(ExitCodeError)
			}
			if err := updateVimrc(conf.VimrcPath, vimrc_content); err != nil {
				fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail remove plugin.\n")
				os.Exit(ExitCodeError)
			}
		case "NeoBundle":
			line := scanRemoveLineForNeoBundle(vimrc_file, plugin_name)

			_, err := vimrc_file.Seek(0, 0)
			if err != nil {
				fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail change file offset.\n")
				os.Exit(ExitCodeError)
			}

			vimrc_content, err := createRemovePluginContentForNeoBundle(vimrc_file, plugin_name, line)
			if err != nil {
				fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Can't read .vimrc file.\n")
				os.Exit(ExitCodeError)
			}
			if err := updateVimrc(conf.VimrcPath, vimrc_content); err != nil {
				fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail remove plugin.\n")
				os.Exit(ExitCodeError)
			}
		case "dein.vim":
			line := scanRemoveLineForDein(vimrc_file, plugin_name)

			_, err := vimrc_file.Seek(0, 0)
			if err != nil {
				fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail change file offset.\n")
				os.Exit(ExitCodeError)
			}

			vimrc_content, err := createRemovePluginContentForDein(vimrc_file, plugin_name, line)
			if err != nil {
				fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Can't read .vimrc file.\n")
				os.Exit(ExitCodeError)
			}
			if err := updateVimrc(conf.VimrcPath, vimrc_content); err != nil {
				fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail remove plugin.\n")
				os.Exit(ExitCodeError)
			}
		default:
			fmt.Printf("\x1b[31m%s\x1b[0m", "Error: ManagerType is not specified.\n")
			os.Exit(ExitCodeError)
		}
	} else {
		fmt.Printf("\x1b[31m%s\x1b[0m", "Error: No .vimrc file exists.\n")
		os.Exit(ExitCodeError)
	}

	fmt.Printf("\x1b[32m%s\x1b[0m", "Success: Remove plugin.\n")
	return nil
}
