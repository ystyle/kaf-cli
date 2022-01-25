package kafcli

import (
	"bytes"
	"fmt"
	"github.com/leotaku/mobi"
	"golang.org/x/text/language"
	"image"
	"math/rand"
	"os"
	"time"
)

type Azw3Converter struct{}

func (convert Azw3Converter) Build(book Book) error {
	fmt.Println("使用第三方库生成azw3, 不保证所有样式都能正常显示")
	fmt.Println("正在生成azw3...")
	start := time.Now()
	chunks := SectionSliceChunk(book.SectionList, 2700)
	for i, chunk := range chunks {
		index := i + 1
		title := fmt.Sprintf("%s_%d", book.Bookname, index)
		filename := fmt.Sprintf("%s_%d.azw3", book.Out, index)
		if len(chunks) == 1 {
			title = fmt.Sprintf("%s", book.Bookname)
			filename = fmt.Sprintf("%s.azw3", book.Out)
		}
		mb := mobi.Book{
			Title:       title,
			Authors:     []string{book.Author},
			CreatedDate: time.Now(),
			Chapters:    []mobi.Chapter{},
			Language:    language.MustParse(book.Lang),
			UniqueID:    rand.Uint32(),
		}
		css := fmt.Sprintf(cssContent, book.Align, book.Bottom, book.Indent)
		for _, section := range chunk {
			ch := mobi.Chapter{
				Title:  section.Title,
				Chunks: mobi.Chunks(convert.wrapMobiTitle(section.Title, section.Content, book.Align)),
			}
			mb.Chapters = append(mb.Chapters, ch)
		}

		mb.CSSFlows = []string{css}
		if book.Cover != "" {
			f, err := os.Open(book.Cover)
			if err != nil {
				return fmt.Errorf("添加封面失败: %w", err)
			}
			img, _, err := image.Decode(f)
			if err != nil {
				return fmt.Errorf("添加封面失败: %w", err)
			}
			mb.CoverImage = img
		}

		// Convert book to PalmDB database
		db := mb.Realize()

		// Write database to file
		f, _ := os.Create(filename)
		err := db.Write(f)
		if err != nil {
			return fmt.Errorf("保存失败: %w", err)
		}
	}

	fmt.Println("生成azw3电子书耗时:", time.Now().Sub(start))
	return nil
}

func (convert Azw3Converter) wrapMobiTitle(title, content, align string) string {
	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf(mobiTtmlTitleStart, align))
	buff.WriteString(title)
	buff.WriteString(htmlTitleEnd)
	buff.WriteString(content)
	return buff.String()
}

func SectionSliceChunk(s []Section, size int) [][]Section {
	var ret [][]Section
	for size < len(s) {
		// s[:size:size] 表示 len 为 size，cap 也为 size，第二个冒号后的 size 表示 cap
		s, ret = s[size:], append(ret, s[:size:size])
	}
	ret = append(ret, s)
	return ret
}
