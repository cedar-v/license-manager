package service

import (
	"testing"
	"time"
)

func TestParseCuOrderTimeRange(t *testing.T) {
	loc := time.FixedZone("CST", 8*3600)
	now := time.Date(2026, 1, 14, 15, 4, 5, 0, loc) // Wednesday

	t.Run("empty", func(t *testing.T) {
		start, end, err := parseCuOrderTimeRange("", now)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if start != nil || end != nil {
			t.Fatalf("expected nil start/end, got %v %v", start, end)
		}
	})

	t.Run("today", func(t *testing.T) {
		start, end, err := parseCuOrderTimeRange("today", now)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		wantStart := time.Date(2026, 1, 14, 0, 0, 0, 0, loc)
		wantEnd := time.Date(2026, 1, 15, 0, 0, 0, 0, loc)
		if start == nil || end == nil || !start.Equal(wantStart) || !end.Equal(wantEnd) {
			t.Fatalf("unexpected range: start=%v end=%v", start, end)
		}
	})

	t.Run("this week", func(t *testing.T) {
		start, end, err := parseCuOrderTimeRange("week", now)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		wantStart := time.Date(2026, 1, 12, 0, 0, 0, 0, loc) // Monday
		wantEnd := time.Date(2026, 1, 19, 0, 0, 0, 0, loc)
		if start == nil || end == nil || !start.Equal(wantStart) || !end.Equal(wantEnd) {
			t.Fatalf("unexpected range: start=%v end=%v", start, end)
		}
	})

	t.Run("this month", func(t *testing.T) {
		start, end, err := parseCuOrderTimeRange("month", now)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		wantStart := time.Date(2026, 1, 1, 0, 0, 0, 0, loc)
		wantEnd := time.Date(2026, 2, 1, 0, 0, 0, 0, loc)
		if start == nil || end == nil || !start.Equal(wantStart) || !end.Equal(wantEnd) {
			t.Fatalf("unexpected range: start=%v end=%v", start, end)
		}
	})

	t.Run("last three months", func(t *testing.T) {
		start, end, err := parseCuOrderTimeRange("three_months", now)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		wantStart := time.Date(2025, 10, 14, 0, 0, 0, 0, loc)
		wantEnd := time.Date(2026, 1, 15, 0, 0, 0, 0, loc)
		if start == nil || end == nil || !start.Equal(wantStart) || !end.Equal(wantEnd) {
			t.Fatalf("unexpected range: start=%v end=%v", start, end)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		_, _, err := parseCuOrderTimeRange("nope", now)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
}
