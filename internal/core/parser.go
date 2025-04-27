package core

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
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
		line = strings.ReplaceAll(line, "<", "&lt;")
		line = strings.ReplaceAll(line, ">", "&gt;")
		// 空行直接跳过
		if len(line) == 0 {
			continue
		}
		// 处理标题（优先匹配卷）
		if utf8.RuneCountInString(line) <= int(book.Max) {
			isVolume := book.VolumeReg.MatchString(line)
			isChapter := book.Reg.MatchString(line)

			if isVolume || isChapter {
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
