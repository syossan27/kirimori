package main

import (
	"errors"
	"os"

	"github.com/urfave/cli"
)

func cmdRemove(c *cli.Context) error {
	// 設定ファイルの読み込み
	pluginName := c.Args().First()
	if pluginName == "" {
		return errors.New("plguin name required")
	}

	conf := config()

	vimrcFile, err := os.OpenFile(conf.VimrcPath, os.O_RDONLY, 0644)
	if err != nil {
		fatal("Error: Can't open .vimrc file.")
	}
	defer vimrcFile.Close()

	// true: プラグインマネージャーの種類を取得し、case文でそれぞれ処理
	switch conf.ManagerType {
	case "Vundle":
		line := scanRemoveLineForVundle(vimrcFile, pluginName)

		_, err := vimrcFile.Seek(0, 0)
		if err != nil {
			fatal("Error: Fail change file offset.")
		}

		vimrcContent, err := createRemovePluginContentForVundle(vimrcFile, pluginName, line)
		if err != nil {
			fatal("Error: Can't read .vimrc file.")
		}
		if err := updateVimrc(conf.VimrcPath, vimrcContent); err != nil {
			fatal("Error: Fail remove plugin.")
		}
	case "NeoBundle":
		line := scanRemoveLineForNeoBundle(vimrcFile, pluginName)

		_, err := vimrcFile.Seek(0, 0)
		if err != nil {
			fatal("Error: Fail change file offset.")
		}

		vimrcContent, err := createRemovePluginContentForNeoBundle(vimrcFile, pluginName, line)
		if err != nil {
			fatal("Error: Can't read .vimrc file.")
		}
		if err := updateVimrc(conf.VimrcPath, vimrcContent); err != nil {
			fatal("Error: Fail remove plugin.")
		}
	case "dein.vim":
		line := scanRemoveLineForDein(vimrcFile, pluginName)

		_, err := vimrcFile.Seek(0, 0)
		if err != nil {
			fatal("Error: Fail change file offset.")
		}

		vimrcContent, err := createRemovePluginContentForDein(vimrcFile, pluginName, line)
		if err != nil {
			fatal("Error: Can't read .vimrc file.")
		}
		if err := updateVimrc(conf.VimrcPath, vimrcContent); err != nil {
			fatal("Error: Fail remove plugin.")
		}
	default:
		fatal("Error: ManagerType is not specified.")
	}

	success("Success: Remove plugin.")
	return nil
}
