package main

import (
	"bufio"
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
	match    string // 正则
	author   string // 正则
)

// 解析程序参数
func init() {
	flag.StringVar(&filename, "filename", "", "txt 文件名")
	flag.StringVar(&author, "author", "YSTYLE", "作者")
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
	start := time.Now()
	bookName := strings.Split(filepath.Base(filename), ".")[0]

	fmt.Println("转换信息:")
	fmt.Println("文件名:", filename)
	if author != "" {
		fmt.Println("作者:", author)
	}
	fmt.Println("匹配条件:", match)
	fmt.Println()

	// Create a ne EPUB
	e := epub.NewEpub(bookName)

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
	var content string

	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				e.AddSection(content, title, "", "")
				break
			}
			fmt.Println("读取文件出错:", err.Error())
			os.Exit(1)
		}
		line = strings.TrimSpace(line)
		if reg.MatchString(line) {
			if title == "" {
				title = "说明"
			}
			e.AddSection(content, title, "", "")
			title = line
			content = line + "<br/>"
		} else {
			content = content + line + "<br/>"
		}
	}
	end := time.Now().Sub(start)
	fmt.Println("耗时:", end)

	// Write the EPUB
	fmt.Println("正在生成电子书...")
	err = e.Write(bookName + ".epub")
	if err != nil {
		// handle error
	}
	// 计算耗时
	end = time.Now().Sub(start)
	fmt.Println("耗时:", end)
	fmt.Println("\n转换完成! 总耗时:", end)
}
