package main

import (
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"
)

func cmdList(c *cli.Context) error {
	// 設定ファイルの読み込み
	var conf Config
	if _, err := toml.DecodeFile(settingFilePath, &conf); err != nil {
		fatal("Error: Can't read setting file.")
	}
	conf.VimrcPath = strings.Replace(conf.VimrcPath, "~", homePath, 1)
	// .vimrcのパスにファイルが存在するかどうか判定
	if !fileExists(conf.VimrcPath) {
		fatal("Error: No .vimrc file exists.\n")
	}
	vimrcFile, err := os.OpenFile(conf.VimrcPath, os.O_RDONLY, 0644)
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
