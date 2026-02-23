package product

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client interface {
	RestockItem(ctx context.Context, productID string, quantity int) error
}

type client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string, httpClient *http.Client) Client {
	return &client{baseURL: baseURL, httpClient: httpClient}
}

func (c *client) RestockItem(ctx context.Context, productID string, quantity int) error {
	body, _ := json.Marshal(RestockRequest{Quantity: quantity})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/products/%s/restock", c.baseURL, productID), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("calling product service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("product service returned status %d", resp.StatusCode)
	}
	return nil
}
