package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/BurntSushi/toml"
	vimlparser "github.com/haya14busa/go-vimlparser"
	"github.com/haya14busa/go-vimlparser/ast"
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

var (
	opt = &vimlparser.ParseOption{}

	addPlugin      = false
	removePlugin   = false
	installPlugins []string
	addLine        int
	removeLine     int
	plugin_name    string
)

type AddVundleVisitor struct {
}

func (v *AddVundleVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch n := node.(type) {
		case *ast.Excmd:
			if n.Cmd().Name == "Bundle" {
				addLine = n.Pos().Line
			}
		}
	}
	return v
}

type RemoveVundleVisitor struct {
}

func (v *RemoveVundleVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch n := node.(type) {
		case *ast.Excmd:
			if n.Cmd().Name == "Bundle" {
				if strings.Contains(n.Command, plugin_name) {
					removeLine = n.Pos().Line
				}
			}
		}
	}
	return v
}

type ListVundleVisitor struct {
}

func (v *ListVundleVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch n := node.(type) {
		case *ast.Excmd:
			if n.Cmd().Name == "Bundle" {
				command := n.Command
				start := n.ExArg.Argpos.Offset - n.ExArg.Cmdpos.Offset
				end := utf8.RuneCountInString(n.Command)
				plugin_name = strings.Replace(command[start:end], "'", "", -1)
				installPlugins = append(installPlugins, plugin_name)
			}
		}
	}
	return v
}

type AddNeoBundleVisitor struct {
}

func (v *AddNeoBundleVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch n := node.(type) {
		case *ast.Excmd:
			if n.Cmd().Name == "NeoBundle" {
				addLine = n.Pos().Line
			}
		}
	}
	return v
}

type RemoveNeoBundleVisitor struct {
}

func (v *RemoveNeoBundleVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch n := node.(type) {
		case *ast.Excmd:
			if n.Cmd().Name == "NeoBundle" {
				if strings.Contains(n.Command, plugin_name) {
					removeLine = n.Pos().Line
				}
			}
		}
	}
	return v
}

type ListNeoBundleVisitor struct {
}

func (v *ListNeoBundleVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch n := node.(type) {
		case *ast.Excmd:
			if n.Cmd().Name == "NeoBundle" {
				command := n.Command
				start := n.ExArg.Argpos.Offset - n.ExArg.Cmdpos.Offset
				end := utf8.RuneCountInString(n.Command)
				plugin_name = strings.Replace(command[start:end], "'", "", -1)
				installPlugins = append(installPlugins, plugin_name)
			}
		}
	}
	return v
}

type AddDeinVisitor struct {
}

func (v *AddDeinVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch n := node.(type) {
		case *ast.Ident:
			if n.Name == "dein#add" {
				addLine = n.Pos().Line
			}
		}
	}
	return v
}

type RemoveDeinVisitor struct {
}

func (v *RemoveDeinVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch n := node.(type) {
		case *ast.Ident:
			if n.Name == "dein#add" {
				removePlugin = true
			}
		case *ast.BasicLit:
			if removePlugin {
				if strings.Contains(n.Value, plugin_name) {
					removeLine = n.Pos().Line
					removePlugin = false
				}
			}
		}
	}
	return v
}

type ListDeinVisitor struct {
}

func (v *ListDeinVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch n := node.(type) {
		case *ast.Ident:
			if n.Name == "dein#add" {
				addPlugin = true
			}
		case *ast.BasicLit:
			if addPlugin {
				plugin_name = strings.Replace(n.Value, "'", "", -1)
				installPlugins = append(installPlugins, plugin_name)
				addPlugin = false
			}
		}
	}
	return v
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
				}

				var vimrc_file_path string
				fmt.Println("Type your .vimrc path. (default: ~/.vimrc)")
				fmt.Print("> ")
				fmt.Scanln(&vimrc_file_path)
				if vimrc_file_path == "" {
					vimrc_file_path = "~/.vimrc"
				}

				var manager_type string
				fmt.Println("Choose a your vim bundle plugin. (default: 1)")
				fmt.Println("\t1) Vundle")
				fmt.Println("\t2) NeoBundle")
				fmt.Println("\t3) dein.vim")
				fmt.Print("Type number > ")
				fmt.Scanln(&manager_type)
				switch manager_type {
				case "1":
					manager_type = "Vundle"
				case "2":
					manager_type = "NeoBundle"
				case "3":
					manager_type = "dein.vim"
				default:
					manager_type = "Vundle"
				}

				file, err := os.Create(setting_file_path)
				if err != nil {
					fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Setting file exist.\n")
					os.Exit(ExitCodeError)
				}
				defer file.Close()

				writer := bufio.NewWriter(file)
				writer.Write(createSettingFileText(vimrc_file_path, manager_type))
				writer.Flush()

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
				plugin_name = c.Args().First()
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
						scanAddLineForVundle(vimrc_file)

						_, err := vimrc_file.Seek(0, 0)
						if err != nil {
							fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail change file offset.\n")
							os.Exit(ExitCodeError)
						}

						vimrc_content, err := createAddPluginContentForVundle(vimrc_file, plugin_name)
						if err != nil {
							fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Can't read .vimrc file.\n")
							os.Exit(ExitCodeError)
						}
						if err := updateVimrc(conf.VimrcPath, vimrc_content); err != nil {
							fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail add plugin.\n")
							os.Exit(ExitCodeError)
						}
					case "NeoBundle":
						scanAddLineForNeoBundle(vimrc_file)

						_, err := vimrc_file.Seek(0, 0)
						if err != nil {
							fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail change file offset.\n")
							os.Exit(ExitCodeError)
						}

						vimrc_content, err := createAddPluginContentForNeoBundle(vimrc_file, plugin_name)
						if err != nil {
							fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Can't read .vimrc file.\n")
							os.Exit(ExitCodeError)
						}
						if err := updateVimrc(conf.VimrcPath, vimrc_content); err != nil {
							fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail add plugin.\n")
							os.Exit(ExitCodeError)
						}
					case "dein.vim":
						scanAddLineForDein(vimrc_file)

						_, err := vimrc_file.Seek(0, 0)
						if err != nil {
							fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail change file offset.\n")
							os.Exit(ExitCodeError)
						}

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
				plugin_name = c.Args().First()
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
						scanRemoveLineForVundle(vimrc_file)

						_, err := vimrc_file.Seek(0, 0)
						if err != nil {
							fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail change file offset.\n")
							os.Exit(ExitCodeError)
						}

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
						scanRemoveLineForNeoBundle(vimrc_file)

						_, err := vimrc_file.Seek(0, 0)
						if err != nil {
							fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail change file offset.\n")
							os.Exit(ExitCodeError)
						}

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
						scanRemoveLineForDein(vimrc_file)

						_, err := vimrc_file.Seek(0, 0)
						if err != nil {
							fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail change file offset.\n")
							os.Exit(ExitCodeError)
						}

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
						scanListPluginForVundle(vimrc_file)
						listPlugin()
					case "NeoBundle":
						scanListPluginForNeoBundle(vimrc_file)
						listPlugin()
					case "dein.vim":
						scanListPluginForDein(vimrc_file)
						listPlugin()
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

func createSettingFileText(vimrc_file_path string, manager_type string) []byte {
	body := []string{
		"VimrcPath = \"" + vimrc_file_path + "\"",
		"ManagerType = \"" + manager_type + "\"",
	}

	return []byte(strings.Join(body, "\n"))
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func scanAddLineForVundle(vimrc_file *os.File) {
	f, err := vimlparser.ParseFile(vimrc_file, "", opt)
	if err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail parse .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	v := new(AddVundleVisitor)
	ast.Walk(v, f)

	return
}

func scanAddLineForNeoBundle(vimrc_file *os.File) {
	f, err := vimlparser.ParseFile(vimrc_file, "", opt)
	if err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail parse .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	v := new(AddNeoBundleVisitor)
	ast.Walk(v, f)

	return
}

func scanAddLineForDein(vimrc_file *os.File) {
	f, err := vimlparser.ParseFile(vimrc_file, "", opt)
	if err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail parse .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	v := new(AddDeinVisitor)
	ast.Walk(v, f)

	return
}

func createAddPluginContentForVundle(vimrc_file *os.File, plugin_name string) ([]byte, error) {
	var rows []string
	var index int = 1
	scanner := bufio.NewScanner(vimrc_file)
	for scanner.Scan() {
		var scan_text = scanner.Text()
		rows = append(rows, scan_text)
		if addLine == index {
			rows = append(rows, "Bundle '"+plugin_name+"'")
		}
		index++
	}
	if addLine == 0 {
		rows = append(rows, "Bundle '"+plugin_name+"'")
	}
	vimrc_content := []byte(strings.Join(rows, "\n"))

	err := scanner.Err()
	return vimrc_content, err
}

func createAddPluginContentForNeoBundle(vimrc_file *os.File, plugin_name string) ([]byte, error) {
	var rows []string
	var index int = 1
	scanner := bufio.NewScanner(vimrc_file)
	for scanner.Scan() {
		var scan_text = scanner.Text()
		rows = append(rows, scan_text)
		if addLine == index {
			rows = append(rows, "NeoBundle '"+plugin_name+"'")
		}
		index++
	}
	if addLine == 0 {
		rows = append(rows, "NeoBundle '"+plugin_name+"'")
	}
	vimrc_content := []byte(strings.Join(rows, "\n"))

	err := scanner.Err()
	return vimrc_content, err
}

func createAddPluginContentForDein(vimrc_file *os.File, plugin_name string) ([]byte, error) {
	var rows []string
	var index int = 1
	scanner := bufio.NewScanner(vimrc_file)
	for scanner.Scan() {
		var scan_text = scanner.Text()
		rows = append(rows, scan_text)
		if addLine == index {
			rows = append(rows, "call dein#add('"+plugin_name+"')")
		}
		index++
	}
	if addLine == 0 {
		rows = append(rows, "call dein#add('"+plugin_name+"')")
	}
	vimrc_content := []byte(strings.Join(rows, "\n"))

	err := scanner.Err()
	return vimrc_content, err
}

func scanRemoveLineForVundle(vimrc_file *os.File) {
	f, err := vimlparser.ParseFile(vimrc_file, "", opt)
	if err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail parse .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	v := new(RemoveVundleVisitor)
	ast.Walk(v, f)

	return
}

func createRemovePluginContentForVundle(vimrc_file *os.File, plugin_name string) ([]byte, error) {
	var rows []string
	var index int = 1
	scanner := bufio.NewScanner(vimrc_file)
	for scanner.Scan() {
		var scan_text = scanner.Text()
		if index == removeLine {
			index++
			continue
		} else {
			rows = append(rows, scan_text)
		}
		index++
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Can't read .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	vimrc_content := []byte(strings.Join(rows, "\n"))
	err := scanner.Err()
	return vimrc_content, err
}

func scanRemoveLineForNeoBundle(vimrc_file *os.File) {
	f, err := vimlparser.ParseFile(vimrc_file, "", opt)
	if err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail parse .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	v := new(RemoveNeoBundleVisitor)
	ast.Walk(v, f)

	return
}

func createRemovePluginContentForNeoBundle(vimrc_file *os.File, plugin_name string) ([]byte, error) {
	var rows []string
	var index int = 1
	scanner := bufio.NewScanner(vimrc_file)
	for scanner.Scan() {
		var scan_text = scanner.Text()
		if index == removeLine {
			index++
			continue
		} else {
			rows = append(rows, scan_text)
		}
		index++
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

func scanRemoveLineForDein(vimrc_file *os.File) {
	f, err := vimlparser.ParseFile(vimrc_file, "", opt)
	if err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail parse .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	v := new(RemoveDeinVisitor)
	ast.Walk(v, f)

	return
}

func createRemovePluginContentForDein(vimrc_file *os.File, plugin_name string) ([]byte, error) {
	var rows []string
	var index int = 1
	scanner := bufio.NewScanner(vimrc_file)
	for scanner.Scan() {
		var scan_text = scanner.Text()
		if index == removeLine {
			index++
			continue
		} else {
			rows = append(rows, scan_text)
		}
		index++
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

func scanListPluginForVundle(vimrc_file *os.File) {
	f, err := vimlparser.ParseFile(vimrc_file, "", opt)
	if err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail parse .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	v := new(ListVundleVisitor)
	ast.Walk(v, f)

	return
}

func scanListPluginForNeoBundle(vimrc_file *os.File) {
	f, err := vimlparser.ParseFile(vimrc_file, "", opt)
	if err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail parse .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	v := new(ListNeoBundleVisitor)
	ast.Walk(v, f)

	return
}

func scanListPluginForDein(vimrc_file *os.File) {
	f, err := vimlparser.ParseFile(vimrc_file, "", opt)
	if err != nil {
		fmt.Printf("\x1b[31m%s\x1b[0m", "Error: Fail parse .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	v := new(ListDeinVisitor)
	ast.Walk(v, f)

	return
}

func listPlugin() {
	if len(installPlugins) == 0 {
		fmt.Printf("\x1b[31m%s\x1b[0m", "Nothing install plugin.\n")
	}
	for _, install_plugin := range installPlugins {
		println(install_plugin)
	}
}

func main() {
	app := makeApp()
	app.Run(os.Args)
}
