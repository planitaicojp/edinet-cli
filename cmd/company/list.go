package company

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/planitaicojp/edinet-cli/cmd/cmdutil"
	"github.com/planitaicojp/edinet-cli/internal/api"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "EDINETコードリストを表示",
	Long:  `EDINETコードリスト（企業一覧）をローカルキャッシュから表示します。`,
	Example: `  edinet company list
  edinet company list --update
  edinet company list --format json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		update, _ := cmd.Flags().GetBool("update")
		cachePath := api.CacheFilePath()

		if update {
			if err := downloadAndCache(cmd, cachePath); err != nil {
				return err
			}
		}

		companies, err := api.LoadCodeList(cachePath)
		if err != nil {
			return err
		}

		formatter := cmdutil.GetFormatter(cmd)
		return formatter.Format(os.Stdout, companies)
	},
}

func init() {
	listCmd.Flags().Bool("update", false, "公式サイトから最新コードリストを再ダウンロード")
}

func downloadAndCache(cmd *cobra.Command, cachePath string) error {
	client, err := cmdutil.NewClient(cmd)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "EDINETコードリストをダウンロード中...")
	body, err := client.DownloadCodeList()
	if err != nil {
		return err
	}
	defer func() { _ = body.Close() }()

	// Read ZIP into memory
	zipData, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("ダウンロードエラー: %w", err)
	}

	// Extract CSV from ZIP
	csvData, err := extractCSVFromZip(zipData)
	if err != nil {
		return err
	}

	// Write CSV to cache
	dir := filepath.Dir(cachePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("キャッシュディレクトリ作成エラー: %w", err)
	}
	if err := os.WriteFile(cachePath, csvData, 0600); err != nil {
		return fmt.Errorf("キャッシュ書き込みエラー: %w", err)
	}

	fmt.Fprintf(os.Stderr, "コードリストを保存しました: %s\n", cachePath)
	return nil
}

func extractCSVFromZip(data []byte) ([]byte, error) {
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("ZIP展開エラー: %w", err)
	}

	for _, f := range r.File {
		if strings.HasSuffix(strings.ToLower(f.Name), ".csv") {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("CSV読み込みエラー: %w", err)
			}
			defer func() { _ = rc.Close() }()
			return io.ReadAll(rc)
		}
	}

	return nil, fmt.Errorf("ZIP内にCSVファイルが見つかりません")
}
