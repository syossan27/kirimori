package main

import (
	"bufio"
	"fmt"
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
					fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Setting file exist.\n")
					os.Exit(ExitCodeError)
				} else {
					file, err := os.Create(setting_file_path)
					if err != nil {
						fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Setting file exist.\n")
						os.Exit(ExitCodeError)
					}
					defer file.Close()

					writer := bufio.NewWriter(file)
					writer.Write(defaultSettingFileText())
					writer.Flush()
				}
				fmt.Printf("\x1b[32m%s\x1b[0m", "Success: Create setting file.\n")
				return nil
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add plugin",
			Action: func(c *cli.Context) error {
				// 設定ファイルの読み込み
				var plugin_name = c.Args().First()
				var conf Config
				if _, err := toml.DecodeFile(setting_file_path, &conf); err != nil {
					println("Error: Can't read setting file.")
					os.Exit(ExitCodeError)
				}
				conf.VimrcPath = strings.Replace(conf.VimrcPath, "~", home_path, 1)
				// .vimrcのパスにファイルが存在するかどうか判定
				if fileExists(conf.VimrcPath) {
					// true: プラグインマネージャーの種類を取得し、case文でそれぞれ処理
					vimrc_file, err := os.OpenFile(conf.VimrcPath, os.O_RDWR|os.O_APPEND, 0666)
					if err != nil {
						fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Can't open .vimrc file.\n")
						os.Exit(ExitCodeError)
					}
					defer vimrc_file.Close()

					switch conf.ManagerType {
					case "Vundle":
						if err := addPluginForVundle(vimrc_file, plugin_name); err != nil {
							fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail add plugin.\n")
							os.Exit(ExitCodeError)
						}
					case "NeoBundle":
						if err := addPluginForNeoBundle(vimrc_file, plugin_name); err != nil {
							fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail add plugin.\n")
							os.Exit(ExitCodeError)
						}
					case "dein.vim":
						vimrc_content, err := createAddPluginContentForDein(vimrc_file, plugin_name)
						if err != nil {
							fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Can't read .vimrc file.\n")
							os.Exit(ExitCodeError)
						}
						if err := updateVimrc(conf.VimrcPath, vimrc_content); err != nil {
							fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail add plugin.\n")
							os.Exit(ExitCodeError)
						}
					default:
						fmt.Printf("\x1b[31m%s\x1b[0m", "Error: ManagerType is not specified.\n")
						os.Exit(ExitCodeError)
					}
				} else {
					fmt.Printf("\x1b[31m%s\x1b[0m", "Error: No .vimrc file exists.\n")
					os.Exit(ExitCodeError)
				}

				fmt.Printf("\x1b[32m%s\x1b[0m", "Success: Add plugin.\n")
				return nil
			},
		},
		{
			Name:    "remove",
			Aliases: []string{"r"},
			Usage:   "remove plugin",
			Action: func(c *cli.Context) error {
				// 設定ファイルの読み込み
				var plugin_name = c.Args().First()
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
						vimrc_content, err := createRemovePluginContentForVundle(vimrc_file, plugin_name)
						if err != nil {
							fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Can't read .vimrc file.\n")
							os.Exit(ExitCodeError)
						}
						if err := updateVimrc(conf.VimrcPath, vimrc_content); err != nil {
							fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail remove plugin.\n")
							os.Exit(ExitCodeError)
						}
					case "NeoBundle":
						vimrc_content, err := createRemovePluginContentForNeoBundle(vimrc_file, plugin_name)
						if err != nil {
							fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Can't read .vimrc file.\n")
							os.Exit(ExitCodeError)
						}
						if err := updateVimrc(conf.VimrcPath, vimrc_content); err != nil {
							fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail remove plugin.\n")
							os.Exit(ExitCodeError)
						}
					case "dein.vim":
						vimrc_content, err := createRemovePluginContentForDein(vimrc_file, plugin_name)
						if err != nil {
							fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Can't read .vimrc file.\n")
							os.Exit(ExitCodeError)
						}
						if err := updateVimrc(conf.VimrcPath, vimrc_content); err != nil {
							fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail remove plugin.\n")
							os.Exit(ExitCodeError)
						}
					default:
						fmt.Printf("\x1b[31m%s\x1b[0m", "Error: ManagerType is not specified.\n")
						os.Exit(ExitCodeError)
					}
				} else {
					fmt.Printf("\x1b[31m%s\x1b[0m", "Error: No .vimrc file exists.\n")
					os.Exit(ExitCodeError)
				}

				fmt.Printf("\x1b[32m%s\x1b[0m", "Success: Remove plugin.\n")
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list plugin",
			Action: func(c *cli.Context) error {
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
						err := listPluginForVundle(vimrc_file)
						if err != nil {
							fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Can't read .vimrc file.\n")
							os.Exit(ExitCodeError)
						}
					case "NeoBundle":
						err := listPluginForNeoBundle(vimrc_file)
						if err != nil {
							fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Can't read .vimrc file.\n")
							os.Exit(ExitCodeError)
						}
					case "dein.vim":
						err := listPluginForDein(vimrc_file)
						if err != nil {
							fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Can't read .vimrc file.\n")
							os.Exit(ExitCodeError)
						}
					default:
						fmt.Printf("\x1b[31m%s\x1b[0m", "Error: ManagerType is not specified.\n")
						os.Exit(ExitCodeError)
					}
				} else {
					fmt.Printf("\x1b[31m%s\x1b[0m", "Error: No .vimrc file exists.\n")
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

func addPluginForVundle(vimrc_file *os.File, plugin_name string) error {
	writer := bufio.NewWriter(vimrc_file)
	_, err := writer.WriteString("\nBundle '" + plugin_name + "'")
	writer.Flush()
	return err
}

func createRemovePluginContentForVundle(vimrc_file *os.File, plugin_name string) ([]byte, error) {
	var rows []string
	scanner := bufio.NewScanner(vimrc_file)
	for scanner.Scan() {
		var scan_text = scanner.Text()
		if strings.Contains(scan_text, "Bundle '"+plugin_name+"'") {
			continue
		} else {
			rows = append(rows, scan_text)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Can't read .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	vimrc_content := []byte(strings.Join(rows, "\n"))
	err := scanner.Err()
	return vimrc_content, err
}

func addPluginForNeoBundle(vimrc_file *os.File, plugin_name string) error {
	writer := bufio.NewWriter(vimrc_file)
	_, err := writer.WriteString("\nNeoBundle '" + plugin_name + "'")
	writer.Flush()
	return err
}

func createRemovePluginContentForNeoBundle(vimrc_file *os.File, plugin_name string) ([]byte, error) {
	var rows []string
	scanner := bufio.NewScanner(vimrc_file)
	for scanner.Scan() {
		var scan_text = scanner.Text()
		rows = append(rows, scan_text)
		if strings.Contains(scan_text, "NeoBundle '"+plugin_name+"'") {
			continue
		} else {
			rows = append(rows, scan_text)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Can't read .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	vimrc_content := []byte(strings.Join(rows, "\n"))
	err := scanner.Err()
	return vimrc_content, err
}

func createAddPluginContentForDein(vimrc_file *os.File, plugin_name string) ([]byte, error) {
	var rows []string
	scanner := bufio.NewScanner(vimrc_file)
	for scanner.Scan() {
		var scan_text = scanner.Text()
		rows = append(rows, scan_text)
		if strings.Contains(scan_text, "call dein#begin") {
			rows = append(rows, "call dein#add('"+plugin_name+"')")
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Can't read .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	vimrc_content := []byte(strings.Join(rows, "\n"))
	err := scanner.Err()
	return vimrc_content, err
}

func createRemovePluginContentForDein(vimrc_file *os.File, plugin_name string) ([]byte, error) {
	var rows []string
	scanner := bufio.NewScanner(vimrc_file)
	for scanner.Scan() {
		var scan_text = scanner.Text()
		if strings.Contains(scan_text, "call dein#add('"+plugin_name+"')") {
			continue
		} else {
			rows = append(rows, scan_text)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Can't read .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	vimrc_content := []byte(strings.Join(rows, "\n"))
	err := scanner.Err()
	return vimrc_content, err
}

func updateVimrc(vimrc_file_path string, vimrc_content []byte) error {
	vimrc_file, err := os.Create(vimrc_file_path)
	if err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Can't open .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	writer := bufio.NewWriter(vimrc_file)
	writer.Write(vimrc_content)
	writer.Flush()
	return err
}

func listPluginForVundle(vimrc_file *os.File) error {
	scanner := bufio.NewScanner(vimrc_file)
	for scanner.Scan() {
		var scan_text = scanner.Text()
		if strings.Contains(scan_text, "Bundle '") {
			scan_text = strings.Replace(scan_text, "Bundle", "", 1)
			scan_text = strings.Replace(scan_text, "'", "", -1)
			println(scan_text)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Can't read .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	err := scanner.Err()
	return err
}

func listPluginForNeoBundle(vimrc_file *os.File) error {
	scanner := bufio.NewScanner(vimrc_file)
	for scanner.Scan() {
		var scan_text = scanner.Text()
		if strings.Contains(scan_text, "NeoBundle '") {
			scan_text = strings.Replace(scan_text, "NeoBundle", "", 1)
			scan_text = strings.Replace(scan_text, "'", "", -1)
			println(scan_text)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Can't read .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	err := scanner.Err()
	return err
}

func listPluginForDein(vimrc_file *os.File) error {
	scanner := bufio.NewScanner(vimrc_file)
	for scanner.Scan() {
		var scan_text = scanner.Text()
		if strings.Contains(scan_text, "call dein#add") {
			scan_text = strings.Replace(scan_text, "call dein#add", "", 1)
			scan_text = strings.Replace(scan_text, "('", "", 1)
			scan_text = strings.Replace(scan_text, "')", "", 1)
			println(scan_text)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Can't read .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	err := scanner.Err()
	return err
}

func main() {
	app := makeApp()
	app.Run(os.Args)
}
