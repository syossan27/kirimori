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

// AddPlugVisitor is walker
type AddPlugVisitor struct {
	Line int
	Name string
}

// Visit implement ast.Walker
func (v *AddPlugVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch n := node.(type) {
		case *ast.Ident:
			if n.Name == "plug#begin" {
				v.Line = n.Pos().Line
			}
		case *ast.Excmd:
			if n.Cmd().Name == "Plug" {
				if strings.Contains(n.Command, v.Name) {
					fatal("Error: Installed plugin.")
				}
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
			if n.Cmd().Name == "Plug" {
				if strings.Contains(n.Command, v.Name) {
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
			if n.Cmd().Name == "Plug" {
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
func (p *PluginPlug) AddLine(r io.Reader, pluginName string) int {
	f, err := vimlparser.ParseFile(r, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(AddPlugVisitor)
	v.Name = pluginName
	ast.Walk(v, f)

	return v.Line
}

func (p *PluginPlug) InstallExCmd() {
	cmd := exec.Command("vim", "-c", "PlugInstall", "-c", "qa")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fatal("Error: Fail install plugin.")
	}
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

func (p *PluginPlug) RemoveExCmd() {
	cmd := exec.Command("vim", "-c", "PlugClean", "-c", "qa")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fatal("Error: Fail remove plugin.")
	}
}

// ListPlugin implement PluginManager.ListPlugin
func (p *PluginPlug) ListPlugin(r io.Reader) []string {
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
