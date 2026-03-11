// Package main is the entry point for the Setokin API server.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/setokin/api/internal/config"
	"github.com/setokin/api/internal/database"
	"github.com/setokin/api/internal/handlers"
	"github.com/setokin/api/internal/middleware"
	minioclient "github.com/setokin/api/internal/minio"
	"github.com/setokin/api/internal/repositories"
	"github.com/setokin/api/internal/routes"
	"github.com/setokin/api/internal/services"
	"github.com/setokin/api/pkg/logger"
)

func main() {
	// Initialize logger
	logger.Init()
	defer logger.Sync()

	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg.DB)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to connect to database: %v", err))
		os.Exit(1)
	}
	logger.Info("Database connection established")

	// Initialize MinIO client
	minioClient, err := minioclient.NewClient(cfg.MinIO)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to initialize MinIO client: %v", err))
		os.Exit(1)
	}
	if err := minioClient.EnsureBucket(context.Background()); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to ensure MinIO bucket: %v", err))
		os.Exit(1)
	}
	logger.Info("MinIO connection established")

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ErrorHandler: customErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(middleware.SecurityHeaders())
	app.Use(middleware.CORSMiddleware(cfg.CORS.AllowedOrigins))
	app.Use(middleware.LoggerMiddleware())
	app.Use(middleware.RateLimitMiddleware())

	// --- Initialize repositories ---
	userRepo := repositories.NewUserRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	unitRepo := repositories.NewUnitRepository(db)
	itemRepo := repositories.NewItemRepository(db)
	supplierRepo := repositories.NewSupplierRepository(db)
	stockRepo := repositories.NewStockRepository(db)
	inventoryRepo := repositories.NewInventoryRepository(db)
	reportRepo := repositories.NewReportRepository(db)
	uploadRepo := repositories.NewUploadRepository(db)

	// --- Initialize services ---
	authService := services.NewAuthService(userRepo, cfg.JWT)
	categoryService := services.NewCategoryService(categoryRepo)
	unitService := services.NewUnitService(unitRepo)
	itemService := services.NewItemService(itemRepo, categoryRepo, unitRepo)
	supplierService := services.NewSupplierService(supplierRepo)
	stockService := services.NewStockService(stockRepo, itemRepo, db)
	uploadService := services.NewUploadService(uploadRepo, minioClient)

	// --- Initialize handlers ---
	h := &routes.Handlers{
		Auth:      handlers.NewAuthHandler(authService),
		Category:  handlers.NewCategoryHandler(categoryService),
		Unit:      handlers.NewUnitHandler(unitService),
		Item:      handlers.NewItemHandler(itemService),
		Supplier:  handlers.NewSupplierHandler(supplierService),
		StockIn:   handlers.NewStockInHandler(stockService),
		StockOut:  handlers.NewStockOutHandler(stockService),
		Batch:     handlers.NewBatchHandler(stockService),
		Inventory: handlers.NewInventoryHandler(inventoryRepo),
		Report:    handlers.NewReportHandler(reportRepo),
		Upload:    handlers.NewUploadHandler(uploadService),
	}

	// Setup routes
	routes.Setup(app, db, cfg.JWT.Secret, h)

	// Start server
	port := cfg.App.Port
	logger.Info(fmt.Sprintf("Server starting on port %s", port))
	if err := app.Listen(":" + port); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to start server: %v", err))
		os.Exit(1)
	}
}

// customErrorHandler provides a global error handler for unhandled errors.
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": fiber.Map{
			"code":    "internal_error",
			"message": err.Error(),
		},
	})
}
