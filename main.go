package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/bmaupin/go-epub"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

var (
	filename string // 目录
	bookname string // 书名
	match    string // 正则
	author   string // 作者
	max      uint   // 标题最大字数
	decoder  *encoding.Decoder
	Tips     bool
)

const (
	htmlPStart       = `<p  style="text-indent: 2em">`
	htmlPEnd         = "</p>"
	htmlTitleStart   = `<h3 style="text-align:center">`
	htmlTitleEnd     = "</h3>"
	DefaultMatchTips = "自动匹配,可自定义"
)

// 解析程序参数
func init() {
	if len(os.Args) == 2 && strings.HasSuffix(os.Args[1], ".txt") {
		fmt.Println("傻瓜模式开启...")
		filename = os.Args[1]
		author = "YSTYLE"
		max = 35
		match = DefaultMatchTips
		Tips = true
	} else {
		flag.StringVar(&filename, "filename", "", "txt 文件名")
		flag.StringVar(&author, "author", "YSTYLE", "作者")
		flag.StringVar(&bookname, "bookname", "", "书名: 默认为txt文件名")
		flag.UintVar(&max, "max", 35, "标题最大字数")
		flag.StringVar(&match, "match", DefaultMatchTips, "匹配标题的正则表达式, 不写可以自动识别, 如果没生成章节就参考教程。例: -match 第.{1,8}章 表示第和章字之间可以有1-8个任意文字")
		flag.BoolVar(&Tips, "tips", true, "添加本软件教程")
		flag.Parse()
	}
	decoder = simplifiedchinese.GBK.NewDecoder()
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

func preNUm(data byte) int {
	str := fmt.Sprintf("%b", data)
	var i int = 0
	for i < len(str) {
		if str[i] != '1' {
			break
		}
		i++
	}
	return i
}
func isUtf8(data []byte) bool {
	for i := 0; i < len(data); {
		if data[i]&0x80 == 0x00 {
			i++
			continue
		} else if num := preNUm(data[i]); num > 2 {
			i++
			for j := 0; j < num-1; j++ {
				if data[i]&0xc0 != 0x80 {
					return false
				}
				i++
			}
		} else {
			return false
		}
	}
	return true
}

func conver(str string) string {
	buff, err := decoder.Bytes([]byte(str))
	if err != nil {
		return str
	}
	return string(buff)
}

func main() {
	if filename == "" {
		fmt.Println("文件名不能为空")
		os.Exit(1)
	}

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

	// 编译正则表达式
	if match == "" || match == DefaultMatchTips {
		match = "^.{0,8}(第.{1,20}(章|节)|(S|s)ection.{1,20}|(C|c)hapter.{1,20})"
	}
	reg, err := regexp.Compile(match)
	if err != nil {
		fmt.Printf("生成匹配规则出错: %s\n%s\n", match, err.Error())
		return
	}

	start := time.Now()
	// Create a ne EPUB
	e := epub.NewEpub(bookname)

	// Set the author
	e.SetAuthor(author)
	fmt.Println("正在读取txt文件...")

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
		if !isUtf8([]byte(line)) {
			line = conver(line)
		}
		line = strings.TrimSpace(line)
		// 空行直接跳过
		if len(line) == 0 {
			continue
		}
		// 处理标题
		if len(line) <= int(max) && reg.MatchString(line) {
			if title == "" {
				title = "说明"
				if Tips {
					content.WriteString(htmlPStart)
					content.WriteString("本书由TmdTextEpub生成: <br/>")
					content.WriteString("制作教程: <a href='https://ystyle.top/2019/12/31/txt-converto-epub-and-mobi/'>https://ystyle.top/2019/12/31/txt-converto-epub-and-mobi</a>")
					content.WriteString(htmlPEnd)
				}
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
	fmt.Println("读取文件耗时:", end)

	// Write the EPUB
	fmt.Println("正在生成电子书...")
	epubName := bookname + ".epub"
	err = e.Write(epubName)
	if err != nil {
		// handle error
	}
	// 计算耗时
	end = time.Now().Sub(start)
	fmt.Println("生成EPUB电子书耗时:", end)
	// 生成kindle格式
	ConverToMobi(epubName)
	end = time.Now().Sub(start)
	fmt.Println("\n转换完成! 总耗时:", end)
}

func Run(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ConverToMobi(bookname string) {
	command := "kindlegen"
	if runtime.GOOS == "windows" {
		command = "kindlegen.exe"
	}
	kindlegen, err := exec.LookPath(command)
	if err != nil {
		return
	}
	fmt.Printf("\n检测到Kindle格式转换器: %s，正在把书籍转换成Kindle格式...\n", command)
	fmt.Println("转换mobi比较花时间, 请等待...")
	start := time.Now()
	Run(kindlegen, "-dont_append_source", bookname)
	// 计算耗时
	end := time.Now().Sub(start)
	fmt.Println("转换为Kindle格式耗时:", end)
}
