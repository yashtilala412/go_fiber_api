package models

// UserReview represents a user review for an app.
type UserReview struct {
	AppName               string  `json:"app_name"`
	TranslatedReview      string  `json:"translated_review"`
	Sentiment             string  `json:"sentiment"`
	SentimentPolarity     *string `json:"sentiment_polarity,omitempty"`
	SentimentSubjectivity *string `json:"sentiment_subjectivity,omitempty"`
}
