package main

import (
	"os"

	"github.com/urfave/cli"
)

func cmdList(c *cli.Context) error {
	// 設定ファイルの読み込み
	conf := config()

	// true: プラグインマネージャーの種類を取得し、case文でそれぞれ処理
	f, err := os.Open(conf.VimrcPath)
	if err != nil {
		fatal("Error: Can't open .vimrc file.")
	}
	defer f.Close()

	// true: プラグインマネージャーの種類を取得し、case文でそれぞれ処理
	listPlugin(conf.Manager().ListPlugins(f))

	return nil
}
