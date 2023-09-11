package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
	"rekber/dto"
	"rekber/helper"
	"rekber/model"
	"rekber/repository"
	"time"
)

type AuthService struct {
	repo     *repository.UserRepository
	authRepo *repository.AuthRepository
}

func NewAuthService(repo *repository.UserRepository, authRepo *repository.AuthRepository) *AuthService {
	return &AuthService{repo, authRepo}
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

func (u *AuthService) RegisterLoginV2(ctx context.Context, email string) (model.Auth, error) {
	if email == "" {
		return model.Auth{}, errors.New("email can't be empty")
	}

	//Check Email Exists?
	user, err := u.repo.Find(ctx, "email", email)
	if user.ID.IsZero() {
		return model.Auth{}, errors.New("user not found")
	}

	if err != nil {
		return model.Auth{}, err
	}

	loginV2 := model.Auth{
		ID:         primitive.NewObjectID(),
		AuthID:     uuid.NewString(),
		UserID:     user.UserID,
		CodeLink:   fmt.Sprintf("type%dv2%d", time.Now().Unix(), time.Now().Second()),
		IsLoggedin: false,
		CreatedAt:  time.Now().Unix(),
		ExpiredAt:  time.Now().Add(1 * time.Minute).Unix(),
	}

	//Add To Service
	newLoginV2, err := u.authRepo.Create(ctx, loginV2)
	if err != nil {
		return model.Auth{}, err
	}

	return newLoginV2, nil
}

func (u *AuthService) LoginV2(ctx context.Context, codeLink string) (model.User, error) {
	if codeLink == "" {
		return model.User{}, errors.New("codelink can't be empty")
	}

	//Check Codelink is Exists?
	loginV2, err := u.authRepo.Find(ctx, bson.M{"code_link": codeLink})
	if loginV2.ID.IsZero() {
		return model.User{}, errors.New("codelink not found")
	}
	if err != nil {
		return model.User{}, err
	}

	//Check IsLoggedin
	if loginV2.IsLoggedin {
		return model.User{}, errors.New("link has login")
	}

	//Check Expired
	timeNow := time.Now().Unix()
	if loginV2.ExpiredAt <= timeNow {
		return model.User{}, errors.New("expired link")
	}

	//check user
	user, err := u.repo.Find(ctx, "user_id", loginV2.UserID)
	if user.ID.IsZero() {
		return model.User{}, errors.New("user not found")
	}

	if err != nil {
		return model.User{}, err
	}

	//Update Is Loggedin
	if err := u.authRepo.Update(ctx, bson.M{"code_link": codeLink}, bson.M{"is_loggedin": true}); err != nil {
		return model.User{}, err
	}

	return user, nil
}
