package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
	"rekber/constant"
	"rekber/helper"
	"rekber/model"
	"rekber/repository"
	"time"
)

type PaymentService struct {
	userRepo          *repository.UserRepository
	orderRepo         *repository.OrderRepository
	paymentRepository *repository.PaymentRepository
}

func NewPaymentService(userRepo *repository.UserRepository, orderRepo *repository.OrderRepository, paymentRepo *repository.PaymentRepository) *PaymentService {
	return &PaymentService{userRepo, orderRepo, paymentRepo}
}

func (p *PaymentService) NewTransaction(ctx context.Context, orderID, userID string) (model.Payment, error) {
	//Check User Is Exists
	checkUser, err := p.userRepo.Find(ctx, "user_id", userID)
	if err != nil {
		return model.Payment{}, err
	}

	if checkUser.ID.IsZero() {
		return model.Payment{}, errors.New("user not found")
	}

	//Check Have Payment Or Not In Current Order
	paymentCheck, _ := p.paymentRepository.Find(ctx, bson.M{"order_id": orderID, "status": constant.PENDING_STATUS})
	if !paymentCheck.ID.IsZero() {
		return model.Payment{}, errors.New(fmt.Sprintf("you have payment pending in : %s", paymentCheck.PaymentID))
	}

	//Check Order
	order, err := p.orderRepo.Find(ctx, bson.M{"order_id": orderID})
	if err != nil {
		return model.Payment{}, err
	}

	if order.ID.IsZero() {
		return model.Payment{}, errors.New("order not found")
	}

	//Check Buyer Or Seller
	if userID != order.BuyerID && userID != order.SellerID {
		return model.Payment{}, errors.New("access denied")
	}

	//If Seller Throw Error
	if userID == order.SellerID {
		return model.Payment{}, errors.New("only buyer")
	}

	//Check If Order Canceled
	if order.Status == constant.CANCELED_STATUS {
		return model.Payment{}, errors.New("order has cancel")
	}

	//Insert To Database

	//Generate Link
	link, orderLinkID, err := helper.GenerateLinkPayment(checkUser, int(order.Amount))
	if err != nil {
		return model.Payment{}, err
	}

	payment := model.Payment{
		ID:        primitive.NewObjectID(),
		PaymentID: uuid.NewString(),
		OrderID:   order.OrderID,
		Status:    constant.PENDING_STATUS,
		Link:      link,
		LinkID:    orderLinkID,
		CreatedAt: time.Now().Unix(),
	}

	newPayment, err := p.paymentRepository.Insert(ctx, payment)
	if err != nil {
		return model.Payment{}, err
	}

	return newPayment, nil
}

func (p *PaymentService) ProcessPayment(ctx context.Context, input model.PaymentNotification) error {

	payment, err := p.paymentRepository.Find(ctx, bson.M{"link_id": input.OrderID})
	if err != nil {
		return err
	}

	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		payment.Status = constant.PAID_STATUS
	} else if input.TransactionStatus == "settlement" {
		payment.Status = constant.PAID_STATUS
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		payment.Status = constant.CANCELED_STATUS
	}

	err = p.paymentRepository.Update(ctx, payment.PaymentID, bson.M{"status": payment.Status})
	if err != nil {
		return err
	}

	if payment.Status == constant.PAID_STATUS {
		err := p.orderRepo.Update(ctx, bson.M{"order_id": payment.OrderID}, bson.M{"status": constant.PAID_STATUS})
		if err != nil {
			return err
		}
	}

	return nil
}
