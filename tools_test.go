package kafcli

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestLook(t *testing.T) {
	kindlegen, err := exec.LookPath("kindlegen")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(kindlegen)
}

func TestExe(t *testing.T) {
	fmt.Println(os.Args)
	path, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(path)
}
