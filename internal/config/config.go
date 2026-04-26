package config

type SearchConfig struct {
	Pattern     string
	Path        string
	IgnoreCase  bool
	FixedString bool
	Workers     int
	Includes    []string
	Excludes    []string
}
