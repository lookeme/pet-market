package controller

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"pet-market/api"
	"pet-market/internal/logger"
	"pet-market/internal/security"
	"pet-market/internal/service"
	"pet-market/internal/utils"

	"github.com/jackc/pgx/v5"
)

const bearer = "Bearer "

type Controller struct {
	Authorization  security.Authorization
	UserService    service.UserService
	BalanceService service.BalanceService
	OrderService   service.OrderService
	log            logger.Logger
}

func NewController(
	auth security.Authorization,
	usrService service.UserService,
	orderService service.OrderService,
	balanceService service.BalanceService,
	log *logger.Logger,
) *Controller {
	return &Controller{
		Authorization:  auth,
		UserService:    usrService,
		OrderService:   orderService,
		BalanceService: balanceService,
		log:            *log,
	}
}

func (s *Controller) GetBalance(w http.ResponseWriter, r *http.Request) {
	if !s.Authorization.Authorize(w, r) {
		s.writeResponse(w, r, http.StatusUnauthorized, nil)
		return
	}
	token := r.Header.Get("Authorization")
	token, tokenErr := security.GetToken(token)
	if tokenErr != nil {
		s.writeResponse(w, r, http.StatusInternalServerError, tokenErr)
	}
	ctx := context.Background()
	userID := security.GetUserID(token)
	balance, err := s.BalanceService.GetBalance(ctx, userID)
	if err != nil {
		s.writeResponse(w, r, http.StatusInternalServerError, err)
	} else {
		s.writeResponse(w, r, http.StatusOK, balance)
	}
}

func (s *Controller) WithdrawBalance(w http.ResponseWriter, r *http.Request) {
	if !s.Authorization.Authorize(w, r) {
		s.writeResponse(w, r, http.StatusUnauthorized, nil)
		return
	}
	token := r.Header.Get("Authorization")
	token, tokenErr := security.GetToken(token)
	if tokenErr != nil {
		s.writeResponse(w, r, http.StatusInternalServerError, tokenErr)
	}
	ctx := context.Background()
	b, errBody := io.ReadAll(r.Body)
	defer r.Body.Close()
	if errBody != nil {
		s.writeResponse(w, r, http.StatusInternalServerError, errBody)
	}
	var withdraw api.RequestWithdraw
	if err := json.Unmarshal(b, &withdraw); err != nil {
		s.writeResponse(w, r, http.StatusBadRequest, err)
	}
	userID := security.GetUserID(token)
	err := s.BalanceService.AddWithdraw(ctx, userID, withdraw)
	if err == nil {
		s.writeResponse(w, r, http.StatusOK, err)
	} else if errors.Is(err, utils.ErrInvalidOrderNum) {
		s.writeResponse(w, r, http.StatusUnprocessableEntity, err)
	} else if errors.Is(err, utils.ErrInsufficientFunds) {
		s.writeResponse(w, r, http.StatusPaymentRequired, err)
	} else {
		s.writeResponse(w, r, http.StatusInternalServerError, err)
	}

}

func (s *Controller) AuthorizeUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var user api.User
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err := json.Unmarshal(body, &user); err != nil {
		s.writeResponse(w, r, http.StatusBadRequest, err)
	}
	usr, queryErr := s.UserService.GetUserByName(ctx, user.Login)
	if errors.Is(queryErr, pgx.ErrNoRows) {
		s.writeResponse(w, r, http.StatusUnauthorized, errors.New("user is not registered"))
		return
	}
	if s.Authorization.VerifyPassword(user.Password, usr.Password) {
		s.writeToken(w, r, usr.Login, *usr.Id)
	} else {
		s.writeResponse(w, r, http.StatusUnauthorized, errors.New("wrong password"))
	}
}

func (s *Controller) OrderList(w http.ResponseWriter, r *http.Request) {
	if !s.Authorization.Authorize(w, r) {
		s.writeResponse(w, r, http.StatusUnauthorized, nil)
		return
	}
	token := r.Header.Get("Authorization")
	token, tokenErr := security.GetToken(token)
	if tokenErr != nil {
		s.writeResponse(w, r, http.StatusInternalServerError, tokenErr)
	}
	ctx := context.Background()
	userID := security.GetUserID(token)
	orders, err := s.OrderService.GetUserOrders(ctx, userID)

	if err != nil {
		s.writeResponse(w, r, http.StatusInternalServerError, err)
	} else if len(orders) == 0 {
		s.writeResponse(w, r, http.StatusNoContent, orders)
	} else {
		s.writeResponse(w, r, http.StatusOK, orders)
	}
}

func (s *Controller) UploadOrder(w http.ResponseWriter, r *http.Request) {
	if !s.Authorization.Authorize(w, r) {
		s.writeResponse(w, r, http.StatusUnauthorized, nil)
		return
	}
	ctx := context.Background()
	b, bodyErr := io.ReadAll(r.Body)
	defer r.Body.Close()
	if bodyErr != nil {
		s.writeResponse(w, r, http.StatusBadRequest, bodyErr)
	}
	orderNum := string(b)
	token := r.Header.Get("Authorization")
	token, tokenErr := security.GetToken(token)
	if tokenErr != nil {
		s.writeResponse(w, r, http.StatusInternalServerError, tokenErr)
	}
	userID := security.GetUserID(token)
	ordErr := s.OrderService.CreateOrder(ctx, orderNum, userID)
	if errors.Is(ordErr, utils.ErrInvalidOrderNum) {
		s.writeResponse(w, r, http.StatusUnprocessableEntity, tokenErr)
	} else if errors.Is(ordErr, utils.ErrOrderNumAttachedToAnotherUser) {
		s.writeResponse(w, r, http.StatusConflict, tokenErr)
	} else if errors.Is(ordErr, utils.ErrOrderNumIsAlreadyRegistered) {
		s.writeResponse(w, r, http.StatusOK, tokenErr)
	}
	if ordErr == nil {
		s.writeResponse(w, r, http.StatusAccepted, nil)
	} else {
		s.writeResponse(w, r, http.StatusInternalServerError, ordErr)
	}
}
func (s *Controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user api.User
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err := json.Unmarshal(body, &user); err != nil {
		s.writeResponse(w, r, http.StatusBadRequest, err)
	}
	ctx := context.Background()
	usr, queryErr := s.UserService.GetUserByName(ctx, user.Login)
	if errors.Is(queryErr, pgx.ErrNoRows) {
		err := s.UserService.CreateUser(ctx, user)
		if err != nil {
			s.writeResponse(w, r, http.StatusInternalServerError, err)
		} else {
			s.writeToken(w, r, usr.Login, *usr.Id)
		}
	}
	if queryErr != nil {
		s.writeResponse(w, r, http.StatusInternalServerError, queryErr)
	}
	if user.Login == usr.Login {
		s.writeResponse(w, r, http.StatusConflict, errors.New("user already exists"))
	}
}

func (s *Controller) WithdrawalsList(w http.ResponseWriter, r *http.Request) {
	if !s.Authorization.Authorize(w, r) {
		s.writeResponse(w, r, http.StatusUnauthorized, nil)
		return
	}
	ctx := context.Background()
	token := r.Header.Get("Authorization")
	token, tokenErr := security.GetToken(token)
	if tokenErr != nil {
		s.writeResponse(w, r, http.StatusInternalServerError, tokenErr)
	}
	userID := security.GetUserID(token)
	withdrawals, err := s.BalanceService.GetAllWithdraws(ctx, userID)
	if err != nil {
		s.writeResponse(w, r, http.StatusInternalServerError, err)
	}
	if len(withdrawals) == 0 {
		s.writeResponse(w, r, http.StatusNoContent, nil)
	} else {
		s.writeResponse(w, r, http.StatusOK, withdrawals)
	}
}

func (s *Controller) writeToken(w http.ResponseWriter, r *http.Request, userName string, userID int) {
	token, err := s.Authorization.BuildJWTString(userName, userID)
	if err != nil {
		s.writeResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Authorization", bearer+token)
	s.writeResponse(w, r, http.StatusOK, err)
}
func (s *Controller) writeResponse(w http.ResponseWriter, _ *http.Request, code int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if response == nil {
		w.WriteHeader(code)
		return
	}
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		s.log.Log.Error(err.Error())
	}
}
