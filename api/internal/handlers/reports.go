package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/setokin/api/internal/repositories"
	"github.com/setokin/api/internal/utils"
	apperrors "github.com/setokin/api/pkg/errors"
)

// ReportHandler handles report endpoints.
type ReportHandler struct {
	repo repositories.ReportRepository
}

// NewReportHandler creates a new ReportHandler.
func NewReportHandler(repo repositories.ReportRepository) *ReportHandler {
	return &ReportHandler{repo: repo}
}

// Daily handles GET /reports/daily.
func (h *ReportHandler) Daily(c *fiber.Ctx) error {
	dateStr := c.Query("date", time.Now().Format("2006-01-02"))
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid date format (YYYY-MM-DD)"))
	}

	report, err := h.repo.GetDailyReport(c.Context(), date)
	if err != nil {
		return utils.ErrorResponse(c, apperrors.ErrInternal.WithMessage("Failed to generate daily report"))
	}

	if report.StockIns == nil {
		report.StockIns = []repositories.DailyEntry{}
	}
	if report.StockOuts == nil {
		report.StockOuts = []repositories.DailyEntry{}
	}

	return utils.SuccessResponse(c, fiber.StatusOK, report)
}

// Weekly handles GET /reports/weekly.
func (h *ReportHandler) Weekly(c *fiber.Ctx) error {
	dateStr := c.Query("start_date")
	var startDate time.Time
	if dateStr != "" {
		parsed, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid date format (YYYY-MM-DD)"))
		}
		startDate = parsed
	} else {
		// Default: start of current week (Monday)
		now := time.Now()
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		startDate = now.AddDate(0, 0, -(weekday - 1)).Truncate(24 * time.Hour)
	}

	rows, err := h.repo.GetWeeklyUsage(c.Context(), startDate)
	if err != nil {
		return utils.ErrorResponse(c, apperrors.ErrInternal.WithMessage("Failed to generate weekly report"))
	}

	if rows == nil {
		rows = []repositories.UsageRow{}
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"period":     startDate.Format("2006-01-02") + " to " + startDate.Add(7*24*time.Hour).Format("2006-01-02"),
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   startDate.Add(7 * 24 * time.Hour).Format("2006-01-02"),
		"usage":      rows,
	})
}

// Monthly handles GET /reports/monthly.
func (h *ReportHandler) Monthly(c *fiber.Ctx) error {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	if y := c.Query("year"); y != "" {
		parsed, err := strconv.Atoi(y)
		if err == nil {
			year = parsed
		}
	}
	if m := c.Query("month"); m != "" {
		parsed, err := strconv.Atoi(m)
		if err == nil && parsed >= 1 && parsed <= 12 {
			month = parsed
		}
	}

	rows, err := h.repo.GetMonthlyUsage(c.Context(), year, month)
	if err != nil {
		return utils.ErrorResponse(c, apperrors.ErrInternal.WithMessage("Failed to generate monthly report"))
	}

	if rows == nil {
		rows = []repositories.UsageRow{}
	}

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"period":     startDate.Format("January 2006"),
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
		"usage":      rows,
	})
}

// StockMovement handles GET /reports/stock-movement/:item_id.
func (h *ReportHandler) StockMovement(c *fiber.Ctx) error {
	itemID, err := uuid.Parse(c.Params("item_id"))
	if err != nil {
		return utils.ErrorResponse(c, apperrors.ErrValidation.WithMessage("Invalid item ID"))
	}

	now := time.Now()
	startDate := now.AddDate(0, -1, 0).Truncate(24 * time.Hour) // default: last 30 days
	endDate := now.Truncate(24 * time.Hour).Add(24 * time.Hour)

	if s := c.Query("start_date"); s != "" {
		parsed, err := time.Parse("2006-01-02", s)
		if err == nil {
			startDate = parsed
		}
	}
	if e := c.Query("end_date"); e != "" {
		parsed, err := time.Parse("2006-01-02", e)
		if err == nil {
			endDate = parsed.Add(24 * time.Hour)
		}
	}

	report, err := h.repo.GetStockMovement(c.Context(), itemID, startDate, endDate)
	if err != nil {
		return utils.ErrorResponse(c, apperrors.ErrInternal.WithMessage("Failed to generate stock movement report"))
	}

	if report.Movements == nil {
		report.Movements = []repositories.MovementEntry{}
	}

	return utils.SuccessResponse(c, fiber.StatusOK, report)
}
