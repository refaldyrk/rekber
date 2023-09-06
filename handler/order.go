package handler

import (
	"fmt"
	"net/http"
	"rekber/dto"
	"rekber/helper"
	"rekber/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(s *service.OrderService) *OrderHandler {
	return &OrderHandler{service: s}
}

func (o *OrderHandler) NewOrder(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, helper.ResponseAPI(false, http.StatusUnauthorized, "unauthorized", gin.H{}))
		return
	}

	var req dto.NewOrderReq
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.ResponseAPI(false, http.StatusBadRequest, err.Error(), gin.H{}))
		return
	}

	newOrder, err := o.service.Insert(c, req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.ResponseAPI(false, http.StatusInternalServerError, err.Error(), gin.H{}))
		return
	}

	c.JSON(http.StatusOK, helper.ResponseAPI(true, http.StatusOK, "success create new order", newOrder))
	return
}

func (o *OrderHandler) FindAllOrderByRole(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, helper.ResponseAPI(false, http.StatusUnauthorized, "invalid user id", gin.H{}))
		return
	}

	queryRole := c.Query("role")
	if queryRole == "" {
		c.JSON(http.StatusBadRequest, helper.ResponseAPI(false, http.StatusBadRequest, "role can't be empty", gin.H{}))
		return
	}

	orders, err := o.service.FindAllOrderByRole(c, userID, queryRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.ResponseAPI(false, http.StatusInternalServerError, err.Error(), gin.H{}))
		return
	}

	c.JSON(http.StatusOK, helper.ResponseAPI(true, http.StatusOK, fmt.Sprintf("success get all orders: %d", len(orders)), orders))
	return
}

func (o *OrderHandler) GetOrderDetailByOrderID(c *gin.Context) {
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

	detailOrder, err := o.service.FindByOrderID(c, orderID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.ResponseAPI(false, http.StatusInternalServerError, err.Error(), gin.H{}))
		return
	}

	c.JSON(http.StatusOK, helper.ResponseAPI(true, http.StatusOK, "success get detail order", detailOrder))
	return
}