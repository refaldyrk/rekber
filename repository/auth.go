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

func (a *AuthRepository) InsertLogout(ctx context.Context, data model.Logout) (model.Logout, error) {
	_, err := a.db.Collection("Logout").InsertOne(ctx, data)
	if err != nil {
		return model.Logout{}, err
	}

	return data, nil
}

func (a *AuthRepository) FindLogoutByToken(ctx context.Context, token string) (model.Logout, error) {
	var result model.Logout
	if err := a.db.Collection("Logout").Find(ctx, bson.M{"token": token}).One(&result); err != nil {
		return model.Logout{}, err
	}

	return result, nil
}

func (a *AuthRepository) InsertLogin(ctx context.Context, data model.Login) (model.Login, error) {
	_, err := a.db.Collection("Login").InsertOne(ctx, data)
	if err != nil {
		return model.Login{}, err
	}

	return data, nil
}

func (a *AuthRepository) FindLogin(ctx context.Context, filter bson.M) (model.Login, error) {
	var result model.Login
	if err := a.db.Collection("Login").Find(ctx, filter).One(&result); err != nil {
		return model.Login{}, err
	}

	return result, nil
}

func (a *AuthRepository) FindAllLogin(ctx context.Context, filter bson.M) ([]model.Login, error) {
	var result []model.Login
	if err := a.db.Collection("Login").Find(ctx, filter).All(&result); err != nil {
		return []model.Login{}, err
	}

	return result, nil
}

func (a *AuthRepository) CountLoginData(ctx context.Context, filter bson.M) (int64, error) {
	result, err := a.db.Collection("Login").Find(ctx, filter).Count()
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (a *AuthRepository) DeleteLogin(ctx context.Context, filter bson.M) error {
	if err := a.db.Collection("Login").Remove(ctx, filter); err != nil {
		return err
	}

	return nil
}
