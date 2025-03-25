package models

type Review struct {
	App               string  `json:"app"`
	TranslatedReview  string  `json:"translated_review"`
	Sentiment         string  `json:"sentiment"`
	SentimentPolarity float64 `json:"sentiment_polarity"`
	SentimentSubject  float64 `json:"sentiment_subjectivity"`
}
