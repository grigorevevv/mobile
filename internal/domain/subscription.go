package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subscription struct {
	ID          uuid.UUID  `gorm:"primaryKey;type:uuid" json:"id"`
	ServiceName string     `gorm:"not null" json:"service_name"` // название сервиса
	Price       int        `gorm:"not null" json:"price"`        // стоимость в рублях
	UserID      uuid.UUID  `gorm:"not null" json:"user_id"`      // уникальный ID пользователя
	StartDate   time.Time  `gorm:"not null" json:"start_date"`   // дата начала подписки
	EndDate     *time.Time `gorm:"" json:"end_date,omitempty"`   // опциональная дата окончания подписки
	CreatedAt   time.Time  `gorm:"not null" json:"created_at"`   // дата создания записи
	UpdatedAt   time.Time  `gorm:"not null" json:"updated_at"`   // дата последнего обновления
}

func (s *Subscription) BeforeCreate(db *gorm.DB) (err error) {
	if s == nil {
		return fmt.Errorf("entity does not exist")
	}

	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
