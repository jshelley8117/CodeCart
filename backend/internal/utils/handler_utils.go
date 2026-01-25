package utils

import (
	"context"
	"math"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

func ParsePaginationInput(ctx context.Context, r *http.Request) (int, int, error) {
	zLog := FromContext(ctx, zap.NewNop())
	page := 1
	pageSize := 20

	if pageParam := r.URL.Query().Get("page"); pageParam != "" {
		parsed, err := strconv.Atoi(pageParam)
		if err != nil {
			zLog.Error("invalid page parameter", zap.Error(err))
			return 0, 0, err
		}
		if parsed <= 0 {
			zLog.Error("page must be greater than 0", zap.Error(err))
			return 0, 0, err
		}
		page = parsed
	}

	if pageSizeParam := r.URL.Query().Get("page_size"); pageSizeParam != "" {
		parsed, err := strconv.Atoi(pageSizeParam)
		if err != nil {
			zLog.Error("invalid page_size parameter", zap.Error(err))
			return 0, 0, err
		}
		if parsed <= 0 {
			zLog.Error("page_size must be greater than 0", zap.Error(err))
			return 0, 0, err
		}
		if parsed > 100 {
			zLog.Error("page_size must not exceed 100", zap.Error(err))
			return 0, 0, err
		}
		pageSize = parsed
	}

	return page, pageSize, nil
}

func CalculateTotalPages(total, pageSize int) int {
	return int(math.Ceil(float64(total) / float64(pageSize)))
}
