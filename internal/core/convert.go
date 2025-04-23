package core

import (
	"errors"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ystyle/kaf-cli/internal/model"
	"github.com/ystyle/kaf-cli/internal/utils"
)

func Check(book *model.Book, version string) error {
	book.Version = version
	if err := validateInput(book); err != nil {
		return err
	}
	parseBookInfoFromFilename(book)
	setDefaultValues(book)
	if err := handleCover(book); err != nil {
		return err
	}
	if err := compileRegex(book); err != nil {
		return err
	}
	return nil
}

func validateInput(book *model.Book) error {
	if !strings.HasSuffix(book.Filename, ".txt") {
		return errors.New("不是txt文件")
	}
	return nil
}

func parseBookInfoFromFilename(book *model.Book) {
	reg, _ := regexp.Compile(`《(.*)》.*作者[：:](.*).txt`)
	if reg.MatchString(book.Filename) {
		group := reg.FindAllStringSubmatch(book.Filename, -1)
		if len(group) == 1 && len(group[0]) >= 3 {
			if book.Bookname == "" {
				book.Bookname = group[0][1]
			}
			if book.Author == "" || book.Author == "YSTYLE" {
				book.Author = group[0][2]
			}
		}
	}
	if book.Bookname == "" {
		book.Bookname = strings.Split(filepath.Base(book.Filename), ".")[0]
	}
}

func setDefaultValues(book *model.Book) {
	if book.Out == "" {
		book.Out = book.Bookname
	}
	book.Lang = utils.ParseLang(book.Lang)
}

func handleCover(book *model.Book) error {
	switch book.Cover {
	case "none":
		book.Cover = ""
	case "gen", "orly":
		cover, err := utils.GenCover(book.Bookname, book.Author, book.CoverOrlyColor, book.CoverOrlyIdx)
		if err != nil {
			return err
		}
		book.Cover = cover
	default:
		if exists, _ := utils.IsExists(book.Cover); !exists {
			book.Cover = ""
		}
	}
	return nil
}

func compileRegex(book *model.Book) error {
	if book.Match == "" {
		book.Match = model.DefaultMatchTips
	}
	reg, err := regexp.Compile(book.Match)
	if err != nil {
		return fmt.Errorf("生成匹配规则出错: %s\n%s\n", book.Match, err.Error())
	}
	book.Reg = reg

	reg2, err := regexp.Compile(book.VolumeMatch)
	if err != nil {
		return fmt.Errorf("生成匹配规则出错: %s\n%s\n", book.VolumeMatch, err.Error())
	}
	book.VolumeReg = reg2
	return nil
}
