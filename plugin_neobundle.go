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

// AddNeoBundleVisitor is walker
type AddNeoBundleVisitor struct {
	Line int
	Name string
}

// Visit implement ast.Walker
func (v *AddNeoBundleVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch n := node.(type) {
		case *ast.Ident:
			if n.Name == "neobundle#begin" {
				v.Line = n.Pos().Line
			}
		case *ast.Excmd:
			name := n.Cmd().Name
			if name == "NeoBundle" || name == "NeoBundleFetch" {
				if strings.Contains(n.Command, v.Name) {
					fatal("Error: Installed plugin.")
				}
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
			name := n.Cmd().Name
			if name == "NeoBundle" || name == "NeoBundleFetch" {
				if strings.Contains(n.Command, v.Name) {
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
			name := n.Cmd().Name
			if name == "NeoBundle" || name == "NeoBundleFetch" {
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

// PluginNeoBundle is implement of PluginManager
type PluginNeoBundle struct{}

// AddLine implement PluginManager.AddLine
func (p *PluginNeoBundle) AddLine(r io.Reader, pluginName string) int {
	f, err := vimlparser.ParseFile(r, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(AddNeoBundleVisitor)
	v.Name = pluginName
	ast.Walk(v, f)

	return v.Line
}

func (p *PluginNeoBundle) InstallExCmd() {
	cmd := exec.Command("vim", "-c", "NeoBundleInstall", "-c", "qa")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fatal("Error: Fail install plugin.")
	}
}

// RemoveLine implement PluginManager.RemoveLine
func (p *PluginNeoBundle) RemoveLine(r io.Reader, pluginName string) int {
	f, err := vimlparser.ParseFile(r, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(RemoveNeoBundleVisitor)
	v.Name = pluginName
	ast.Walk(v, f)

	return v.Line
}

func (p *PluginNeoBundle) RemoveExCmd() {
	// Noop
	// Note: https://github.com/Shougo/neobundle.vim/issues/356
}

// ListPlugin implement PluginManager.ListPlugin
func (p *PluginNeoBundle) ListPlugin(r io.Reader) []string {
	f, err := vimlparser.ParseFile(r, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(ListNeoBundleVisitor)
	ast.Walk(v, f)

	return v.InstallPlugins
}

// Format implement PluginManager.Format
func (p *PluginNeoBundle) Format(name string) string {
	return fmt.Sprintf("NeoBundle '%s'", name)
}
