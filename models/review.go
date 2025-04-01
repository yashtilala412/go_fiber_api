package models

import (
	"bytes"
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
	App                   string  `csv:"App"`
	TranslatedReview      string  `csv:"Translated_Review"`
	Sentiment             string  `csv:"Sentiment"`
	SentimentPolarity     float64 `csv:"Sentiment_Polarity"`
	SentimentSubjectivity float64 `csv:"Sentiment_Subjectivity"`
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

// Helper function to check if float is NaN
func isNaN(f float64) bool {
	return f != f
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
func (rm *ReviewModel) DeleteReview(appName string) error {
	reviewMutex.Lock()
	defer reviewMutex.Unlock()

	// 1. Read all reviews from CSV
	reviews, err := rm.ParseReviews()
	if err != nil {
		return err
	}

	// 2. Filter out the reviews with matching app name
	var updatedReviews []Review
	found := false
	for _, review := range reviews {

		if !strings.EqualFold(strings.TrimSpace(review.App), strings.TrimSpace(appName)) {
			updatedReviews = append(updatedReviews, review)
		} else {
			found = true
		}
	}

	if !found {
		return errors.New("App not found")
	}

	// 3. Rewrite the CSV file
	file, err := os.Create(rm.config.ReviewFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	csvBytes, err := csvutil.Marshal(updatedReviews)
	if err != nil {
		return err
	}

	r := csv.NewReader(bytes.NewReader(csvBytes))
	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	err = writer.WriteAll(records)
	if err != nil {
		return err
	}

	// 4. Update the in-memory cache
	reviewCache = updatedReviews

	return nil
}
