package mcp

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/ystyle/kaf-cli/internal/converter"
	"github.com/ystyle/kaf-cli/internal/core"
	"github.com/ystyle/kaf-cli/internal/model"
	"github.com/ystyle/kaf-cli/pkg/analytics"
)

var (
	secret      string
	measurement string
)

type ConverterService struct {
	version string
}

func (s *ConverterService) RegisterTools(srv *server.MCPServer, version string) {
	s.version = version
	tool := mcp.NewTool("kaf_convert",
		mcp.WithDescription("电子书格式转换器，支持把txt文件转换成epub/azw3/mobi电子书格式\n若转换成功，AI助手在返回结果给用户时应该使用markdown: `[/home/user/documents/book.epub](/home/user/documents/book.epub)`格式, 两个URI都使用完整路径，以方便用户查看和点击跳转"),
		mcp.WithString("filename",
			mcp.Required(),
			mcp.Description("txt小说文件, 支持相对路径，相对路径默认会从配置目录读取小说文件"),
			mcp.Pattern(`\.txt$`),
		),
		mcp.WithString("bookname",
			mcp.Description("书名, 为空可以自动识别的文件名格式: 《(.*)》.*作者[：:](.*).txt"),
		),
		mcp.WithString("author",
			mcp.Description("作者, 为空可以自动识别的文件名格式: 《(.*)》.*作者[：:](.*).txt"),
		),
		mcp.WithString("match",
			mcp.Description("章节匹配规则的正则表达式, 不填会自动识别常见章节名称, 例: 第.{1,8}章"),
		),
		mcp.WithString("format",
			mcp.Description("输出格式: epub、mobi、azw3、all(全部格式)。默认epub"),
			mcp.Enum("epub", "mobi", "azw3", "all"),
		),
		mcp.WithString("cover",
			mcp.Description("封面图片: 本地图片路径(jpg/png)、orly(在线生成O'Rly风格封面,需联网)、none(无封面)。默认使用cover.png"),
		),
		mcp.WithString("cover-orly-color",
			mcp.Description("O'Rly封面主题色, 可为1-16编号或hex颜色代码(如#ff6600), 不填随机"),
		),
		mcp.WithNumber("cover-orly-idx",
			mcp.Description("O'Rly封面动物图案编号 0-41, 不填随机, 查看: https://orly.nanmu.me"),
		),
		mcp.WithString("font",
			mcp.Description("嵌入字体文件路径, 嵌入后epub正文将使用该字体"),
		),
		mcp.WithString("custom-css-file",
			mcp.Description("自定义CSS文件路径, 覆盖默认样式。可用选择器: h2.volume(卷名), h3.title(章节标题), .content(正文), body(整体)"),
		),
		mcp.WithString("lang",
			mcp.Description("书籍语言: en,de,fr,it,es,zh,ja,pt,ru,nl。默认zh"),
			mcp.Enum("en", "de", "fr", "it", "es", "zh", "ja", "pt", "ru", "nl"),
		),
		mcp.WithNumber("indent",
			mcp.Description("段落缩进字数, 默认2"),
		),
		mcp.WithString("align",
			mcp.Description("标题对齐方式: left、center、right。默认center"),
			mcp.Enum("left", "center", "right"),
		),
		mcp.WithString("line-height",
			mcp.Description("行高/行间距, 如: 1.5rem、2。默认1.5rem"),
		),
		mcp.WithString("bottom",
			mcp.Description("段落间距, 单位可为em、px。默认1em"),
		),
		mcp.WithNumber("max",
			mcp.Description("标题最大字数, 超过的标题会被忽略。默认35"),
		),
		mcp.WithString("out",
			mcp.Description("输出文件名, 不需要包含格式后缀。默认为书名"),
		),
		mcp.WithString("volume-match",
			mcp.Description("卷匹配正则表达式, 设为false禁用卷识别。默认自动匹配"),
		),
		mcp.WithString("exclude",
			mcp.Description("排除无效章节/卷的正则表达式"),
		),
		mcp.WithString("unknow-title",
			mcp.Description("未知章节默认名称。默认: 章节正文"),
		),
		mcp.WithBoolean("separate-chapter-number",
			mcp.Description("是否分离章节序号和标题样式(序号单独一行显示)。默认false"),
		),
		mcp.WithBoolean("tips",
			mcp.Description("是否在书中添加软件教程。默认true"),
		),
	)
	srv.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logger.Info("request received",
			"tool", "kaf_convert",
			"params", req.Params.Arguments,
		)
		filename, ok := req.Params.Arguments["filename"].(string)
		if !ok {
			err := errors.New("filename is required")
			logger.Error("parameter validation failed", "error", err)
			return nil, err
		}

		book, _ := model.NewBookSimple(filename)

		if v, _ := req.Params.Arguments["author"].(string); v != "" {
			book.Author = v
		}
		if v, _ := req.Params.Arguments["bookname"].(string); v != "" {
			book.Bookname = v
		}
		if v, _ := req.Params.Arguments["match"].(string); v != "" {
			book.Match = v
		}
		if v, _ := req.Params.Arguments["format"].(string); v != "" {
			book.Format = v
		} else {
			book.Format = "epub"
		}
		if v, _ := req.Params.Arguments["cover"].(string); v != "" {
			book.Cover = v
		}
		if v, _ := req.Params.Arguments["cover-orly-color"].(string); v != "" {
			book.CoverOrlyColor = v
		}
		if v, ok := req.Params.Arguments["cover-orly-idx"].(float64); ok {
			book.CoverOrlyIdx = int(v)
		}
		if v, _ := req.Params.Arguments["font"].(string); v != "" {
			book.Font = v
		}
		if v, _ := req.Params.Arguments["custom-css-file"].(string); v != "" {
			book.CustomCSSFile = v
		}
		if v, _ := req.Params.Arguments["lang"].(string); v != "" {
			book.Lang = v
		}
		if v, ok := req.Params.Arguments["indent"].(float64); ok && v > 0 {
			book.Indent = uint(v)
		}
		if v, _ := req.Params.Arguments["align"].(string); v != "" {
			book.Align = v
		}
		if v, _ := req.Params.Arguments["line-height"].(string); v != "" {
			book.LineHeight = v
		}
		if v, _ := req.Params.Arguments["bottom"].(string); v != "" {
			book.Bottom = v
		}
		if v, ok := req.Params.Arguments["max"].(float64); ok && v > 0 {
			book.Max = uint(v)
		}
		if v, _ := req.Params.Arguments["out"].(string); v != "" {
			book.Out = v
		}
		if v, _ := req.Params.Arguments["volume-match"].(string); v != "" {
			book.VolumeMatch = v
		}
		if v, _ := req.Params.Arguments["exclude"].(string); v != "" {
			book.ExclusionPattern = v
		}
		if v, _ := req.Params.Arguments["unknow-title"].(string); v != "" {
			book.UnknowTitle = v
		}
		if v, ok := req.Params.Arguments["separate-chapter-number"].(bool); ok {
			book.SeparateChapterNumber = v
		}
		if v, ok := req.Params.Arguments["tips"].(bool); ok {
			book.Tips = v
		}

		logger.Info("check start")
		if err := core.Check(book, s.version); err != nil {
			logger.Error("check failed", "error", err, "filename", filename)
			return nil, err
		}

		logger.Info("parse start")
		if err := core.Parse(book); err != nil {
			logger.Error("parse failed", "error", err, "filename", filename)
			return nil, err
		}

		conv := converter.Dispatcher{
			Book: book,
		}
		logger.Info("convert start")
		if err := conv.Convert(); err != nil {
			logger.Error("convert failed", "error", err, "filename", filename)
			return nil, err
		}
		analytics.Analytics(s.version, secret, measurement, book.Format)

		bookfile := fmt.Sprintf("%s.%s", book.Out, book.Format)
		if book.Format == "all" {
			bookfile = fmt.Sprintf("%s.epub", book.Out)
		}
		bookfile, _ = filepath.Abs(bookfile)
		return mcp.NewToolResultResource(book.Bookname, mcp.BlobResourceContents{
			URI: fmt.Sprintf("%s", bookfile),
		}), nil

	})

	websitetool := mcp.NewTool("repo-url",
		mcp.WithDescription("获取kaf-cli的代码仓库地址"),
	)
	srv.AddTool(websitetool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("https://github.com/ystyle/kaf-cli"), nil
	})
}
