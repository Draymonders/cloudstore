package main

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	// dirPath     = "/data/tmp/"
	dirPath = "d:\\tmp\\"
)

func main() {
	var cmd *exec.Cmd
	filepath := dirPath + "/root5d0dd1ec/"
	filestore := dirPath + "1111111.pdf"

	cmd = exec.Command(dirPath+"mergeAll.sh", filepath, filestore)
	// cmd.Run()
	if _, err := cmd.Output(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println(filestore, " has been merge complete")
	}
}
