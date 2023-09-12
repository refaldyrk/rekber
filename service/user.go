package service

import (
	"context"
	"errors"
	"rekber/helper"
	"rekber/model"
	"rekber/repository"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo}
}

func (u *UserService) Register(ctx context.Context, req model.User) (model.User, error) {
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return model.User{}, errors.New("body invalid request")
	}

	if len(req.Password) < 8 {
		return model.User{}, errors.New("password invalid request, must be 8 character")
	}

	//Check Username
	userCheck, _ := u.repo.Find(ctx, "username", req.Username)

	if !userCheck.ID.IsZero() {
		return model.User{}, errors.New("username already exist")
	}

	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return model.User{}, err
	}

	user := model.User{
		ID:            primitive.NewObjectID(),
		UserID:        uuid.NewString(),
		Username:      req.Username,
		Email:         req.Email,
		DeviceConnect: 0,
		Password:      hashedPassword,
		CreatedAt:     time.Now().Unix(),
		UpdatedAt:     time.Now().Unix(),
	}

	newUser, err := u.repo.Create(ctx, user)
	if err != nil {
		return model.User{}, err
	}

	return newUser, nil
}

func (u *UserService) MySelf(ctx context.Context, userID string) (model.User, error) {
	if userID == "" {
		return model.User{}, errors.New("user id can't be empty")
	}

	user, err := u.repo.Find(ctx, "user_id", userID)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
