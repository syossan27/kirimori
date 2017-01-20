package main

import (
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"
)

func cmdAdd(c *cli.Context) error {
	// 設定ファイルの読み込み
	plugin_name := c.Args().First()
	var conf Config
	if _, err := toml.DecodeFile(setting_file_path, &conf); err != nil {
		fatal("Error: Can't read setting file.")
	}
	conf.VimrcPath = strings.Replace(conf.VimrcPath, "~", home_path, 1)
	// .vimrcのパスにファイルが存在するかどうか判定
	if fileExists(conf.VimrcPath) {
		// true: プラグインマネージャーの種類を取得し、case文でそれぞれ処理
		vimrc_file, err := os.OpenFile(conf.VimrcPath, os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			fatal("Error: Can't open .vimrc file.")
		}
		defer vimrc_file.Close()

		switch conf.ManagerType {
		case "Vundle":
			line := scanAddLineForVundle(vimrc_file)

			_, err := vimrc_file.Seek(0, 0)
			if err != nil {
				fatal("Error: Fail change file offset.")
			}

			vimrc_content, err := createAddPluginContentForVundle(vimrc_file, plugin_name, line)
			if err != nil {
				fatal("Error: Can't read .vimrc file.")
			}
			if err := updateVimrc(conf.VimrcPath, vimrc_content); err != nil {
				fatal("Error: Fail add plugin.")
			}
		case "NeoBundle":
			line := scanAddLineForNeoBundle(vimrc_file)

			_, err := vimrc_file.Seek(0, 0)
			if err != nil {
				fatal("Error: Fail change file offset.")
			}

			vimrc_content, err := createAddPluginContentForNeoBundle(vimrc_file, plugin_name, line)
			if err != nil {
				fatal("Error: Can't read .vimrc file.")
			}
			if err := updateVimrc(conf.VimrcPath, vimrc_content); err != nil {
				fatal("Error: Fail add plugin.")
			}
		case "dein.vim":
			line := scanAddLineForDein(vimrc_file)

			_, err := vimrc_file.Seek(0, 0)
			if err != nil {
				fatal("Error: Fail change file offset.")
			}

			vimrc_content, err := createAddPluginContentForDein(vimrc_file, plugin_name, line)
			if err != nil {
				fatal("Error: Can't read .vimrc file.")
			}
			if err := updateVimrc(conf.VimrcPath, vimrc_content); err != nil {
				fatal("Error: Fail add plugin.")
			}
		default:
			fatal("Error: ManagerType is not specified.")
		}
	} else {
		fatal("Error: No .vimrc file exists.")
	}

	success("Success: Add plugin.")

	return nil
}
