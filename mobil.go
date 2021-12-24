package kafcli

import (
	"fmt"
	mobi2 "github.com/766b/mobi"
)

type MobiConverter struct{}

func (convert MobiConverter) Build(book Book) error {
	m, err := mobi2.NewWriter(fmt.Sprintf("%s.mobi", book.Out))
	if err != nil {
		panic(err)
	}
	m.Title(book.Bookname)
	m.Compression(mobi2.CompressionNone)
	if book.Cover != "" {
		m.AddCover(book.Cover, book.Cover)
	}
	m.NewExthRecord(mobi2.EXTH_DOCTYPE, "EBOK")
	m.NewExthRecord(mobi2.EXTH_AUTHOR, book.Author)

	for _, section := range book.SectionList {
		m.NewChapter(section.Title, []byte(section.Content))
	}
	m.Write()
	return nil
}
