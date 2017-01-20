package main

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	vimlparser "github.com/haya14busa/go-vimlparser"
	"github.com/haya14busa/go-vimlparser/ast"
)

// AddPlugVisitor is walker
type AddPlugVisitor struct {
	Line int
}

// Visit implement ast.Walker
func (v *AddPlugVisitor) Visit(node ast.Node) (w ast.Visitor) {
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

// RemovePlugVisitor is walker
type RemovePlugVisitor struct {
	Line int
	Name string
}

// Visit implement ast.Walker
func (v *RemovePlugVisitor) Visit(node ast.Node) (w ast.Visitor) {
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

// ListPlugVisitor is walker
type ListPlugVisitor struct {
	InstallPlugins []string
}

// Visit implement ast.Walker
func (v *ListPlugVisitor) Visit(node ast.Node) (w ast.Visitor) {
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

// PluginPlug is implement of PluginManager
type PluginPlug struct{}

// AddLine implement PluginManager.AddLine
func (p *PluginPlug) AddLine(r io.Reader) int {
	f, err := vimlparser.ParseFile(r, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(AddPlugVisitor)
	ast.Walk(v, f)

	return v.Line
}

// RemoveLine implement PluginManager.RemoveLine
func (p *PluginPlug) RemoveLine(r io.Reader, pluginName string) int {
	f, err := vimlparser.ParseFile(r, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(RemovePlugVisitor)
	v.Name = pluginName
	ast.Walk(v, f)

	return v.Line
}

// ListPlugins implement PluginManager.ListPlugins
func (p *PluginPlug) ListPlugins(r io.Reader) []string {
	f, err := vimlparser.ParseFile(r, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(ListPlugVisitor)
	ast.Walk(v, f)

	return v.InstallPlugins
}

// Format implement PluginManager.Format
func (p *PluginPlug) Format(name string) string {
	return fmt.Sprintf("Plug '%s'", name)
}
