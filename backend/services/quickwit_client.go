package services

import (
	"bytes"
	"context"
	"encoding/json"
	"exporter/models"
	"fmt"
	"net/http"
	"time"
)

type QuickwitClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewQuickwitClient(baseURL string) *QuickwitClient {
	// Reuse connections (HTTP keep-alive) for better perf.
	// Go's http.Transport auto-decompresses gzip responses by default.
	transport := &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 5,
		IdleConnTimeout:     90 * time.Second,
	}
	return &QuickwitClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
	}
}

// quickwitSearchResponse matches the actual Quickwit search response shape:
//
//	{
//	  "hits": [
//	    { "timestamp": 1782322199000, "message": "...", "source_ip": "...", ... },
//	    ...
//	  ],
//	  "num_hits": 1234,
//	  "took_ms": 45
//	}
type quickwitSearchResponse struct {
	Hits    []map[string]interface{} `json:"hits"`
	NumHits uint64                   `json:"num_hits"`
}

func (c *QuickwitClient) Search(ctx context.Context, query string, max_hits *int, index_id *string) (*models.SearchResponse, error) {
	payload := models.QuickwitSearchRequest{
		Query:   query,
		MaxHits: *max_hits,
		SortBy:  "index_timestamp", // sort by fast-field timestamp (newest-first by doc-id default)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/api/v1/"+*index_id+"/search",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("quickwit API error: status %d", resp.StatusCode)
	}

	// Go's http.Transport already decompresses gzip transparently,
	// so resp.Body is always plain JSON here. Stream-decode directly.
	var qwResp quickwitSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&qwResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &models.SearchResponse{
		Hits:  qwResp.Hits,
		Total: qwResp.NumHits,
	}, nil
}
