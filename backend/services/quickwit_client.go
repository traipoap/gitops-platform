package services

import (
	"bytes"
	"context"
	"encoding/json"
	"exporter/models"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

type QuickwitClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewQuickwitClient(baseURL string) *QuickwitClient {
	return &QuickwitClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *QuickwitClient) Search(ctx context.Context, query string, max_hits *int, index_id *string) (*models.SearchResponse, error) {
	payload := models.QuickwitSearchRequest{
		Query:       query,
		MaxHits:     *max_hits,
		SortByField: "index_timestamp",
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

	fmt.Println("request body:", string(body))

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("quickwit API error: status %d", resp.StatusCode)
	}

	// Read body as []byte don't unmarshal
	body, _ = io.ReadAll(resp.Body)

	hits := gjson.GetBytes(body, "hits").Array()
	total := gjson.GetBytes(body, "num_hits").Uint()

	// Convert []gjson.Result -> []map[string]interface{}
	hitMaps := make([]map[string]interface{}, len(hits))
	for i, h := range hits {
		hitMaps[i] = h.Value().(map[string]interface{})
	}

	return &models.SearchResponse{
		Hits:  hitMaps,
		Total: total,
	}, nil
}
