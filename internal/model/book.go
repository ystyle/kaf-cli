package model

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/ystyle/kaf-cli/internal/utils"
	"golang.org/x/text/encoding"
)

var (
	ErrInvalidFile   = errors.New("invalid input file")
	ErrMissingConfig = errors.New("missing required configuration")
)

const (
	VolumeMatch      = "^第[0-9一二三四五六七八九十零〇百千两 ]+[卷部]"
	DefaultMatchTips = "^第[0-9一二三四五六七八九十零〇百千两 ]+[章回节集幕卷部]|^[Ss]ection.{1,20}$|^[Cc]hapter.{1,20}$|^[Pp]age.{1,20}$|^\\d{1,4}$|^\\d+、$|^引子$|^楔子$|^章节目录|^章节|^序章|^最终章 \\w{1,20}$|^番外\\d?\\w{0,20}|^完本感言.{0,4}$"
	DefaultExclusion = "^第[0-9一二三四五六七八九十零〇百千两 ]+(部门|部队|部属|部分|部件|部落|部.*：$)"
	Tutorial         = `本书由kaf-cli生成: <br/>
制作教程: <a href='https://ystyle.top/2019/12/31/txt-converto-epub-and-mobi/'>https://ystyle.top/2019/12/31/txt-converto-epub-and-mobi</a>
`
)

func NewBookSimple(filename string) (*Book, error) {
	book := Book{
		Filename: filename,
	}
	SetDefault(&book)
	return &book, nil
}

type Book struct {
	Filename         string    // 目录
	Bookname         string    // 书名
	Author           string    // 作者
	SectionList      []Section // 章节
	Match            string    // 正则
	VolumeMatch      string    // 卷匹配规则
	ExclusionPattern string    // 用户自定义的排除规则（正则）
	Max              uint      // 标题最大字数
	Indent           uint      // 段落缩进字段
	Align            string    // 标题对齐方式
	UnknowTitle      string    // 未知章节名称
	Cover            string    // 封面图片
	CoverOrlyColor   string    // 生成封面图片的颜色
	CoverOrlyIdx     int       // 生成封面图片的动物
	Font             string    // 嵌入字体
	Bottom           string    // 段阿落间距
	LineHeight       string    // 行高
	Tips             bool      // 是否添加教程文本
	Lang             string    // 设置语言
	Out              string    // 输出文件名
	Format           string    // 书籍格式
	Decoder          *encoding.Decoder
	PageStylesFile   string
	Reg              *regexp.Regexp
	VolumeReg        *regexp.Regexp
	ExclusionReg     *regexp.Regexp // 动态生成的正则，用于排除无效标题
	Version          string
}

type Section struct {
	Title    string
	Content  string
	Sections []Section
}

func SectionCount(sections []Section) int {
	var count int
	for _, section := range sections {
		count += 1 + len(section.Sections)
	}
	return count
}

func SetDefault(book *Book) {
	book.Match = utils.DefaultString(book.Match, DefaultMatchTips)
	book.VolumeMatch = utils.DefaultString(book.VolumeMatch, VolumeMatch)
	book.Author = utils.DefaultString(book.Author, "YSTYLE")
	book.UnknowTitle = utils.DefaultString(book.UnknowTitle, "章节正文")
	book.Max = utils.DefalutInt(book.Max, 35)
	book.Indent = utils.DefalutInt(book.Indent, 2)
	book.Align = utils.DefaultString(book.Align, utils.GetEnv("KAF_CLI_ALIGN", "center"))
	book.Cover = utils.DefaultString(book.Cover, "cover.png")
	book.Bottom = utils.DefaultString(book.Bottom, "1em")
	book.Lang = utils.DefaultString(book.Lang, utils.GetEnv("KAF_CLI_LANG", "zh"))
	book.Format = utils.DefaultString(book.Format, utils.GetEnv("KAF_CLI_FORMAT", "all"))
	book.CoverOrlyIdx = utils.DefalutInt(book.CoverOrlyIdx, -1)
	book.ExclusionPattern = utils.DefaultString(book.ExclusionPattern, DefaultExclusion) // 默认排除规则
}

func (book *Book) ToString() {
	fmt.Println("转换信息:")
	fmt.Println("软件版本:", book.Version)
	fmt.Println("文件名:\t", book.Filename)
	fmt.Println("书籍书名:", book.Bookname)
	fmt.Println("书籍作者:", book.Author)
	if book.Cover != "" {
		fmt.Println("书籍封面:", book.Cover)
	}
	fmt.Println("书籍语言:", book.Lang)
	if book.Match == DefaultMatchTips {
		fmt.Println("匹配条件:", "自动匹配")
	} else {
		fmt.Println("匹配条件:", book.Match)
	}
	fmt.Println("卷匹配条件:", book.VolumeMatch)
	fmt.Println("转换格式:", book.Format)
	fmt.Println()
}
