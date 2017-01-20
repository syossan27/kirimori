package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func testPlugin(t *testing.T, name, key string) {
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

func TestPlugins(t *testing.T) {
	for _, pm := range pluginManagers {
		testPlugin(t, pm.Name, pm.Key)
	}
}
