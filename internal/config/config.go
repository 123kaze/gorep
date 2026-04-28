package config

type SearchConfig struct {
	Pattern     string
	Path        string
	IgnoreCase  bool
	FixedString bool
	Workers     int
	Includes    []string
	Excludes    []string
	All         bool
	NoColor     bool
	BeforeCtx   int
	AfterCtx    int
	CtxLines    int
	NoGitIgnore bool
}
