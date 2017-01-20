package main

import (
	"errors"
	"os"

	"github.com/urfave/cli"
)

func cmdAdd(c *cli.Context) error {
	// 設定ファイルの読み込み
	name := c.Args().First()
	if name == "" {
		return errors.New("plguin name required")
	}

	conf := config()

	// true: プラグインマネージャーの種類を取得し、case文でそれぞれ処理
	f, err := os.OpenFile(conf.VimrcPath, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fatal("Error: Can't open .vimrc file.")
	}
	defer f.Close()

	var line int
	var format string

	switch conf.ManagerType {
	case "Vundle":
		line = scanAddLineForVundle(f)
		format = "Bundle '%s'"
	case "NeoBundle":
		line = scanAddLineForNeoBundle(f)
		format = "NeoBundle '%s'"
	case "dein.vim":
		line = scanAddLineForDein(f)
		format = "call dein#add('%s')"
	default:
		fatal("Error: ManagerType is not specified.")
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		fatal("Error: Fail change file offset.")
	}

	content, err := createAddPluginContent(f, format, name, line)
	if err != nil {
		fatal("Error: Can't read .vimrc file.")
	}
	if err := updateVimrc(conf.VimrcPath, content); err != nil {
		fatal("Error: Fail add plugin.")
	}

	success("Success: Add plugin.")

	return nil
}
