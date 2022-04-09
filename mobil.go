package kafcli

import (
	"fmt"
	"github.com/766b/mobi"
	"time"
)

type MobiConverter struct{}

func (convert MobiConverter) Build(book Book) error {
	fmt.Println("使用第三方库生成mobi, 不保证所有样式都能正常显示")
	fmt.Println("正在生成mobi...")
	start := time.Now()
	m, err := mobi.NewWriter(fmt.Sprintf("%s.mobi", book.Out))
	if err != nil {
		panic(err)
	}
	m.Title(book.Bookname)
	m.Compression(mobi.CompressionNone)
	if book.Cover != "" {
		m.AddCover(book.Cover, book.Cover)
	}
	m.NewExthRecord(mobi.EXTH_DOCTYPE, "EBOK")
	m.NewExthRecord(mobi.EXTH_AUTHOR, book.Author)
	for _, section := range book.SectionList {
		m.NewChapter(section.Title, []byte(section.Content))
	}
	m.Write()
	fmt.Println("生成mobi电子书耗时:", time.Now().Sub(start))
	return nil
}
