package main

import (
	"errors"
	"os"

	"github.com/urfave/cli"
)

func cmdRemove(c *cli.Context) error {
	// 設定ファイルの読み込み
	name := c.Args().First()
	if name == "" {
		return errors.New("plguin name required")
	}

	conf := config()

	f, err := os.OpenFile(conf.VimrcPath, os.O_RDONLY, 0644)
	if err != nil {
		fatal("Error: Can't open .vimrc file.")
	}
	defer f.Close()

	var line int

	// true: プラグインマネージャーの種類を取得し、case文でそれぞれ処理
	switch conf.ManagerType {
	case "Vundle":
		line = scanRemoveLineForVundle(f, name)
	case "NeoBundle":
		line = scanRemoveLineForNeoBundle(f, name)
	case "dein.vim":
		line = scanRemoveLineForDein(f, name)
	default:
		fatal("Error: ManagerType is not specified.")
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		fatal("Error: Fail change file offset.")
	}

	content, err := createRemovePluginContent(f, name, line)
	if err != nil {
		fatal("Error: Can't read .vimrc file.")
	}
	if err := updateVimrc(conf.VimrcPath, content); err != nil {
		fatal("Error: Fail remove plugin.")
	}

	success("Success: Remove plugin.")
	return nil
}
