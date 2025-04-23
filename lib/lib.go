package main

import "C"
import (
	"encoding/json"
	"fmt"

	"github.com/ystyle/kaf-cli/internal/converter"
	"github.com/ystyle/kaf-cli/internal/core"
	"github.com/ystyle/kaf-cli/internal/model"
	"github.com/ystyle/kaf-cli/pkg/analytics"
)

var (
	secret      string
	measurement string
	version     string
)

//export KafConvert
func KafConvert(params *C.char) int64 {
	var book model.Book
	err := json.Unmarshal([]byte(C.GoString(params)), &book)
	if err != nil {
		return 1
	}
	if err := core.Check(&book, version); err != nil {
		return 2
	}
	analytics.Analytics(version, secret, measurement, book.Format)
	if err := core.Parse(&book); err != nil {
		return 3
	}
	conv := converter.Dispatcher{
		Book: &book,
	}
	if err := conv.Convert(); err != nil {
		return 4
	}
	return 0
}

//export KafPreview
func KafPreview(params *C.char) *C.char {
	var bookArg model.Book
	err := json.Unmarshal([]byte(C.GoString(params)), &bookArg)
	if err != nil {
		return C.CString(fmt.Sprintf("ERROR: 参数错误, %s", err.Error()))
	}
	if err := core.Check(&bookArg, version); err != nil {
		return C.CString(fmt.Sprintf("ERROR: 参数错误, %s", err.Error()))
	}
	if err := core.Parse(&bookArg); err != nil {
		return C.CString(fmt.Sprintf("ERROR: 解析错误, %s", err.Error()))
	}
	bs, _ := json.Marshal(bookArg.SectionList)
	return C.CString(string(bs))

}
func main() {

}
