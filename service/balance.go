package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (b *BalanceService) InsertNewBalance(ctx context.Context, orderID string, amount int64) error {
	if orderID == "" || amount == 0 {
		return errors.New("invalid request")
	}

	//Insert New Data
	balance := model.Balance{
		ID:           primitive.NewObjectID(),
		BalanceID:    fmt.Sprintf("BAL%sANCE", uuid.NewString()),
		OrderID:      orderID,
		Amount:       amount,
		IsWithdrawal: false,
		PaidedAt:     time.Now().Unix(),
	}

	err := b.repo.Insert(ctx, balance)
	if err != nil {
		return err
	}

	return nil
}
