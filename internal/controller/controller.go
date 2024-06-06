package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"pet-market/api"
	"pet-market/internal/logger"
	"pet-market/internal/security"
	"pet-market/internal/service"
)

const bearer = "Bearer "

type Controller struct {
	Authorization  security.Authorization
	UserService    service.UserService
	BalanceService service.BalanceService
	OrderService   service.OrderService
	log            logger.Logger
}

func NewController() *Controller {
	return &Controller{}
}
func (s *Controller) GetOrdersNumber(w http.ResponseWriter, r *http.Request, number string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Controller) GetBalance(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Controller) WithdrawBalance(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Controller) AuthorizeUser(w http.ResponseWriter, r *http.Request) {
	var user api.User
	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &user); err != nil {
		s.writeResponse(w, r, http.StatusBadRequest, err)
	}
	usr, err := s.UserService.GetUserByName(user.Name)
	if err != nil {
		s.writeResponse(w, r, http.StatusInternalServerError, err)
	}
	if s.Authorization.VerifyPassword(user.Password, usr.Password) {
		s.writeToken(w, r, usr.Name)
	} else {
		s.writeResponse(w, r, http.StatusUnauthorized, err)
	}
}

func (s *Controller) OrderList(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Controller) UploadOrder(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (s *Controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user api.User
	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &user); err != nil {
		s.writeResponse(w, r, http.StatusBadRequest, err)
	}
	usr, err := s.UserService.GetUserByName(user.Name)
	if err != nil {
		s.writeResponse(w, r, http.StatusInternalServerError, err)
	}
	if usr != nil {
		s.writeResponse(w, r, http.StatusConflict, err)
	} else {
		err = s.UserService.CreateUser(user)
		if err != nil {
			s.writeResponse(w, r, http.StatusInternalServerError, err)
		} else {
			s.writeToken(w, r, usr.Name)
		}
	}
}

func (s *Controller) writeToken(w http.ResponseWriter, r *http.Request, userName string) {
	token, err := s.Authorization.BuildJWTString(userName)
	if err != nil {
		s.writeResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Authorization", bearer+token)
	s.writeResponse(w, r, http.StatusOK, err)
}
func (s *Controller) writeResponse(w http.ResponseWriter, _ *http.Request, code int, response interface{}) {
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
