package main

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	vimlparser "github.com/haya14busa/go-vimlparser"
	"github.com/haya14busa/go-vimlparser/ast"
)

type AddVundleVisitor struct {
	Line int
}

func (v *AddVundleVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch n := node.(type) {
		case *ast.Excmd:
			if n.Cmd().Name == "Bundle" {
				v.Line = n.Pos().Line
			}
		}
	}
	return v
}

type RemoveVundleVisitor struct {
	Line int
	Name string
}

func (v *RemoveVundleVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch n := node.(type) {
		case *ast.Excmd:
			if n.Cmd().Name == "Bundle" {
				if v.Name != "" && strings.Contains(n.Command, v.Name) {
					v.Line = n.Pos().Line
				}
			}
		}
	}
	return v
}

type ListVundleVisitor struct {
	InstallPlugins []string
}

func (v *ListVundleVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch n := node.(type) {
		case *ast.Excmd:
			if n.Cmd().Name == "Bundle" {
				command := n.Command
				start := n.ExArg.Argpos.Offset - n.ExArg.Cmdpos.Offset
				end := utf8.RuneCountInString(n.Command)
				name := strings.Replace(command[start:end], "'", "", -1)
				v.InstallPlugins = append(v.InstallPlugins, name)
			}
		}
	}
	return v
}

type AddNeoBundleVisitor struct {
	Line int
}

func (v *AddNeoBundleVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch n := node.(type) {
		case *ast.Excmd:
			if n.Cmd().Name == "NeoBundle" {
				v.Line = n.Pos().Line
			}
		}
	}
	return v
}

type RemoveNeoBundleVisitor struct {
	Line int
	Name string
}

func (v *RemoveNeoBundleVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch n := node.(type) {
		case *ast.Excmd:
			if n.Cmd().Name == "NeoBundle" {
				if v.Name != "" && strings.Contains(n.Command, v.Name) {
					v.Line = n.Pos().Line
				}
			}
		}
	}
	return v
}

type ListNeoBundleVisitor struct {
	InstallPlugins []string
}

func (v *ListNeoBundleVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch n := node.(type) {
		case *ast.Excmd:
			if n.Cmd().Name == "NeoBundle" {
				command := n.Command
				start := n.ExArg.Argpos.Offset - n.ExArg.Cmdpos.Offset
				end := utf8.RuneCountInString(n.Command)
				name := strings.Replace(command[start:end], "'", "", -1)
				v.InstallPlugins = append(v.InstallPlugins, name)
			}
		}
	}
	return v
}

type AddDeinVisitor struct {
	Line int
}

func (v *AddDeinVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch n := node.(type) {
		case *ast.Ident:
			if n.Name == "dein#add" {
				v.Line = n.Pos().Line
			}
		}
	}
	return v
}

type RemoveDeinVisitor struct {
	Line    int
	Name    string
	Removed bool
}

func (v *RemoveDeinVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch n := node.(type) {
		case *ast.Ident:
			if n.Name == "dein#add" {
				v.Removed = true
			}
		case *ast.BasicLit:
			if v.Removed {
				if v.Name != "" && strings.Contains(n.Value, v.Name) {
					v.Line = n.Pos().Line
					v.Removed = false
				}
			}
		}
	}
	return v
}

type ListDeinVisitor struct {
	Added          bool
	InstallPlugins []string
}

func (v *ListDeinVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch n := node.(type) {
		case *ast.Ident:
			if n.Name == "dein#add" {
				v.Added = true
			}
		case *ast.BasicLit:
			if v.Added {
				name := strings.Replace(n.Value, "'", "", -1)
				v.InstallPlugins = append(v.InstallPlugins, name)
				v.Added = false
			}
		}
	}
	return v
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
