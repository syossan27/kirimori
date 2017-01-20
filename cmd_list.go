package main

import (
	"os"

	"github.com/urfave/cli"
)

func cmdList(c *cli.Context) error {
	conf := config()

	// true: プラグインマネージャーの種類を取得し、case文でそれぞれ処理
	vimrcFile, err := os.Open(conf.VimrcPath)
	if err != nil {
		fatal("Error: Can't open .vimrc file.")
	}
	defer vimrcFile.Close()

	// true: プラグインマネージャーの種類を取得し、case文でそれぞれ処理
	switch conf.ManagerType {
	case "Vundle":
		listPlugin(scanListPluginForVundle(vimrcFile))
	case "NeoBundle":
		listPlugin(scanListPluginForNeoBundle(vimrcFile))
	case "dein.vim":
		listPlugin(scanListPluginForDein(vimrcFile))
	default:
		fatal("Error: ManagerType is not specified.")
	}
	return nil
}
