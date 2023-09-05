package service

import (
	"context"
	"errors"
	"rekber/dto"
	"rekber/helper"
	"rekber/model"
	"rekber/repository"
)

type AuthService struct {
	repo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{repo}
}

func (u *AuthService) Login(ctx context.Context, req dto.LoginReq) (model.User, error) {
	if req.Username == "" {
		return model.User{}, errors.New("username is empty")
	}

	user, err := u.repo.Find(ctx, "username", req.Username)
	if err != nil {
		return model.User{}, errors.New("user not found")
	}

	if user.Username == "" {
		return model.User{}, errors.New("username or password is wrong")
	}

	// compare password
	ok := helper.CheckPasswordHash(req.Password, user.Password)
	if !ok {
		return model.User{}, errors.New("username or password is wrong")
	}

	return user, nil
}
