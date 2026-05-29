package kafcli

import (
	"github.com/ystyle/kaf-cli/internal/converter"
	"github.com/ystyle/kaf-cli/internal/core"
	"github.com/ystyle/kaf-cli/internal/model"
)

type Book = model.Book
type Section = model.Section

var NewSimpleBook = model.NewBookSimple

// SetDefault fills default values for a book configuration.
func SetDefault(book *Book) {
	model.SetDefault(book)
}

// Check validates a book, parses metadata from filename, and compiles regex patterns.
func Check(book *Book, version string) error {
	return core.Check(book, version)
}

// Parse reads the txt file, detects encoding, and splits content into sections.
func Parse(book *Book) error {
	return core.Parse(book)
}

// Build generates epub/azw3/mobi output files from an already-parsed book.
// The caller must have called Check and Parse (or filled SectionList manually) before calling Build.
func Build(book *Book) error {
	conv := converter.Dispatcher{Book: book}
	return conv.Convert()
}

// Convert runs the full pipeline: Check → Parse → Build.
func Convert(book *Book) error {
	if err := core.Check(book, "v1.0.0"); err != nil {
		return err
	}
	if err := core.Parse(book); err != nil {
		return err
	}
	return Build(book)
}

// Preview is like Convert but only returns the parsed section list without generating output files.
func Preview(book *Book) ([]Section, error) {
	if err := core.Check(book, "v1.0.0"); err != nil {
		return nil, err
	}
	if err := core.Parse(book); err != nil {
		return nil, err
	}
	return book.SectionList, nil
}