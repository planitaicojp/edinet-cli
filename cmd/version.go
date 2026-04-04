package cmd

import (
	"github.com/spf13/cobra"
)

var version = "dev"

const banner = `  _____ ____ ___ _   _ _____ _____
 | ____|  _ \_ _| \ | | ____|_   _|
 |  _| | | | | ||  \| |  _|   | |
 | |___| |_| | || |\  | |___  | |
 |_____|____/___|_| \_|_____| |_|
  edinet-cli %s
  Author:  Tonghyun Kim
  License: Apache-2.0
  Home:    https://github.com/planitaicojp/edinet-cli

  This is an unofficial tool and is not affiliated with
  or endorsed by FSA (Financial Services Agency of Japan).
`

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "バージョン情報を表示",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf(banner, version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
