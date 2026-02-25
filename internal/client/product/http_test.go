package product

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRestockItem_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/products/prod-001/restock" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected application/json content type")
		}

		body, _ := io.ReadAll(r.Body)
		var req RestockRequest
		json.Unmarshal(body, &req)

		if req.Quantity != 5 {
			t.Errorf("expected quantity=5, got %d", req.Quantity)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, srv.Client())
	err := client.RestockItem(context.Background(), "prod-001", 5)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRestockItem_Non200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, srv.Client())
	err := client.RestockItem(context.Background(), "prod-001", 5)

	if err == nil {
		t.Fatal("expected error for 503 response")
	}
}

func TestRestockItem_ServerDown(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))
	srv.Close()

	client := NewHTTPClient(srv.URL, srv.Client())
	err := client.RestockItem(context.Background(), "prod-001", 5)

	if err == nil {
		t.Fatal("expected error when server is down")
	}
}

func TestRestockItem_CancelledContext(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	client := NewHTTPClient(srv.URL, srv.Client())
	err := client.RestockItem(ctx, "prod-001", 5)

	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}
