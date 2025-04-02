package utils

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type ErrorResp struct {
	Message string `json:"message"`
}
type StatusResponse struct {
	Status string `json:"status"`
}

func NewResponseError(w http.ResponseWriter, statusCode int, message string) {
	logrus.Error(message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(ErrorResp{Message: message}); err != nil {
		logrus.Error(err)
		return
	}
}

func BadRequest(w http.ResponseWriter, err error, logger *zap.Logger) {
	logger.Info("the user entered incorrect data", zap.Error(err))
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func InternalServerError(w http.ResponseWriter, err error, logger *zap.Logger) {
	logger.Error(http.StatusText(http.StatusInternalServerError), zap.Error(err))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func Forbidden(w http.ResponseWriter, err error, logger *zap.Logger) {
	logger.Info(http.StatusText(http.StatusForbidden), zap.Error(err))
	http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
}

func NotFoundErr(w http.ResponseWriter, err error, logger *zap.Logger) {
	logger.Info(http.StatusText(http.StatusNotFound), zap.Error(err))
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func ResetContentServerError(w http.ResponseWriter, err error, logger *zap.Logger) {
	logger.Info(http.StatusText(http.StatusResetContent), zap.Error(err))
	http.Error(w, http.StatusText(http.StatusResetContent), http.StatusResetContent)
}

func UnauthorizedError(w http.ResponseWriter, err error, logger *zap.Logger) {
	logger.Info(http.StatusText(http.StatusUnauthorized), zap.Error(err))
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

func ResponseServer(resp interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		NewResponseError(w, http.StatusInternalServerError, err.Error())
		logrus.Println(err)
		return
	}
}

func GetUserId(w http.ResponseWriter, r *http.Request) (int, error) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		NewResponseError(w, http.StatusBadRequest, err.Error())
		logrus.Println(err)
	}
	return id, nil
}

func InitConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("configs")
	return viper.ReadInConfig()
}
