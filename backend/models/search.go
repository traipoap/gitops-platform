package models

type SearchParams struct {
	FromTimestamp     *string `form:"from_timestamp"`
	ToTimestamp       *string `form:"to_timestamp"`
	SourceIP          *string `form:"source_ip"`
	Message           *string `form:"message"`
	MaxHits           *int    `form:"max_hits"`
	EndIndexTimestamp *string `form:"index_timestamp"`
}

type SearchResponse struct {
	Hits  []map[string]interface{} `json:"hits"`
	Total uint64                   `json:"total"`
}

// Quickwit API request/response (internal use)
type QuickwitSearchRequest struct {
	Query       string `json:"query"`
	MaxHits     int    `json:"max_hits"`
	SortByField string `json:"sort_by_field"`
	SortOrder   string `json:"sort_order,omitempty"`
}

type QuickwitSearchResponse struct {
	NumHits uint64 `json:"num_hits"`
	Hits    []struct {
		Source map[string]interface{} `json:"_source"`
	} `json:"hits"`
}
