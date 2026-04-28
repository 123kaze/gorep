package model

type LineMatch struct {
	LineNum   int
	Content   string
	Ranges    [][2]int // 匹配发生在第几位到第几位
	BeforeCtx []string
	AfterCtx  []string
}

type FileMatch struct {
	FilePath string
	Lines    []LineMatch
}
