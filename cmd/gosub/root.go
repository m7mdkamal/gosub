package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/m7mdkamal/gosub"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gosub",
	Short: "GoSub is subtitle downloader.",
	Long:  "GoSub is subtitle downloader.",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := os.Getwd()

		if len(args) == 0 {
			args = append(args, path)
		}

		for _, arg := range args {

			if !filepath.IsAbs(arg) {
				path = filepath.Join(path, arg)
			} else {
				path = arg
			}
			lang, _ := cmd.Flags().GetString("language")
			query, _ := cmd.Flags().GetString("query")
			gosub.Run(path, lang, query)

		}
	},
}

func init() {
	rootCmd.Flags().StringP("language", "l", "english", "Set subtitle language")
	rootCmd.Flags().StringP("query", "q", "", "Set custom search query")
}

// Execute to start the cmd cobra
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
