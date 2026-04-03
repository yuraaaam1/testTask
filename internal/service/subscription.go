package service

import (
	"fmt"
	"time"

	"github.com/yuraaaam1/testTask/internal/logger"
	"github.com/yuraaaam1/testTask/internal/model"
	"github.com/yuraaaam1/testTask/internal/repository"
	"go.uber.org/zap"
)

type SubscriptionService struct {
	repo *repository.SubscriptionRepository
}

func NewSubscriptionService(repo *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) Create(input model.CreateUpdateSubscriptionInput) (*model.Subscription, error) {
	sub, err := s.repo.Create(input)
	if err != nil {
		logger.Log.Error("failed to create subscription", zap.Error(err))
		return nil, err
	}

	logger.Log.Info("subscription created", zap.String("id", sub.ID))
	return sub, nil
}

func (s *SubscriptionService) GetByID(id string) (*model.Subscription, error) {
	sub, err := s.repo.GetByID(id)
	if err != nil {
		logger.Log.Error("failed to get subscription", zap.Error(err))
		return nil, err
	}

	logger.Log.Info("subscription succesfully get", zap.String("id", id))
	return sub, nil
}

func (s *SubscriptionService) List() ([]*model.Subscription, error) {
	subs, err := s.repo.List()
	if err != nil {
		logger.Log.Error("failed to get subscriptions", zap.Error(err))
		return nil, err
	}

	logger.Log.Info("subscriptions succesfully get", zap.Int("count", len(subs)))
	return subs, nil
}

func (s *SubscriptionService) Update(id string, input model.CreateUpdateSubscriptionInput) (*model.Subscription, error) {
	sub, err := s.repo.Update(id, input)
	if err != nil {
		logger.Log.Error("failet to update subscription", zap.Error(err))
		return nil, err
	}

	logger.Log.Info("subscription succesfully update", zap.String("id", id))
	return sub, nil
}

func (s *SubscriptionService) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		logger.Log.Error("failed to delete subscription", zap.Error(err))
		return err
	}

	logger.Log.Info("subscription succesfully delete")
	return nil
}

func (s *SubscriptionService) CalculateTotalCost(userID, serviceName, dateFrom, dateTo string) (*model.TotalCostResult, error) {
	subs, err := s.repo.GetByFilter(userID, serviceName, dateFrom, dateTo)
	if err != nil {
		logger.Log.Error("failed to get subscriptions by filter", zap.Error(err))
		return nil, err
	}

	from, err := time.Parse("01-2006", dateFrom)
	if err != nil {
		return nil, fmt.Errorf("invalid date_from format: %w", err)
	}

	to, err := time.Parse("01-2006", dateTo)
	if err != nil {
		return nil, fmt.Errorf("invalid date_to format: %w", err)
	}

	total := 0

	for _, sub := range subs {
		start := sub.StartDate.ToTime()
		var end time.Time
		if sub.EndDate != nil {
			end = *sub.EndDate.ToTime()
		} else {
			end = to
		}

		effectiveStart := maxTime(*start, from)
		effectiveEnd := minTime(end, to)

		months := monthsBeetween(effectiveStart, effectiveEnd)
		total += sub.Price * months
	}

	logger.Log.Info("total cost calculated",
		zap.String("user_id", userID),
		zap.String("service_name", serviceName),
		zap.String("date_from", dateFrom),
		zap.String("date_to", dateTo),
		zap.Int("total_cost", total),
	)

	return &model.TotalCostResult{TotalCost: total}, nil
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func monthsBeetween(from, to time.Time) int {
	months := ((to.Year()-from.Year())*12 + int(to.Month()) - int(from.Month()) + 1)
	if months < 0 {
		return 0
	}
	return months
}
