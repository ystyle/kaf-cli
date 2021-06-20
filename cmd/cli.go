package main

import (
	"github.com/ystyle/kaf-cli"
)

var (
	secret      string
	measurement string
	version     string
)

func main() {
	kafcli.Convert()
	kafcli.Analytics(version, secret, measurement)
}
