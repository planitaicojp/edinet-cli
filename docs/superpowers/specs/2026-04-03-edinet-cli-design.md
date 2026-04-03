# edinet-cli Design Spec

## Overview

EDINET (Electronic Disclosure for Investors' NETwork) API v2를 조작하기 위한 Go CLI 도구.
금융庁が提供する有価証券報告書等の開示書類を検索・取得する。

- **対象API**: EDINET API v2 (`https://api.edinet-fsa.go.jp/api/v2`)
- **認証**: Subscription-Key (API キー、クエリパラメータ)
- **参考プロジェクト**: `~/dev/crowdy/conoha-cli` (Cobra + internal/ 構造)

## API Endpoints

EDINET API v2は2つのエンドポイントのみ提供:

### 1. 書類一覧API
```
GET /api/v2/documents.json?date={YYYY-MM-DD}&type={1|2}&Subscription-Key={key}
```
- `date` (必須): ファイル日付
- `type` (任意): 1=メタデータのみ(デフォルト), 2=提出書類一覧+メタデータ
- レスポンス: JSON (metadata + results[])
- Document: 29フィールド (docID, edinetCode, secCode, JCN, filerName, docTypeCode, xbrlFlag, pdfFlag, csvFlag, legalStatus 等)

### 2. 書類取得API
```
GET /api/v2/documents/{docID}?type={1-5}&Subscription-Key={key}
```
- `type` (必須): 1=XBRL(ZIP), 2=PDF, 3=代替書面(ZIP), 4=英文(ZIP), 5=CSV(ZIP)
- レスポンス: バイナリ (ZIP or PDF)

### ステータスコード
| HTTP Status | 意味 |
|-------------|------|
| 200 | 成功 |
| 400 | パラメータエラー |
| 403 | 認証失敗 |
| 404 | 書類なし |
| 500 | サーバーエラー |

## Architecture

### Directory Structure

```
edinet-cli/
├── main.go                          # cmd.Execute() のみ
├── cmd/
│   ├── root.go                      # グローバルフラグ: --format, --api-key, --verbose, --no-color
│   ├── version.go                   # edinet version
│   ├── completion.go                # edinet completion {bash,zsh,fish}
│   ├── document/
│   │   ├── document.go              # edinet document (グループコマンド)
│   │   ├── list.go                  # edinet document list
│   │   ├── show.go                  # edinet document show <docID>
│   │   └── download.go             # edinet document download <docID>
│   ├── company/
│   │   ├── company.go               # edinet company (グループコマンド)
│   │   ├── list.go                  # edinet company list
│   │   └── search.go               # edinet company search <query>
│   ├── config/
│   │   ├── config.go                # edinet config (グループコマンド)
│   │   ├── set.go                   # edinet config set api-key <key>
│   │   └── show.go                  # edinet config show
│   └── cmdutil/
│       ├── client.go                # NewClient(cmd) → api.Client 生成
│       ├── format.go                # GetFormat(cmd) → output.Formatter
│       └── flags.go                 # グローバルフラグヘルパー
├── internal/
│   ├── api/
│   │   ├── client.go                # HTTP client (BaseURL, APIキー注入, リトライ, デバッグログ)
│   │   ├── documents.go             # ListDocuments(), GetDocument(), DownloadDocument()
│   │   └── codelist.go             # DownloadCodeList()
│   ├── config/
│   │   ├── config.go                # Load/Save config.yaml
│   │   └── env.go                   # 環境変数定義
│   ├── model/
│   │   ├── document.go              # Document, Metadata, ResultSet 構造体
│   │   ├── company.go               # Company (EDINETコードリスト項目)
│   │   └── response.go             # APIレスポンスラッパー
│   ├── output/
│   │   ├── formatter.go             # Formatter インターフェース + New(format)
│   │   ├── table.go                 # TableFormatter (tabwriter)
│   │   ├── json.go                  # JSONFormatter
│   │   └── csv.go                   # CSVFormatter
│   ├── errors/
│   │   ├── errors.go                # APIError, AuthError, NotFoundError 等
│   │   └── exitcodes.go            # 終了コード定数
│   └── xbrl/
│       ├── parser.go                # XBRL/iXBRLインスタンス文書パーサー
│       ├── taxonomy.go              # 日本GAAP/IFRSタクソノミマッピング
│       └── financial.go             # 財務項目抽出
└── test/
    └── fixtures/                    # テスト用JSON/XBRLサンプル
```

### Go Module

- Module: `github.com/planitaicojp/edinet-cli`
- Binary: `edinet`

## Authentication & Config

### Config Priority (高い順)
1. `--api-key` CLIフラグ
2. `EDINET_API_KEY` 環境変数
3. `~/.config/edinet/config.yaml`

### config.yaml
```yaml
api_key: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
default_format: table
```

### Environment Variables
| 変数名 | 説明 | デフォルト |
|--------|------|-----------|
| `EDINET_API_KEY` | APIキー | - |
| `EDINET_FORMAT` | 出力形式 | `table` |
| `EDINET_CONFIG_DIR` | configディレクトリ | `~/.config/edinet` |
| `EDINET_DEBUG` | デバッグログ | `false` |

### API Key Handling
- EDINET APIは `Subscription-Key` クエリパラメータで認証
- `cmdutil.NewClient(cmd)` がフラグ→環境変数→config順で解決し `api.Client` に注入
- `config set api-key` でconfig.yamlに保存 (ファイルパーミッション 0600)
- `config show` でAPIキーはマスキング表示 (`xxxx...xxxx`)

## Subcommands

### edinet document list
```
edinet document list --date 2024-04-01 [--type 1|2] [--format json|table|csv]
```
- `--date` (必須): 対象日付 (YYYY-MM-DD). `today`, `yesterday` もサポート
- `--type`: 1=メタデータのみ(デフォルト), 2=提出書類一覧+メタデータ
- type=2時: docID, filerName, docDescription, submitDateTime 等の主要フィールドをテーブル表示

### edinet document show
```
edinet document show <docID>
```
- 書類一覧API(type=2)から該当docIDの詳細情報を表示
- 29フィールド全体をフォーマットに合わせて出力

### edinet document download
```
edinet document download <docID> [--type xbrl|pdf|attach|english|csv] [--output dir/] [--extract]
```
- `--type`: xbrl(デフォルト), pdf, attach, english, csv → API type 1~5にマッピング
- `--output`: 保存ディレクトリ (デフォルト: カレントディレクトリ)
- `--extract`: ZIP自動展開
- ファイル名: `{docID}_{type}.{zip|pdf}`

### edinet company list
```
edinet company list [--update]
```
- EDINETコードリストCSVをローカルキャッシュから表示
- `--update`: 公式サイトから最新CSVを再ダウンロード
- キャッシュ: `~/.config/edinet/codelist.csv`

### edinet company search
```
edinet company search <query> [--by name|code|edinet-code|all]
```
- ローカルキャッシュのコードリストを検索
- `--by`: 検索対象 (デフォルト: `all`)
- 部分一致検索

### edinet config set/show
```
edinet config set api-key <value>
edinet config show
```

## API Client

### HTTP Client
```go
type Client struct {
    BaseURL    string       // https://api.edinet-fsa.go.jp/api/v2
    APIKey     string
    HTTP       *http.Client // 30s timeout
    Debug      bool
}
```

- リトライ: GETリクエストに対し429/5xx時最大3回、指数バックオフ
- デバッグモード: `--verbose` or `EDINET_DEBUG=true` で要求/応答ログ (APIキーマスキング)
- Subscription-Key: 全リクエストのクエリパラメータに自動付与

### API Methods
```go
// documents.go
func (c *Client) ListDocuments(date string, typ int) (*model.DocumentListResponse, error)
func (c *Client) DownloadDocument(docID string, typ int) (io.ReadCloser, string, error)

// codelist.go
func (c *Client) DownloadCodeList() (io.ReadCloser, error)
```

## Error Handling

### Error Types
| エラー型 | Exit Code | 状況 |
|---------|-----------|------|
| `ValidationError` | 2 | 不正な日付形式、必須パラメータ未指定 |
| `AuthError` | 3 | APIキー未設定/無効 (403) |
| `APIError` | 4 | APIレスポンスエラー (400, 500等) |
| `NetworkError` | 5 | ネットワーク接続失敗 |
| `NotFoundError` | 6 | 書類なし (404) |

### HTTP Status → CLI Error Mapping
| HTTP Status | CLI Error |
|-------------|-----------|
| 200 | - |
| 400 | `ValidationError` |
| 403 | `AuthError` |
| 404 | `NotFoundError` |
| 500 | `APIError` |

## Output Formatting

### Formatter Interface
```go
type Formatter interface {
    Format(w io.Writer, data any) error
}
```

### Implementations
- **Table**: `text/tabwriter`, 構造体の `json` タグをカラムヘッダーに使用
- **JSON**: `json.MarshalIndent` (pretty print)
- **CSV**: `encoding/csv`, 構造体フィールド順

### Format Selection Priority
1. `--format` フラグ (最優先)
2. `EDINET_FORMAT` 環境変数
3. config.yamlの `default_format`
4. デフォルト: `table`

## XBRL Parsing

### Strategy
ダウンロードしたZIP → 展開 → XBRLインスタンス文書パース

### Files
- `parser.go` — ZIP展開 + XBRL/iXBRLインスタンス文書ロード
- `taxonomy.go` — 日本GAAP/IFRSタクソノミネームスペースマッピング
- `financial.go` — 財務項目抽出 (コンテキスト別値解釈)

### Supported Formats
- XBRL 2.1 インスタンス文書
- Inline XBRL (iXBRL) — EDINETの最新提出はほぼiXBRL

### Extraction Targets (例)
- 売上高/営業収益、営業利益、経常利益、当期純利益
- 総資産、純資産、自己資本比率
- 発行済株式数、EPS、BPS、配当金

### Future Extension
`edinet xbrl parse <file>` サブコマンド追加の余地あり

## Dependencies (go.mod)

| Package | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLIフレームワーク |
| `gopkg.in/yaml.v3` | config.yamlパース |
| `golang.org/x/net/html` | iXBRL (Inline XBRL) パース |

標準ライブラリでカバー可能なものは外部依存なし (`encoding/xml`, `encoding/csv`, `text/tabwriter` 等).

## Notes

- GitHub非公式OpenAPI spec (gabu/edinet-api-spec) は v1ベースで使用不可
- 公式PDF仕様書 (2026年1月版) を正とする
- EDINETコードリストCSVは公式サイトから別途ダウンロード
