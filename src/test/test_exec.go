package main

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	bashLinux   = "/bin/sh"
	bashWindows = "C:\\Program Files\\Git\\cmd\\git.exe"
)

func main() {
	var cmd *exec.Cmd
	filepath := "/d/tmp/root5d0dd1ec/"
	filestore := "/d/tmp/2333.pdf"
	command := "/d/tmp/mergeAll.sh " + filepath + " " + filestore

	cmd = exec.Command(bashWindows, "-c", command)
	cmd.Run()
	if data, err := cmd.Output(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		for v := range data {
			fmt.Println(v)
		}
	}
}
