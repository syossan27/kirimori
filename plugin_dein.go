package main

import (
	"os"
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

type PluginDein struct{}

func (p *PluginDein) AddLine(vimrcFile *os.File) int {
	f, err := vimlparser.ParseFile(vimrcFile, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(AddDeinVisitor)
	ast.Walk(v, f)

	return v.Line
}

func (p *PluginDein) RemoveLine(vimrcFile *os.File, pluginName string) int {
	f, err := vimlparser.ParseFile(vimrcFile, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(RemoveDeinVisitor)
	v.Name = pluginName
	ast.Walk(v, f)

	return v.Line
}

func (p *PluginDein) ListPlugins(vimrcFile *os.File) []string {
	f, err := vimlparser.ParseFile(vimrcFile, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(ListDeinVisitor)
	ast.Walk(v, f)

	return v.InstallPlugins
}

func (p *PluginDein) Format() string {
	return "call dein#add('%s')"
}
