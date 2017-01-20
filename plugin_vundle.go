package main

import (
	"os"
	"strings"
	"unicode/utf8"

	vimlparser "github.com/haya14busa/go-vimlparser"
	"github.com/haya14busa/go-vimlparser/ast"
)

// AddVundleVisitor is walker
type AddVundleVisitor struct {
	Line int
}

// Visit implement ast.Walker
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

// RemoveVundleVisitor is walker
type RemoveVundleVisitor struct {
	Line int
	Name string
}

// Visit implement ast.Walker
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

// ListVundleVisitor is walker
type ListVundleVisitor struct {
	InstallPlugins []string
}

// Visit implement ast.Walker
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

// PluginVundle is implement of PluginManager
type PluginVundle struct{}

// AddLine implement PluginManager.AddLine
func (p *PluginVundle) AddLine(vimrcFile *os.File) int {
	f, err := vimlparser.ParseFile(vimrcFile, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(AddVundleVisitor)
	ast.Walk(v, f)

	return v.Line
}

// RemoveLine implement PluginManager.RemoveLine
func (p *PluginVundle) RemoveLine(vimrcFile *os.File, pluginName string) int {
	f, err := vimlparser.ParseFile(vimrcFile, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(RemoveVundleVisitor)
	v.Name = pluginName
	ast.Walk(v, f)

	return v.Line
}

// ListPlugins implement PluginManager.ListPlugins
func (p *PluginVundle) ListPlugins(vimrcFile *os.File) []string {
	f, err := vimlparser.ParseFile(vimrcFile, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(ListVundleVisitor)
	ast.Walk(v, f)

	return v.InstallPlugins
}

// Format implement PluginManager.Format
func (p *PluginVundle) Format() string {
	return "Bundle '%s'"
}