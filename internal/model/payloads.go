package model

type InsightPayload struct {
	Text string `json:"text"`
}

type ChartPayload struct {
	Title string `json:"title"`
	Axes  struct {
		X string `json:"x"`
		Y string `json:"y"`
	} `json:"axes"`
	// Keep it flexible; data structure can vary
	Data any `json:"data"`
}

type AudiencePayload struct {
	Gender           []string `json:"gender,omitempty"`
	BirthCountry     []string `json:"birthCountry,omitempty"`
	AgeGroups        []string `json:"ageGroups,omitempty"`
	HoursSocialMedia []string `json:"hoursSocialMedia,omitempty"`
	Purchases        []string `json:"purchases,omitempty"`
}
