package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"lang/internal/repository"
	"lang/pkg/models"
	"time"
)

const (
	tokenTTl   = 12 * time.Hour
	signingKey = "nv4q[hq3083)(_#mkamcsn4"
)

type AuthService struct {
	repo repository.Authorization
}

func (s *AuthService) UserDataValidation(user models.User) error {
	if validateUsername(user.Username) != nil {
		return errors.New("username is invalid")
	}

	if validateEmail(user.Email) != nil {
		return errors.New("email is invalid")
	}

	//if validatePassword(user.Password) != nil {
	//	return errors.New(fmt.Sprintf("password is invalid %s", user.Password))
	//}

	if !validateGender(user.Gender) {
		return errors.New("gender is invalid")
	}

	//if validateProfilePhoto() != nil {
	//return errors.New("photo is invalid")
	//}

	return nil
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user models.User) (int, error) {
	user.Password = generationPasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func generationPasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash)
}

type TokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func (s *AuthService) GenerationToken(username, password string) (string, error) {
	//get user from db
	user, err := s.repo.GetUser(username, generationPasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTl).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return 0, errors.New("token claims aren't of type *TokenClaims")
	}

	return claims.UserId, nil
}
