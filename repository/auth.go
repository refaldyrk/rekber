package repository

import (
	"context"
	"github.com/qiniu/qmgo"
	"gopkg.in/mgo.v2/bson"
	"rekber/model"
)

type AuthRepository struct {
	db *qmgo.Database
}

func NewAuthRepository(db *qmgo.Database) *AuthRepository {
	return &AuthRepository{
		db,
	}
}

func (a *AuthRepository) Create(ctx context.Context, data model.Auth) (model.Auth, error) {
	_, err := a.db.Collection("Auth").InsertOne(ctx, data)
	if err != nil {
		return model.Auth{}, err
	}

	return data, nil
}

func (a *AuthRepository) Find(ctx context.Context, filter bson.M) (model.Auth, error) {
	var result model.Auth
	if err := a.db.Collection("Auth").Find(ctx, filter).One(&result); err != nil {
		return model.Auth{}, err
	}

	return result, nil
}

func (a *AuthRepository) Update(ctx context.Context, filter bson.M, update bson.M) error {
	if err := a.db.Collection("Auth").UpdateOne(ctx, filter, bson.M{
		"$set": update,
	}); err != nil {
		return err
	}

	return nil
}
