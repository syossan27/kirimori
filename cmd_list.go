package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"
)

func cmdList(c *cli.Context) error {
	// 設定ファイルの読み込み
	var conf Config
	if _, err := toml.DecodeFile(setting_file_path, &conf); err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Can't read setting file.\n")
		os.Exit(ExitCodeError)
	}
	conf.VimrcPath = strings.Replace(conf.VimrcPath, "~", home_path, 1)
	// .vimrcのパスにファイルが存在するかどうか判定
	if fileExists(conf.VimrcPath) {
		vimrc_file, err := os.OpenFile(conf.VimrcPath, os.O_RDONLY, 0644)
		if err != nil {
			fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Can't open .vimrc file.\n")
			os.Exit(ExitCodeError)
		}
		defer vimrc_file.Close()

		// true: プラグインマネージャーの種類を取得し、case文でそれぞれ処理
		switch conf.ManagerType {
		case "Vundle":
			listPlugin(scanListPluginForVundle(vimrc_file))
		case "NeoBundle":
			listPlugin(scanListPluginForNeoBundle(vimrc_file))
		case "dein.vim":
			listPlugin(scanListPluginForDein(vimrc_file))
		default:
			fmt.Printf("\x1b[31m%s\x1b[0m", "Error: ManagerType is not specified.\n")
			os.Exit(ExitCodeError)
		}
	} else {
		fmt.Printf("\x1b[31m%s\x1b[0m", "Error: No .vimrc file exists.\n")
		os.Exit(ExitCodeError)
	}
	return nil
}
