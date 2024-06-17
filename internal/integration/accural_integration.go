package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"pet-market/internal/logger"
	"time"

	"go.uber.org/zap"
)

var url = "/api/orders/"

type AccuralIntegration struct {
	host   string
	Client *http.Client
	log    logger.Logger
}

func New(host string) *AccuralIntegration {
	return &AccuralIntegration{
		host:   host,
		Client: &http.Client{},
	}
}

func (a *AccuralIntegration) GetOrder(orderNumber string, timeout time.Duration) (*OrderAccural, error) {
	a.log.Log.Info("create request to accural...")
	req, err := http.NewRequest(http.MethodGet, a.host+url+orderNumber, nil)
	if err != nil {
		a.log.Log.Error(err.Error())
		return nil, err
	}
	res, getErr := a.Client.Do(req)
	if getErr != nil {
		a.log.Log.Error(err.Error())
		return nil, getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		a.log.Log.Error(readErr.Error())
		return nil, readErr
	}
	order := OrderAccural{}
	jsonErr := json.Unmarshal(body, &order)
	if jsonErr != nil {
		return nil, jsonErr
	}
	a.log.Log.Info("return order",
		zap.String("order", order.Order),
		zap.String("status", order.Status),
		zap.Int("accural", order.Accrual))
	return &order, nil
}

type OrderAccural struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}
