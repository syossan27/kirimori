package main

import (
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"
)

func cmdAdd(c *cli.Context) error {
	// 設定ファイルの読み込み
	pluginName := c.Args().First()
	var conf Config
	if _, err := toml.DecodeFile(settingFilePath, &conf); err != nil {
		fatal("Error: Can't read setting file.")
	}
	conf.VimrcPath = strings.Replace(conf.VimrcPath, "~", homePath, 1)
	// .vimrcのパスにファイルが存在するかどうか判定
	if fileExists(conf.VimrcPath) {
		// true: プラグインマネージャーの種類を取得し、case文でそれぞれ処理
		vimrcFile, err := os.OpenFile(conf.VimrcPath, os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			fatal("Error: Can't open .vimrc file.")
		}
		defer vimrcFile.Close()

		switch conf.ManagerType {
		case "Vundle":
			line := scanAddLineForVundle(vimrcFile)

			_, err := vimrcFile.Seek(0, 0)
			if err != nil {
				fatal("Error: Fail change file offset.")
			}

			vimrcContent, err := createAddPluginContentForVundle(vimrcFile, pluginName, line)
			if err != nil {
				fatal("Error: Can't read .vimrc file.")
			}
			if err := updateVimrc(conf.VimrcPath, vimrcContent); err != nil {
				fatal("Error: Fail add plugin.")
			}
		case "NeoBundle":
			line := scanAddLineForNeoBundle(vimrcFile)

			_, err := vimrcFile.Seek(0, 0)
			if err != nil {
				fatal("Error: Fail change file offset.")
			}

			vimrcContent, err := createAddPluginContentForNeoBundle(vimrcFile, pluginName, line)
			if err != nil {
				fatal("Error: Can't read .vimrc file.")
			}
			if err := updateVimrc(conf.VimrcPath, vimrcContent); err != nil {
				fatal("Error: Fail add plugin.")
			}
		case "dein.vim":
			line := scanAddLineForDein(vimrcFile)

			_, err := vimrcFile.Seek(0, 0)
			if err != nil {
				fatal("Error: Fail change file offset.")
			}

			vimrcContent, err := createAddPluginContentForDein(vimrcFile, pluginName, line)
			if err != nil {
				fatal("Error: Can't read .vimrc file.")
			}
			if err := updateVimrc(conf.VimrcPath, vimrcContent); err != nil {
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
