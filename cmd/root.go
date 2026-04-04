package cmd

import (
	"github.com/planitaicojp/edinet-cli/cmd/company"
	configCmd "github.com/planitaicojp/edinet-cli/cmd/config"
	"github.com/planitaicojp/edinet-cli/cmd/document"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "edinet",
	Short: "EDINET API CLI - 金融庁 開示書類の検索・取得ツール",
	Long: `edinet は、金融庁が提供する EDINET API v2 を操作するための CLI ツールです。
有価証券報告書等の開示書類を検索・取得できます。`,
	Version:       version,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().StringP("format", "f", "", "出力形式 (table, json, csv)")
	rootCmd.PersistentFlags().String("api-key", "", "EDINET API キー")
	rootCmd.PersistentFlags().Bool("verbose", false, "デバッグ出力を有効化")
	rootCmd.PersistentFlags().Bool("no-color", false, "カラー出力を無効化")

	rootCmd.AddCommand(document.Cmd)
	rootCmd.AddCommand(company.Cmd)
	rootCmd.AddCommand(configCmd.Cmd)
}

func Execute() error {
	return rootCmd.Execute()
}
