package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"rekber/helper"
	"rekber/model"
	"rekber/service"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService}
}

func (p *PaymentHandler) NewPayment(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, helper.ResponseAPI(false, http.StatusUnauthorized, "unauthorized", gin.H{}))
		return
	}

	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, helper.ResponseAPI(false, http.StatusBadRequest, "id param can't be empty", gin.H{}))
		return
	}

	//Service
	payment, err := p.paymentService.NewTransaction(c, orderID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.ResponseAPI(false, http.StatusInternalServerError, err.Error(), gin.H{}))
		return
	}

	c.JSON(http.StatusOK, helper.ResponseAPI(true, http.StatusOK, "success create payment", payment))
	return

}

func (p *PaymentHandler) NotificationPayment(c *gin.Context) {
	var input model.PaymentNotification

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, helper.ResponseAPI(false, http.StatusUnprocessableEntity, err.Error(), gin.H{}))
		return
	}

	err = p.paymentService.ProcessPayment(c, input)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, helper.ResponseAPI(false, http.StatusInternalServerError, err.Error(), gin.H{}))
		return
	}

	c.JSON(http.StatusOK, input)
}
