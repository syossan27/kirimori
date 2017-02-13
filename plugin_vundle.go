package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"unicode/utf8"

	vimlparser "github.com/haya14busa/go-vimlparser"
	"github.com/haya14busa/go-vimlparser/ast"
)

// AddVundleVisitor is walker
type AddVundleVisitor struct {
	Line int
	Name string
}

// Visit implement ast.Walker
func (v *AddVundleVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch n := node.(type) {
		case *ast.Ident:
			if n.Name == "vundle#begin" {
				v.Line = n.Pos().Line
			}
		case *ast.Excmd:
			if n.Cmd().Name == "Plugin" {
				if strings.Contains(n.Command, v.Name) {
					fatal("Error: Installed plugin.")
				}
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
			if n.Cmd().Name == "Plugin" {
				if strings.Contains(n.Command, v.Name) {
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
			if n.Cmd().Name == "Plugin" {
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
func (p *PluginVundle) AddLine(r io.Reader, pluginName string) int {
	f, err := vimlparser.ParseFile(r, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(AddVundleVisitor)
	v.Name = pluginName
	ast.Walk(v, f)

	return v.Line
}

// InstallExCmd implement PluginManager.InstallExCmd
func (p *PluginVundle) InstallExCmd() {
	cmd := exec.Command("vim", "-c", "PluginInstall", "-c", "qa")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fatal("Error: Fail install plugin.")
	}
}

// RemoveLine implement PluginManager.RemoveLine
func (p *PluginVundle) RemoveLine(r io.Reader, pluginName string) int {
	f, err := vimlparser.ParseFile(r, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(RemoveVundleVisitor)
	v.Name = pluginName
	ast.Walk(v, f)

	return v.Line
}

// RemoveExCmd implement PluginManager.RemoveExCmd
func (p *PluginVundle) RemoveExCmd() {
	cmd := exec.Command("vim", "-c", "PluginClean", "-c", "qa")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fatal("Error: Fail remove plugin.")
	}
}

// ListPlugin implement PluginManager.ListPlugin
func (p *PluginVundle) ListPlugin(r io.Reader) []string {
	f, err := vimlparser.ParseFile(r, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(ListVundleVisitor)
	ast.Walk(v, f)

	return v.InstallPlugins
}

// Format implement PluginManager.Format
func (p *PluginVundle) Format(name string) string {
	return fmt.Sprintf("Plugin '%s'", name)
}
