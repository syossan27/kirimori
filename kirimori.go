package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

const (
	ExitCodeOK = iota
	ExitCodeError
)

type Config struct {
	VimrcPath   string
	ManagerType string
}

func makeApp() *cli.App {
	app := cli.NewApp()

	app.Name = "kirimori"
	app.Usage = "Add Vim Plugin Tool"
	app.Version = "1.0"
	home_path, _ := homedir.Dir()
	var setting_file_path string = home_path + "/.kirimori.toml"

	app.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "create setting file",
			Action: func(c *cli.Context) error {
				if fileExists(setting_file_path) {
					println("Error: Setting file exist.")
					os.Exit(ExitCodeError)
				} else {
					file, err := os.Create(setting_file_path)
					if err != nil {
						println("Error: Can't create setting file.")
						os.Exit(ExitCodeError)
					}
					defer file.Close()

					writer := bufio.NewWriter(file)
					writer.Write(defaultSettingFileText())
					writer.Flush()
				}
				println("Success: Create setting file.")
				return nil
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add vim plugin",
			Action: func(c *cli.Context) error {
				// 設定ファイルの読み込み
				var package_name = c.Args().First()
				var conf Config
				if _, err := toml.DecodeFile(setting_file_path, &conf); err != nil {
					println("Error: Can't read setting file.")
					os.Exit(ExitCodeError)
				}
				conf.VimrcPath = strings.Replace(conf.VimrcPath, "~", home_path, 1)
				// .vimrcのパスにファイルが存在するかどうか判定
				if fileExists(conf.VimrcPath) {
					vimrc_file, err := os.OpenFile(conf.VimrcPath, os.O_RDWR|os.O_APPEND, 0666)
					if err != nil {
						println("Error: Can't open .vimrc file.")
						os.Exit(ExitCodeError)
					}
					defer vimrc_file.Close()
					// true: プラグインマネージャーの種類を取得し、case文でそれぞれ処理
					switch conf.ManagerType {
					case "Vundle":
						if err := addPackageForVundle(vimrc_file, package_name); err != nil {
							println("Error: Fail add package.")
							os.Exit(ExitCodeError)
						}
					case "NeoBundle":
						if err := addPackageForNeoBundle(vimrc_file, package_name); err != nil {
							println("Error: Fail add package.")
							os.Exit(ExitCodeError)
						}
					case "dein.vim":
						if err := addPackageForDein(vimrc_file); err != nil {
							println("Error: Fail add package.")
							os.Exit(ExitCodeError)
						}
					default:
						println("Error: ManagerType is not specified.")
						os.Exit(ExitCodeError)
					}
				} else {
					println("Error: No .vimrc file exists.")
					os.Exit(ExitCodeError)
				}
				return nil
			},
		},
	}

	return app
}

func defaultSettingFileText() []byte {
	body := []string{
		"# VimrcPath = \"~/.vimrc\"",
		"",
		"# ManagerType = \"NeoBundle\"",
		"# ManagerType = \"Vundle\"",
		"# ManagerType = \"dein.vim\"",
	}

	return []byte(strings.Join(body, "\n"))
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func addPackageForVundle(vimrc_file *os.File, package_name string) error {
	writer := bufio.NewWriter(vimrc_file)
	_, err := writer.WriteString("Bundle '" + package_name + "'")
	writer.Flush()
	return err
}

func addPackageForNeoBundle(vimrc_file *os.File) error {
	writer := bufio.NewWriter(vimrc_file)
	_, err := writer.WriteString("NeoBundle '" + package_name + "'")
	writer.Flush()
	return nil
}

func addPackageForDein(vimrc_file *os.File) error {
	return nil
}

func main() {
	app := makeApp()
	app.Run(os.Args)
}
