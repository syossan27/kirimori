package main

import (
	"os"

	"github.com/urfave/cli"
)

func cmdList(c *cli.Context) error {
	// 設定ファイルの読み込み
	conf := config()

	f, err := os.Open(conf.VimrcPath)
	if err != nil {
		fatal("Error: Can't open .vimrc file.")
	}
	defer f.Close()

	printLines(conf.Manager().ListPlugin(f))

	return nil
}
