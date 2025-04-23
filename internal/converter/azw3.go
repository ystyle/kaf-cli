package converter

import (
	"bytes"
	"fmt"
	"image"
	"math/rand"
	"os"
	"time"

	"github.com/leotaku/mobi"
	"github.com/ystyle/kaf-cli/internal/model"
	"golang.org/x/text/language"
)

type Azw3Converter struct {
	MobiTtmlTitleStart string // AZW3专属标题标签
	HTMLTitleEnd       string
	CSSContent         string
}

func NewAzw3Converter() *Azw3Converter {
	return &Azw3Converter{
		MobiTtmlTitleStart: `<h3 style="text-align:%s;">`,
		HTMLTitleEnd:       "</h3>",
		CSSContent: `
            .title {text-align: %s}
            .content { margin-bottom: %s; text-indent: %dem; %s }
        `,
	}
}

func (convert Azw3Converter) Build(book model.Book) error {
	fmt.Println("使用第三方库生成azw3, 不保证所有样式都能正常显示")
	fmt.Println("正在生成azw3...")
	start := time.Now()
	chunks := SectionSliceChunk(book.SectionList, 2000)
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
		var excss string
		if book.LineHeight != "" {
			excss = fmt.Sprintf("line-height: %s;", book.LineHeight)
		}
		css := fmt.Sprintf(convert.CSSContent, book.Align, book.Bottom, book.Indent, excss)
		for _, section := range chunk {
			ch := mobi.Chapter{
				Title:  section.Title,
				Chunks: mobi.Chunks(convert.wrapTitle(section.Title, section.Content, book.Align)),
			}
			mb.Chapters = append(mb.Chapters, ch)
			if len(section.Sections) > 0 {
				for _, subsection := range section.Sections {
					ch := mobi.Chapter{
						Title:  subsection.Title,
						Chunks: mobi.Chunks(convert.wrapTitle(subsection.Title, subsection.Content, book.Align)),
					}
					mb.Chapters = append(mb.Chapters, ch)
				}
			}
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

func (convert Azw3Converter) wrapTitle(title, content, align string) string {
	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf(convert.MobiTtmlTitleStart, align))
	buff.WriteString(title)
	buff.WriteString(convert.HTMLTitleEnd)
	buff.WriteString(content)
	return buff.String()
}

func SectionSliceChunk(s []model.Section, size int) [][]model.Section {
	var ret [][]model.Section
	for size < len(s) {
		// s[:size:size] 表示 len 为 size，cap 也为 size，第二个冒号后的 size 表示 cap
		s, ret = s[size:], append(ret, s[:size:size])
	}
	ret = append(ret, s)
	return ret
}
