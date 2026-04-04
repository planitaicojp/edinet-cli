package company

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "company",
	Short: "企業情報の検索・一覧表示 (EDINETコードリスト)",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(searchCmd)
}
