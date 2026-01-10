package core

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ystyle/kaf-cli/internal/model"
	"github.com/ystyle/kaf-cli/internal/utils"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func readBuffer(book *model.Book, filename string) *bufio.Reader {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("读取文件出错: ", err.Error())
		os.Exit(1)
	}
	temBuf := bufio.NewReader(f)
	bs, _ := temBuf.Peek(1024)
	encodig, encodename, _ := charset.DetermineEncoding(bs, "text/plain")
	if encodename != "utf-8" {
		f.Seek(0, 0)
		bs, err := io.ReadAll(f)
		if err != nil {
			fmt.Println("读取文件出错: ", err.Error())
			os.Exit(1)
		}
		var buf bytes.Buffer
		book.Decoder = encodig.NewDecoder()
		if encodename == "windows-1252" {
			book.Decoder = simplifiedchinese.GB18030.NewDecoder()
		}
		bs, _, _ = transform.Bytes(book.Decoder, bs)
		buf.Write(bs)
		return bufio.NewReader(&buf)
	} else {
		f.Seek(0, 0)
		buf := bufio.NewReader(f)
		return buf
	}
}

// sanitizeHTMLTags 智能处理 HTML 标签
// 保留 epub 支持的标签（如 img, br, hr 等），转义其他标签
func sanitizeHTMLTags(line string) string {
	// 快速检查：如果没有 < 符号，直接返回
	if !strings.Contains(line, "<") {
		return line
	}

	// 定义允许的标签（EPUB 3.0 支持的常用标签）
	// 只保留最常用的标签以提升性能
	allowedTags := []string{
		"img", "br", "hr", // 单标签（最常用）
		"p", "span", "div",
		"b", "i", "u", "s", "strong", "em",
		"a", "table", "tr", "td", "th",
	}

	// 构建正则表达式来匹配允许的标签
	var tagPatterns []string
	for _, tag := range allowedTags {
		tagPatterns = append(tagPatterns, fmt.Sprintf(`<%s\b[^>]*>`, tag))
		tagPatterns = append(tagPatterns, fmt.Sprintf(`</%s>`, tag))
		tagPatterns = append(tagPatterns, fmt.Sprintf(`<%s\b[^>]*/>`, tag))
	}

	pattern := strings.Join(tagPatterns, "|")
	re := regexp.MustCompile(pattern)

	// 找到所有匹配的标签
	matches := re.FindAllStringIndex(line, -1)

	// 如果没有匹配到任何标签，直接全部转义
	if len(matches) == 0 {
		line = strings.ReplaceAll(line, "<", "&lt;")
		line = strings.ReplaceAll(line, ">", "&gt;")
		return line
	}

	// 有匹配的标签，需要保护
	var result strings.Builder
	lastEnd := 0

	for _, match := range matches {
		start, end := match[0], match[1]

		// 转义匹配标签之前的文本
		segment := line[lastEnd:start]
		segment = strings.ReplaceAll(segment, "<", "&lt;")
		segment = strings.ReplaceAll(segment, ">", "&gt;")
		result.WriteString(segment)

		// 保留标签本身
		result.WriteString(line[start:end])
		lastEnd = end
	}

	// 处理最后一部分
	segment := line[lastEnd:]
	segment = strings.ReplaceAll(segment, "<", "&lt;")
	segment = strings.ReplaceAll(segment, ">", "&gt;")
	result.WriteString(segment)

	return result.String()
}

func Parse(book *model.Book) error {
	if book == nil {
		return fmt.Errorf("book参数不能为nil")
	}
	var contentList []model.Section
	fmt.Println("正在读取txt文件...")
	start := time.Now()
	buf := readBuffer(book, book.Filename)
	var title string
	var content bytes.Buffer
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if line != "" {
					if line = strings.TrimSpace(line); line != "" {
						utils.AddPart(&content, line)
					}
				}
				contentList = append(contentList, model.Section{
					Title:   title,
					Content: content.String(),
				})
				content.Reset()
				break
			}
			return fmt.Errorf("读取文件出错: %w", err)
		}
		line = strings.TrimSpace(line)
		// 智能处理 HTML 标签：保留 epub 支持的标签，转义其他标签
		line = sanitizeHTMLTags(line)
		// 空行直接跳过
		if len(line) == 0 {
			continue
		}
		// 处理标题（优先匹配卷）
		if utf8.RuneCountInString(line) <= int(book.Max) {
			isVolume := book.VolumeReg.MatchString(line)
			isChapter := book.Reg.MatchString(line)
			isExclusion := false
			if book.ExclusionReg != nil && book.ExclusionReg.MatchString(line) {
				isExclusion = true
			}

			if !isExclusion && (isVolume || isChapter) {
				if title == "" {
					title = book.UnknowTitle
				}
				if content.Len() > 0 || title != book.UnknowTitle {
					contentList = append(contentList, model.Section{
						Title:   title,
						Content: content.String(),
					})
				}
				title = line
				content.Reset()
				continue
			}
		}
		utils.AddPart(&content, line)
	}
	// 没识别到章节又没识别到 EOF 时，把所有的内容写到最后一章
	if content.Len() != 0 {
		if title == "" {
			title = "章节正文"
		}
		contentList = append(contentList, model.Section{
			Title:   title,
			Content: content.String(),
		})
	}
	var sectionList []model.Section
	var volumeSection *model.Section
	for _, section := range contentList {
		if book.VolumeMatch != "false" && book.VolumeReg.MatchString(section.Title) {
			if volumeSection != nil {
				sectionList = append(sectionList, *volumeSection)
				volumeSection = nil
			}
			temp := section
			volumeSection = &temp
		} else if strings.HasPrefix(section.Title, "完本感言") || strings.HasPrefix(section.Title, "番外") {
			if volumeSection != nil {
				sectionList = append(sectionList, *volumeSection)
				volumeSection = nil
			}
			sectionList = append(sectionList, section)
		} else {
			if volumeSection == nil {
				sectionList = append(sectionList, section)
			} else {
				volumeSection.Sections = append(volumeSection.Sections, section)
			}
		}
	}
	// 如果有最后一卷,添加到章节列表
	if volumeSection != nil {
		sectionList = append(sectionList, *volumeSection)
		volumeSection = nil
	}
	end := time.Now().Sub(start)
	fmt.Println("读取文件耗时:", end)
	fmt.Println("匹配章节:", model.SectionCount(sectionList))
	// 添加提示
	if book.Tips {
		tuorialSection := model.Section{
			Title:   "制作说明",
			Content: model.Tutorial,
		}
		sectionList = append([]model.Section{tuorialSection}, sectionList...)
		sectionList = append(sectionList, tuorialSection)
	}
	book.SectionList = sectionList
	return nil
}
