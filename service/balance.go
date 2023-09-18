package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
	"rekber/model"
	"rekber/repository"
	"time"
)

type BalanceService struct {
	repo *repository.BalanceRepository
}

func NewBalanceService(repo *repository.BalanceRepository) *BalanceService {
	return &BalanceService{
		repo: repo,
	}
}

func (b *BalanceService) InsertNewBalance(ctx context.Context, orderID, sellerID string, amount, fee int64) error {
	if orderID == "" || amount == 0 {
		return errors.New("invalid request")
	}

	//Insert New Data
	balance := model.Balance{
		ID:           primitive.NewObjectID(),
		BalanceID:    fmt.Sprintf("BAL%sANCE", uuid.NewString()),
		SellerID:     sellerID,
		OrderID:      orderID,
		Amount:       amount - fee,
		IsWithdrawal: false,
		PaidedAt:     time.Now().Unix(),
	}

	err := b.repo.Insert(ctx, balance)
	if err != nil {
		return err
	}

	return nil
}

func (b *BalanceService) FindAllBalanceByUserID(ctx context.Context, userID string) ([]model.Balance, error) {
	if userID == "" {
		return []model.Balance{}, errors.New("unauthorized")
	}

	//Service
	balances, err := b.repo.FindAll(ctx, bson.M{"seller_id": userID})
	if err != nil {
		return nil, err
	}

	return balances, nil
}

func (b *BalanceService) GetDetailBalanceByID(ctx context.Context, balanceID string) (model.Balance, error) {
	if balanceID == "" {
		return model.Balance{}, errors.New("order id can't be empty")
	}

	//Check Order ID
	balance, err := b.repo.Find(ctx, bson.M{"balance_id": balanceID})
	if balance.ID.IsZero() {
		if err != nil {
			return model.Balance{}, qmgo.ErrNoSuchDocuments
		}
	}

	if err != nil {
		return model.Balance{}, err
	}

	//Get Service
	byBalanceID, err := b.repo.GetBalanceByBalanceID(ctx, balanceID)
	if err != nil {
		return model.Balance{}, err
	}

	return byBalanceID, nil
}
