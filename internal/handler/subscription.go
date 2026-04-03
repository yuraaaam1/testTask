package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuraaaam1/testTask/internal/model"
	"github.com/yuraaaam1/testTask/internal/service"
)

type SubscriptionHandler struct {
	service *service.SubscriptionService
}

func NewSubscriptionHandler(service *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

// @Summary Создать подписку
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body model.CreateUpdateSubscriptionInput true "Данные подписки"
// @Success 201 {object} model.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func (h *SubscriptionHandler) CreateHandler(c *gin.Context) {
	var input model.CreateUpdateSubscriptionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sub, err := h.service.Create(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sub)
}

// @Summary Получить подписку по ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} model.Subscription
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByIDHandler(c *gin.Context) {
	id := c.Param("id")

	sub, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if sub == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	c.JSON(http.StatusOK, sub)

}

// @Summary Получить список подписок
// @Tags subscriptions
// @Produce json
// @Success 200 {array} model.Subscription
// @Failure 500 {object} map[string]string
// @Router /subscriptions [get]
func (h *SubscriptionHandler) ListHandler(c *gin.Context) {
	subs, err := h.service.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subs)
}

// @Summary Обновить подписку
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Param input body model.CreateUpdateSubscriptionInput true "Данные подписки"
// @Success 200 {object} model.Subscription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateHandler(c *gin.Context) {
	id := c.Param("id")

	var input model.CreateUpdateSubscriptionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.service.Update(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if sub == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	c.JSON(http.StatusOK, sub)
}

// @Summary Удалить подписку
// @Tags subscriptions
// @Param id path string true "ID подписки"
// @Success 204
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteHandler(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(id); err != nil {
		if err.Error() == "subscription not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Подсчёт суммарной стоимости подписок
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "ID пользователя"
// @Param service_name query string false "Название сервиса"
// @Param date_from query string true "Начало периода (MM-YYYY)"
// @Param date_to query string true "Конец периода (MM-YYYY)"
// @Success 200 {object} model.TotalCostResult
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/cost [get]
func (h *SubscriptionHandler) CalculateTotalCostHandler(c *gin.Context) {
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	if dateFrom == "" || dateTo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date_from and date_to required"})
		return
	}

	result, err := h.service.CalculateTotalCost(userID, serviceName, dateFrom, dateTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
