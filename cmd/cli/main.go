package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/ystyle/kaf-cli/internal/converter"
	"github.com/ystyle/kaf-cli/internal/core"
	"github.com/ystyle/kaf-cli/internal/model"
	"github.com/ystyle/kaf-cli/internal/utils"
	"github.com/ystyle/kaf-cli/pkg/analytics"
)

var (
	secret      string
	measurement string
	version     string
)

func NewBookArgs() *model.Book {
	var book model.Book
	flag.StringVar(&book.Filename, "filename", "", "txt 文件名")
	flag.StringVar(&book.Bookname, "bookname", "", "书名: 默认为txt文件名")
	flag.StringVar(&book.Author, "author", "YSTYLE", "作者")
	flag.StringVar(&book.Match, "match", "", "匹配标题的正则表达式, 不写可以自动识别, 如果没生成章节就参考教程。例: -match 第.{1,8}章 表示第和章字之间可以有1-8个任意文字")
	flag.StringVar(&book.VolumeMatch, "volume-match", model.VolumeMatch, "卷匹配规则,设置为false可以禁用卷识别")
	flag.StringVar(&book.UnknowTitle, "unknow-title", "章节正文", "未知章节默认名称")
	flag.StringVar(&book.Cover, "cover", "cover.png", "封面图片可为: 本地图片, 和orly。 设置为orly时生成orly风格的封面, 需要连接网络。")
	flag.StringVar(&book.CoverOrlyColor, "cover-orly-color", "", "orly封面的主题色, 可以为1-16和hex格式的颜色代码, 不填时随机")
	flag.IntVar(&book.CoverOrlyIdx, "cover-orly-idx", -1, "orly封面的动物, 可以为0-41, 不填时随机, 具体图案可以查看: https://orly.nanmu.me")
	flag.UintVar(&book.Max, "max", 35, "标题最大字数")
	flag.UintVar(&book.Indent, "indent", 2, "段落缩进字数")
	flag.StringVar(&book.Align, "align", utils.GetEnv("KAF_CLI_ALIGN", "center"), "标题对齐方式: left、center、righ。环境变量KAF_CLI_ALIGN可修改默认值")
	flag.StringVar(&book.Bottom, "bottom", "1em", "段落间距(单位可以为em、px)")
	flag.StringVar(&book.LineHeight, "line-height", "", "行高(用于设置行间距, 默认为1.5rem)")
	flag.StringVar(&book.Font, "font", "", "嵌入字体, 之后epub的正文都将使用该字体")
	flag.StringVar(&book.Lang, "lang", utils.GetEnv("KAF_CLI_LANG", "zh"), "设置语言: en,de,fr,it,es,zh,ja,pt,ru,nl。环境变量KAF_CLI_LANG可修改默认值")
	flag.StringVar(&book.Format, "format", utils.GetEnv("KAF_CLI_FORMAT", "all"), "书籍格式: all、epub、mobi、azw3。环境变量KAF_CLI_FORMAT可修改默认值")
	flag.StringVar(&book.Out, "out", "", "输出文件名，不需要包含格式后缀")
	flag.BoolVar(&book.Tips, "tips", true, "添加本软件教程")
	flag.Parse()
	return &book
}

func printHelp(version string) {
	fmt.Println("错误: 文件名不能为空")
	fmt.Println("软件版本: \t", version)
	fmt.Println("简洁模式: \t把文件拖放到kaf-cli上")
	fmt.Println("命令行简单模式: kaf-cli ebook.txt")
	fmt.Println("\n以下为kaf-cli的全部参数")
	flag.PrintDefaults()
	if runtime.GOOS == "windows" {
		time.Sleep(time.Second * 10)
	}
}

func main() {
	var book *model.Book
	var err error
	if len(os.Args) == 2 && strings.HasSuffix(os.Args[1], ".txt") {
		book, err = model.NewBookSimple(os.Args[1])
		if err != nil {
			fmt.Printf("错误: %s\n", err.Error())
			os.Exit(1)
		}
	} else {
		book = NewBookArgs()
	}
	if err := core.Check(book, version); err != nil {
		if err.Error() == "不是txt文件" {
			fmt.Printf("错误: %s\n", err.Error())
			os.Exit(1)
		}
		printHelp(version)
		os.Exit(1)
	}
	analytics.Analytics(version, secret, measurement, book.Format)
	book.ToString()
	if err := core.Parse(book); err != nil {
		fmt.Printf("错误: %s\n", err.Error())
		os.Exit(2)
	}
	conv := converter.Dispatcher{
		Book: book,
	}
	if err := conv.Convert(); err != nil {
		fmt.Printf("错误: %s\n", err.Error())
		os.Exit(1)
	}
}
