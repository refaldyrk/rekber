package service

import (
	"context"
	"errors"
	"fmt"
	"rekber/constant"
	"rekber/dto"
	"rekber/helper"
	"rekber/model"
	"rekber/repository"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

type OrderService struct {
	repo     *repository.OrderRepository
	userRepo *repository.UserRepository
}

func NewOrderService(repo *repository.OrderRepository, userRepo *repository.UserRepository) *OrderService {
	return &OrderService{repo: repo, userRepo: userRepo}
}

func (o *OrderService) Insert(ctx context.Context, req dto.NewOrderReq, userID string) (model.Order, error) {
	if req.BuyerIdentity == "" || req.SellerIdentity == "" || req.Type == "" || req.Amount < constant.MINIMUM_PRICE {
		return model.Order{}, errors.New("invalid request")
	}

	buyer, _ := o.userRepo.FindByUsernameOrEmail(ctx, req.BuyerIdentity)
	if buyer.ID.IsZero() {
		return model.Order{}, errors.New("buyer id not found")
	}

	seller, _ := o.userRepo.FindByUsernameOrEmail(ctx, req.SellerIdentity)
	if seller.ID.IsZero() {
		return model.Order{}, errors.New("seller id not found")
	}

	if userID != seller.UserID && userID != buyer.UserID {
		return model.Order{}, errors.New("invalid request")
	}

	//Calculate Fee
	fee := helper.CalculateFee(int(req.Amount))

	order := model.Order{
		ID:        primitive.NewObjectID(),
		OrderID:   uuid.NewString(),
		SellerID:  seller.UserID,
		BuyerID:   buyer.UserID,
		Type:      req.Type,
		Amount:    req.Amount + int64(fee),
		Fee:       fee,
		Status:    constant.PENDING_STATUS,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	newOrder, err := o.repo.Insert(ctx, order)
	if err != nil {
		return model.Order{}, err
	}

	return newOrder, nil
}

func (o *OrderService) FindAllOrderByRole(ctx context.Context, userID string, role string) ([]model.Order, error) {
	if role == "" {
		return []model.Order{}, errors.New("role isn't be empty")
	}

	if role == constant.BUYER {
		orders, err := o.repo.FindAll(ctx, bson.M{"buyer_id": userID})
		if err != nil {
			return []model.Order{}, err
		}

		return orders, nil
	} else if role == constant.SELLER {
		orders, err := o.repo.FindAll(ctx, bson.M{"seller_id": userID})
		if err != nil {
			return []model.Order{}, err
		}

		return orders, nil
	}

	return []model.Order{}, errors.New("undefined role, only buyer and seller")
}

func (o *OrderService) FindByOrderID(ctx context.Context, orderID, userID string) (model.Order, error) {
	if userID == "" || orderID == "" {
		return model.Order{}, errors.New("invalid request")
	}

	order, err := o.repo.Find(ctx, bson.M{"order_id": orderID})
	if order.ID.IsZero() {
		return model.Order{}, errors.New("not found")
	}

	if err != nil {
		return model.Order{}, err
	}

	if order.BuyerID != userID && order.SellerID != userID {
		return model.Order{}, errors.New("access denied")
	}

	orderDataDetail, err := o.repo.GetOrderByOrderID(ctx, orderID)
	if err != nil {
		return model.Order{}, err
	}

	return orderDataDetail, nil
}

func (o *OrderService) FindAllOrderByStatus(ctx context.Context, status, userID string) ([]model.Order, error) {
	if userID == "" || status == "" {
		return []model.Order{}, errors.New("invalid request")
	}

	//Match Status
	switch status {
	case constant.CANCELED_STATUS:
		status = constant.CANCELED_STATUS
		break
	case constant.PENDING_STATUS:
		status = constant.PENDING_STATUS
		break
	case constant.SUCCESS_STATUS:
		status = constant.SUCCESS_STATUS
		break
	default:
		return []model.Order{}, errors.New("uknown status")
	}

	//Check User Exists
	user, err := o.userRepo.Find(ctx, "user_id", userID)
	if user.ID.IsZero() {
		return []model.Order{}, errors.New("user not found")
	}

	if err != nil {
		return []model.Order{}, err
	}

	filter := bson.M{
		"$or": []bson.M{
			{"buyer_id": user.UserID},
			{"seller_id": user.UserID},
		},
		"status": status,
	}

	allOrder, err := o.repo.FindAll(ctx, filter)
	if err != nil {
		return []model.Order{}, err
	}

	return allOrder, nil
}

func (o *OrderService) SetStatusByOrderID(ctx context.Context, orderID, userID, status string) (bool, error) {
	if userID == "" || orderID == "" {
		return false, errors.New("invalid request")
	}

	//Check Status?
	switch status {
	case constant.CANCELED_STATUS:
		status = constant.CANCELED_STATUS
		break
	case constant.SUCCESS_STATUS:
		status = constant.SUCCESS_STATUS
		break
	default:
		return false, errors.New("uknown status")
	}

	order, err := o.repo.Find(ctx, bson.M{"order_id": orderID})
	if order.ID.IsZero() {
		return false, errors.New("not found")
	}

	if err != nil {
		return false, err
	}

	//Is Already Or Not?
	if status == order.Status {
		return false, errors.New(fmt.Sprintf("already status: %s", order.Status))
	}

	//Check User, Seller Or Buyer?
	if order.BuyerID != userID && order.SellerID != userID {
		return false, errors.New("access denied")
	}

	//Update Order Data
	isSuccesUpdate, err := o.repo.SetStatusOrderByOrderID(ctx, orderID, status)
	if err != nil {
		return isSuccesUpdate, err
	}

	return isSuccesUpdate, nil
}
