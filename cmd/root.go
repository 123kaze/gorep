package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gorep/internal/config"
	"gorep/internal/engine"
	"os"
	"runtime"
)

var (
	ignoreCase  bool
	fixedString bool
	workers     int
)
var rootCmd = &cobra.Command{
	Use:   "gorep <pattern> <path>",
	Short: "A fast file search CLI for Go",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.SearchConfig{
			Pattern:     args[0],
			Path:        args[1],
			IgnoreCase:  ignoreCase,
			FixedString: fixedString,
			Workers:     workers,
		}

		matcher, err := engine.NewMatcher(cfg.Pattern, cfg.FixedString, cfg.IgnoreCase)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}

		matchs, err := engine.SearchPath(cfg, matcher)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		if len(matchs) == 0 {
			fmt.Println("no matched content was found")
			return
		}
		for _, match := range matchs {
			fmt.Println(match.FilePath)
			for _, line := range match.Lines {
				fmt.Printf("  %d:  %s\n", line.LineNum, line.Content)
			}
		}
		// 初版代码
		//if match == nil {
		//	fmt.Println("no matched content in", cfg.Path)
		//	return
		//}
		//fmt.Println(match.FilePath)
		//for _, line := range match.Lines {
		//	fmt.Printf("  %d:%s\n", line.LineNum, line.Content)
		//}
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&ignoreCase, "ignore-case", "i", false, "ignore case distinction")
	rootCmd.Flags().BoolVarP(&fixedString, "fixed-string", "F", false, "target Pattern as a literal string")
	rootCmd.Flags().IntVarP(&workers, "workers", "j", runtime.NumCPU(), "number of worker goroutines")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
