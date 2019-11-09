package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/bmaupin/go-epub"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	filename string // 目录
	bookname string // 书名
	match    string // 正则
	author   string // 正则
)

const (
	htmlPStart     = `<p  style="text-indent: 2em">`
	htmlPEnd       = "</p>"
	htmlTitleStart = `<h3 style="text-align:center">`
	htmlTitleEnd   = "</h3>"
)

// 解析程序参数
func init() {
	flag.StringVar(&filename, "filename", "", "txt 文件名")
	flag.StringVar(&author, "author", "YSTYLE", "作者")
	flag.StringVar(&bookname, "bookname", "", "书名: 默认为txt文件名")
	flag.StringVar(&match, "match", "第.{1,8}章", "匹配标题的正则表达式, 例: -match 第.{1,8}章 表示第和章字之间可以有1-8个任意文字")
	flag.Parse()
}

func readBuffer(filename string) *bufio.Reader {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("读取文件出错: ", err.Error())
		os.Exit(1)
	}
	buf := bufio.NewReader(f)
	return buf
}

func main() {
	if filename == "" {
		fmt.Println("文件名不能为空")
		os.Exit(1)
	}

	start := time.Now()
	if bookname == "" {
		bookname = strings.Split(filepath.Base(filename), ".")[0]
	}

	fmt.Println("转换信息:")
	fmt.Println("文件名:", filename)
	fmt.Println("书名:", bookname)
	if author != "" {
		fmt.Println("作者:", author)
	}
	fmt.Println("匹配条件:", match)
	fmt.Println()

	// Create a ne EPUB
	e := epub.NewEpub(bookname)

	// Set the author
	e.SetAuthor(author)
	fmt.Println("正在读取txt文件...")

	// 编译正则表达式
	reg, err := regexp.Compile(match)
	if err != nil {
		fmt.Printf("生成匹配规则出错: %s\n%s\n", match, err.Error())
		return
	}

	buf := readBuffer(filename)
	var title string
	var content bytes.Buffer
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				e.AddSection(content.String(), title, "", "")
				break
			}
			fmt.Println("读取文件出错:", err.Error())
			os.Exit(1)
		}
		line = strings.TrimSpace(line)
		// 处理标题
		if reg.MatchString(line) {
			if title == "" {
				title = "说明"
			}
			e.AddSection(content.String(), title, "", "")
			title = line
			content.Reset()
			content.WriteString(htmlTitleStart)
			content.WriteString(title)
			content.WriteString(htmlTitleEnd)
			continue
		}
		if strings.HasSuffix(line, "==") ||
			strings.HasSuffix(line, "**") ||
			strings.HasSuffix(line, "--") ||
			strings.HasSuffix(line, "//") {
			content.WriteString(line)
			continue
		}
		content.WriteString(htmlPStart)
		content.WriteString(line)
		content.WriteString(htmlPEnd)
	}
	end := time.Now().Sub(start)
	fmt.Println("耗时:", end)

	// Write the EPUB
	fmt.Println("正在生成电子书...")
	err = e.Write(bookname + ".epub")
	if err != nil {
		// handle error
	}
	// 计算耗时
	end = time.Now().Sub(start)
	fmt.Println("耗时:", end)
	fmt.Println("\n转换完成! 总耗时:", end)
}
