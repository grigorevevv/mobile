package api_error

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	ErrNotFound      = errors.New("Подписка не найдена")
	ErrAlreadyExists = errors.New("Подписка уже существует")
	ErrEntity        = errors.New("Подписки не существует")
	Erruuid          = errors.New("Неверный формат UUID")
	//ErrDatabaseIsEmpty = errors.New("database is empty")
)

type ResponseError struct {
	Status      int       `json:"status"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
	Err         string    `json:"error"`
}

func NewResponseError(status int, desc string, err error) *ResponseError {
	newresponseerr := &ResponseError{
		status,
		"error: " + desc,
		time.Now(),
		err.Error(),
	}
	return newresponseerr
}

func (r *ResponseError) Error() string {
	return fmt.Sprintf("%d: %v", r.Status, r.Error)
}

func ErrorHandler(c *gin.Context, err error) {
	// 	Ошибка преобразования из текста в число
	var numError *strconv.NumError
	if errors.As(err, &numError) {
		err := NewResponseError(http.StatusBadRequest, "handler", fmt.Errorf("error: convert params is invalid"))
		c.JSON(err.Status, err)
		return
	}

	// 	Ошибка преобразования время в текст
	var timeError *time.ParseError
	if errors.As(err, &timeError) {
		err := NewResponseError(http.StatusBadRequest, "handler", fmt.Errorf("error: convert params times is invalid"))
		c.JSON(err.Status, err)
		return
	}

	// 	Ошибка преобразования JSON
	var syntaxError *json.SyntaxError
	if errors.As(err, &syntaxError) {
		err := NewResponseError(http.StatusBadRequest, "handler", fmt.Errorf("error: json is invalid"))
		c.JSON(err.Status, err)
		return
	}

	//	Прочие ошибки
	if errors.Is(err, ErrNotFound) {
		err := NewResponseError(http.StatusNotFound, "service", err)
		c.JSON(err.Status, err)
	} else if errors.Is(err, ErrAlreadyExists) {
		err := NewResponseError(http.StatusConflict, "service", err)
		c.JSON(err.Status, err)
	} else if errors.Is(err, Erruuid) {
		err := NewResponseError(http.StatusBadRequest, "handler", err)
		c.JSON(err.Status, err)
	} else {
		err := NewResponseError(http.StatusInternalServerError, "other", err)
		c.JSON(err.Status, err)
	}
}
