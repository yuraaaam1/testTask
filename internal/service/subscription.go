package service

import (
	"github.com/yuraaaam1/testTask/internal/model"
	"github.com/yuraaaam1/testTask/internal/repository"
)

type SubscriptionService struct {
	repo *repository.SubscriptionRepository
}

func NewSubscriptionService(repo *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) Create(input model.CreateUpdateSubscriptionInput) (*model.Subscription, error) {
	return s.repo.Create(input)
}

func (s *SubscriptionService) GetByID(id string) (*model.Subscription, error) {
	return s.repo.GetByID(id)
}

func (s *SubscriptionService) List() ([]*model.Subscription, error) {
	return s.repo.List()
}

func (s *SubscriptionService) Update(id string, input model.CreateUpdateSubscriptionInput) (*model.Subscription, error) {
	return s.repo.Update(id, input)
}

func (s *SubscriptionService) Delete(id string) error {
	return s.repo.Delete(id)
}
