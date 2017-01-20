package main

import (
	"os"
	"strings"
	"unicode/utf8"

	vimlparser "github.com/haya14busa/go-vimlparser"
	"github.com/haya14busa/go-vimlparser/ast"
)

// AddNeoBundleVisitor is walker
type AddNeoBundleVisitor struct {
	Line int
}

// Visit implement ast.Walker
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

// RemoveNeoBundleVisitor is walker
type RemoveNeoBundleVisitor struct {
	Line int
	Name string
}

// Visit implement ast.Walker
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

// ListNeoBundleVisitor is walker
type ListNeoBundleVisitor struct {
	InstallPlugins []string
}

// Visit implement ast.Walker
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

type PluginNeoBundle struct{}

func (p *PluginNeoBundle) AddLine(vimrcFile *os.File) int {
	f, err := vimlparser.ParseFile(vimrcFile, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(AddNeoBundleVisitor)
	ast.Walk(v, f)

	return v.Line
}

func (p *PluginNeoBundle) RemoveLine(vimrcFile *os.File, pluginName string) int {
	f, err := vimlparser.ParseFile(vimrcFile, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(RemoveNeoBundleVisitor)
	v.Name = pluginName
	ast.Walk(v, f)

	return v.Line
}

func (p *PluginNeoBundle) ListPlugins(vimrcFile *os.File) []string {
	f, err := vimlparser.ParseFile(vimrcFile, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(ListNeoBundleVisitor)
	ast.Walk(v, f)

	return v.InstallPlugins
}

func (p *PluginNeoBundle) Format() string {
	return "NeoBundle '%s'"
}
