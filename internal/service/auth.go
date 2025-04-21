package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"lang/internal/repository"
	"lang/pkg/models"
	"net/http"
	"regexp"
	"time"
)

const (
	tokenTTl   = 12 * time.Hour
	signingKey = "nv4q[hq3083)(_#mlnbnkamcsn4"
)

type AuthService struct {
	repo repository.Authorization
	ver  repository.Verification
}

func NewAuthService(repo repository.Authorization, ver repository.Verification) *AuthService {
	return &AuthService{repo: repo, ver: ver}
}

func (s *AuthService) CreateUser(user models.User) (int, models.ErrorResponse) {
	var errResp models.ErrorResponse

	// Валидация входных данных
	switch {
	case !IsValidEmail(user.Email):
		errResp = models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid email address",
		}
	case len(user.Username) < 3:
		errResp = models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Username must be at least 3 characters long",
		}
	case user.Gender >= 3:
		errResp = models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid gender value",
		}
	case !IsValidBirth(user.Birthday):
		errResp = models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid birthday date",
		}
	case !IsValidPassword(user.Password):
		errResp = models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Password must be at least 8 characters long and contain at least one number and special character (!@#$%^&*)",
		}
	}

	if errResp.Code != 0 {
		return 0, errResp
	}

	// Проверка существования языка
	ok, err := s.repo.CheckLangId(user.Language)
	if err != nil {
		return 0, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to validate language",
			Error:   err,
		}
	}
	if !ok {
		return 0, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid language ID",
		}
	}

	// Проверка уникальности email
	answer, err := s.repo.IsEmailFree(user.Email)
	if err != nil {
		return 0, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to check email availability",
			Error:   err,
		}
	}
	if answer {
		return 0, models.ErrorResponse{
			Code:    http.StatusConflict,
			Message: "User with this email already exists",
		}
	}

	// Хеширование пароля и создание пользователя
	user.Password = generationPasswordHash(user.Password)
	userID, err := s.repo.CreateUser(user)
	if err != nil {
		return 0, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create user",
			Error:   err,
		}
	}

	return userID, models.ErrorResponse{}
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

func (s *AuthService) GenerationToken(username, password string) (string, models.ErrorResponse) {
	var errUserNotFound = errors.New("user not found")

	// Валидация входных параметров
	if username == "" || password == "" {
		return "", models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Username and password are required",
		}
	}

	// Получаем пользователя из БД
	hashedPassword := generationPasswordHash(password)
	user, err := s.repo.GetUser(username, hashedPassword)
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			return "", models.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid username or password",
			}
		}
		return "", models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get user data",
			Error:   err,
		}
	}

	// Проверяем активность пользователя
	active, err := s.ver.GetUserActive(user.Id, username)
	if err != nil {
		return "", models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to check user status",
			Error:   err,
		}
	}

	if !active {
		return "", models.ErrorResponse{
			Code:    http.StatusForbidden,
			Message: "User account is not active",
		}
	}

	// Создаем JWT токен
	claims := &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTl).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId:   user.Id,
		Username: user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate token",
			Error:   err,
		}
	}

	return signedToken, models.ErrorResponse{}
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
