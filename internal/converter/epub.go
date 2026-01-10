package converter

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/go-shiori/go-epub"
	"github.com/ystyle/kaf-cli/internal/model"
	"github.com/ystyle/kaf-cli/internal/utils"
)

type EpubConverter struct {
	HTMLPStart       string // EPUB专属段落标签
	HTMLPEnd         string
	HTMLTitleStart   string
	HTMLTitleEnd     string
	HTMLVolumeStart  string
	HTMLVolumeEnd    string
	CSSContent       string
}

func NewEpubConverter() *EpubConverter {
	return &EpubConverter{
		HTMLPStart:      `<p class="content">`,
		HTMLPEnd:        "</p>",
		HTMLTitleStart:  `<h3 class="title">`,
		HTMLTitleEnd:    "</h3>",
		HTMLVolumeStart: `<h2 class="volume">`,
		HTMLVolumeEnd:   "</h2>",
		CSSContent: `
            h2.volume {
                text-align: center;
                font-size: 2.2em;
                margin: 1.5em 0 1em 0;
                padding: 0.5em 0;
                border-top: 3px double #666;
                border-bottom: 3px double #666;
                background: linear-gradient(to bottom, #f9f9f9, #ffffff);
                font-weight: bold;
            }
            h3.title {
                text-align: %s;
                font-size: 1.8em;
                margin: 1em 0;
                border-bottom: 2px solid #ccc;
            }
            h3.title span.chapter-number {
                display: block;
                font-size: 0.65em;
            }
            .content { margin-bottom: %s; text-indent: %dem; %s }
        `,
	}
}

func (convert EpubConverter) wrapTitle(title, content string, separateNumber bool, isVolume bool) string {
	var buff bytes.Buffer

	if isVolume {
		// 卷名使用专门的样式
		buff.WriteString(convert.HTMLVolumeStart)
		buff.WriteString(title)
		buff.WriteString(convert.HTMLVolumeEnd)
		buff.WriteString(content)
		return buff.String()
	}

	if separateNumber {
		// 尝试分离章节序号和标题
		number, text := parseChapterTitle(title)
		buff.WriteString(convert.HTMLTitleStart)
		if number != "" {
			// 有序号，将序号和标题分开显示
			buff.WriteString(fmt.Sprintf(`<span class="chapter-number">%s</span>`, number))
			if text != "" {
				buff.WriteString(text)
			}
		} else {
			// 无序号，直接显示标题
			buff.WriteString(title)
		}
		buff.WriteString(convert.HTMLTitleEnd)
	} else {
		// 不分离，直接显示标题
		buff.WriteString(convert.HTMLTitleStart)
		buff.WriteString(title)
		buff.WriteString(convert.HTMLTitleEnd)
	}
	buff.WriteString(content)
	return buff.String()
}

// parseChapterTitle 解析章节标题，返回序号和标题
// 支持的格式：
//   "第一章 标题" -> number="第一章", text="标题"
//   "第1章 标题" -> number="第1章", text="标题"
//   "1. 标题" -> number="1.", text="标题"
//   "一、标题" -> number="一、", text="标题"
//   "引子" -> number="引子", text=""
//   "卷名" -> number="", text="卷名"（没有匹配到序号）
func parseChapterTitle(title string) (number, text string) {
	// 匹配 "第X章/回/节/集" 格式
	re := regexp.MustCompile(`^(第[0-9一二三四五六七八九十零〇百千两 ]+[章回节集])\s*(.*)$`)
	if matches := re.FindStringSubmatch(title); matches != nil {
		return matches[1], matches[2]
	}

	// 匹配 "数字." 或 "数字、" 格式（使用字符串拼接来支持中文顿号）
	re = regexp.MustCompile(`^(\d+[.` + string(rune(0x3001)) + `])\s*(.*)$`)
	if matches := re.FindStringSubmatch(title); matches != nil {
		return matches[1], matches[2]
	}

	// 匹配 "中文数字、" 格式
	re = regexp.MustCompile(`^([一二三四五六七八九十]+[.` + string(rune(0x3001)) + `])\s*(.*)$`)
	if matches := re.FindStringSubmatch(title); matches != nil {
		return matches[1], matches[2]
	}

	// 匹配特殊章节名（引子、楔子、序章等）
	re = regexp.MustCompile(`^(引子|楔子|序章|最终章|完本感言|番外)\s*(.*)$`)
	if matches := re.FindStringSubmatch(title); matches != nil {
		if matches[2] != "" {
			return matches[1], matches[2]
		}
		return matches[1], ""
	}

	// 没有匹配到序号格式，返回空序号
	return "", title
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
			// 这是一个卷（包含子章节）
			internalFilename, _ := e.AddSection(
				convert.wrapTitle(section.Title, section.Content, book.SeparateChapterNumber, true),
				section.Title,
				"",
				css,
			)
			for _, subsecton := range section.Sections {
				e.AddSubSection(
					internalFilename,
					convert.wrapTitle(subsecton.Title, subsecton.Content, book.SeparateChapterNumber, false),
					subsecton.Title,
					"",
					css,
				)
			}
		} else {
			e.AddSection(convert.wrapTitle(section.Title, section.Content, book.SeparateChapterNumber, false), section.Title, "", css)
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
