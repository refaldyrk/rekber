package repository

import (
	"context"
	"github.com/qiniu/qmgo"
	"gopkg.in/mgo.v2/bson"
	"rekber/helper"
	"rekber/model"
)

type BalanceRepository struct {
	db *qmgo.Database
}

func NewBalanceRepository(db *qmgo.Database) *BalanceRepository {
	return &BalanceRepository{
		db: db,
	}
}

func (b *BalanceRepository) Insert(ctx context.Context, data model.Balance) error {
	_, err := b.db.Collection("Balance").InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (b *BalanceRepository) Find(ctx context.Context, filter bson.M) (model.Balance, error) {
	var balance model.Balance
	if err := b.db.Collection("Balance").Find(ctx, filter).One(&balance); err != nil {
		return model.Balance{}, err
	}

	return balance, nil
}

func (b *BalanceRepository) FindAll(ctx context.Context, filter bson.M) ([]model.Balance, error) {
	var balances []model.Balance
	if err := b.db.Collection("Balance").Find(ctx, filter).All(&balances); err != nil {
		return []model.Balance{}, err
	}

	return balances, nil
}

func (b *BalanceRepository) Update(ctx context.Context, filter bson.M, update bson.M) error {
	err := b.db.Collection("Balance").UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return err
	}

	return nil
}

func (b *BalanceRepository) GetBalanceByBalanceID(ctx context.Context, balanceID string) (model.Balance, error) {
	var balance model.Balance

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"$and": []bson.M{
					{
						"balance_id": balanceID,
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "Order",
				"localField":   "order_id",
				"foreignField": "order_id",
				"as":           "order",
			},
		},
		{
			"$unwind": "$order",
		},
		{
			"$project": helper.GetBSONTagMap(&model.Balance{}, bson.M{
				"order": "$order",
			}),
		},
	}

	err := b.db.Collection("Balance").Aggregate(ctx, pipeline).One(&balance)
	if err != nil {
		return model.Balance{}, err
	}

	return balance, nil
}
