package lms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client interface {
	GetLoan(ctx context.Context, loanID string) (*Loan, error)
	GetInstallments(ctx context.Context, userID string) ([]Installment, error)
	GetUpcomingInstallments(ctx context.Context) ([]Installment, error)
	GetOverdueInstallments(ctx context.Context) ([]Installment, error)
	UpdateLoanStatus(ctx context.Context, loanID, status string) error
	RecordPayment(ctx context.Context, req RecordPaymentRequest) error
}

type client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string, httpClient *http.Client) Client {
	return &client{baseURL: baseURL, httpClient: httpClient}
}

func (c *client) GetLoan(ctx context.Context, loanID string) (*Loan, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/loans/%s", c.baseURL, loanID), nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling LMS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LMS returned status %d", resp.StatusCode)
	}

	var loan Loan
	if err := json.NewDecoder(resp.Body).Decode(&loan); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &loan, nil
}

func (c *client) GetInstallments(ctx context.Context, userID string) ([]Installment, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/loans?user_id=%s", c.baseURL, userID), nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling LMS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LMS returned status %d", resp.StatusCode)
	}

	var installments []Installment
	if err := json.NewDecoder(resp.Body).Decode(&installments); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return installments, nil
}

func (c *client) GetUpcomingInstallments(ctx context.Context) ([]Installment, error) {
	return c.listInstallments(ctx, "upcoming")
}

func (c *client) GetOverdueInstallments(ctx context.Context) ([]Installment, error) {
	return c.listInstallments(ctx, "overdue")
}

func (c *client) listInstallments(ctx context.Context, status string) ([]Installment, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/installments?status=%s", c.baseURL, status), nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling LMS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LMS returned status %d", resp.StatusCode)
	}

	var installments []Installment
	if err := json.NewDecoder(resp.Body).Decode(&installments); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return installments, nil
}

func (c *client) UpdateLoanStatus(ctx context.Context, loanID, status string) error {
	body, _ := json.Marshal(map[string]string{"status": status})
	req, err := http.NewRequestWithContext(ctx, http.MethodPut,
		fmt.Sprintf("%s/loans/%s/status", c.baseURL, loanID), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("calling LMS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("LMS returned status %d", resp.StatusCode)
	}
	return nil
}

func (c *client) RecordPayment(ctx context.Context, reqBody RecordPaymentRequest) error {
	body, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/loans/%s/payments", c.baseURL, reqBody.LoanID), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("calling LMS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("LMS returned status %d", resp.StatusCode)
	}
	return nil
}
