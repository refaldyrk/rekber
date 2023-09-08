package repository

import (
	"context"
	"github.com/qiniu/qmgo"
	"gopkg.in/mgo.v2/bson"
	"rekber/helper"
	"rekber/model"
)

type PaymentRepository struct {
	db *qmgo.Database
}

func NewPaymentRepository(db *qmgo.Database) *PaymentRepository {
	return &PaymentRepository{
		db,
	}
}

func (p *PaymentRepository) Insert(ctx context.Context, payment model.Payment) (model.Payment, error) {
	_, err := p.db.Collection("Payment").InsertOne(ctx, payment)
	if err != nil {
		return model.Payment{}, err
	}

	return payment, nil
}

func (p *PaymentRepository) Find(ctx context.Context, filter bson.M) (model.Payment, error) {
	var payment model.Payment

	err := p.db.Collection("Payment").Find(ctx, filter).One(&payment)
	if err != nil {
		return model.Payment{}, err
	}

	return payment, nil
}

func (p *PaymentRepository) FindAll(ctx context.Context, filter bson.M) ([]model.Payment, error) {
	var payments []model.Payment

	err := p.db.Collection("Payment").Find(ctx, filter).All(&payments)
	if err != nil {
		return []model.Payment{}, err
	}

	return payments, nil
}

func (p *PaymentRepository) Update(ctx context.Context, paymentID string, update bson.M) error {
	err := p.db.Collection("Payment").UpdateOne(ctx, bson.M{"payment_id": paymentID}, bson.M{"$set": update})
	if err != nil {
		return err
	}

	return nil
}

func (p *PaymentRepository) GetPaymentByID(ctx context.Context, paymentID string) (model.Payment, error) {
	var payment model.Payment

	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "Order",
				"localField":   "order_id",
				"foreignField": "order_id",
				"as":           "order",
			},
		},
		{
			"$match": bson.M{
				"$and": []bson.M{
					{
						"payment_id": paymentID,
					},
				},
			},
		},
		{
			"$unwind": "$order",
		},
		{
			"$project": helper.GetBSONTagMap(&model.Order{}, bson.M{
				"order": "$order",
			}),
		},
	}

	err := p.db.Collection("Order").Aggregate(ctx, pipeline).One(&payment)
	if err != nil {
		return model.Payment{}, err
	}

	return payment, nil
}
