package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"encoding/json"

	"text/tabwriter"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/haya14busa/go-vimlparser"
	"github.com/koron/go-dproxy"
	"github.com/kyokomi/emoji"
	"github.com/mattn/go-colorable"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
	"gopkg.in/resty.v1"
)

const (
	// ExitCodeOK is exit code for OK
	ExitCodeOK = iota
	// ExitCodeError is exit code for Error
	ExitCodeError
)

// Manager manage the PluginManager
var pluginManagers = []struct {
	Name    string
	Key     string
	Manager PluginManager
	URL     string
}{
	{
		Name:    "Vundle",
		Key:     "vundle",
		Manager: new(PluginVundle),
		URL:     "https://github.com/VundleVim/Vundle.vim",
	},
	{
		Name:    "NeoBundle",
		Key:     "neobundle",
		Manager: new(PluginNeoBundle),
		URL:     "https://github.com/Shougo/neobundle.vim",
	},
	{
		Name:    "dein.vim",
		Key:     "dein",
		Manager: new(PluginDein),
		URL:     "https://github.com/Shougo/dein.vim",
	},
	{
		Name:    "vim-plug",
		Key:     "plug",
		Manager: new(PluginPlug),
		URL:     "https://github.com/junegunn/vim-plug/",
	},
}

// Config hold the path and type for vimrc
type Config struct {
	VimrcPath   string
	ManagerType string
}

// Manager return PluginManager for the ManagerType
func (c *Config) Manager() PluginManager {
	for _, pm := range pluginManagers {
		if pm.Name == c.ManagerType {
			return pm.Manager
		}
	}
	fatal("Error: ManagerType is not specified.")
	return nil
}

var (
	opt             = &vimlparser.ParseOption{}
	homePath, _     = homedir.Dir()
	settingFilePath = filepath.Join(homePath, ".kirimori.toml")
	stdout          = colorable.NewColorableStdout()
	stderr          = colorable.NewColorableStderr()
)

// PluginManager is common interface of the plugin manages
type PluginManager interface {
	AddLine(io.Reader, string) int
	InstallExCmd()
	ListPlugin(io.Reader) []string
	RemoveLine(io.Reader, string) int
	RemoveExCmd()
	Format(string) string
}

func success(msg string) {
	fmt.Fprintf(stdout, "\x1b[32m%s\x1b[0m\n", msg)
}

func fatal(msg string) {
	fmt.Fprintf(stderr, "\x1b[31m%s\x1b[0m\n", msg)
	os.Exit(ExitCodeError)
}

func makeApp() *cli.App {
	app := cli.NewApp()

	app.Name = "kirimori"
	app.Usage = "Manage Vim Plugin Tool"
	app.Version = "0.0.4"

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
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "search plugin in vimawesome.com",
			Action:  cmdSearch,
		},
		{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "configure",
			Action:  cmdConfig,
		},
	}

	return app
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func createAddPluginContent(r io.Reader, line string, addLine int) ([]byte, error) {
	var rows []string
	var index = 1
	scanner := bufio.NewScanner(r)
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
	b := []byte(strings.Join(rows, "\n") + "\n")

	err := scanner.Err()
	return b, err
}

func createRemovePluginContent(r io.Reader, removeLine int) ([]byte, error) {
	var rows []string
	var index = 1
	scanner := bufio.NewScanner(r)
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
	b := []byte(strings.Join(rows, "\n") + "\n")
	err := scanner.Err()
	return b, err
}

func updateVimrc(filename string, b []byte) error {
	f, err := os.Create(filename)
	if err != nil {
		fatal("Error: Can't open .vimrc file.")
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	writer.Write(b)
	writer.Flush()
	return err
}

func searchPlugin(pluginName string) error {
	// Get search results
	result, err := resty.SetRedirectPolicy(resty.FlexibleRedirectPolicy(20)).
		R().
		SetQueryParams(map[string]string{
			"query": pluginName,
			"page":  "1",
		}).
		Get("http://vimawesome.com/api/plugins")

	// Parse search results for print
	var resultInterface interface{}
	err = json.Unmarshal(result.Body(), &resultInterface)
	proxy := dproxy.New(resultInterface)
	pluginsInterface, err := proxy.M("plugins").Array()
	totalResults, err := proxy.M("total_results").Int64()

	// Set PrintColor
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	// Print search results
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 0, 8, ' ', 0)
	fmt.Printf("Total: %d\n\n", totalResults)
	fmt.Fprintf(w, "%s\t%s\n", red("Plugin Name"), red("Short Description"))
	for _, pluginInterface := range pluginsInterface {
		proxy = dproxy.New(pluginInterface)
		githubURL, _ := proxy.M("github_url").String()
		pluginName := strings.Replace(githubURL, "https://github.com/", "", 1)
		shortDesc, _ := proxy.M("short_desc").String()
		fmt.Fprintf(w, "%s\t%s\n", green(pluginName), yellow(emoji.Sprint(shortDesc)))
	}
	w.Flush()

	return err
}

func config() *Config {
	var conf Config
	if _, err := toml.DecodeFile(settingFilePath, &conf); err != nil {
		fatal("Error: Can't read setting file.")
	}
	conf.VimrcPath = regexp.MustCompile(`^~[/\\]`).ReplaceAllString(conf.VimrcPath, homePath)

	if !fileExists(conf.VimrcPath) {
		fatal("Error: No .vimrc file exists.\n")
	}
	return &conf
}

func printLines(plugins []string) {
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
