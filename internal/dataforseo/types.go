package dataforseo

// TaskResponse is the top-level response envelope from the DataForSEO API.
type TaskResponse struct {
	Version        string `json:"version"`
	StatusCode     int    `json:"status_code"`
	StatusMessage  string `json:"status_message"`
	Time           string `json:"time"`
	Cost           float64 `json:"cost"`
	TasksCount     int    `json:"tasks_count"`
	TasksError     int    `json:"tasks_error"`
	Tasks          []Task `json:"tasks"`
}

// Task is a single task in the DataForSEO response.
type Task struct {
	ID            string          `json:"id"`
	StatusCode    int             `json:"status_code"`
	StatusMessage string          `json:"status_message"`
	Time          string          `json:"time"`
	Cost          float64         `json:"cost"`
	ResultCount   int             `json:"result_count"`
	Path          []string        `json:"path"`
	Data          any             `json:"data"`
	Result        []TaskResult    `json:"result"`
}

// TaskResult holds result data; the shape varies by endpoint.
type TaskResult struct {
	Keyword         string           `json:"keyword,omitempty"`
	LocationCode    int              `json:"location_code,omitempty"`
	LanguageCode    string           `json:"language_code,omitempty"`
	SearchVolume    int              `json:"search_volume,omitempty"`
	Competition     float64          `json:"competition,omitempty"`
	CompetitionLevel string          `json:"competition_level,omitempty"`
	CPC             float64          `json:"cpc,omitempty"`
	MonthlySearches []MonthlySearch  `json:"monthly_searches,omitempty"`
	KeywordDifficulty int            `json:"keyword_difficulty,omitempty"`
	SeedKeyword     string           `json:"seed_keyword,omitempty"`
	ItemsCount      int              `json:"items_count,omitempty"`
	Items           []any            `json:"items,omitempty"`
	TotalCount      int              `json:"total_count,omitempty"`
	CheckURL        string           `json:"check_url,omitempty"`
	Datetime        string           `json:"datetime,omitempty"`
	SpellCorrection any              `json:"spell,omitempty"`
	RefinementChips any              `json:"refinement_chips,omitempty"`
	ItemTypes       []string         `json:"item_types,omitempty"`
	SEDomain        string           `json:"se_domain,omitempty"`
	Type            string           `json:"type,omitempty"`
}

// MonthlySearch represents monthly search volume data.
type MonthlySearch struct {
	Year         int `json:"year"`
	Month        int `json:"month"`
	SearchVolume int `json:"search_volume"`
}

// KeywordData represents keyword search volume data for display.
type KeywordData struct {
	Keyword          string  `json:"keyword"`
	SearchVolume     int     `json:"search_volume"`
	Competition      float64 `json:"competition"`
	CompetitionLevel string  `json:"competition_level"`
	CPC              float64 `json:"cpc"`
	LocationCode     int     `json:"location_code"`
	LanguageCode     string  `json:"language_code"`
}

// DifficultyData represents keyword difficulty data for display.
type DifficultyData struct {
	Keyword           string `json:"keyword"`
	KeywordDifficulty int    `json:"keyword_difficulty"`
	LocationCode      int    `json:"location_code"`
	LanguageCode      string `json:"language_code"`
}

// RelatedKeywordItem represents a related keyword.
type RelatedKeywordItem struct {
	Keyword          string  `json:"keyword"`
	SearchVolume     int     `json:"search_volume"`
	Competition      float64 `json:"competition"`
	CompetitionLevel string  `json:"competition_level"`
	CPC              float64 `json:"cpc"`
}

// SerpItem represents a single SERP result.
type SerpItem struct {
	Type        string `json:"type"`
	RankGroup   int    `json:"rank_group"`
	RankAbsolute int   `json:"rank_absolute"`
	Position    string `json:"position"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Domain      string `json:"domain"`
	Description string `json:"description"`
	Breadcrumb  string `json:"breadcrumb"`
}

// PAAItem represents a People Also Ask item.
type PAAItem struct {
	Type        string `json:"type"`
	RankGroup   int    `json:"rank_group"`
	RankAbsolute int   `json:"rank_absolute"`
	Position    string `json:"position"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Domain      string `json:"domain"`
	Description string `json:"description"`
}
