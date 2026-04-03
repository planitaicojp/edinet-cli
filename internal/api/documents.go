package api

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/planitaicojp/edinet-cli/internal/model"
)

func (c *Client) ListDocuments(date string, typ int) (*model.DocumentListResponse, error) {
	params := map[string]string{
		"date": date,
		"type": strconv.Itoa(typ),
	}
	body, err := c.Get("/documents.json", params)
	if err != nil {
		return nil, err
	}

	var resp model.DocumentListResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &resp, nil
}

func (c *Client) DownloadDocument(docID string, typ int) (io.ReadCloser, string, error) {
	params := map[string]string{
		"type": strconv.Itoa(typ),
	}
	return c.GetBinary(fmt.Sprintf("/documents/%s", docID), params)
}
