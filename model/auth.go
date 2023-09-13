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

type Logout struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	LogoutID string             `json:"logout_id" bson:"logout_id"`
	Token    string             `json:"token" bson:"token"`
	LogoutAt int64              `json:"logout_at" bson:"logout_at"`
}

type Login struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	LoginID    string             `json:"login_id" bson:"login_id"`
	UserID     string             `json:"user_id" bson:"user_id"`
	Token      string             `json:"token" bson:"token"`
	LoggedinAt int64              `json:"loggedin_at" bson:"loggedin_at"`
}
