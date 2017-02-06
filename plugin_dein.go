package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	vimlparser "github.com/haya14busa/go-vimlparser"
	"github.com/haya14busa/go-vimlparser/ast"
)

// AddDeinVisitor is walker
type AddDeinVisitor struct {
	Line int
}

// Visit implement ast.Walker
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

// RemoveDeinVisitor is walker
type RemoveDeinVisitor struct {
	Line  int
	Name  string
	Found bool
}

// Visit implement ast.Walker
func (v *RemoveDeinVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch n := node.(type) {
		case *ast.Ident:
			if n.Name == "dein#add" {
				v.Found = true
			}
		case *ast.BasicLit:
			if v.Found {
				if v.Name != "" && strings.Contains(n.Value, v.Name) {
					v.Line = n.Pos().Line
					v.Found = false
				}
			}
		}
	}
	return v
}

// ListDeinVisitor is walker
type ListDeinVisitor struct {
	Found          bool
	InstallPlugins []string
}

// Visit implement ast.Walker
func (v *ListDeinVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch n := node.(type) {
		case *ast.Ident:
			if n.Name == "dein#add" {
				v.Found = true
			}
		case *ast.BasicLit:
			if v.Found {
				name := strings.Replace(n.Value, "'", "", -1)
				v.InstallPlugins = append(v.InstallPlugins, name)
				v.Found = false
			}
		}
	}
	return v
}

// PluginDein is implement of PluginManager
type PluginDein struct{}

// AddLine implement PluginManager.AddLine
func (p *PluginDein) AddLine(r io.Reader) int {
	f, err := vimlparser.ParseFile(r, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(AddDeinVisitor)
	ast.Walk(v, f)

	return v.Line
}

func (p *PluginDein) InstallExCmd() {
	cmd := exec.Command("vim", "-c", "call dein#install()", "-c", "qa")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fatal("Error: Fail install plugin.")
	}
}

// RemoveLine implement PluginManager.RemoveLine
func (p *PluginDein) RemoveLine(r io.Reader, pluginName string) int {
	f, err := vimlparser.ParseFile(r, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(RemoveDeinVisitor)
	v.Name = pluginName
	ast.Walk(v, f)

	return v.Line
}

func (p *PluginDein) RemoveExCmd() {
	// Noop
	// Note: https://github.com/Shougo/dein.vim/issues/156
}

// ListPlugin implement PluginManager.ListPlugin
func (p *PluginDein) ListPlugin(r io.Reader) []string {
	f, err := vimlparser.ParseFile(r, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(ListDeinVisitor)
	ast.Walk(v, f)

	return v.InstallPlugins
}

// Format implement PluginManager.Format
func (p *PluginDein) Format(name string) string {
	return fmt.Sprintf("call dein#add('%s')", name)
}
