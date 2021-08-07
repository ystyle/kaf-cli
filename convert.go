package kafcli

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/bmaupin/go-epub"
	"github.com/leotaku/mobi"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/language"
	"golang.org/x/text/transform"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

type Section struct {
	Title   string
	Content string
}

var (
	filename       string // 目录
	bookname       string // 书名
	match          string // 正则
	author         string // 作者
	max            uint   // 标题最大字数
	indent         uint   // 段落缩进字段
	align          string // 标题对齐方式
	cover          string // 封面图片
	bottom         string // 段阿落间距
	Tips           bool   // 是否添加教程文本
	lang           string // 设置语言
	out            string // 输出文件名
	format         string // 书籍格式
	decoder        *encoding.Decoder
	pageStylesFile string
	reg            *regexp.Regexp
)

const (
	htmlPStart         = `<p class="content">`
	htmlPEnd           = "</p>"
	htmlTitleStart     = `<h3 class="title">`
	mobiTtmlTitleStart = `<h3 style="text-align:%s;">`
	htmlTitleEnd       = "</h3>"
	DefaultMatchTips   = "自动匹配,可自定义"
	cssContent         = `
.title {text-align:%s}
.content {
  margin-bottom: %s;
  margin-top: 0;
  text-indent: %dem;
}
`
	Tutorial = `本书由kaf-cli生成: <br/>
制作教程: <a href='https://ystyle.top/2019/12/31/txt-converto-epub-and-mobi/'>https://ystyle.top/2019/12/31/txt-converto-epub-and-mobi</a>
`
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
		lang = "zh"
		indent = 2
		align = "center"
		cover = "cover.png"
		bottom = "1em"
		format = "both"
	} else {
		flag.StringVar(&filename, "filename", "", "txt 文件名")
		flag.StringVar(&author, "author", "YSTYLE", "作者")
		flag.StringVar(&bookname, "bookname", "", "书名: 默认为txt文件名")
		flag.UintVar(&max, "max", 35, "标题最大字数")
		flag.StringVar(&match, "match", DefaultMatchTips, "匹配标题的正则表达式, 不写可以自动识别, 如果没生成章节就参考教程。例: -match 第.{1,8}章 表示第和章字之间可以有1-8个任意文字")
		flag.UintVar(&indent, "indent", 2, "段落缩进字数")
		flag.StringVar(&align, "align", "center", "标题对齐方式: left、center、righ")
		flag.StringVar(&cover, "cover", "cover.png", "封面图片")
		flag.StringVar(&bottom, "bottom", "1em", "段落间距(单位可以为em、px)")
		flag.StringVar(&format, "format", "all", "书籍格式: all、epub、mobi、azw3")
		flag.StringVar(&lang, "lang", "zh", "设置语言: en,de,fr,it,es,zh,ja,pt,ru,nl。 支持使用环境变量KAF-CLI-LANG设置")
		flag.BoolVar(&Tips, "tips", true, "添加本软件教程")
		flag.Parse()
	}

	if filename == "" {
		fmt.Println("文件名不能为空")
		fmt.Println("简洁模式: 直接把文件播放到kaf-cli上")
		fmt.Println("命令行简单模式: kaf-cli ebook.txt")
		fmt.Println("查看命令行参数: kaf-cli -h")
		fmt.Println("以下为kaf-cli的全部参数")
		flag.PrintDefaults()
		time.Sleep(time.Second * 10)
		os.Exit(1)
	}

	if bookname == "" {
		bookname = strings.Split(filepath.Base(filename), ".")[0]
	}

	if out == "" {
		out = bookname
	}

	if l := os.Getenv("KAF-CLI-LANG"); l != "" {
		lang = l
	}
	lang = parseLang(lang)

	if exists, _ := isExists(cover); !exists {
		cover = ""
	}

	fmt.Println("转换信息:")
	fmt.Println("文件名:", filename)
	fmt.Println("书名:", bookname)
	if author != "" {
		fmt.Println("作者:", author)
	}
	if cover != "" {
		fmt.Println("封面:", cover)
	}
	fmt.Println("匹配条件:", match)
	fmt.Println("书籍语言:", lang)
	fmt.Println()

	// 编译正则表达式
	if match == "" || match == DefaultMatchTips {
		match = "^.{0,8}(第.{1,20}(章|节)|(S|s)ection.{1,20}|(C|c)hapter.{1,20}|(P|p)age.{1,20})|^\\d{1,4}.{0,20}$|^引子|^楔子|^章节目录"
	}
	var err error
	reg, err = regexp.Compile(match)
	if err != nil {
		fmt.Printf("生成匹配规则出错: %s\n%s\n", match, err.Error())
		return
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

func parseBook() []Section {
	var contentList []Section
	fmt.Println("正在读取txt文件...")
	start := time.Now()
	buf := readBuffer(filename)
	var title string
	var content bytes.Buffer
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if line != "" {
					if line = strings.TrimSpace(line); line != "" {
						addPart(&content, line)
					}
				}
				contentList = append(contentList, Section{
					Title:   title,
					Content: content.String(),
				})
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
					addPart(&content, Tutorial)
				}
			}
			if content.Len() == 0 {
				continue
			}
			contentList = append(contentList, Section{
				Title:   title,
				Content: content.String(),
			})
			title = line
			content.Reset()
			continue
		}
		addPart(&content, line)
	}
	// 没识别到章节又没识别到 EOF 时，把所有的内容写到最后一章
	if content.Len() != 0 {
		if title == "" {
			title = "章节正文"
		}
		contentList = append(contentList, Section{
			Title:   title,
			Content: content.String(),
		})
	}
	end := time.Now().Sub(start)
	fmt.Println("读取文件耗时:", end)

	// 添加提示
	if Tips {
		contentList = append(contentList, Section{
			Title:   "制作说明",
			Content: Tutorial,
		})
	}
	return contentList
}

func Convert() {
	start := time.Now()
	// 解析文本
	sectionList := parseBook()
	fmt.Println()

	// 判断要生成的格式
	var isEpub, isMobi, isAzw3 bool
	switch format {
	case "epub":
		isEpub = true
	case "mobi":
		isEpub = true
		isMobi = true
	case "azw3":
		isAzw3 = true
	default:
		isEpub = true
		isMobi = true
		isAzw3 = true
	}

	hasKinldegen := lookKindlegen()
	if isMobi && hasKinldegen == "" {
		isEpub = false
	}
	// 生成epub
	if isEpub {
		buildEpub(sectionList)
		fmt.Println()
	}
	// 生成azw3格式
	if isAzw3 {
		// 生成kindle格式
		buildAzw3(sectionList)
	}
	// 生成mobi格式
	if isMobi {
		if hasKinldegen != "" {
			converToMobi(fmt.Sprintf("%s.epub", out))
		}
	}
	end := time.Now().Sub(start)
	fmt.Println("\n转换完成! 总耗时:", end)
}

func wrapMobiTitle(title, content string) string {
	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf(mobiTtmlTitleStart, align))
	buff.WriteString(title)
	buff.WriteString(htmlTitleEnd)
	buff.WriteString(content)
	return buff.String()
}

func buildAzw3(sectionList []Section) {
	fmt.Println("使用第三方库生成azw3, 不保证所有样式都能正常显示")
	fmt.Println("正在生成azw3...")
	start := time.Now()
	mb := mobi.Book{
		Title:       bookname,
		Authors:     []string{author},
		CreatedDate: time.Now(),
		Chapters:    []mobi.Chapter{},
		Language:    language.MustParse(lang),
		UniqueID:    rand.Uint32(),
	}
	css := fmt.Sprintf(cssContent, align, bottom, indent)
	for _, section := range sectionList {
		ch := mobi.Chapter{
			Title:  section.Title,
			Chunks: mobi.Chunks(wrapMobiTitle(section.Title, section.Content)),
		}
		mb.Chapters = append(mb.Chapters, ch)
	}

	mb.CSSFlows = []string{css}
	if cover != "" {
		f, err := os.Open(cover)
		if err != nil {
			panic(err)
		}
		img, _, err := image.Decode(f)
		if err != nil {
			panic(err)
		}
		mb.CoverImage = img
	}

	// Convert book to PalmDB database
	db := mb.Realize()

	// Write database to file
	f, _ := os.Create(fmt.Sprintf("%s.azw3", out))
	err := db.Write(f)
	if err != nil {
		panic(err)
	}
	fmt.Println("生成azw3电子书耗时:", time.Now().Sub(start))
}

func wrapEpubTitle(title, content string) string {
	var buff bytes.Buffer
	buff.WriteString(htmlTitleStart)
	buff.WriteString(title)
	buff.WriteString(htmlTitleEnd)
	buff.WriteString(content)
	return buff.String()
}

func buildEpub(sectionList []Section) {
	fmt.Println("正在生成epub")
	start := time.Now()
	// 写入样式
	tempDir, err := ioutil.TempDir("", "kaf-cli")
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			panic(fmt.Sprintf("创建临时文件夹失败: %s", err))
		}
	}()
	pageStylesFile = path.Join(tempDir, "page_styles.css")

	err = ioutil.WriteFile(pageStylesFile, []byte(fmt.Sprintf(cssContent, align, bottom, indent)), 0666)
	if err != nil {
		panic(fmt.Sprintf("无法写入样式文件: %s", err))
	}
	// Create a ne EPUB
	e := epub.NewEpub(bookname)
	e.SetLang(lang)
	// Set the author
	e.SetAuthor(author)
	css, err := e.AddCSS(pageStylesFile, "")
	if err != nil {
		panic(fmt.Sprintf("无法写入样式文件: %s", err))
	}

	if cover != "" {
		img, err := e.AddImage(cover, cover)
		if err != nil {
			panic(err)
		}
		e.SetCover(img, "")
	}

	for _, section := range sectionList {
		e.AddSection(wrapEpubTitle(section.Title, section.Content), section.Title, "", css)
	}

	// Write the EPUB
	fmt.Println("正在生成电子书...")
	epubName := out + ".epub"
	err = e.Write(epubName)
	if err != nil {
		// handle error
	}
	// 计算耗时
	end := time.Now().Sub(start)
	fmt.Println("生成EPUB电子书耗时:", end)
}

func addPart(buff *bytes.Buffer, content string) {
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
