package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Auth struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	AuthID     string             `json:"auth_id" bson:"auth_id"`
	UserID     string             `json:"user_id" bson:"user_id"`
	CodeLink   string             `json:"code_link" bson:"code_link"`
	IsLoggedin bool               `json:"is_loggedin" bson:"is_loggedin"`
	CreatedAt  int64              `json:"created_at" bson:"created_at"`
	ExpiredAt  int64              `json:"expired_at" bson:"expired_at"`
}
