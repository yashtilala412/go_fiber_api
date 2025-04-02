package models

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/config"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jszwec/csvutil"
)

type Review struct {
	App                   string  `csv:"App" validate:"required"`                       // Required field
	TranslatedReview      string  `csv:"Translated_Review" validate:"required"`         // Required field
	Sentiment             string  `csv:"Sentiment" validate:"required"`                 // Required field
	SentimentPolarity     float64 `csv:"Sentiment_Polarity" validate:"gte=-1,lte=1"`    // Must be between -1 and 1
	SentimentSubjectivity float64 `csv:"Sentiment_Subjectivity" validate:"gte=0,lte=1"` // Must be between 0 and 1
}

type ReviewModel struct {
	config config.AppConfig
}

// NewReviewModel initializes a new ReviewModel instance
// models/review.go

func NewReviewModel(config config.AppConfig) *ReviewModel {
	return &ReviewModel{
		config: config,
	}
}

// Global Cache Variables
var (
	reviewCache []Review
	reviewMutex sync.RWMutex
	reviewOnce  sync.Once
)

// loadCache: Loads review data into cache
func (rm *ReviewModel) loadCache() error {
	reviews, err := rm.ParseReviews()
	if err != nil {
		return err
	}
	reviewCache = reviews
	return nil
}

// ListReviewsFromCache: Returns data from cache or loads it if expired
func (rm *ReviewModel) ListReviewsFromCache() ([]Review, error) {
	// First-time cache load
	reviewOnce.Do(func() {
		_ = rm.loadCache()
	})
	reviewMutex.RLock()
	defer reviewMutex.RUnlock()
	if len(reviewCache) == 0 {
		err := rm.loadCache()
		if err != nil {
			return nil, err
		}
	}
	return reviewCache, nil
}

// ParseReviews: Reads and parses reviews from CSV using csvutils.Unmarshal
func (rm *ReviewModel) ParseReviews() ([]Review, error) {
	if rm.config.ReviewFilePath == "" {
		return nil, errors.New("REview file path is not configured")
	}
	records, err := utils.ReadCSV(rm.config.ReviewFilePath)
	if err != nil {
		return nil, errors.New("could not read CSV file")
	}

	// Unmarshal CSV into struct
	var reviews []Review
	if err := csvutil.Unmarshal(records, &reviews); err != nil {
		return nil, err
	}

	// Filter out invalid reviews
	var validReviews []Review
	for _, review := range reviews {
		if review.App != "" && review.App != "nan" &&
			review.Sentiment != "nan" &&
			!isNaN(review.SentimentPolarity) {
			validReviews = append(validReviews, review)
		}
	}
	return validReviews, nil
}
func (r *Review) Validate() error {
	return validate.Struct(r)
}

// Helper function to check if float is NaN
func isNaN(f float64) bool {
	return f != f
}
func (r *Review) ValidateApp() error {
	return validate.Struct(r)
}

// ListReviews: Fetches reviews based on filters
func (rm *ReviewModel) ListReviews(c *fiber.Ctx, appName, sentiment string, polarityMin, polarityMax float64) ([]Review, error) {
	reviews, err := rm.ListReviewsFromCache()
	if err != nil {
		return nil, err
	}
	var filteredReviews []Review
	for _, review := range reviews {
		matchesApp := appName == "" ||
			strings.EqualFold(strings.TrimSpace(review.App), strings.TrimSpace(appName))

		matchesSentiment := sentiment == "" ||
			strings.EqualFold(strings.TrimSpace(review.Sentiment), strings.TrimSpace(sentiment))

		matchesPolarity := review.SentimentPolarity >= polarityMin &&
			review.SentimentPolarity <= polarityMax

		if matchesApp && matchesSentiment && matchesPolarity {
			filteredReviews = append(filteredReviews, review)
		}
	}

	if len(filteredReviews) == 0 {
		return nil, fmt.Errorf(
			"No reviews found matching: App=%s, Sentiment=%s, Polarity=%f-%f",
			appName, sentiment, polarityMin, polarityMax,
		)
	}
	return filteredReviews, nil
}
func (rm *ReviewModel) AddReview(review Review) error {
	if err := review.ValidateApp(); err != nil {
		return err
	}
	reviewMutex.Lock()
	defer reviewMutex.Unlock()

	file, err := os.OpenFile(rm.config.ReviewFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{
		review.App,
		review.TranslatedReview,
		review.Sentiment,
		fmt.Sprintf("%f", review.SentimentPolarity),
		fmt.Sprintf("%f", review.SentimentSubjectivity),
	}

	if err := writer.Write(record); err != nil {
		return err
	}

	// Append to the in-memory cache
	reviewCache = append(reviewCache, review)

	return nil
}
