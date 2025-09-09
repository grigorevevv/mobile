package mapper

import (
	"mobile/internal/domain"
	"time"

	"github.com/google/uuid"
)

type SubscriptionMap struct {
	ServiceName string    `gorm:"" json:"service_name"`
	Price       int       `gorm:"" json:"price"`
	UserID      uuid.UUID `gorm:"" json:"user_id"`
	StartDate   string    `gorm:"" json:"start_date"`
	EndDate     string    `gorm:"" json:"end_date"`
}

func SubscriptionMapper(sm SubscriptionMap) (*domain.Subscription, error) {
	stDate, err := timeParse(sm.StartDate)
	if err != nil {
		return nil, err
	}

	endDate, err := timeParse(sm.EndDate)
	if err != nil {
		return nil, err
	}

	newSubscr := &domain.Subscription{
		ID:          uuid.New(),
		ServiceName: sm.ServiceName,
		Price:       sm.Price,
		UserID:      sm.UserID,
		StartDate:   stDate,
		EndDate:     &endDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return newSubscr, nil
}

func timeParse(dates string) (time.Time, error) {
	if dates == "" {
		return time.Time{}, nil
	}

	startDate, err := time.Parse("01-2006", dates)
	if err != nil {
		return time.Time{}, err
	}
	return startDate, err
}
