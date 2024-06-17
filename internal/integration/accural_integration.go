package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"pet-market/internal/logger"
	"time"
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

	req, err := http.NewRequest(http.MethodGet, a.host+url+orderNumber, nil)
	if err != nil {
		return nil, err
	}
	res, getErr := a.Client.Do(req)
	if getErr != nil {
		return nil, getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	order := OrderAccural{}
	jsonErr := json.Unmarshal(body, &order)
	if jsonErr != nil {
		return nil, jsonErr
	}
	return &order, nil
}

type OrderAccural struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}
