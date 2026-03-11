package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReportRepository defines the interface for report data access.
type ReportRepository interface {
	GetDailyReport(ctx context.Context, date time.Time) (*DailyReport, error)
	GetWeeklyUsage(ctx context.Context, startDate time.Time) ([]UsageRow, error)
	GetMonthlyUsage(ctx context.Context, year int, month int) ([]UsageRow, error)
	GetStockMovement(ctx context.Context, itemID uuid.UUID, startDate, endDate time.Time) (*StockMovementReport, error)
}

// DailyReport represents a daily activity report.
type DailyReport struct {
	Date     string        `json:"date"`
	StockIns []DailyEntry  `json:"stock_ins"`
	StockOuts []DailyEntry `json:"stock_outs"`
}

// DailyEntry represents a stock in/out entry for daily report.
type DailyEntry struct {
	ItemID   uuid.UUID `json:"item_id"`
	ItemName string    `json:"item_name"`
	Unit     string    `json:"unit"`
	Quantity float64   `json:"quantity"`
}

// UsageRow represents an item usage row for weekly/monthly reports.
type UsageRow struct {
	ItemID       uuid.UUID `json:"item_id"`
	ItemName     string    `json:"item_name"`
	CategoryName string    `json:"category_name"`
	Unit         string    `json:"unit"`
	TotalUsed    float64   `json:"total_used"`
}

// StockMovementReport represents stock movement for a single item.
type StockMovementReport struct {
	ItemID       uuid.UUID           `json:"item_id"`
	ItemName     string              `json:"item_name"`
	Unit         string              `json:"unit"`
	Period       string              `json:"period"`
	TotalIn      float64             `json:"total_in"`
	TotalOut     float64             `json:"total_out"`
	NetChange    float64             `json:"net_change"`
	Movements    []MovementEntry     `json:"movements"`
}

// MovementEntry represents a single stock movement.
type MovementEntry struct {
	Date      time.Time `json:"date"`
	Type      string    `json:"type"` // "in" or "out"
	Quantity  float64   `json:"quantity"`
	Notes     *string   `json:"notes,omitempty"`
	CreatedBy string    `json:"created_by"`
}

type reportRepository struct {
	db *gorm.DB
}

// NewReportRepository creates a new ReportRepository.
func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) GetDailyReport(ctx context.Context, date time.Time) (*DailyReport, error) {
	nextDay := date.Add(24 * time.Hour)
	report := &DailyReport{Date: date.Format("2006-01-02")}

	// Stock Ins
	var stockIns []DailyEntry
	err := r.db.WithContext(ctx).
		Table("data.stock_in si").
		Select("si.item_id, i.name AS item_name, u.abbreviation AS unit, si.quantity").
		Joins("JOIN data.items i ON si.item_id = i.id").
		Joins("JOIN data.units u ON i.unit_id = u.id").
		Where("si.created_at >= ? AND si.created_at < ?", date, nextDay).
		Order("si.created_at ASC").
		Scan(&stockIns).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch daily stock ins: %w", err)
	}
	report.StockIns = stockIns

	// Stock Outs
	var stockOuts []DailyEntry
	err = r.db.WithContext(ctx).
		Table("data.stock_out so").
		Select("so.item_id, i.name AS item_name, u.abbreviation AS unit, so.quantity").
		Joins("JOIN data.items i ON so.item_id = i.id").
		Joins("JOIN data.units u ON i.unit_id = u.id").
		Where("so.created_at >= ? AND so.created_at < ?", date, nextDay).
		Order("so.created_at ASC").
		Scan(&stockOuts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch daily stock outs: %w", err)
	}
	report.StockOuts = stockOuts

	return report, nil
}

func (r *reportRepository) GetWeeklyUsage(ctx context.Context, startDate time.Time) ([]UsageRow, error) {
	endDate := startDate.Add(7 * 24 * time.Hour)
	return r.getUsageReport(ctx, startDate, endDate)
}

func (r *reportRepository) GetMonthlyUsage(ctx context.Context, year int, month int) ([]UsageRow, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)
	return r.getUsageReport(ctx, startDate, endDate)
}

func (r *reportRepository) getUsageReport(ctx context.Context, startDate, endDate time.Time) ([]UsageRow, error) {
	var rows []UsageRow
	err := r.db.WithContext(ctx).
		Table("data.stock_out so").
		Select("so.item_id, i.name AS item_name, c.name AS category_name, u.abbreviation AS unit, SUM(so.quantity) AS total_used").
		Joins("JOIN data.items i ON so.item_id = i.id").
		Joins("JOIN data.categories c ON i.category_id = c.id").
		Joins("JOIN data.units u ON i.unit_id = u.id").
		Where("so.created_at >= ? AND so.created_at < ?", startDate, endDate).
		Group("so.item_id, i.name, c.name, u.abbreviation").
		Order("total_used DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch usage report: %w", err)
	}
	return rows, nil
}

func (r *reportRepository) GetStockMovement(ctx context.Context, itemID uuid.UUID, startDate, endDate time.Time) (*StockMovementReport, error) {
	// Get item info
	var itemInfo struct {
		Name string
		Unit string
	}
	err := r.db.WithContext(ctx).
		Table("data.items i").
		Select("i.name, u.abbreviation AS unit").
		Joins("JOIN data.units u ON i.unit_id = u.id").
		Where("i.id = ?", itemID).
		Scan(&itemInfo).Error
	if err != nil {
		return nil, fmt.Errorf("item not found: %w", err)
	}

	report := &StockMovementReport{
		ItemID:   itemID,
		ItemName: itemInfo.Name,
		Unit:     itemInfo.Unit,
		Period:   fmt.Sprintf("%s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")),
	}

	// Stock In movements
	var inMovements []struct {
		Date      time.Time
		Quantity  float64
		Notes     *string
		FullName  string
	}
	err = r.db.WithContext(ctx).
		Table("data.stock_in si").
		Select("si.created_at AS date, si.quantity, si.notes, u.full_name").
		Joins("JOIN data.users u ON si.created_by = u.id").
		Where("si.item_id = ? AND si.created_at >= ? AND si.created_at < ?", itemID, startDate, endDate).
		Order("si.created_at ASC").
		Scan(&inMovements).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock in movements: %w", err)
	}

	for _, m := range inMovements {
		report.TotalIn += m.Quantity
		report.Movements = append(report.Movements, MovementEntry{
			Date: m.Date, Type: "in", Quantity: m.Quantity, Notes: m.Notes, CreatedBy: m.FullName,
		})
	}

	// Stock Out movements
	var outMovements []struct {
		Date     time.Time
		Quantity float64
		Notes    *string
		FullName string
	}
	err = r.db.WithContext(ctx).
		Table("data.stock_out so").
		Select("so.created_at AS date, so.quantity, so.notes, u.full_name").
		Joins("JOIN data.users u ON so.created_by = u.id").
		Where("so.item_id = ? AND so.created_at >= ? AND so.created_at < ?", itemID, startDate, endDate).
		Order("so.created_at ASC").
		Scan(&outMovements).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stock out movements: %w", err)
	}

	for _, m := range outMovements {
		report.TotalOut += m.Quantity
		report.Movements = append(report.Movements, MovementEntry{
			Date: m.Date, Type: "out", Quantity: m.Quantity, Notes: m.Notes, CreatedBy: m.FullName,
		})
	}

	report.NetChange = report.TotalIn - report.TotalOut
	return report, nil
}
