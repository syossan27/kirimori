package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	vimlparser "github.com/haya14busa/go-vimlparser"
	"github.com/haya14busa/go-vimlparser/ast"
	"github.com/mattn/go-colorable"
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
	opt                      = &vimlparser.ParseOption{}
	home_path, _             = homedir.Dir()
	setting_file_path string = filepath.Join(home_path, "/.kirimori.toml")
	stdout                   = colorable.NewColorableStdout()
	stderr                   = colorable.NewColorableStderr()
)

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

func scanAddLineForVundle(vimrc_file *os.File) int {
	f, err := vimlparser.ParseFile(vimrc_file, "", opt)
	if err != nil {
		fmt.Fprintf(stdout, "\x1b[31m%s\x1b[0m", "Error: Fail parse .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	v := new(AddVundleVisitor)
	ast.Walk(v, f)

	return v.Line
}

func scanAddLineForNeoBundle(vimrc_file *os.File) int {
	f, err := vimlparser.ParseFile(vimrc_file, "", opt)
	if err != nil {
		fmt.Fprintf(stdout, "\x1b[31m%s\x1b[0m", "Error: Fail parse .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	v := new(AddNeoBundleVisitor)
	ast.Walk(v, f)

	return v.Line
}

func scanAddLineForDein(vimrc_file *os.File) int {
	f, err := vimlparser.ParseFile(vimrc_file, "", opt)
	if err != nil {
		fmt.Fprintf(stdout, "\x1b[31m%s\x1b[0m", "Error: Fail parse .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	v := new(AddDeinVisitor)
	ast.Walk(v, f)

	return v.Line
}

func createAddPluginContentForVundle(vimrc_file *os.File, plugin_name string, addLine int) ([]byte, error) {
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

func createAddPluginContentForNeoBundle(vimrc_file *os.File, plugin_name string, addLine int) ([]byte, error) {
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

func createAddPluginContentForDein(vimrc_file *os.File, plugin_name string, addLine int) ([]byte, error) {
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

func scanRemoveLineForVundle(vimrc_file *os.File, plugin_name string) int {
	f, err := vimlparser.ParseFile(vimrc_file, "", opt)
	if err != nil {
		fmt.Fprintf(stdout, "\x1b[31m%s\x1b[0m", "Error: Fail parse .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	v := new(RemoveVundleVisitor)
	v.Name = plugin_name
	ast.Walk(v, f)

	return v.Line
}

func createRemovePluginContentForVundle(vimrc_file *os.File, plugin_name string, removeLine int) ([]byte, error) {
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
		fmt.Fprintf(stdout, "\x1b[31m%s\x1b[0m", "Error: Can't read .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	vimrc_content := []byte(strings.Join(rows, "\n"))
	err := scanner.Err()
	return vimrc_content, err
}

func scanRemoveLineForNeoBundle(vimrc_file *os.File, plugin_name string) int {
	f, err := vimlparser.ParseFile(vimrc_file, "", opt)
	if err != nil {
		fmt.Fprintf(stdout, "\x1b[31m%s\x1b[0m", "Error: Fail parse .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	v := new(RemoveNeoBundleVisitor)
	v.Name = plugin_name
	ast.Walk(v, f)

	return v.Line
}

func createRemovePluginContentForNeoBundle(vimrc_file *os.File, plugin_name string, removeLine int) ([]byte, error) {
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
		fmt.Fprintf(stdout, "\x1b[31m%s\x1b[0m", "Error: Can't read .vimrc file.\n")
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

func scanRemoveLineForDein(vimrc_file *os.File, plugin_name string) int {
	f, err := vimlparser.ParseFile(vimrc_file, "", opt)
	if err != nil {
		fmt.Fprintf(stdout, "\x1b[31m%s\x1b[0m", "Error: Fail parse .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	v := new(RemoveDeinVisitor)
	v.Name = plugin_name
	ast.Walk(v, f)

	return v.Line
}

func createRemovePluginContentForDein(vimrc_file *os.File, plugin_name string, removeLine int) ([]byte, error) {
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
		fmt.Fprintf(stdout, "\x1b[31m%s\x1b[0m", "Error: Can't read .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	vimrc_content := []byte(strings.Join(rows, "\n"))
	err := scanner.Err()
	return vimrc_content, err
}

func updateVimrc(vimrc_file_path string, vimrc_content []byte) error {
	vimrc_file, err := os.Create(vimrc_file_path)
	if err != nil {
		fmt.Fprintf(stdout, "\x1b[31m%s\x1b[0m", "Error: Can't open .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	writer := bufio.NewWriter(vimrc_file)
	writer.Write(vimrc_content)
	writer.Flush()
	return err
}

func scanListPluginForVundle(vimrc_file *os.File) []string {
	f, err := vimlparser.ParseFile(vimrc_file, "", opt)
	if err != nil {
		fmt.Fprintf(stdout, "\x1b[31m%s\x1b[0m", "Error: Fail parse .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	v := new(ListVundleVisitor)
	ast.Walk(v, f)

	return v.InstallPlugins
}

func scanListPluginForNeoBundle(vimrc_file *os.File) []string {
	f, err := vimlparser.ParseFile(vimrc_file, "", opt)
	if err != nil {
		fmt.Fprintf(stdout, "\x1b[31m%s\x1b[0m", "Error: Fail parse .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	v := new(ListNeoBundleVisitor)
	ast.Walk(v, f)

	return v.InstallPlugins
}

func scanListPluginForDein(vimrc_file *os.File) []string {
	f, err := vimlparser.ParseFile(vimrc_file, "", opt)
	if err != nil {
		fmt.Fprintf(stdout, "\x1b[31m%s\x1b[0m", "Error: Fail parse .vimrc file.\n")
		os.Exit(ExitCodeError)
	}
	v := new(ListDeinVisitor)
	ast.Walk(v, f)

	return v.InstallPlugins
}

func listPlugin(plugins []string) {
	if len(plugins) == 0 {
		fmt.Fprintf(stdout, "\x1b[31m%s\x1b[0m", "Nothing install plugin.\n")
	}
	for _, install_plugin := range plugins {
		fmt.Println(install_plugin)
	}
}

func main() {
	app := makeApp()
	app.Run(os.Args)
}
