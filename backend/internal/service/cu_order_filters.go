package service

import (
	"fmt"
	"strings"
	"time"
)

func parseCuOrderTimeRange(value string, now time.Time) (*time.Time, *time.Time, error) {
	raw := strings.TrimSpace(value)
	if raw == "" {
		return nil, nil, nil
	}

	loc := now.Location()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	startOfTomorrow := startOfToday.AddDate(0, 0, 1)

	switch strings.ToLower(raw) {
	case "today", "今天":
		return &startOfToday, &startOfTomorrow, nil
	case "week", "this_week", "本周":
		// 默认按周一作为一周开始
		daysSinceMonday := (int(now.Weekday()) + 6) % 7
		start := startOfToday.AddDate(0, 0, -daysSinceMonday)
		end := start.AddDate(0, 0, 7)
		return &start, &end, nil
	case "month", "this_month", "本月":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		end := start.AddDate(0, 1, 0)
		return &start, &end, nil
	case "three_months", "3months", "last_3_months", "近三个月":
		threeMonthsAgo := now.AddDate(0, -3, 0)
		start := time.Date(threeMonthsAgo.Year(), threeMonthsAgo.Month(), threeMonthsAgo.Day(), 0, 0, 0, 0, loc)
		return &start, &startOfTomorrow, nil
	default:
		return nil, nil, fmt.Errorf("invalid time filter: %s", raw)
	}
}
