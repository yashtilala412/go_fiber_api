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

// ListReviewsFromCache: Returns data from cache or loads it if expired
func (rm *ReviewModel) ListReviewsFromCache() ([]Review, error) {
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
		defer reviewMutex.RUnlock()
	}

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

// ListReviews: Fetches reviews based on filters
func (rm *ReviewModel) ListReviews(c *fiber.Ctx, appName, sentiment string, polarityMin, polarityMax float64) ([]Review, error) {
	reviews, err := rm.ListReviewsFromCache()
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
		utils.JSONSuccess(c, fiber.StatusNotFound, nil)
		return nil, errors.New(fmt.Sprintf(
			"No reviews found matching: App=%s, Sentiment=%s, Polarity=%f-%f",
			appName, sentiment, polarityMin, polarityMax,
		))
	}

	return filteredReviews, nil
}

// AddReview adds a new review to the CSV and updates the cache.
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

// DeleteReviewByAppName deletes all reviews for a given app name.
func (rm *ReviewModel) DeleteReviewByAppName(appName string) error {
	reviewMutex.Lock()
	defer reviewMutex.Unlock()

	// 1. Read all reviews from CSV
	reviews, err := rm.ParseReviews()
	if err != nil {
		return err
	}

	// 2. Filter out reviews to be deleted
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
		return errors.New("No reviews found for app: " + appName)
	}

	// 3. Write the updated review list back to CSV using csvutil.Marshal
	file, err := os.Create(rm.config.ReviewFilePath) // Create overwrites the file
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
