package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	cerrors "github.com/planitaicojp/edinet-cli/internal/errors"
)

const DefaultBaseURL = "https://api.edinet-fsa.go.jp/api/v2"

type Client struct {
	BaseURL string
	APIKey  string
	HTTP    *http.Client
	Debug   bool
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

func (c *Client) Get(path string, params map[string]string) ([]byte, error) {
	u, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("Subscription-Key", c.APIKey)
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	resp, err := c.HTTP.Get(u.String())
	if err != nil {
		return nil, &cerrors.NetworkError{Err: err}
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
	u, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return nil, "", err
	}

	q := u.Query()
	q.Set("Subscription-Key", c.APIKey)
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	resp, err := c.HTTP.Get(u.String())
	if err != nil {
		return nil, "", &cerrors.NetworkError{Err: err}
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

func MaskAPIKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "..." + key[len(key)-4:]
}
