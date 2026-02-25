package lms

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLoan_Success(t *testing.T) {
	expected := Loan{
		ID:          "loan-001",
		UserID:      "user-123",
		Status:      "active",
		PaidAmount:  15000,
		TotalAmount: 60000,
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/loans/loan-001" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, srv.Client())
	loan, err := client.GetLoan(context.Background(), "loan-001")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loan.ID != expected.ID {
		t.Errorf("expected loan ID %q, got %q", expected.ID, loan.ID)
	}
	if loan.PaidAmount != expected.PaidAmount {
		t.Errorf("expected paid amount %d, got %d", expected.PaidAmount, loan.PaidAmount)
	}
}

func TestGetLoan_Non200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, srv.Client())
	_, err := client.GetLoan(context.Background(), "missing")

	if err == nil {
		t.Fatal("expected error for 404 response")
	}
}

func TestGetLoan_MalformedJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`not json`))
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, srv.Client())
	_, err := client.GetLoan(context.Background(), "loan-001")

	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}

func TestGetInstallments_Success(t *testing.T) {
	expected := []Installment{
		{ID: "inst-001", LoanID: "loan-001", Amount: 15000, Status: "paid", DueDate: "2026-01-15"},
		{ID: "inst-002", LoanID: "loan-001", Amount: 15000, Status: "upcoming", DueDate: "2026-02-15"},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Query().Get("user_id") != "user-123" {
			t.Errorf("unexpected user_id query param: %s", r.URL.Query().Get("user_id"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, srv.Client())
	installments, err := client.GetInstallments(context.Background(), "user-123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(installments) != 2 {
		t.Fatalf("expected 2 installments, got %d", len(installments))
	}
	if installments[0].ID != "inst-001" {
		t.Errorf("expected first installment ID %q, got %q", "inst-001", installments[0].ID)
	}
}

func TestGetUpcomingInstallments_Success(t *testing.T) {
	expected := []Installment{
		{ID: "inst-005", LoanID: "loan-002", Amount: 20000, Status: "upcoming", DueDate: "2026-03-01"},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("status") != "upcoming" {
			t.Errorf("expected status=upcoming, got %s", r.URL.Query().Get("status"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, srv.Client())
	installments, err := client.GetUpcomingInstallments(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(installments) != 1 {
		t.Fatalf("expected 1 installment, got %d", len(installments))
	}
}

func TestGetOverdueInstallments_Success(t *testing.T) {
	expected := []Installment{
		{ID: "inst-009", LoanID: "loan-003", Amount: 12000, Status: "overdue", DueDate: "2026-01-01"},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("status") != "overdue" {
			t.Errorf("expected status=overdue, got %s", r.URL.Query().Get("status"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expected)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, srv.Client())
	installments, err := client.GetOverdueInstallments(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(installments) != 1 {
		t.Fatalf("expected 1 installment, got %d", len(installments))
	}
	if installments[0].Status != "overdue" {
		t.Errorf("expected status %q, got %q", "overdue", installments[0].Status)
	}
}

func TestUpdateLoanStatus_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/loans/loan-001/status" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected application/json content type")
		}

		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["status"] != "refunded" {
			t.Errorf("expected status=refunded in body, got %q", body["status"])
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, srv.Client())
	err := client.UpdateLoanStatus(context.Background(), "loan-001", "refunded")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateLoanStatus_Non200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, srv.Client())
	err := client.UpdateLoanStatus(context.Background(), "loan-001", "refunded")

	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestRecordPayment_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/loans/loan-001/payments" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var req RecordPaymentRequest
		json.Unmarshal(body, &req)

		if req.LoanID != "loan-001" {
			t.Errorf("expected loan_id=loan-001, got %q", req.LoanID)
		}
		if req.Amount != 15000 {
			t.Errorf("expected amount=15000, got %d", req.Amount)
		}
		if req.TransactionID != "txn-abc" {
			t.Errorf("expected transaction_id=txn-abc, got %q", req.TransactionID)
		}

		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, srv.Client())
	err := client.RecordPayment(context.Background(), RecordPaymentRequest{
		LoanID:        "loan-001",
		InstallmentID: "inst-001",
		Amount:        15000,
		TransactionID: "txn-abc",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRecordPayment_Non200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, srv.Client())
	err := client.RecordPayment(context.Background(), RecordPaymentRequest{
		LoanID:        "loan-001",
		InstallmentID: "inst-001",
		Amount:        15000,
		TransactionID: "txn-abc",
	})

	if err == nil {
		t.Fatal("expected error for 400 response")
	}
}

func TestGetLoan_ServerDown(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))
	srv.Close() // close immediately to simulate unreachable server

	client := NewHTTPClient(srv.URL, srv.Client())
	_, err := client.GetLoan(context.Background(), "loan-001")

	if err == nil {
		t.Fatal("expected error when server is down")
	}
}

func TestGetLoan_CancelledContext(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Loan{ID: "loan-001"})
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before making the call

	client := NewHTTPClient(srv.URL, srv.Client())
	_, err := client.GetLoan(ctx, "loan-001")

	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}
