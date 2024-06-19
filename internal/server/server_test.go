package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"pet-market/api"
	"pet-market/internal/controller"
	"pet-market/internal/integration"
	"pet-market/internal/logger"
	"pet-market/internal/mocks"
	"pet-market/internal/models"
	"pet-market/internal/security"
	"pet-market/internal/service"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var hash, _ = security.HashPassword("qwerty")

const (
	currenBalance    = float32(350)
	withdrawnBalance = float32(60)
)

func Pointer(val float32) *float32 { return &val }
func TestShop(t *testing.T) {
	log, _ := zap.NewDevelopment()
	zlog := logger.Logger{
		Log: log,
	}
	auth := security.New(&zlog)
	mockCtrl := gomock.NewController(t)

	orderAccural := integration.OrderAccural{
		Accrual: Pointer(200),
		Order:   "12345678903",
		Status:  "PROCESS",
	}

	order := models.Order{
		Accrual: Pointer(200),
		OrderID: "12345678903",
		Status:  "PROCESS",
		UserID:  1,
	}
	ID := 1
	defer mockCtrl.Finish()
	userRepoMock := mocks.NewMockUserRepository(mockCtrl)
	userRepoMock.EXPECT().Save(context.Background(), "login", "qwerty").Return(1, nil).AnyTimes()
	userRepoMock.EXPECT().GetUserByLogin(context.Background(), "login").Return(api.User{
		Id:       &ID,
		Login:    "login",
		Password: hash,
	}, nil).AnyTimes()
	userRepoMock.EXPECT().GetUserByLogin(context.Background(), "user").Return(api.User{}, pgx.ErrNoRows).AnyTimes()

	orderRepoMock := mocks.NewMockOrderRepository(mockCtrl)

	orderRepoMock.EXPECT().GetByOrderNumber(context.Background(), "12345678903").Return(
		models.Order{}, pgx.ErrNoRows).AnyTimes()

	orderRepoMock.EXPECT().GetAll(context.Background(), 1).Return([]models.Order{
		{
			OrderID:    "12345678903",
			Status:     "INVALID",
			UploadedAt: time.Now(),
			UserID:     1,
		},
		{
			OrderID:    "12345678904",
			Accrual:    Pointer(200),
			Status:     "PROCESSED",
			UploadedAt: time.Now(),
			UserID:     1,
		},
		{
			OrderID:    "12345678904",
			Accrual:    Pointer(200),
			Status:     "PROCESSED",
			UploadedAt: time.Now(),
			UserID:     1,
		},
		{
			OrderID:    "12345678906",
			Accrual:    Pointer(200),
			Status:     "PROCESSED",
			UploadedAt: time.Now(),
			UserID:     1,
		},
	}, nil).AnyTimes()

	orderRepoMock.EXPECT().Save(context.Background(), order, 1).Return(nil).AnyTimes()
	balanceRepoMock := mocks.NewMockBalanceRepository(mockCtrl)
	balanceRepoMock.EXPECT().GetBalance(context.Background(), 1).Return(api.Balance{
		Current:   350,
		Withdrawn: 60,
	}, nil).AnyTimes()

	withdrawnRepoMock := mocks.NewMockWithdrawalsRepository(mockCtrl)
	withdrawnRepoMock.EXPECT().GetAllByUserID(context.Background(), 1).Return([]models.Withdraw{
		{
			OrderNum:    "12345678903",
			ProcessedAt: time.Now(),
			Sum:         20,
			UserID:      1,
		},
		{
			OrderNum:    "12345678904",
			ProcessedAt: time.Now(),
			Sum:         20,
			UserID:      1,
		},

		{
			OrderNum:    "12345678905",
			ProcessedAt: time.Now(),
			Sum:         20,
			UserID:      1,
		},
	}, nil).AnyTimes()

	client := mocks.NewTestClient(func(req *http.Request) *http.Response {
		jsonBody, _ := json.Marshal(orderAccural)
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(string(jsonBody))),
			Header:     make(http.Header),
		}
	})
	accural := integration.AccrualClient{
		Client: client,
		Log:    zlog,
	}
	userService := service.NewUserService(userRepoMock, auth, &zlog)
	balanceService := service.NewBalanceService(balanceRepoMock, withdrawnRepoMock)
	orderService := service.NewOrderService(&accural, orderRepoMock)
	ctr := controller.NewController(
		auth,
		userService,
		orderService,
		balanceService,
		&zlog,
	)
	t.Run("get balance test #1", func(t *testing.T) {
		token, _ := auth.BuildJWTString("login", 1)
		req := httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rctx := chi.NewRouteContext()
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()
		ctr.GetBalance(w, req)
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		body, readErr := io.ReadAll(res.Body)
		require.NoError(t, readErr)
		balance := api.Balance{}
		jsonErr := json.Unmarshal(body, &balance)
		require.NoError(t, jsonErr)
		checkBalance := currenBalance - withdrawnBalance
		assert.Equal(t, balance.Current, checkBalance)
		assert.Equal(t, balance.Withdrawn, withdrawnBalance)
		err := res.Body.Close()
		require.NoError(t, err)
	})

	t.Run("get user orders #2", func(t *testing.T) {
		token, _ := auth.BuildJWTString("login", 1)
		req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		ctr.OrderList(w, req)
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		body, readErr := io.ReadAll(res.Body)
		require.NoError(t, readErr)
		var orders []api.OrderResponse
		jsonErr := json.Unmarshal(body, &orders)
		require.NoError(t, jsonErr)
		assert.Equal(t, len(orders), 4)
		err := res.Body.Close()
		require.NoError(t, err)
	})

	t.Run("get user withdraws #3", func(t *testing.T) {
		token, _ := auth.BuildJWTString("login", 1)
		req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rctx := chi.NewRouteContext()
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()
		ctr.WithdrawalsList(w, req)
		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		body, readErr := io.ReadAll(res.Body)
		require.NoError(t, readErr)
		var withdraws []api.ResponseWithdraw
		jsonErr := json.Unmarshal(body, &withdraws)
		require.NoError(t, jsonErr)
		assert.Equal(t, len(withdraws), 3)
		err := res.Body.Close()
		require.NoError(t, err)
	})

	t.Run("create user test #4", func(t *testing.T) {
		user := api.User{
			Id:       &ID,
			Login:    "login",
			Password: "qwerty",
		}
		jsonData, err := json.Marshal(&user)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBuffer(jsonData))
		w := httptest.NewRecorder()
		ctr.AuthorizeUser(w, req)
		res := w.Result()
		res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		token := res.Header.Get("Authorization")
		token, tokenErr := security.GetToken(token)
		assert.True(t, len(token) > 0)
		require.NoError(t, tokenErr)
	})

	t.Run("create user oder user  #5", func(t *testing.T) {
		token, _ := auth.BuildJWTString("login", 1)
		req := httptest.NewRequest(http.MethodPost, "/api/user/order", bytes.NewBuffer([]byte("12345678903")))
		req.Header.Set("Authorization", "Bearer "+token)
		rctx := chi.NewRouteContext()
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()
		ctr.UploadOrder(w, req)
		res := w.Result()
		res.Body.Close()
		assert.Equal(t, http.StatusAccepted, res.StatusCode)
	})

}
