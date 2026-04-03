package document

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:     "document",
	Aliases: []string{"doc"},
	Short:   "開示書類の一覧取得・ダウンロード",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(downloadCmd)
}
