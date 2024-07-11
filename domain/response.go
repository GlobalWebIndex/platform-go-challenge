package domain

type AssetsResponse struct {
	Charts    []Chart    `json:"charts"`
	Insights  []Insight  `json:"insights"`
	Audiences []Audience `json:"audiences"`
}
