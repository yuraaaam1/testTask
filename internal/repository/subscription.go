package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/yuraaaam1/testTask/internal/model"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(input model.CreateUpdateSubscriptionInput) (*model.Subscription, error) {
	startDate, err := time.Parse("01-2006", input.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format: %w", err)
	}

	var endDate *time.Time
	if input.EndDate != "" {
		t, err := time.Parse("01-2006", input.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date format: %w", err)
		}
		endDate = &t
	}

	var start time.Time
	var end *time.Time
	sub := &model.Subscription{}
	err = r.db.QueryRow(
		`INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, service_name, price, user_id, start_date, end_date`,
		input.ServiceName, input.Price, input.UserID, startDate, endDate,
	).Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &start, &end)

	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	sub.StartDate = model.MonthDate(start)
	if end != nil {
		md := model.MonthDate(*end)
		sub.EndDate = &md
	}

	return sub, nil
}

func (r *SubscriptionRepository) GetByID(id string) (*model.Subscription, error) {
	var start time.Time
	var end *time.Time
	sub := &model.Subscription{}
	err := r.db.QueryRow(
		`SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions WHERE id = $1`, id,
	).Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &start, &end)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	sub.StartDate = model.MonthDate(start)
	if end != nil {
		md := model.MonthDate(*end)
		sub.EndDate = &md
	}

	return sub, nil

}

func (r *SubscriptionRepository) List() ([]*model.Subscription, error) {
	rows, err := r.db.Query(
		`SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions`)

	if err != nil {
		return nil, fmt.Errorf("failed to get list subscriptions: %w", err)
	}
	defer rows.Close()

	var start time.Time
	var end *time.Time
	var subs []*model.Subscription
	for rows.Next() {
		sub := &model.Subscription{}
		if err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &start, &end); err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		sub.StartDate = model.MonthDate(start)
		if end != nil {
			md := model.MonthDate(*end)
			sub.EndDate = &md
		}

		subs = append(subs, sub)
	}

	return subs, nil
}

func (r *SubscriptionRepository) Update(id string, input model.CreateUpdateSubscriptionInput) (*model.Subscription, error) {
	existing, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, nil
	}

	if input.ServiceName != "" {
		existing.ServiceName = input.ServiceName
	}

	if input.Price != 0 {
		existing.Price = input.Price
	}

	if input.UserID != "" {
		existing.UserID = input.UserID
	}

	if input.StartDate != "" {
		t, err := time.Parse("01-2006", input.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start_date formate: %w", err)
		}
		existing.StartDate = model.MonthDate(t)
	}

	if input.EndDate != "" {
		t, err := time.Parse("01-2006", input.EndDate)
		if err != nil {
			return nil, fmt.Errorf("ivalid end_date formate: %w", err)
		}
		md := model.MonthDate(t)
		existing.EndDate = &md
	}

	var start time.Time
	var end *time.Time
	sub := &model.Subscription{}
	err = r.db.QueryRow(
		`UPDATE subscriptions SET
		service_name = $1,
		price = $2,
		user_id = $3,
		start_date = $4,
		end_date = $5
		WHERE id = $6
		RETURNING id, service_name, price, user_id, start_date, end_date`,
		existing.ServiceName, existing.Price, existing.UserID, *existing.StartDate.ToTime(), existing.EndDate.ToTime(), id,
	).Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &start, &end)

	if err != nil {
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	sub.StartDate = model.MonthDate(start)
	if end != nil {
		md := model.MonthDate(*end)
		sub.EndDate = &md
	}

	return sub, nil
}

func (r *SubscriptionRepository) Delete(id string) error {
	result, err := r.db.Exec(
		`DELETE FROM subscriptions
		WHERE id = $1`, id)

	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

func (r *SubscriptionRepository) GetByFilter(userID, serviceName, dateFrom, dateTo string) ([]*model.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date
			FROM subscriptions WHERE 1=1`

	args := []interface{}{}
	argNum := 1

	if userID != "" {
		query += fmt.Sprintf(" AND user_id = $%d", argNum)
		args = append(args, userID)
		argNum++
	}

	if serviceName != "" {
		query += fmt.Sprintf(" AND service_name = $%d", argNum)
		args = append(args, serviceName)
		argNum++
	}

	if dateFrom != "" {
		from, err := time.Parse("01-2006", dateFrom)
		if err != nil {
			return nil, fmt.Errorf("invalid date_from format: %w", err)
		}

		query += fmt.Sprintf(" AND (end_date IS NULL OR end_date >= $%d)", argNum)
		args = append(args, from)
		argNum++
	}

	if dateTo != "" {
		to, err := time.Parse("01-2006", dateTo)
		if err != nil {
			return nil, fmt.Errorf("invalid date_to format: %w", err)
		}

		query += fmt.Sprintf("  AND start_date <= $%d", argNum)
		args = append(args, to)
		argNum++
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions by filter: %w", err)
	}
	defer rows.Close()

	var subs []*model.Subscription
	for rows.Next() {
		var start time.Time
		var end *time.Time
		sub := &model.Subscription{}

		if err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &start, &end); err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		sub.StartDate = model.MonthDate(start)
		if end != nil {
			md := model.MonthDate(*end)
			sub.EndDate = &md
		}

		subs = append(subs, sub)
	}

	return subs, nil
}
