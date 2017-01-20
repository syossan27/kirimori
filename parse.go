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
	Line    int
	Name    string
	Removed bool
}

// Visit implement ast.Walker
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

// ListDeinVisitor is walker
type ListDeinVisitor struct {
	Added          bool
	InstallPlugins []string
}

// Visit implement ast.Walker
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

func scanAddLineForVundle(vimrcFile *os.File) int {
	f, err := vimlparser.ParseFile(vimrcFile, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(AddVundleVisitor)
	ast.Walk(v, f)

	return v.Line
}

func scanAddLineForNeoBundle(vimrcFile *os.File) int {
	f, err := vimlparser.ParseFile(vimrcFile, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(AddNeoBundleVisitor)
	ast.Walk(v, f)

	return v.Line
}

func scanAddLineForDein(vimrcFile *os.File) int {
	f, err := vimlparser.ParseFile(vimrcFile, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(AddDeinVisitor)
	ast.Walk(v, f)

	return v.Line
}

func scanRemoveLineForVundle(vimrcFile *os.File, pluginName string) int {
	f, err := vimlparser.ParseFile(vimrcFile, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(RemoveVundleVisitor)
	v.Name = pluginName
	ast.Walk(v, f)

	return v.Line
}

func scanRemoveLineForNeoBundle(vimrcFile *os.File, pluginName string) int {
	f, err := vimlparser.ParseFile(vimrcFile, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(RemoveNeoBundleVisitor)
	v.Name = pluginName
	ast.Walk(v, f)

	return v.Line
}

func scanRemoveLineForDein(vimrcFile *os.File, pluginName string) int {
	f, err := vimlparser.ParseFile(vimrcFile, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(RemoveDeinVisitor)
	v.Name = pluginName
	ast.Walk(v, f)

	return v.Line
}

func scanListPluginForVundle(vimrcFile *os.File) []string {
	f, err := vimlparser.ParseFile(vimrcFile, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(ListVundleVisitor)
	ast.Walk(v, f)

	return v.InstallPlugins
}

func scanListPluginForNeoBundle(vimrcFile *os.File) []string {
	f, err := vimlparser.ParseFile(vimrcFile, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(ListNeoBundleVisitor)
	ast.Walk(v, f)

	return v.InstallPlugins
}

func scanListPluginForDein(vimrcFile *os.File) []string {
	f, err := vimlparser.ParseFile(vimrcFile, "", opt)
	if err != nil {
		fatal("Error: Fail parse .vimrc file.")
	}
	v := new(ListDeinVisitor)
	ast.Walk(v, f)

	return v.InstallPlugins
}
