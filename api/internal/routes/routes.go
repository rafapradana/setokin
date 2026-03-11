// Package routes defines all API route definitions.
package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/setokin/api/internal/handlers"
	"github.com/setokin/api/internal/middleware"
	"gorm.io/gorm"
)

// Handlers holds all handler instances for route registration.
type Handlers struct {
	Auth      *handlers.AuthHandler
	Category  *handlers.CategoryHandler
	Unit      *handlers.UnitHandler
	Item      *handlers.ItemHandler
	Supplier  *handlers.SupplierHandler
	StockIn   *handlers.StockInHandler
	StockOut  *handlers.StockOutHandler
	Batch     *handlers.BatchHandler
	Inventory *handlers.InventoryHandler
	Report    *handlers.ReportHandler
	Upload    *handlers.UploadHandler
}

// Setup registers all routes for the application.
func Setup(app *fiber.App, db *gorm.DB, jwtSecret string, h *Handlers) {
	// Health check (no auth)
	app.Get("/v1/health", handlers.HealthCheck(db))

	// API v1 group
	v1 := app.Group("/v1")

	// --- Auth routes (no auth required) ---
	auth := v1.Group("/auth")
	auth.Post("/register", h.Auth.Register)
	auth.Post("/login", h.Auth.Login)
	auth.Post("/refresh", h.Auth.Refresh)

	// Auth routes (auth required)
	authProtected := auth.Group("", middleware.AuthMiddleware(jwtSecret))
	authProtected.Post("/logout", h.Auth.Logout)
	authProtected.Get("/me", h.Auth.GetMe)

	// --- Protected routes ---
	protected := v1.Group("", middleware.AuthMiddleware(jwtSecret))

	// Categories
	categories := protected.Group("/categories")
	categories.Get("/", h.Category.List)
	categories.Get("/:id", h.Category.Get)
	categories.Post("/", middleware.RoleMiddleware("owner", "manager"), h.Category.Create)
	categories.Put("/:id", middleware.RoleMiddleware("owner", "manager"), h.Category.Update)
	categories.Delete("/:id", middleware.RoleMiddleware("owner", "manager"), h.Category.Delete)

	// Units
	protected.Get("/units", h.Unit.List)

	// Items
	items := protected.Group("/items")
	items.Get("/", h.Item.List)
	items.Get("/:id", h.Item.Get)
	items.Post("/", middleware.RoleMiddleware("owner", "manager"), h.Item.Create)
	items.Put("/:id", middleware.RoleMiddleware("owner", "manager"), h.Item.Update)
	items.Delete("/:id", middleware.RoleMiddleware("owner", "manager"), h.Item.Delete)

	// Item batches (nested under items)
	items.Get("/:item_id/batches", h.Batch.ListByItem)

	// Batches (direct access)
	protected.Get("/batches/:id", h.Batch.Get)

	// Suppliers
	suppliers := protected.Group("/suppliers")
	suppliers.Get("/", h.Supplier.List)
	suppliers.Post("/", middleware.RoleMiddleware("owner", "manager"), h.Supplier.Create)
	suppliers.Put("/:id", middleware.RoleMiddleware("owner", "manager"), h.Supplier.Update)
	suppliers.Delete("/:id", middleware.RoleMiddleware("owner", "manager"), h.Supplier.Delete)

	// Stock In
	stockIn := protected.Group("/stock-in")
	stockIn.Post("/", h.StockIn.Create)
	stockIn.Get("/", h.StockIn.List)

	// Stock Out
	stockOut := protected.Group("/stock-out")
	stockOut.Post("/", h.StockOut.Create)
	stockOut.Get("/", h.StockOut.List)
	stockOut.Get("/:id/details", h.StockOut.GetDetails)

	// Inventory
	inventory := protected.Group("/inventory")
	inventory.Get("/", h.Inventory.GetCurrentInventory)
	inventory.Get("/expiring", h.Inventory.GetExpiring)

	// Reports
	reports := protected.Group("/reports")
	reports.Get("/daily", h.Report.Daily)
	reports.Get("/weekly", h.Report.Weekly)
	reports.Get("/monthly", h.Report.Monthly)
	reports.Get("/stock-movement/:item_id", h.Report.StockMovement)

	// Uploads
	uploads := protected.Group("/uploads")
	uploads.Post("/request", h.Upload.RequestUpload)
	uploads.Post("/:upload_id/confirm", h.Upload.ConfirmUpload)
	uploads.Get("/:upload_id/download", h.Upload.GetDownloadURL)
}
