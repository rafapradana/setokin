package utils

import (
	"encoding/base64"
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// PaginationParams holds parsed pagination query parameters.
type PaginationParams struct {
	Limit  int
	Cursor string
}

// ParsePaginationParams extracts pagination parameters from a Fiber context.
func ParsePaginationParams(c *fiber.Ctx) PaginationParams {
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return PaginationParams{
		Limit:  limit,
		Cursor: c.Query("cursor", ""),
	}
}

// EncodeCursor encodes a map of values into a base64 cursor string.
func EncodeCursor(data map[string]string) string {
	bytes, _ := json.Marshal(data)
	return base64.URLEncoding.EncodeToString(bytes)
}

// DecodeCursor decodes a base64 cursor string into a map of values.
func DecodeCursor(cursor string) (map[string]string, error) {
	bytes, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}
	var data map[string]string
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}
	return data, nil
}
