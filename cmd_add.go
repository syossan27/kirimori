package main

import (
	"errors"
	"os"

	"github.com/urfave/cli"
)

func cmdAdd(c *cli.Context) error {
	name := c.Args().First()
	if name == "" {
		return errors.New("plguin name required")
	}

	// 設定ファイルの読み込み
	conf := config()

	// true: プラグインマネージャーの種類を取得し、case文でそれぞれ処理
	f, err := os.OpenFile(conf.VimrcPath, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fatal("Error: Can't open .vimrc file.")
	}
	defer f.Close()

	manager := conf.Manager()
	line := manager.AddLine(f)

	_, err = f.Seek(0, 0)
	if err != nil {
		fatal("Error: Fail change file offset.")
	}

	b, err := createAddPluginContent(f, manager.Format(name), line)
	if err != nil {
		fatal("Error: Can't read .vimrc file.")
	}
	if err := updateVimrc(conf.VimrcPath, b); err != nil {
		fatal("Error: Fail add plugin.")
	}

	manager.InstallExCmd()

	success("Success: Add plugin.")

	return nil
}
