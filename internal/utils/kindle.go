package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func LookKindlegen() string {
	command := "kindlegen"
	if runtime.GOOS == "windows" {
		command = "kindlegen.exe"
	}
	kindlegen, err := exec.LookPath(command)
	if err != nil {
		currentDir, err := os.Executable()
		if err != nil {
			return ""
		}
		kindlegen = filepath.Join(filepath.Dir(currentDir), command)
		if exist, _ := IsExists(kindlegen); !exist {
			return ""
		}
		fmt.Println("kindlegen: ", kindlegen)
	}
	return kindlegen
}
