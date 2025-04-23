package kafcli

import (
	"github.com/ystyle/kaf-cli/internal/converter"
	"github.com/ystyle/kaf-cli/internal/core"
	"github.com/ystyle/kaf-cli/internal/model"
)

type Book = model.Book

// NewSimpleBook 是 model.NewBookSimple 的别名，用于快速创建简单书籍实例。
var NewSimpleBook = model.NewBookSimple

func Convert(book *Book) error {
	if err := core.Check(book, "v1.0.0"); err != nil {
		return err
	}
	if err := core.Parse(book); err != nil {
		return err
	}
	conv := converter.Dispatcher{Book: book}
	return conv.Convert()
}
