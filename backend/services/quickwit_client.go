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
	// Reuse connections (HTTP keep-alive) for better perf
	transport := &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 5,
		IdleConnTimeout:     90 * time.Second,
		// Go's Transport auto-handshake for gzip by default (DisableCompression: false).
		// It will decompress response automatically AND strip Content-Encoding header.
	}
	return &QuickwitClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
	}
}

func (c *QuickwitClient) Search(ctx context.Context, query string, max_hits *int, index_id *string) (*models.SearchResponse, error) {
	payload := models.QuickwitSearchRequest{
		Query:       query,
		MaxHits:     *max_hits,
		SortByField: "index_timestamp",
		SortOrder:   "desc", // newest first — matches UI expectation
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
	// so resp.Body is always plain JSON here. Just stream-decode it.
	var quickwitResp struct {
		NumHits uint64 `json:"num_hits"`
		Hits    []struct {
			Source map[string]interface{} `json:"_source"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&quickwitResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	hits := make([]map[string]interface{}, len(quickwitResp.Hits))
	for i, h := range quickwitResp.Hits {
		hits[i] = h.Source
	}

	return &models.SearchResponse{
		Hits:  hits,
		Total: quickwitResp.NumHits,
	}, nil
}
