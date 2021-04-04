package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/bmaupin/go-epub"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	filename string // 目录
	bookname string // 书名
	match    string // 正则
	author   string // 作者
	max      uint   // 标题最大字数
	Tips     bool
	lang     string // 设置语言
	decoder  *encoding.Decoder
)

const (
	htmlPStart       = `<p class="content">`
	htmlPEnd         = "</p>"
	htmlTitleStart   = `<h3 class="title">`
	htmlTitleEnd     = "</h3>"
	DefaultMatchTips = "自动匹配,可自定义"
	cssContent       = `
.title {text-align:center}
.content {text-indent: 2em}
`
	Tutorial = `本书由TmdTextEpub生成: <br/>
制作教程: <a href='https://ystyle.top/2019/12/31/txt-converto-epub-and-mobi/'>https://ystyle.top/2019/12/31/txt-converto-epub-and-mobi</a>
`
)

func parseLang(lang string) string {
	var langs = "en,de,fr,it,es,zh,ja,pt,ru,nl"
	if strings.Contains(langs, lang) {
		return lang
	}
	return "en"
}

// 解析程序参数
func init() {
	if len(os.Args) == 2 && strings.HasSuffix(os.Args[1], ".txt") {
		fmt.Println("傻瓜模式开启...")
		filename = os.Args[1]
		author = "YSTYLE"
		max = 35
		match = DefaultMatchTips
		Tips = true
		lang = "zh"
	} else {
		flag.StringVar(&filename, "filename", "", "txt 文件名")
		flag.StringVar(&author, "author", "YSTYLE", "作者")
		flag.StringVar(&bookname, "bookname", "", "书名: 默认为txt文件名")
		flag.UintVar(&max, "max", 35, "标题最大字数")
		flag.StringVar(&match, "match", DefaultMatchTips, "匹配标题的正则表达式, 不写可以自动识别, 如果没生成章节就参考教程。例: -match 第.{1,8}章 表示第和章字之间可以有1-8个任意文字")
		flag.StringVar(&lang, "lang", "zh", "设置语言: en,de,fr,it,es,zh,ja,pt,ru,nl。 支持使用环境变量KAF-CLI-LANG设置")
		flag.BoolVar(&Tips, "tips", true, "添加本软件教程")
		flag.Parse()
	}
}

func readBuffer(filename string) *bufio.Reader {
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
		bs, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println("读取文件出错: ", err.Error())
			os.Exit(1)
		}
		var buf bytes.Buffer
		decoder = encodig.NewDecoder()
		if encodename == "windows-1252" {
			decoder = simplifiedchinese.GB18030.NewDecoder()
		}
		bs, _, _ = transform.Bytes(decoder, bs)
		buf.Write(bs)
		return bufio.NewReader(&buf)
	} else {
		f.Seek(0, 0)
		buf := bufio.NewReader(f)
		return buf
	}
}

func main() {
	if filename == "" {
		fmt.Println("文件名不能为空")
		os.Exit(1)
	}

	if bookname == "" {
		bookname = strings.Split(filepath.Base(filename), ".")[0]
	}
	if l := os.Getenv("KAF-CLI-LANG"); l != "" {
		lang = l
	}
	lang = parseLang(lang)

	fmt.Println("转换信息:")
	fmt.Println("文件名:", filename)
	fmt.Println("书名:", bookname)
	if author != "" {
		fmt.Println("作者:", author)
	}
	fmt.Println("匹配条件:", match)
	fmt.Println("书籍语言:", lang)
	fmt.Println()

	// 编译正则表达式
	if match == "" || match == DefaultMatchTips {
		match = "^.{0,8}(第.{1,20}(章|节)|(S|s)ection.{1,20}|(C|c)hapter.{1,20}|(P|p)age.{1,20})|^\\d{1,4}.{0,20}$|^引子|^楔子|^章节目录"
	}
	reg, err := regexp.Compile(match)
	if err != nil {
		fmt.Printf("生成匹配规则出错: %s\n%s\n", match, err.Error())
		return
	}

	// 写入样式
	tempDir, err := ioutil.TempDir("", "TmdTextEpub")
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			panic(fmt.Sprintf("创建临时文件夹失败: %s", err))
		}
	}()
	pageStylesFile := path.Join(tempDir, "page_styles.css")
	err = ioutil.WriteFile(pageStylesFile, []byte(cssContent), 0666)
	if err != nil {
		panic(fmt.Sprintf("无法写入样式文件: %s", err))
	}

	start := time.Now()
	// Create a ne EPUB
	e := epub.NewEpub(bookname)
	e.SetLang(lang)
	// Set the author
	e.SetAuthor(author)
	css, err := e.AddCSS(pageStylesFile, "")
	if err != nil {
		panic(fmt.Sprintf("无法写入样式文件: %s", err))
	}

	fmt.Println("正在读取txt文件...")

	buf := readBuffer(filename)
	var title string
	var content bytes.Buffer

	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if line != "" {
					if line = strings.TrimSpace(line); line != "" {
						AddPart(&content, line)
					}
				}
				e.AddSection(content.String(), title, "", css)
				content.Reset()
				break
			}
			fmt.Println("读取文件出错:", err.Error())
			os.Exit(1)
		}
		line = strings.TrimSpace(line)
		line = strings.ReplaceAll(line, "<", "&lt;")
		line = strings.ReplaceAll(line, ">", "&gt;")
		// 空行直接跳过
		if len(line) == 0 {
			continue
		}
		// 处理标题
		if utf8.RuneCountInString(line) <= int(max) && reg.MatchString(line) {
			if title == "" {
				title = "说明"
				if Tips {
					AddPart(&content, Tutorial)
				}
			}
			if content.Len() == 0 {
				continue
			}
			e.AddSection(content.String(), title, "", css)
			title = line
			content.Reset()
			content.WriteString(htmlTitleStart)
			content.WriteString(title)
			content.WriteString(htmlTitleEnd)
			continue
		}
		AddPart(&content, line)
	}
	// 没识别到章节又没识别到 EOF 时，把所有的内容写到最后一章
	if content.Len() != 0 {
		if title == "" {
			title = "章节正文"
		}
		e.AddSection(content.String(), title, "", "")
	}
	end := time.Now().Sub(start)
	fmt.Println("读取文件耗时:", end)

	// 添加提示
	if Tips {
		e.AddSection(Tutorial, "制作说明", "", "")
	}

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

func AddPart(buff *bytes.Buffer, content string) {
	if strings.HasSuffix(content, "==") ||
		strings.HasSuffix(content, "**") ||
		strings.HasSuffix(content, "--") ||
		strings.HasSuffix(content, "//") {
		buff.WriteString(content)
		return
	}
	buff.WriteString(htmlPStart)
	buff.WriteString(content)
	buff.WriteString(htmlPEnd)
}

func Run(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ConverToMobi(bookname string) {
	command := "kindlegen"
	if runtime.GOOS == "windows" {
		command = "kindlegen.exe"
	}
	kindlegen, _ := exec.LookPath(command)
	if kindlegen == "" {
		currentDir := path.Dir(os.Args[0])
		kindlegen = path.Join(currentDir, command)
		if exist, _ := IsExists(kindlegen); !exist {
			return
		}
	}
	fmt.Printf("\n检测到Kindle格式转换器: %s，正在把书籍转换成Kindle格式...\n", command)
	fmt.Println("转换mobi比较花时间, 大约耗时1-10分钟, 请等待...")
	start := time.Now()
	Run(kindlegen, "-dont_append_source", "-locale", lang, "-c1", bookname)
	// 计算耗时
	end := time.Now().Sub(start)
	fmt.Println("转换为Kindle格式耗时:", end)
}

func IsExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
