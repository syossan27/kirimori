package main

import (
	"errors"
	"os"

	"github.com/urfave/cli"
)

func cmdRemove(c *cli.Context) error {
	name := c.Args().First()
	if name == "" {
		return errors.New("plguin name required")
	}

	// 設定ファイルの読み込み
	conf := config()

	f, err := os.OpenFile(conf.VimrcPath, os.O_RDONLY, 0644)
	if err != nil {
		fatal("Error: Can't open .vimrc file.")
	}
	defer f.Close()

	// true: プラグインマネージャーの種類を取得し、case文でそれぞれ処理
	line := conf.Manager().RemoveLine(f, name)

	_, err = f.Seek(0, 0)
	if err != nil {
		fatal("Error: Fail change file offset.")
	}

	b, err := createRemovePluginContent(f, line)
	if err != nil {
		fatal("Error: Can't read .vimrc file.")
	}
	if err := updateVimrc(conf.VimrcPath, b); err != nil {
		fatal("Error: Fail remove plugin.")
	}

	success("Success: Remove plugin.")
	return nil
}
