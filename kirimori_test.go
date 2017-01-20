package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testAddPlugin(t *testing.T, name, key string) {
	settingFilePath = filepath.Join("testdir", key+".toml")
	conf := config()
	if conf.ManagerType != name {
		t.Fatalf("expected %v but %v", name, conf.ManagerType)
	}

	manager := conf.Manager()

	f, err := os.Open(conf.VimrcPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	line := conf.Manager().AddLine(f)

	_, err = f.Seek(0, 0)
	if err != nil {
		t.Fatal(err)
	}

	b, err := createAddPluginContent(f, manager.Format("mattn/emmet-vim"), line)
	if err != nil {
		t.Fatal(err)
	}
	fb, err := ioutil.ReadFile(filepath.Join("testdir", key+".vimrc.add"))
	if err != nil {
		t.Fatal(err)
	}
	expected := string(fb)
	got := string(b)
	if expected != got {
		t.Fatalf("expected %s but %s", expected, got)
	}
}

func testListPlugin(t *testing.T, name, key string) {
	settingFilePath = filepath.Join("testdir", key+".toml")
	conf := config()
	if conf.ManagerType != name {
		t.Fatalf("expected %v but %v", name, conf.ManagerType)
	}

	f, err := os.Open(conf.VimrcPath + ".add")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	plugins := conf.Manager().ListPlugin(f)

	fb, err := ioutil.ReadFile(filepath.Join("testdir", key+".list"))
	if err != nil {
		t.Fatal(err)
	}
	expected := strings.TrimSpace(string(fb))
	got := strings.Join(plugins, "\n")
	if expected != got {
		t.Fatalf("expected %s but %s", expected, got)
	}
}

func testRemovePlugin(t *testing.T, name, key string) {
	settingFilePath = filepath.Join("testdir", key+".toml")
	conf := config()
	if conf.ManagerType != name {
		t.Fatalf("expected %v but %v", name, conf.ManagerType)
	}

	f, err := os.Open(conf.VimrcPath + ".add")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	line := conf.Manager().RemoveLine(f, "mattn/emmet-vim")

	_, err = f.Seek(0, 0)
	if err != nil {
		t.Fatal(err)
	}

	b, err := createRemovePluginContent(f, line)
	if err != nil {
		t.Fatal(err)
	}
	fb, err := ioutil.ReadFile(filepath.Join("testdir", key+".vimrc"))
	if err != nil {
		t.Fatal(err)
	}
	expected := string(fb)
	got := string(b)
	if expected != got {
		t.Fatalf("expected %s but %s", expected, got)
	}
}

func TestPlugins(t *testing.T) {
	for _, pm := range pluginManagers {
		testAddPlugin(t, pm.Name, pm.Key)
		testListPlugin(t, pm.Name, pm.Key)
		testRemovePlugin(t, pm.Name, pm.Key)
	}
}
