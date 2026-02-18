package ipresolver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResolveIP_WithIPAddressField(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get(apiKeyHeader); got != "test-key" {
			t.Fatalf("expected API key header test-key, got %q", got)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ip_address":"203.0.113.10"}`))
	}))
	defer server.Close()

	ip, err := ResolveIP(context.Background(), server.URL, "test-key")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if ip != "203.0.113.10" {
		t.Fatalf("expected IP 203.0.113.10, got %q", ip)
	}
}

func TestResolveIP_WithLegacyIPField(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ip":"198.51.100.20"}`))
	}))
	defer server.Close()

	ip, err := ResolveIP(context.Background(), server.URL, "unused")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if ip != "198.51.100.20" {
		t.Fatalf("expected IP 198.51.100.20, got %q", ip)
	}
}

func TestResolveIP_Non200(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	}))
	defer server.Close()

	_, err := ResolveIP(context.Background(), server.URL, "unused")
	if err == nil {
		t.Fatalf("expected error for non-200 response")
	}
}

func TestResolveIP_InvalidIP(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ip_address":"not-an-ip"}`))
	}))
	defer server.Close()

	_, err := ResolveIP(context.Background(), server.URL, "unused")
	if err == nil {
		t.Fatalf("expected validation error for invalid IP")
	}
}
