package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	vimlparser "github.com/haya14busa/go-vimlparser"
	"github.com/mattn/go-colorable"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

const (
	// ExitCodeOK is exit code for OK
	ExitCodeOK = iota
	// ExitCodeError is exit code for Error
	ExitCodeError
)

// Config hold the path and type for vimrc
type Config struct {
	VimrcPath   string
	ManagerType string
}

func (c *Config) Manager() PluginManager {
	switch c.ManagerType {
	case "Vundle":
		return new(PluginVundle)
	case "NeoBundle":
		return new(PluginNeoBundle)
	case "dein.vim":
		return new(PluginDein)
	case "plug.vim":
		return new(PluginPlug)
	}
	fatal("Error: ManagerType is not specified.")
	return nil
}

var (
	opt                    = &vimlparser.ParseOption{}
	homePath, _            = homedir.Dir()
	settingFilePath string = filepath.Join(homePath, ".kirimori.toml")
	stdout                 = colorable.NewColorableStdout()
	stderr                 = colorable.NewColorableStderr()
)

type PluginManager interface {
	AddLine(*os.File) int
	ListPlugins(*os.File) []string
	RemoveLine(*os.File, string) int
	Format() string
}

func success(msg string) {
	fmt.Fprintf(stderr, "\x1b[32m%s\x1b[0m\n", msg)
}

func fatal(msg string) {
	fmt.Fprintf(stderr, "\x1b[31m%s\x1b[0m\n", msg)
	os.Exit(ExitCodeError)
}

func makeApp() *cli.App {
	app := cli.NewApp()

	app.Name = "kirimori"
	app.Usage = "Add Vim Plugin Tool"
	app.Version = "1.0"

	app.Commands = []cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "create setting file",
			Action:  cmdInit,
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add plugin",
			Action:  cmdAdd,
		},
		{
			Name:    "remove",
			Aliases: []string{"r"},
			Usage:   "remove plugin",
			Action:  cmdRemove,
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list plugin",
			Action:  cmdList,
		},
	}

	return app
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func createAddPluginContent(vimrcFile *os.File, format string, pluginName string, addLine int) ([]byte, error) {
	var rows []string
	var index = 1
	scanner := bufio.NewScanner(vimrcFile)
	line := fmt.Sprintf(format, pluginName)
	for scanner.Scan() {
		var scanText = scanner.Text()
		rows = append(rows, scanText)
		if addLine == index {
			rows = append(rows, line)
		}
		index++
	}
	if addLine == 0 {
		rows = append(rows, line)
	}
	vimrcContent := []byte(strings.Join(rows, "\n"))

	err := scanner.Err()
	return vimrcContent, err
}

func createRemovePluginContent(vimrcFile *os.File, pluginName string, removeLine int) ([]byte, error) {
	var rows []string
	var index = 1
	scanner := bufio.NewScanner(vimrcFile)
	for scanner.Scan() {
		var scanText = scanner.Text()
		if index == removeLine {
			index++
			continue
		} else {
			rows = append(rows, scanText)
		}
		index++
	}
	if err := scanner.Err(); err != nil {
		fatal("Error: Can't read .vimrc file.")
	}
	vimrcContent := []byte(strings.Join(rows, "\n"))
	err := scanner.Err()
	return vimrcContent, err
}

func updateVimrc(vimrcFilePath string, vimrcContent []byte) error {
	vimrcFile, err := os.Create(vimrcFilePath)
	if err != nil {
		fatal("Error: Can't open .vimrc file.")
	}
	defer vimrcFile.Close()

	writer := bufio.NewWriter(vimrcFile)
	writer.Write(vimrcContent)
	writer.Flush()
	return err
}

func config() *Config {
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
	return &conf
}

func listPlugin(plugins []string) {
	if len(plugins) == 0 {
		success("Nothing install plugin.")
		return
	}
	for _, plugin := range plugins {
		fmt.Println(plugin)
	}
}

func main() {
	app := makeApp()
	app.Run(os.Args)
}
