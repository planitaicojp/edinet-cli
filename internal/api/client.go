package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	cerrors "github.com/planitaicojp/edinet-cli/internal/errors"
)

const (
	DefaultBaseURL = "https://api.edinet-fsa.go.jp/api/v2"
	maxRetries     = 3
)

type Client struct {
	BaseURL     string
	APIKey      string
	HTTP        *http.Client
	Debug       bool
	BackoffFunc func(attempt int) time.Duration // override for testing
}

func NewClient(apiKey string) *Client {
	return &Client{
		BaseURL: DefaultBaseURL,
		APIKey:  apiKey,
		HTTP: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) buildURL(path string, params map[string]string) (string, error) {
	u, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("Subscription-Key", c.APIKey)
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (c *Client) doGet(rawURL string) (*http.Response, error) {
	if c.Debug {
		SetDebug(true)
	}

	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}

	debugLogRequest(req)
	start := time.Now()

	var resp *http.Response
	for attempt := 0; attempt <= maxRetries; attempt++ {
		resp, err = c.HTTP.Do(req)
		if err != nil {
			if attempt == maxRetries {
				return nil, &cerrors.NetworkError{Err: err}
			}
			wait := c.backoff(attempt)
			debugLogRetry(attempt+1, 0, wait)
			time.Sleep(wait)
			continue
		}
		if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			if attempt < maxRetries {
				_ = resp.Body.Close()
				wait := c.backoff(attempt)
				debugLogRetry(attempt+1, resp.StatusCode, wait)
				time.Sleep(wait)
				continue
			}
		}
		break
	}
	elapsed := time.Since(start)

	if resp == nil {
		return nil, &cerrors.NetworkError{Err: fmt.Errorf("no response after %d retries", maxRetries)}
	}

	// Debug log response — only buffer body for text/JSON to avoid OOM on large binaries
	if debugEnabled {
		ct := resp.Header.Get("Content-Type")
		if isTextContentType(ct) {
			body, _ := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			resp.Body = io.NopCloser(bytes.NewReader(body))
			debugLogResponse(resp, elapsed, body)
		} else {
			debugLogResponse(resp, elapsed, nil)
		}
	}

	return resp, nil
}

// backoff returns the wait duration for exponential backoff: 1s, 2s, 4s
func (c *Client) backoff(attempt int) time.Duration {
	if c.BackoffFunc != nil {
		return c.BackoffFunc(attempt)
	}
	return time.Duration(1<<uint(attempt)) * time.Second
}

func (c *Client) Get(path string, params map[string]string) ([]byte, error) {
	rawURL, err := c.buildURL(path, params)
	if err != nil {
		return nil, err
	}

	resp, err := c.doGet(rawURL)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, classifyError(resp.StatusCode, body)
	}

	return body, nil
}

func (c *Client) GetBinary(path string, params map[string]string) (io.ReadCloser, string, error) {
	rawURL, err := c.buildURL(path, params)
	if err != nil {
		return nil, "", err
	}

	resp, err := c.doGet(rawURL)
	if err != nil {
		return nil, "", err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		return nil, "", classifyError(resp.StatusCode, body)
	}

	return resp.Body, resp.Header.Get("Content-Type"), nil
}

func classifyError(statusCode int, body []byte) error {
	var errResp struct {
		Metadata struct {
			Message string `json:"message"`
		} `json:"metadata"`
	}
	msg := string(body)
	if json.Unmarshal(body, &errResp) == nil && errResp.Metadata.Message != "" {
		msg = errResp.Metadata.Message
	}

	switch statusCode {
	case http.StatusBadRequest:
		return &cerrors.ValidationError{Field: "request", Message: msg}
	case http.StatusForbidden:
		return &cerrors.AuthError{Message: msg}
	case http.StatusNotFound:
		return &cerrors.NotFoundError{Resource: "document", ID: ""}
	default:
		return &cerrors.APIError{StatusCode: statusCode, Message: msg}
	}
}

// isTextContentType returns true for content types safe to buffer for debug logging.
func isTextContentType(ct string) bool {
	ct = strings.ToLower(ct)
	return strings.HasPrefix(ct, "application/json") ||
		strings.HasPrefix(ct, "text/") ||
		strings.Contains(ct, "xml")
}

func MaskAPIKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "..." + key[len(key)-4:]
}
