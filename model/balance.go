package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Balance struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	BalanceID    string             `json:"balance_id" bson:"balance_id"`
	OrderID      string             `json:"order_id" bson:"order_id"`
	Amount       int64              `json:"amount" bson:"amount"`
	IsWithdrawal bool               `json:"is_withdrawal" bson:"is_withdrawal"`
	PaidedAt     int64              `json:"paided_at" bson:"paided_at"`

	//Not Save In DB
	Order *Order `json:"order,omitempty" bson:"order,omitempty"`
}
