package repository

import (
	"context"
	"rekber/model"

	"github.com/qiniu/qmgo"
	"gopkg.in/mgo.v2/bson"
)

type UserRepository struct {
	db *qmgo.Database
}

func NewUserRepository(db *qmgo.Database) *UserRepository {
	return &UserRepository{db}
}

func (u *UserRepository) Create(ctx context.Context, user model.User) (model.User, error) {
	_, err := u.db.Collection("User").InsertOne(ctx, user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (u *UserRepository) Find(ctx context.Context, nameFilter, valueFilter string) (model.User, error) {
	var user model.User
	err := u.db.Collection("User").Find(ctx, bson.M{nameFilter: valueFilter}).One(&user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (u *UserRepository) FindByUsernameOrEmail(ctx context.Context, nameFilter string) (model.User, error) {
	var user model.User
	filter := bson.M{
		"$or": []bson.M{
			{"username": nameFilter},
			{"email": nameFilter},
		},
	}

	err := u.db.Collection("User").Find(ctx, filter).One(&user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (u *UserRepository) FindAll(ctx context.Context, nameFilter, valueFilter string) ([]model.User, error) {
	var user []model.User
	err := u.db.Collection("User").Find(ctx, bson.M{nameFilter: valueFilter}).All(&user)
	if err != nil {
		return []model.User{}, err
	}

	return user, nil
}

func (u *UserRepository) Update(ctx context.Context, filter bson.M, update bson.M) error {
	if err := u.db.Collection("User").UpdateOne(ctx, filter, bson.M{
		"$set": update,
	}); err != nil {
		return err
	}
	return nil
}
