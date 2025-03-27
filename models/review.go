package models

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"strings"
	"sync"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/config"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/fiber-csv-app/utils"
	"github.com/jszwec/csvutil"
	"go.uber.org/zap"
)

type Review struct {
	App                   string  `csv:"App"`
	TranslatedReview      string  `csv:"Translated_Review"`
	Sentiment             string  `csv:"Sentiment"`
	SentimentPolarity     float64 `csv:"Sentiment_Polarity"`
	SentimentSubjectivity float64 `csv:"Sentiment_Subjectivity"`
}

type ReviewModel struct {
	logger *zap.Logger
	config config.AppConfig
}

// NewReviewModel initializes a new ReviewModel instance
func NewReviewModel(logger *zap.Logger, config config.AppConfig) *ReviewModel {
	return &ReviewModel{
		logger: logger,
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
	reviewMutex.Lock()
	defer reviewMutex.Unlock()

	reviews, err := rm.ParseReviews()
	if err != nil {
		return err
	}

	reviewCache = reviews
	return nil
}

// GetReviewsFromCache: Returns data from cache or loads it if expired
func (rm *ReviewModel) GetReviewsFromCache() ([]Review, error) {
	// First-time cache load
	reviewOnce.Do(func() {
		_ = rm.loadCache()
	})

	reviewMutex.RLock()
	if len(reviewCache) == 0 {
		reviewMutex.RUnlock()
		err := rm.loadCache()
		if err != nil {
			return nil, err
		}
		reviewMutex.RLock()
	}
	defer reviewMutex.RUnlock()

	return reviewCache, nil
}

// ParseReviews: Reads and parses reviews from CSV using csvutils.Unmarshal
func (rm *ReviewModel) ParseReviews() ([]Review, error) {
	records, err := utils.ReadCSV(rm.config.ReviewFilePath)
	if err != nil {
		rm.logger.Error("Error reading CSV file", zap.Error(err))
		return nil, errors.New("could not read CSV file")
	}

	// Convert records to CSV format
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	err = writer.WriteAll(records)
	if err != nil {
		return nil, err
	}
	writer.Flush()

	// Unmarshal CSV into struct
	var reviews []Review
	if err := csvutil.Unmarshal(buf.Bytes(), &reviews); err != nil {
		rm.logger.Error("Error unmarshaling CSV data", zap.Error(err))
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

	rm.logger.Info("Parsed Reviews",
		zap.Int("total_reviews", len(reviews)),
		zap.Int("valid_reviews", len(validReviews)),
	)

	return validReviews, nil
}

// Helper function to check if float is NaN
func isNaN(f float64) bool {
	return f != f
}

// GetReviews: Fetches reviews based on filters
func (rm *ReviewModel) GetReviews(appName, sentiment string, polarityMin, polarityMax float64) ([]Review, error) {
	reviews, err := rm.GetReviewsFromCache()
	if err != nil {
		return nil, err
	}

	// Extensive logging for debugging
	rm.logger.Info("Review Filtering Debug",
		zap.String("appName", appName),
		zap.String("sentiment", sentiment),
		zap.Float64("polarityMin", polarityMin),
		zap.Float64("polarityMax", polarityMax),
		zap.Int("total_reviews", len(reviews)),
	)

	// Detailed comparison with logging
	var filteredReviews []Review
	for _, review := range reviews {
		// Logging for each review
		debugFields := []zap.Field{
			zap.String("review_app", review.App),
			zap.String("review_sentiment", review.Sentiment),
			zap.Float64("review_polarity", review.SentimentPolarity),
		}

		// Flexible matching with case-insensitive and trimmed comparisons
		matchesApp := appName == "" ||
			strings.EqualFold(strings.TrimSpace(review.App), strings.TrimSpace(appName))

		matchesSentiment := sentiment == "" ||
			strings.EqualFold(strings.TrimSpace(review.Sentiment), strings.TrimSpace(sentiment))

		matchesPolarity := review.SentimentPolarity >= polarityMin &&
			review.SentimentPolarity <= polarityMax

		if matchesApp && matchesSentiment && matchesPolarity {
			filteredReviews = append(filteredReviews, review)
			rm.logger.Debug("Review Matched", append(debugFields,
				zap.Bool("app_match", matchesApp),
				zap.Bool("sentiment_match", matchesSentiment),
				zap.Bool("polarity_match", matchesPolarity),
			)...)
		}
	}

	// More detailed logging
	rm.logger.Info("Filtering Results",
		zap.Int("filtered_review_count", len(filteredReviews)),
	)

	if len(filteredReviews) == 0 {
		return nil, errors.New(fmt.Sprintf(
			"No reviews found matching: App=%s, Sentiment=%s, Polarity=%f-%f",
			appName, sentiment, polarityMin, polarityMax,
		))
	}

	return filteredReviews, nil
}
