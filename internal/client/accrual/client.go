package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/invinciblewest/gophermart/internal/logger"
	"github.com/invinciblewest/gophermart/internal/model"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type Client struct {
	client  *http.Client
	baseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{
		client:  http.DefaultClient,
		baseURL: baseURL,
	}
}

func (c *Client) GetOrderInfo(ctx context.Context, orderNumber string) (*model.AccrualResponse, int, error) {
	path, err := url.JoinPath(c.baseURL, "api", "orders", orderNumber)
	if err != nil {
		return nil, 0, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, 0, err
	}

	response, err := c.client.Do(request)
	if err != nil {
		return nil, 0, err
	}
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			logger.Log.Info("failed to close response body", zap.Error(err))
		}
	}(response.Body)

	if response.StatusCode == http.StatusTooManyRequests {
		retryAfterStr := response.Header.Get("Retry-After")
		seconds, err := strconv.Atoi(retryAfterStr)
		if err != nil {
			return nil, 0, err
		}

		return nil, seconds, nil
	}

	if response.StatusCode == http.StatusNoContent {
		return nil, 0, nil
	}

	if response.StatusCode != http.StatusOK {

		return nil, 0, errors.New("failed to get order info")
	}

	var accrualResponse *model.AccrualResponse
	if err = json.NewDecoder(response.Body).Decode(&accrualResponse); err != nil {
		return nil, 0, err
	}

	return accrualResponse, 0, nil
}
