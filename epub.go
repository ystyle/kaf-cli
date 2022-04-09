package kafcli

import (
	"bytes"
	"fmt"
	"github.com/bmaupin/go-epub"
	"io/ioutil"
	"os"
	"path"
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
	tempDir, err := ioutil.TempDir("", "kaf-cli")
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			panic(fmt.Sprintf("创建临时文件夹失败: %s", err))
		}
	}()
	pageStylesFile := path.Join(tempDir, "page_styles.css")

	err = ioutil.WriteFile(pageStylesFile, []byte(fmt.Sprintf(cssContent, book.Align, book.Bottom, book.Indent)), 0666)
	if err != nil {
		return fmt.Errorf("无法写入样式文件: %w", err)
	}
	// Create a ne EPUB
	e := epub.NewEpub(book.Bookname)
	e.SetLang(book.Lang)
	// Set the author
	e.SetAuthor(book.Author)
	css, err := e.AddCSS(pageStylesFile, "")
	if err != nil {
		return fmt.Errorf("无法写入样式文件: %w", err)
	}

	if book.Cover != "" {
		img, err := e.AddImage(book.Cover, book.Cover)
		if err != nil {
			return fmt.Errorf("添加封面失败: %w", err)
		}
		e.SetCover(img, "")
	}

	for _, section := range book.SectionList {
		e.AddSection(convert.wrapTitle(section.Title, section.Content), section.Title, "", css)
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
