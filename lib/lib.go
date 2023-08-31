package main

import "C"
import (
	"encoding/json"
	kafcli "github.com/ystyle/kaf-cli"
)

var (
	secret      string
	measurement string
	version     string
)

//export KafConvert
func KafConvert(params *C.char) int64 {
	var bookArg kafcli.Book
	err := json.Unmarshal([]byte(C.GoString(params)), &bookArg)
	if err != nil {
		return 1
	}
	bookArg.SetDefault()
	if err := bookArg.Check(version); err != nil {
		return 2
	}
	kafcli.Analytics(version, secret, measurement, bookArg.Format)
	if err := bookArg.Parse(); err != nil {
		return 3
	}
	bookArg.Convert()
	return 0
}

func main() {

}
