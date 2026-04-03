package service

import (
	"testing"
	"time"
)

func TestMonthsBeetween(t *testing.T) {
	tests := []struct {
		name     string
		from     time.Time
		to       time.Time
		expected int
	}{
		{
			name:     "один месяц",
			from:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			to:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		{
			name:     "весь год",
			from:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			to:       time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
			expected: 12,
		},
		{
			name:     "несколько месяцев",
			from:     time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
			to:       time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
			expected: 5,
		},
		{
			name:     "разные годы",
			from:     time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
			to:       time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			expected: 4,
		},
		{
			name:     "отрицательноый период",
			from:     time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
			to:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := monthsBeetween(tt.from, tt.to)
			if result != tt.expected {
				t.Errorf("monthBeetween(%v, %v) = %d, want %d", tt.from, tt.to, result, tt.expected)
			}
		})
	}
}
