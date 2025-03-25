package cli

import (
	"fmt"
	"log"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/routes"
	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/services"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// rootCmd is the main CLI command
var rootCmd = &cobra.Command{
	Use:   "fiber-csv-app",
	Short: "CLI for managing CSV operations",
}

// apiCmd starts the API when running `go run main.go api`
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the API server",
	Run: func(cmd *cobra.Command, args []string) {
		logger, _ := zap.NewProduction()
		defer logger.Sync()

		// Load CSV data before starting the API
		appData, reviewData, err := services.LoadAppData(logger)
		if err != nil {
			logger.Error("Failed to load CSV data", zap.Error(err))
		} else {
			logger.Info("CSV data loaded successfully", zap.Int("appRows", len(appData)), zap.Int("reviewRows", len(reviewData)))
		}

		app := fiber.New()

		// Register API routes
		routes.RegisterRoutes(app, logger)

		port := viper.GetString("PORT")
		log.Println("Starting API on port", port)
		log.Fatal(app.Listen(":" + port))
	},
}

// loadCmd loads CSV data when running `go run main.go load`
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load CSV data",
	Run: func(cmd *cobra.Command, args []string) {
		logger, _ := zap.NewProduction()
		defer logger.Sync()

		appData, reviewData, err := services.LoadAppData(logger)
		if err != nil {
			log.Fatalf("Error loading data: %v", err)
		}

		fmt.Println("App Data:", appData)
		fmt.Println("Review Data:", reviewData)
	},
}

// StartCLI initializes the CLI
func StartCLI() {
	rootCmd.AddCommand(apiCmd)  // Register API command
	rootCmd.AddCommand(loadCmd) // Register Load command

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
