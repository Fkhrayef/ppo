package psp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type httpClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHTTPClient(baseURL string, hc *http.Client) Client {
	return &httpClient{baseURL: baseURL, httpClient: hc}
}

func (c *httpClient) Charge(ctx context.Context, reqBody ChargeRequest) (*ChargeResponse, error) {
	body, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/charges", c.baseURL), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling PSP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("PSP returned status %d", resp.StatusCode)
	}

	var chargeResp ChargeResponse
	if err := json.NewDecoder(resp.Body).Decode(&chargeResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &chargeResp, nil
}

func (c *httpClient) Refund(ctx context.Context, reqBody RefundRequest) (*RefundResponse, error) {
	body, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/refunds", c.baseURL), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling PSP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("PSP returned status %d", resp.StatusCode)
	}

	var refundResp RefundResponse
	if err := json.NewDecoder(resp.Body).Decode(&refundResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &refundResp, nil
}
