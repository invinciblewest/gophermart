package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/invinciblewest/gophermart/internal/helper"
	"github.com/invinciblewest/gophermart/internal/logger"
	"github.com/invinciblewest/gophermart/internal/model"
	"github.com/invinciblewest/gophermart/internal/usecase"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

type Handler struct {
	UserUseCase    usecase.UserUseCase
	OrderUseCase   usecase.OrderUseCase
	BalanceUseCase usecase.BalanceUseCase
}

func NewHandler(
	userUseCase usecase.UserUseCase,
	orderUseCase usecase.OrderUseCase,
	balanceUseCase usecase.BalanceUseCase,
) *Handler {
	return &Handler{
		UserUseCase:    userUseCase,
		OrderUseCase:   orderUseCase,
		BalanceUseCase: balanceUseCase,
	}
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := h.UserUseCase.RegisterAndLogin(r.Context(), &user)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrEmptyLoginOrPassword):
			w.WriteHeader(http.StatusBadRequest)
			return
		case errors.Is(err, model.ErrUserAlreadyExists):
			w.WriteHeader(http.StatusConflict)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			logger.Log.Info("failed to register user", zap.Error(err))
			return
		}
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token, err := h.UserUseCase.Login(r.Context(), user)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrEmptyLoginOrPassword):
			w.WriteHeader(http.StatusBadRequest)
			return
		case errors.Is(err, model.ErrUserNotFound):
			logger.Log.Info("user not found", zap.String("login", user.Login))
			w.WriteHeader(http.StatusUnauthorized)
			return
		case errors.Is(err, model.ErrInvalidPassword):
			logger.Log.Info("invalid password", zap.String("login", user.Login))
			w.WriteHeader(http.StatusUnauthorized)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			logger.Log.Info("failed to login user", zap.Error(err))
			return
		}
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) AddOrder(w http.ResponseWriter, r *http.Request) {
	userID, err := helper.GetUserID(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orderNumber, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			logger.Log.Error("failed to close request body", zap.Error(err))
		}
	}(r.Body)

	if len(orderNumber) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	order := model.Order{
		UserID: userID,
		Number: string(orderNumber),
	}

	if err = h.OrderUseCase.AddOrder(r.Context(), &order); err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidOrderNumber):
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		case errors.Is(err, model.ErrOrderAlreadyExists):
			w.WriteHeader(http.StatusOK)
			return
		case errors.Is(err, model.ErrOrderAlreadyExistsForAnotherUser):
			w.WriteHeader(http.StatusConflict)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			logger.Log.Info("failed to add order", zap.Error(err))
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userID, err := helper.GetUserID(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	orders, err := h.OrderUseCase.GetByUser(r.Context(), userID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		logger.Log.Info("failed to get orders", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(orders); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Log.Info("failed to encode orders", zap.Error(err))
		return
	}
}

func (h *Handler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	userID, err := helper.GetUserID(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	balance, err := h.BalanceUseCase.GetUserBalance(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(balance); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Log.Info("failed to encode balance", zap.Error(err))
		return
	}
}

func (h *Handler) WithdrawBalance(w http.ResponseWriter, r *http.Request) {
	userID, err := helper.GetUserID(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var withdrawRequest model.WithdrawRequest
	if err = json.NewDecoder(r.Body).Decode(&withdrawRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = h.BalanceUseCase.WithdrawBalance(r.Context(), userID, withdrawRequest); err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidWithdrawSum):
			w.WriteHeader(http.StatusPaymentRequired)
			return
		case errors.Is(err, model.ErrInvalidOrderNumber):
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			logger.Log.Info("failed to withdraw balance", zap.Error(err))
			return
		}
	}

}

func (h *Handler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	userID, err := helper.GetUserID(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	withdrawals, err := h.BalanceUseCase.GetWithdrawals(r.Context(), userID)
	if err != nil {
		if errors.Is(err, model.ErrWithdrawalNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		logger.Log.Info("failed to get withdrawals", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(withdrawals); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Log.Info("failed to encode withdrawals", zap.Error(err))
		return
	}
}
