package repository

import (
	"context"
	"errors"
	"rekber/helper"
	"rekber/model"

	"github.com/qiniu/qmgo"
	"gopkg.in/mgo.v2/bson"
)

type OrderRepository struct {
	db *qmgo.Database
}

func NewOrderRepository(db *qmgo.Database) *OrderRepository {
	return &OrderRepository{db}
}

func (o *OrderRepository) Insert(ctx context.Context, order model.Order) (model.Order, error) {
	_, err := o.db.Collection("Order").InsertOne(ctx, order)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (o *OrderRepository) Find(ctx context.Context, filter bson.M) (model.Order, error) {
	var order model.Order

	err := o.db.Collection("Order").Find(ctx, filter).One(&order)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (o *OrderRepository) FindAll(ctx context.Context, filter bson.M) ([]model.Order, error) {
	var orders []model.Order

	err := o.db.Collection("Order").Find(ctx, filter).All(&orders)
	if err != nil {
		return []model.Order{}, err
	}

	return orders, nil
}

func (o *OrderRepository) GetOrderByOrderID(ctx context.Context, orderID string) (model.Order, error) {
	var orders model.Order
	collection := o.db.Collection("Order")

	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "User",
				"localField":   "seller_id",
				"foreignField": "user_id",
				"as":           "seller",
			},
		},
		{
			"$match": bson.M{
				"$and": []bson.M{
					{
						"order_id": orderID,
					},
				},
			},
		},
		{
			"$unwind": "$seller",
		},
		{
			"$lookup": bson.M{
				"from":         "User",
				"localField":   "buyer_id",
				"foreignField": "user_id",
				"as":           "buyer",
			},
		},
		{
			"$unwind": "$buyer",
		},
		{
			"$project": helper.GetBSONTagMap(&model.Order{}, bson.M{
				"seller": "$seller",
				"buyer":  "$buyer",
			}),
		},
	}

	err := collection.Aggregate(ctx, pipeline).One(&orders)
	if err != nil {
		return model.Order{}, err
	}

	return orders, nil
}

func (o *OrderRepository) SetStatusOrderByOrderID(ctx context.Context, orderID, status string) (bool, error) {
	err := o.db.Collection("Order").UpdateOne(ctx, bson.M{"order_id": orderID}, bson.M{
		"$set": bson.M{
			"status": status,
		},
	})

	if err == qmgo.ErrNoSuchDocuments {
		return false, errors.New("not found")
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
