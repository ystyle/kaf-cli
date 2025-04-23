package converter

import (
	"fmt"
	"time"

	"github.com/766b/mobi"
	"github.com/ystyle/kaf-cli/internal/model"
)

type MobiConverter struct {
	HTMLPStart     string // MOBI专属段落标签
	HTMLPEnd       string
	HTMLTitleStart string
	HTMLTitleEnd   string
}

func NewMobiConverter() *MobiConverter {
	return &MobiConverter{
		HTMLPStart:     `<p class="content">`,
		HTMLPEnd:       "</p>",
		HTMLTitleStart: `<h3 class="title">`,
		HTMLTitleEnd:   "</h3>",
	}
}

func (convert MobiConverter) Build(book model.Book) error {
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
		if len(section.Sections) > 0 {
			for _, subsection := range section.Sections {
				m.NewChapter(subsection.Title, []byte(subsection.Content))
			}
		}
	}
	m.Write()
	fmt.Println("生成mobi电子书耗时:", time.Now().Sub(start))
	return nil
}
