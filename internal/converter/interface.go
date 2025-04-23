package converter

import "github.com/ystyle/kaf-cli/internal/model"

type Converter interface {
	Build(book model.Book) error
}
