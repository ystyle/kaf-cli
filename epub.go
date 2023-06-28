package kafcli

import (
	"bytes"
	"fmt"
	"github.com/bmaupin/go-epub"
	"os"
	"path/filepath"
	"time"
)

type EpubConverter struct{}

func (convert EpubConverter) wrapTitle(title, content string) string {
	var buff bytes.Buffer
	buff.WriteString(htmlTitleStart)
	buff.WriteString(title)
	buff.WriteString(htmlTitleEnd)
	buff.WriteString(content)
	return buff.String()
}

func (convert EpubConverter) Build(book Book) error {
	fmt.Println("正在生成epub")
	start := time.Now()
	// 写入样式
	tempDir, err := os.MkdirTemp("", "kaf-cli")
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			panic(fmt.Sprintf("创建临时文件夹失败: %s", err))
		}
	}()

	// Create a ne EPUB
	e := epub.NewEpub(book.Bookname)
	e.SetLang(book.Lang)
	// Set the author
	e.SetAuthor(book.Author)

	pageStylesFile := filepath.Join(tempDir, "page_styles.css")
	var epubcss = cssContent
	var excss string
	if book.LineHeight != "" {
		excss = fmt.Sprintf("line-height: %s;", book.LineHeight)
	}
	if b, _ := isExists(book.Font); b {
		fontfile, _ := e.AddFont(book.Font, "")
		excss += `
font-family: "embedfont";
`
		epubcss += fmt.Sprintf(`
@font-face {
  font-family: "embedfont";
  src: url(%s) format('truetype');
}
`, fontfile)
	}

	err = os.WriteFile(pageStylesFile, []byte(fmt.Sprintf(epubcss, book.Align, book.Bottom, book.Indent, excss)), 0666)
	if err != nil {
		return fmt.Errorf("无法写入样式文件: %w", err)
	}
	css, err := e.AddCSS(pageStylesFile, "")
	if err != nil {
		return fmt.Errorf("无法写入样式文件: %w", err)
	}

	if book.Cover != "" {
		img, err := e.AddImage(book.Cover, filepath.Base(book.Cover))
		if err != nil {
			return fmt.Errorf("添加封面失败: %w", err)
		}
		e.SetCover(img, "")
	}

	for _, section := range book.SectionList {
		if len(section.Sections) > 0 {
			internalFilename, _ := e.AddSection(
				convert.wrapTitle(section.Title, section.Content),
				section.Title,
				"",
				css,
			)
			for _, subsecton := range section.Sections {
				e.AddSubSection(
					internalFilename,
					convert.wrapTitle(subsecton.Title, subsecton.Content),
					subsecton.Title,
					"",
					css,
				)
			}
		} else {
			e.AddSection(convert.wrapTitle(section.Title, section.Content), section.Title, "", css)
		}
	}

	// Write the EPUB
	fmt.Println("正在生成电子书...")
	epubName := book.Out + ".epub"
	err = e.Write(epubName)
	if err != nil {
		// handle error
	}
	// 计算耗时
	end := time.Now().Sub(start)
	fmt.Println("生成EPUB电子书耗时:", end)
	return nil
}
