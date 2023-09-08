package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Payment struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	PaymentID string             `json:"payment_id" bson:"payment_id"`
	OrderID   string             `json:"order_id" bson:"order_id"`
	Status    string             `json:"status" bson:"status"`
	Link      string             `json:"link" bson:"link"`
	LinkID    string             `json:"link_id" bson:"link_id"`
	CreatedAt int64              `json:"created_at" bson:"created_at"`

	//Relation
	Order *Order `json:"order,omitempty" bson:"order,omitempty"`
}

type PaymentNotification struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}
