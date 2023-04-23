package models

type AssetResponse struct {
	UserId    int        `json:"user_id"`
	Charts    []Chart    `json:"charts"`
	Insights  []Insight  `json:"insights"`
	Audiences []Audience `json:"audiences"`
}
