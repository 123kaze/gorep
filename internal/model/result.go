package model

type LineMatch struct {
	LineNum int
	Content string
}

type FileMatch struct {
	FilePath string
	Lines    []LineMatch
}
