package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	OrderID   string             `json:"order_id" bson:"order_id"`
	SellerID  string             `json:"seller_id" bson:"seller_id"`
	BuyerID   string             `json:"buyer_id" bson:"buyer_id"`
	Type      string             `json:"type" bson:"type"`
	Amount    int64              `json:"amount" bson:"amount" `
	Status    string             `json:"status" bson:"status"` //Waiting, Pending, Success, Canceled
	CreatedAt int64              `json:"created_at" bson:"created_at"`
	UpdatedAt int64              `json:"updated_at" bson:"updated_at"`

	//=======> Omitempty
	Fee    int   `json:"fee" bson:"fee"`
	Seller *User `json:"seller,omitempty" bson:"seller,omitempty"`
	Buyer  *User `json:"buyer,omitempty" bson:"buyer,omitempty"`
}
