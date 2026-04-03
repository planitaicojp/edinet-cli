package config

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "config",
	Short: "設定の表示・変更",
}

func init() {
	Cmd.AddCommand(setCmd)
	Cmd.AddCommand(showCmd)
}
