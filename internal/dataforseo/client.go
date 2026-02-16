package dataforseo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/builtbyrobben/dataforseo-cli/internal/api"
)

// Client wraps the REST API client with DataForSEO-specific methods.
type Client struct {
	api *api.Client
}

// NewClient creates a new DataForSEO client.
func NewClient(authBase64 string, opts ...api.ClientOption) *Client {
	allOpts := append([]api.ClientOption{
		api.WithAuthFn(func(r *http.Request) {
			r.Header.Set("Authorization", "Basic "+authBase64)
		}),
		api.WithUserAgent("dataforseo-cli/1.0"),
	}, opts...)

	return &Client{
		api: api.NewClient("https://api.dataforseo.com/v3", allOpts...),
	}
}

func (c *Client) unwrapResult(resp *TaskResponse) ([]TaskResult, error) {
	if len(resp.Tasks) == 0 {
		return nil, fmt.Errorf("no tasks in response")
	}

	task := resp.Tasks[0]
	if task.StatusCode != 20000 {
		return nil, fmt.Errorf("task error (%d): %s", task.StatusCode, task.StatusMessage)
	}

	return task.Result, nil
}

// SearchVolume retrieves search volume data for the given keywords.
func (c *Client) SearchVolume(ctx context.Context, keywords []string, location int, language string) ([]KeywordData, error) {
	body := []map[string]any{
		{
			"keywords":      keywords,
			"location_code": location,
			"language_code": language,
		},
	}

	var resp TaskResponse
	if err := c.api.Post(ctx, "/keywords_data/google_ads/search_volume/live", body, &resp); err != nil {
		return nil, fmt.Errorf("search volume: %w", err)
	}

	results, err := c.unwrapResult(&resp)
	if err != nil {
		return nil, fmt.Errorf("search volume: %w", err)
	}

	var data []KeywordData

	for _, r := range results {
		data = append(data, KeywordData{
			Keyword:          r.Keyword,
			SearchVolume:     r.SearchVolume,
			Competition:      r.Competition,
			CompetitionLevel: r.CompetitionLevel,
			CPC:              r.CPC,
			LocationCode:     r.LocationCode,
			LanguageCode:     r.LanguageCode,
		})
	}

	return data, nil
}

// KeywordDifficulty retrieves difficulty scores for the given keywords.
func (c *Client) KeywordDifficulty(ctx context.Context, keywords []string, location int, language string) ([]DifficultyData, error) {
	body := []map[string]any{
		{
			"keywords":      keywords,
			"location_code": location,
			"language_code": language,
		},
	}

	var resp TaskResponse
	if err := c.api.Post(ctx, "/dataforseo_labs/google/keyword_difficulty/live", body, &resp); err != nil {
		return nil, fmt.Errorf("keyword difficulty: %w", err)
	}

	results, err := c.unwrapResult(&resp)
	if err != nil {
		return nil, fmt.Errorf("keyword difficulty: %w", err)
	}

	var data []DifficultyData

	for _, r := range results {
		data = append(data, DifficultyData{
			Keyword:           r.Keyword,
			KeywordDifficulty: r.KeywordDifficulty,
			LocationCode:      r.LocationCode,
			LanguageCode:      r.LanguageCode,
		})
	}

	return data, nil
}

// RelatedKeywords retrieves related keywords for the given seed keyword.
func (c *Client) RelatedKeywords(ctx context.Context, keyword string, location int, language string, limit int) ([]RelatedKeywordItem, error) {
	body := []map[string]any{
		{
			"keyword":       keyword,
			"location_code": location,
			"language_code": language,
			"limit":         limit,
		},
	}

	var resp TaskResponse
	if err := c.api.Post(ctx, "/dataforseo_labs/google/related_keywords/live", body, &resp); err != nil {
		return nil, fmt.Errorf("related keywords: %w", err)
	}

	results, err := c.unwrapResult(&resp)
	if err != nil {
		return nil, fmt.Errorf("related keywords: %w", err)
	}

	return extractKeywordItems(results)
}

// KeywordSuggestions retrieves keyword suggestions for the given seed keyword.
func (c *Client) KeywordSuggestions(ctx context.Context, keyword string, location int, language string, limit int) ([]RelatedKeywordItem, error) {
	body := []map[string]any{
		{
			"keyword":       keyword,
			"location_code": location,
			"language_code": language,
			"limit":         limit,
		},
	}

	var resp TaskResponse
	if err := c.api.Post(ctx, "/dataforseo_labs/google/keyword_suggestions/live", body, &resp); err != nil {
		return nil, fmt.Errorf("keyword suggestions: %w", err)
	}

	results, err := c.unwrapResult(&resp)
	if err != nil {
		return nil, fmt.Errorf("keyword suggestions: %w", err)
	}

	return extractKeywordItems(results)
}

// SerpGoogle retrieves Google organic SERP results for the given keyword.
func (c *Client) SerpGoogle(ctx context.Context, keyword string, location int, language string, device string) ([]SerpItem, error) {
	body := []map[string]any{
		{
			"keyword":       keyword,
			"location_code": location,
			"language_code": language,
			"device":        device,
		},
	}

	var resp TaskResponse
	if err := c.api.Post(ctx, "/serp/google/organic/live/regular", body, &resp); err != nil {
		return nil, fmt.Errorf("serp google: %w", err)
	}

	results, err := c.unwrapResult(&resp)
	if err != nil {
		return nil, fmt.Errorf("serp google: %w", err)
	}

	return extractSerpItems(results, "")
}

// SerpPAA retrieves People Also Ask items from Google SERP.
func (c *Client) SerpPAA(ctx context.Context, keyword string, location int, language string) ([]PAAItem, error) {
	body := []map[string]any{
		{
			"keyword":       keyword,
			"location_code": location,
			"language_code": language,
		},
	}

	var resp TaskResponse
	if err := c.api.Post(ctx, "/serp/google/organic/live/regular", body, &resp); err != nil {
		return nil, fmt.Errorf("serp paa: %w", err)
	}

	results, err := c.unwrapResult(&resp)
	if err != nil {
		return nil, fmt.Errorf("serp paa: %w", err)
	}

	return extractPAAItems(results)
}

func extractKeywordItems(results []TaskResult) ([]RelatedKeywordItem, error) {
	var items []RelatedKeywordItem

	for _, r := range results {
		for _, rawItem := range r.Items {
			itemBytes, err := json.Marshal(rawItem)
			if err != nil {
				continue
			}

			var parsed struct {
				KeywordData struct {
					Keyword          string  `json:"keyword"`
					SearchVolume     int     `json:"keyword_info.search_volume"`
					Competition      float64 `json:"keyword_info.competition"`
					CompetitionLevel string  `json:"keyword_info.competition_level"`
					CPC              float64 `json:"keyword_info.cpc"`
				} `json:"keyword_data"`
				KeywordInfo struct {
					SearchVolume     int     `json:"search_volume"`
					Competition      float64 `json:"competition"`
					CompetitionLevel string  `json:"competition_level"`
					CPC              float64 `json:"cpc"`
				} `json:"keyword_info"`
				Keyword string `json:"keyword"`
			}

			if err := json.Unmarshal(itemBytes, &parsed); err != nil {
				continue
			}

			kw := parsed.Keyword
			if kw == "" && parsed.KeywordData.Keyword != "" {
				kw = parsed.KeywordData.Keyword
			}

			items = append(items, RelatedKeywordItem{
				Keyword:          kw,
				SearchVolume:     parsed.KeywordInfo.SearchVolume,
				Competition:      parsed.KeywordInfo.Competition,
				CompetitionLevel: parsed.KeywordInfo.CompetitionLevel,
				CPC:              parsed.KeywordInfo.CPC,
			})
		}
	}

	return items, nil
}

func extractSerpItems(results []TaskResult, filterType string) ([]SerpItem, error) {
	var items []SerpItem

	for _, r := range results {
		for _, rawItem := range r.Items {
			itemBytes, err := json.Marshal(rawItem)
			if err != nil {
				continue
			}

			var parsed SerpItem
			if err := json.Unmarshal(itemBytes, &parsed); err != nil {
				continue
			}

			if filterType != "" && parsed.Type != filterType {
				continue
			}

			if parsed.Type == "organic" || filterType == "" {
				items = append(items, parsed)
			}
		}
	}

	return items, nil
}

func extractPAAItems(results []TaskResult) ([]PAAItem, error) {
	var items []PAAItem

	for _, r := range results {
		for _, rawItem := range r.Items {
			itemBytes, err := json.Marshal(rawItem)
			if err != nil {
				continue
			}

			var parsed struct {
				Type         string `json:"type"`
				RankGroup    int    `json:"rank_group"`
				RankAbsolute int    `json:"rank_absolute"`
				Position     string `json:"position"`
				Title        string `json:"title"`
				URL          string `json:"url"`
				Domain       string `json:"domain"`
				Description  string `json:"description"`
				Items        []struct {
					Type        string `json:"type"`
					Title       string `json:"title"`
					URL         string `json:"url"`
					Domain      string `json:"domain"`
					Description string `json:"description"`
				} `json:"items"`
			}

			if err := json.Unmarshal(itemBytes, &parsed); err != nil {
				continue
			}

			if parsed.Type != "people_also_ask" {
				continue
			}

			for _, sub := range parsed.Items {
				items = append(items, PAAItem{
					Type:        sub.Type,
					RankGroup:   parsed.RankGroup,
					RankAbsolute: parsed.RankAbsolute,
					Position:    parsed.Position,
					Title:       sub.Title,
					URL:         sub.URL,
					Domain:      sub.Domain,
					Description: sub.Description,
				})
			}

			if len(parsed.Items) == 0 {
				items = append(items, PAAItem{
					Type:        parsed.Type,
					RankGroup:   parsed.RankGroup,
					RankAbsolute: parsed.RankAbsolute,
					Position:    parsed.Position,
					Title:       parsed.Title,
					URL:         parsed.URL,
					Domain:      parsed.Domain,
					Description: parsed.Description,
				})
			}
		}
	}

	return items, nil
}
