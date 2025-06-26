package converter

import (
	"fmt"
	"time"

	"github.com/ystyle/kaf-cli/internal/model"
	"github.com/ystyle/kaf-cli/internal/utils"
)

type Dispatcher struct {
	Book *model.Book
}

func (d *Dispatcher) Convert() error {
	start := time.Now()
	// 解析文本
	fmt.Println()
	// 判断要生成的格式
	var isEpub, isMobi, isAzw3 bool
	switch d.Book.Format {
	case "epub":
		isEpub = true
	case "mobi":
		isEpub = true
		isMobi = true
	case "azw3":
		isAzw3 = true
	default:
		isEpub = true
		isMobi = true
		isAzw3 = true
	}

	hasKinldegen := utils.LookKindlegen()
	if d.Book.Format == "mobi" && hasKinldegen == "" {
		isEpub = false
	}

	var convert Converter
	// 生成epub
	if isEpub {
		convert = NewEpubConverter()
		convert.Build(*d.Book)
		fmt.Println()
	}
	// 生成azw3格式
	if isAzw3 {
		convert = NewAzw3Converter()
		// 生成kindle格式
		convert.Build(*d.Book)
	}
	// 生成mobi格式
	if isMobi {
		if hasKinldegen == "" {
			convert = NewMobiConverter()
			convert.Build(*d.Book)
		} else {
			ConverToMobi(fmt.Sprintf("%s.epub", d.Book.Out), d.Book.Lang)
		}
	}
	end := time.Now().Sub(start)
	fmt.Println("\n转换完成! 总耗时:", end)

	return nil
}
