package main

import (
	"fmt"
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
			fmt.Printf("错误: %s\n", err.Error())
			os.Exit(3)
		}
	} else {
		book = kafcli.NewBookArgs()
	}
	if err := book.Check(version); err != nil {
		fmt.Printf("错误: %s\n", err.Error())
		os.Exit(1)
	}
	kafcli.Analytics(version, secret, measurement, book.Format)
	book.ToString()
	if err := book.Parse(); err != nil {
		fmt.Printf("错误: %s\n", err.Error())
		os.Exit(2)
	}
	book.Convert()
}
