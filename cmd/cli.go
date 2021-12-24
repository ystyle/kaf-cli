package main

import (
	"github.com/ystyle/kaf-cli"
	"os"
	"strings"
)

var (
	secret      string
	measurement string
	version     string
)

func main() {
	var book *kafcli.Book
	var err error
	if len(os.Args) == 2 && strings.HasSuffix(os.Args[1], ".txt") {
		book, err = kafcli.NewBookSimple(os.Args[1])
		if err != nil {
			panic(err)
		}
	} else {
		book = kafcli.NewBookArgs()
	}
	if err := book.Check(); err != nil {
		panic(err)
	}
	book.ToString()
	if err := book.Parse(); err != nil {
		panic(err)
	}
	book.Convert()
	kafcli.Analytics(version, secret, measurement)
}
