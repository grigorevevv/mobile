package repository

import (
	"context"
	"fmt"
	"mobile/internal/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionRepo struct {
	db *gorm.DB
}

func NewSubscriptionRepo(db *gorm.DB) *SubscriptionRepo {
	return &SubscriptionRepo{db: db}
}

func (r *SubscriptionRepo) Create(ctx context.Context, s *domain.Subscription) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *SubscriptionRepo) GetById(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	var Subscr domain.Subscription

	err := r.db.WithContext(ctx).Find(&Subscr, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &Subscr, nil
}

func (r *SubscriptionRepo) Update(ctx context.Context, sub *domain.Subscription, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&domain.Subscription{}).Where("id = ?", sub.ID).Updates(fields).Error
}

func (r *SubscriptionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Subscription{}).Error
}

type ListReport struct {
	UserId      uuid.UUID
	ServiceName string
	Sum         int
}

func (r *SubscriptionRepo) Report(ctx context.Context, params *ListParams) ([]ListReport, error) {
	rows := make([]ListReport, 0)

	err := r.db.WithContext(ctx).
		Raw("select s.user_id, s.service_name, sum(s.price) \n"+
			"from subscriptions s\n"+
			"where s.start_date >= ?\n"+
			"and (s.end_date <= ? or s.end_date is NULL)\n"+
			"and s.user_id  = ?\n"+
			"and s.service_name = ?\n"+
			"group by s.user_id, s.service_name\n",
			params.StartDate, params.EndDate, params.IdUser, params.ServiceName).
		Scan(&rows).Error

	if err != nil {
		return nil, fmt.Errorf("Ошибка формирования отчета: %w", err)
	}

	fmt.Println(rows)
	return rows, nil
}

type ListParams struct {
	StartDate   time.Time
	EndDate     time.Time
	IdUser      uuid.UUID
	ServiceName string
}
