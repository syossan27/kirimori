package main

import (
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
		fatal("Error: Can't read setting file.")
	}
	conf.VimrcPath = strings.Replace(conf.VimrcPath, "~", home_path, 1)
	// .vimrcのパスにファイルが存在するかどうか判定
	if fileExists(conf.VimrcPath) {
		vimrc_file, err := os.OpenFile(conf.VimrcPath, os.O_RDONLY, 0644)
		if err != nil {
			fatal("Error: Can't open .vimrc file.")
		}
		defer vimrc_file.Close()
		// true: プラグインマネージャーの種類を取得し、case文でそれぞれ処理
		switch conf.ManagerType {
		case "Vundle":
			line := scanRemoveLineForVundle(vimrc_file, plugin_name)

			_, err := vimrc_file.Seek(0, 0)
			if err != nil {
				fatal("Error: Fail change file offset.")
			}

			vimrc_content, err := createRemovePluginContentForVundle(vimrc_file, plugin_name, line)
			if err != nil {
				fatal("Error: Can't read .vimrc file.")
			}
			if err := updateVimrc(conf.VimrcPath, vimrc_content); err != nil {
				fatal("Error: Fail remove plugin.")
			}
		case "NeoBundle":
			line := scanRemoveLineForNeoBundle(vimrc_file, plugin_name)

			_, err := vimrc_file.Seek(0, 0)
			if err != nil {
				fatal("Error: Fail change file offset.")
			}

			vimrc_content, err := createRemovePluginContentForNeoBundle(vimrc_file, plugin_name, line)
			if err != nil {
				fatal("Error: Can't read .vimrc file.")
			}
			if err := updateVimrc(conf.VimrcPath, vimrc_content); err != nil {
				fatal("Error: Fail remove plugin.")
			}
		case "dein.vim":
			line := scanRemoveLineForDein(vimrc_file, plugin_name)

			_, err := vimrc_file.Seek(0, 0)
			if err != nil {
				fatal("Error: Fail change file offset.")
			}

			vimrc_content, err := createRemovePluginContentForDein(vimrc_file, plugin_name, line)
			if err != nil {
				fatal("Error: Can't read .vimrc file.")
			}
			if err := updateVimrc(conf.VimrcPath, vimrc_content); err != nil {
				fatal("Error: Fail remove plugin.")
			}
		default:
			fatal("Error: ManagerType is not specified.")
		}
	} else {
		fatal("Error: No .vimrc file exists.")
	}

	success("Success: Remove plugin.")
	return nil
}
