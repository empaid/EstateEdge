package main

import (
	"context"
	"log"

	"github.com/empaid/estateedge/services/auth/internal/repository"
	"github.com/empaid/estateedge/services/common/genproto/auth"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	store *repository.Storage
	auth.UnimplementedAuthServiceServer
}

func (a authService) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {

	user, err := a.store.UserStore.GetUser(ctx, req.Username)
	if err != nil {
		log.Print("Error when logging in user", err)
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password))
	if err != nil {
		log.Print("Password incorrect")
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"id":       user.ID,
	})

	authToken, err := token.SignedString([]byte("temp_secret_key_store_in_env"))

	if err != nil {
		log.Print("Error while Login: ", err)
		return nil, err
	}

	res := &auth.LoginResponse{
		Success:   true,
		AuthToken: authToken,
	}
	return res, nil
}

func (a authService) Validate(ctx context.Context, req *auth.ValidateRequest) (*auth.ValidateResponse, error) {

	res := &auth.ValidateResponse{
		Success: true,
		UserId:  req.AuthToken,
	}
	return res, nil
}

func (a authService) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		log.Print("Error while hashing passowrd", err)
	}

	user := repository.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
	}
	if err := a.store.UserStore.RegisterUser(ctx, &user); err != nil {
		log.Fatal("Error while creating the  user")
		return nil, err
	}

	return &auth.RegisterResponse{
		Success:   true,
		AuthToken: "valid_user_created",
	}, nil

}
