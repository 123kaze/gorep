package printer

import (
	"gorep/internal/model"
)

type Printer interface {
	Print(match model.FileMatch)
}
