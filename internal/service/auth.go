package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"lang/internal/repository"
	"lang/pkg/models"
	"regexp"
	"time"
)

const (
	tokenTTl   = 12 * time.Hour
	signingKey = "nv4q[hq3083)(_#mlnbnkamcsn4"
)

type AuthService struct {
	repo repository.Authorization
	ver  Verification
}

func NewAuthService(repo repository.Authorization, ver repository.Verification) *AuthService {
	return &AuthService{repo: repo, ver: ver}
}

func (s *AuthService) CreateUser(user models.User) (int, error) {

	if ok := IsValidEmail(user.Email); !ok {
		return 0, errors.New("invalid email")
	}
	if len(user.Username) < 3 {
		return 0, errors.New("username too short")
	}
	if user.Gender >= 3 {
		return 0, errors.New("invalid gender")
	}
	if !IsValidBirth(user.Birthday) {
		return 0, errors.New("invalid birthday")
	}
	if !IsValidPassword(user.Password) {
		return 0, errors.New("invalid password(пароль должен быть из 8 и более символов содержать хоть одну цифру и !, @, #, $, %, ^, &, *. ")
	}

	ok, err := s.repo.CheckLangId(user.Language)
	if err != nil {
		return 0, err
	}

	if !ok {
		return 0, errors.New("invalid language id")
	}

	answer, err := s.repo.IsEmailFree(user.Email)
	if err != nil {
		return 0, err
	}
	if answer {
		return 0, errors.New("User with this email already exists!!! ")
	}

	user.Password = generationPasswordHash(user.Password)

	return s.repo.CreateUser(user)
}

func generationPasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(signingKey)))
}

func IsValidEmail(email string) bool {
	// Простая, но эффективная проверка
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	return emailRegex.MatchString(email)
}

func IsValidBirth(b string) bool {
	birthRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	return birthRegex.MatchString(b)
}

func IsValidPassword(password string) bool {
	passwordRegex := regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*]{8,}$`)
	return passwordRegex.MatchString(password)
}

type TokenClaims struct {
	jwt.StandardClaims
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
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
		user.Username,
	})

	active, err := s.ver.GetUserActive(user.Id, username)
	if err != nil {
		return "", err

	}

	if active == false {
		return "", errors.New("user is not active")
	}

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
