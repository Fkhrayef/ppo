package psp

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCharge_Success(t *testing.T) {
	expected := ChargeResponse{TransactionID: "txn-123", Status: "captured"}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/charges" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected application/json content type")
		}

		body, _ := io.ReadAll(r.Body)
		var req ChargeRequest
		json.Unmarshal(body, &req)

		if req.Amount != 25000 {
			t.Errorf("expected amount=25000, got %d", req.Amount)
		}
		if req.Currency != "SAR" {
			t.Errorf("expected currency=SAR, got %q", req.Currency)
		}
		if req.CardToken != "tok-abc" {
			t.Errorf("expected card_token=tok-abc, got %q", req.CardToken)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expected)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, srv.Client())
	resp, err := client.Charge(context.Background(), ChargeRequest{
		Amount:    25000,
		Currency:  "SAR",
		CardToken: "tok-abc",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TransactionID != "txn-123" {
		t.Errorf("expected transaction_id=txn-123, got %q", resp.TransactionID)
	}
	if resp.Status != "captured" {
		t.Errorf("expected status=captured, got %q", resp.Status)
	}
}

func TestCharge_Non200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, srv.Client())
	_, err := client.Charge(context.Background(), ChargeRequest{
		Amount:    25000,
		Currency:  "SAR",
		CardToken: "tok-abc",
	})

	if err == nil {
		t.Fatal("expected error for 422 response")
	}
}

func TestCharge_MalformedJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{bad json`))
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, srv.Client())
	_, err := client.Charge(context.Background(), ChargeRequest{
		Amount: 25000, Currency: "SAR", CardToken: "tok-abc",
	})

	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}

func TestRefund_Success(t *testing.T) {
	expected := RefundResponse{RefundID: "ref-456", Status: "refunded"}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/refunds" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var req RefundRequest
		json.Unmarshal(body, &req)

		if req.OrderID != "order-001" {
			t.Errorf("expected order_id=order-001, got %q", req.OrderID)
		}
		if req.Amount != 15000 {
			t.Errorf("expected amount=15000, got %d", req.Amount)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, srv.Client())
	resp, err := client.Refund(context.Background(), RefundRequest{
		OrderID:   "order-001",
		Amount:    15000,
		Currency:  "SAR",
		CardToken: "tok-abc",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.RefundID != "ref-456" {
		t.Errorf("expected refund_id=ref-456, got %q", resp.RefundID)
	}
	if resp.Status != "refunded" {
		t.Errorf("expected status=refunded, got %q", resp.Status)
	}
}

func TestRefund_Non200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, srv.Client())
	_, err := client.Refund(context.Background(), RefundRequest{
		OrderID:   "order-001",
		Amount:    15000,
		Currency:  "SAR",
		CardToken: "tok-abc",
	})

	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestCharge_ServerDown(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))
	srv.Close()

	client := NewHTTPClient(srv.URL, srv.Client())
	_, err := client.Charge(context.Background(), ChargeRequest{
		Amount: 25000, Currency: "SAR", CardToken: "tok-abc",
	})

	if err == nil {
		t.Fatal("expected error when server is down")
	}
}

func TestRefund_CancelledContext(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(RefundResponse{RefundID: "ref-789"})
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := NewHTTPClient(srv.URL, srv.Client())
	_, err := client.Refund(ctx, RefundRequest{
		OrderID: "order-001", Amount: 15000, Currency: "SAR", CardToken: "tok-abc",
	})

	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}
