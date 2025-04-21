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

func NewResponseError(w http.ResponseWriter, statusCode int, message string, err error, logger *zap.Logger) {
	logger.Error(message, zap.Error(err))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if encodeErr := json.NewEncoder(w).Encode(ErrorResp{Message: message}); encodeErr != nil {
		logger.Error("Failed to encode error response", zap.Error(encodeErr))
		http.Error(w, http.StatusText(statusCode), statusCode)
	}
}

func ResponseServer(resp interface{}, w http.ResponseWriter, logger *zap.Logger) {
	w.Header().Set("Content-Type", "application/json")

	if resp == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		NewResponseError(
			w,
			http.StatusInternalServerError,
			"Failed to encode JSON response",
			err,
			logger,
		)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetUserId(w http.ResponseWriter, r *http.Request) (int, error) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		NewResponseError(w, http.StatusBadRequest, "", err, nil)
		logrus.Println(err)
	}
	return id, nil
}

func InitConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("configs")
	return viper.ReadInConfig()
}
