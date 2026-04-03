package model

import (
	"strings"
	"time"
)

type Subscription struct {
	ID          string     `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserID      string     `json:"user_id"`
	StartDate   MonthDate  `json:"start_date"`
	EndDate     *MonthDate `json:"end_date"`
}

type CreateUpdateSubscriptionInput struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
}

type MonthDate time.Time

type TotalCostResult struct {
	TotalCost int `json:"total_cost"`
}

func (m MonthDate) MarshalJSON() ([]byte, error) {
	t := time.Time(m)
	return []byte(`"` + t.Format("01-2006") + `"`), nil
}

func (m *MonthDate) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	t, err := time.Parse("01-2006", s)
	if err != nil {
		return err
	}
	*m = MonthDate(t)
	return nil
}

func (m *MonthDate) ToTime() *time.Time {
	if m == nil {
		return nil
	}
	t := time.Time(*m)
	return &t
}
