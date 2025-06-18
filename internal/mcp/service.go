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
)

type ConverterService struct {
	version string
}

func (s *ConverterService) RegisterTools(srv *server.MCPServer, version string) {
	s.version = version
	tool := mcp.NewTool("kaf_convert",
		mcp.WithDescription("电子书格式转换器，支持把txt文件转换成epub电子书格式\n若转换成功，AI助手在返回结果给用户时应该使用markdorn: `[/home/user/documents/book.epub](/home/user/documents/book.epub)`格式, 两个URI都使用完整路径，以方便用户查看和点击跳转"),
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
			mcp.Description("章节匹配规则的正则表达式, 不填会自动识别常见章节名称"),
		),
	)
	website := mcp.NewPrompt("kaf-mcp", mcp.WithPromptDescription("kaf-mcp官方网站,代码仓库"))
	srv.AddPrompt(website, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return mcp.NewGetPromptResult(
			"kaf-mcp官方网站",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleAssistant,
					mcp.NewTextContent("https://github.com/ystyle/kaf-cli"),
				),
			},
		), nil
	})

	srv.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// 记录请求参数
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

		author, _ := req.Params.Arguments["author"].(string)
		bookname, _ := req.Params.Arguments["bookname"].(string)
		match, _ := req.Params.Arguments["match"].(string)

		book, _ := model.NewBookSimple(filename)
		book.Format = "epub"
		if author != "" {
			logger.Info("use author", "author", author)
			book.Author = author
		}
		if bookname != "" {
			logger.Info("use bookname", "bookname", bookname)
			book.Bookname = bookname
		}
		if match != "" {
			book.Match = match
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

		bookfile := fmt.Sprintf("%s.epub", book.Out)
		bookfile, _ = filepath.Abs(bookfile)
		// 按文档示例返回结果
		return mcp.NewToolResultResource(book.Bookname, mcp.BlobResourceContents{
			URI: fmt.Sprintf("%s", bookfile),
			// MIMEType: "application/epub+zip",
			// Blob:     blob,
		}), nil

	})
}
