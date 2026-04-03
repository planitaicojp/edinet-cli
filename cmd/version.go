package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "バージョン情報を表示",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("edinet version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
