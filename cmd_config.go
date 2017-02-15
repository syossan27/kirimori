package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/urfave/cli"
)

func cmdConfig(c *cli.Context) error {
	var cmd *exec.Cmd
	var Editor = "vim"
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", fmt.Sprintf("%s %s", Editor, settingFilePath))
	} else {
		cmd = exec.Command("sh", "-c", fmt.Sprintf("%s %s", Editor, settingFilePath))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
