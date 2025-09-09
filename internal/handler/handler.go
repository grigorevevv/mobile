package handler

import (
	"mobile/internal/api_error"
	"mobile/internal/domain"
	"mobile/internal/mapper"
	"mobile/internal/repository"
	"mobile/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type HandlerSrv struct {
	service *service.SubscriptionService
	logger  *logrus.Logger
}

func NewMsHandler(service *service.SubscriptionService, log *logrus.Logger) *HandlerSrv {
	return &HandlerSrv{service: service, logger: log}
}

// Ответ одного типа
type ginResponse[T any] struct {
	Status    int       `json:"status"`
	Data      T         `json:"data"`
	TimeStamp time.Time `json:"timeStamp"`
}

// Создание подписки
func (h *HandlerSrv) CreateSubscription(c *gin.Context) {
	var newSubscriptionMap mapper.SubscriptionMap

	if err := c.ShouldBindJSON(&newSubscriptionMap); err != nil {
		h.logger.Errorf("Ошибка привязки JSON-данных: %v", err)
		api_error.ErrorHandler(c, err)
		return
	}

	newSubscription, err := mapper.SubscriptionMapper(newSubscriptionMap)
	if err != nil {
		h.logger.Errorf("Ошибка маппера: %v", err)
		api_error.ErrorHandler(c, err)
		return
	}

	if err := h.service.SubscriptionCreate(c.Request.Context(), newSubscription); err != nil {
		h.logger.Errorf("Ошибка создания подписки: %v", err)
		api_error.ErrorHandler(c, err)
		return
	}

	pesp := ginResponse[string]{
		Status:    http.StatusOK,
		Data:      "Подписка добавлена",
		TimeStamp: time.Now(),
	}
	c.JSON(http.StatusOK, pesp)
}

// возвращает информацию о подписке по её ID
func (h *HandlerSrv) GetSubscription(c *gin.Context) {
	idStr := c.Query("id")

	// Преобразовываем строку в UUID
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Errorf("Неверный формат UUID: %v", err)
		api_error.ErrorHandler(c, err)
		return
	}

	sub, err := h.service.GetById(c.Request.Context(), id)
	if err != nil {
		h.logger.Errorf("Ошибка получения подписки: %v", err)
		api_error.ErrorHandler(c, err)
		return
	}

	pesp := ginResponse[domain.Subscription]{
		Status:    http.StatusOK,
		Data:      *sub,
		TimeStamp: time.Now(),
	}
	c.JSON(http.StatusOK, pesp)
}

// Обновляет данные у подписки по id
func (h *HandlerSrv) UpdateSubscription(c *gin.Context) {
	idStr := c.Query("id")

	// Преобразовываем строку в UUID
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Errorf("Неверный формат UUID: %v", err)
		api_error.ErrorHandler(c, err)
		return
	}

	var subscriptionUpdate mapper.SubscriptionMap

	if err := c.ShouldBindJSON(&subscriptionUpdate); err != nil {
		h.logger.Errorf("Ошибка привязки JSON-данных: %v", err)
		api_error.ErrorHandler(c, err)
		return
	}

	newSubscription, err := mapper.SubscriptionMapper(subscriptionUpdate)
	if err != nil {
		h.logger.Errorf("Ошибка маппера: %v", err)
		api_error.ErrorHandler(c, err)
		return
	}

	newSubscription.ID = id

	err = h.service.UpdateSubscription(c.Request.Context(), newSubscription)
	if err != nil {
		h.logger.Errorf("Ошибка обновления подписки: %v", err)
		api_error.ErrorHandler(c, err)
		return
	}

	pesp := ginResponse[string]{
		Status:    http.StatusOK,
		Data:      "Подписка обновлена",
		TimeStamp: time.Now(),
	}
	c.JSON(http.StatusOK, pesp)

}

// Удаляем подписку по id
func (h *HandlerSrv) DeleteSubscription(c *gin.Context) {
	idStr := c.Query("id")

	// Преобразовываем строку в UUID
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Errorf("Неверный формат UUID: %v", err)
		api_error.ErrorHandler(c, err)
		return
	}

	err = h.service.DeleteSubscription(c.Request.Context(), id)
	if err != nil {
		h.logger.Errorf("Ошибка удаления подписки: %v", err)
		api_error.ErrorHandler(c, err)
		return
	}

	pesp := ginResponse[string]{
		Status:    http.StatusOK,
		Data:      "Подписка удалена",
		TimeStamp: time.Now(),
	}
	c.JSON(http.StatusOK, pesp)
}

func (h *HandlerSrv) Report(c *gin.Context) {
	params, err := checkingListParams(c)
	if err != nil {
		h.logger.Errorf("Ошибка параметров: %v", err)
		api_error.ErrorHandler(c, err)
		return
	}

	data, err := h.service.ReportSubscription(c, params)
	if err != nil {
		h.logger.Errorf("Ошибка отчета: %v", err)
		api_error.ErrorHandler(c, err)
		return
	}

	pesp := ginResponse[[]repository.ListReport]{
		Status:    http.StatusOK,
		Data:      data,
		TimeStamp: time.Now(),
	}
	c.JSON(http.StatusOK, pesp)

}

func checkingListParams(c *gin.Context) (*repository.ListParams, error) {
	startDateParam := c.Query("start") // например: 2025-07-27
	endDateParam := c.Query("end")
	idStr := c.Query("id")
	serviceName := c.Query("service_name")

	startDate, err := time.Parse("01-2006", startDateParam)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse("01-2006", endDateParam)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	param := &repository.ListParams{
		StartDate:   startDate,
		EndDate:     endDate,
		IdUser:      id,
		ServiceName: serviceName,
	}

	return param, nil
}
