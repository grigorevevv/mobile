package service

import (
	"context"
	"mobile/internal/api_error"
	"mobile/internal/domain"
	"mobile/internal/repository"

	"github.com/google/uuid"
)

type SubscriptionService struct {
	subscriptionRepo *repository.SubscriptionRepo
}

func NewSubscriptionService(r *repository.SubscriptionRepo) *SubscriptionService {
	return &SubscriptionService{subscriptionRepo: r}
}

func (s *SubscriptionService) SubscriptionCreate(ctx context.Context, subscription *domain.Subscription) error {
	return s.subscriptionRepo.Create(ctx, subscription)
}

func (s *SubscriptionService) GetById(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	return s.subscriptionRepo.GetById(ctx, id)
}

// Update обновляет подписку.
func (s *SubscriptionService) UpdateSubscription(ctx context.Context, sub *domain.Subscription) error {
	// Выделяем только изменяемые поля
	fields := map[string]interface{}{}

	if sub.ServiceName != "" {
		fields["service_name"] = sub.ServiceName
	}
	if sub.Price != 0 {
		fields["price"] = sub.Price
	}
	if !sub.StartDate.IsZero() {
		fields["start_date"] = sub.StartDate
	}
	if sub.EndDate != nil && !sub.EndDate.IsZero() {
		fields["end_date"] = sub.EndDate
	}

	// Если нет полей для обновления, ничего не делаем
	if len(fields) == 0 {
		return nil
	}

	return s.subscriptionRepo.Update(ctx, sub, fields)
}

func (s *SubscriptionService) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	sub, err := s.subscriptionRepo.GetById(ctx, id)
	if err != nil {
		return err
	}
	if sub.ID == uuid.Nil {
		return api_error.ErrEntity
	}

	return s.subscriptionRepo.Delete(ctx, id)
}

func (s *SubscriptionService) ReportSubscription(ctx context.Context, params *repository.ListParams) ([]repository.ListReport, error) {
	data, err := s.subscriptionRepo.Report(ctx, params)
	if err != nil {
		return nil, err
	}

	return data, nil
}
