package converter

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/go-shiori/go-epub"
	"github.com/ystyle/kaf-cli/internal/model"
	"github.com/ystyle/kaf-cli/internal/utils"
)

type EpubConverter struct {
	HTMLPStart     string // EPUB专属段落标签
	HTMLPEnd       string
	HTMLTitleStart string
	HTMLTitleEnd   string
	CSSContent     string
}

func NewEpubConverter() *EpubConverter {
	return &EpubConverter{
		HTMLPStart:     `<p class="content">`,
		HTMLPEnd:       "</p>",
		HTMLTitleStart: `<h3 class="title">`,
		HTMLTitleEnd:   "</h3>",
		CSSContent: `
            .title {text-align: %s}
            .content { margin-bottom: %s; text-indent: %dem; %s }
        `,
	}
}

func (convert EpubConverter) wrapTitle(title, content string) string {
	var buff bytes.Buffer
	buff.WriteString(convert.HTMLTitleStart)
	buff.WriteString(title)
	buff.WriteString(convert.HTMLTitleEnd)
	buff.WriteString(content)
	return buff.String()
}

func (convert EpubConverter) Build(book model.Book) error {
	log.Default().SetOutput(io.Discard)
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
	e, err := epub.NewEpub(book.Bookname)
	if err != nil {
		return fmt.Errorf("创建小说文件失败")
	}
	e.SetLang(book.Lang)
	// Set the author
	e.SetAuthor(book.Author)

	pageStylesFile := filepath.Join(tempDir, "page_styles.css")
	var epubcss = convert.CSSContent
	var excss string
	if book.LineHeight != "" {
		excss = fmt.Sprintf("line-height: %s;", book.LineHeight)
	}
	if b, _ := utils.IsExists(book.Font); b {
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

	err = os.WriteFile(pageStylesFile, fmt.Appendf(nil, epubcss, book.Align, book.Bottom, book.Indent, excss), 0666)
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
