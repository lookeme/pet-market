package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"pet-market/internal/logger"
	"pet-market/internal/models"
	"time"

	"go.uber.org/zap"
)

var url = "/api/orders/"

type AccrualClient struct {
	host   string
	Client *http.Client
	Log    logger.Logger
}

func New(host string, log *logger.Logger, timeout time.Duration) *AccrualClient {
	return &AccrualClient{
		host: host,
		Client: &http.Client{
			Timeout: timeout,
		},
		Log: *log,
	}
}

func (a *AccrualClient) GetOrder(orderNumber string) (*OrderAccural, error) {
	a.Log.Log.Info("create request to accural.", zap.String("orderNum", orderNumber))
	req, err := http.NewRequest(http.MethodGet, a.host+url+orderNumber, nil)
	req.Header.Set("Content-Length", "0")
	if err != nil {
		a.Log.Log.Error(err.Error())
		return nil, err
	}
	res, getErr := a.Client.Do(req)

	if getErr != nil {
		a.Log.Log.Error(getErr.Error())
		return nil, getErr
	}

	a.Log.Log.Info("response status",
		zap.String("status", res.Status),
		zap.Int("status", res.StatusCode),
	)
	order := OrderAccural{}
	if res.StatusCode != http.StatusOK {
		order.Order = orderNumber
		order.Status = models.PROCESSING
		return &order, nil
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		a.Log.Log.Error(readErr.Error())
		return nil, readErr
	}
	a.Log.Log.Info("accural response", zap.String("body", string(body)))
	jsonErr := json.Unmarshal(body, &order)
	if jsonErr != nil {
		a.Log.Log.Error(jsonErr.Error())
		return nil, jsonErr
	}
	return &order, nil
}

type OrderAccural struct {
	Order   string   `json:"order"`
	Status  string   `json:"status"`
	Accrual *float32 `json:"accrual"`
}
