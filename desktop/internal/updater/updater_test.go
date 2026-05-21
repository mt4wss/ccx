package updater

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestCompareSemver(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.1", "1.0.0", 1},
		{"1.0.0", "1.0.1", -1},
		{"2.0.0", "1.99.99", 1},
		{"1.0.0", "1.0.0-rc1", 1},
		{"1.0.0-rc2", "1.0.0-rc1", 1},
		{"1.0.0-rc1", "1.0.0", -1},
		{"1.2", "1.2.0", 0},
		{"1.2.3.4", "1.2.3", 1},
	}
	for _, c := range cases {
		got, err := compareSemver(c.a, c.b)
		if err != nil {
			t.Fatalf("compareSemver(%q,%q): unexpected error: %v", c.a, c.b, err)
		}
		if got != c.want {
			t.Errorf("compareSemver(%q,%q) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestCompareSemverError(t *testing.T) {
	if _, err := compareSemver("abc", "1.0.0"); err == nil {
		t.Error("expected error on invalid semver, got nil")
	}
}

func TestPlatformAssetName(t *testing.T) {
	name := PlatformAssetName("2.7.9")
	if name == "" {
		t.Error("PlatformAssetName 返回空")
	}
}

func TestVerifyRequiresSha256URL(t *testing.T) {
	u := New("1.0.0")
	if err := u.Verify(t.Context(), "unused", ""); err != nil {
		t.Errorf("expected nil when sha256 URL is empty (skip), got: %v", err)
	}
}

func TestDoWithRetrySuccess(t *testing.T) {
	var attempts atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := attempts.Add(1)
		if n < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, srv.URL, nil)
	resp, err := doWithRetry(t.Context(), client, req, 3)
	if err != nil {
		t.Fatalf("expected success, got: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if n := attempts.Load(); n != 3 {
		t.Errorf("expected 3 attempts, got %d", n)
	}
}

func TestDoWithRetryExhausted(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer srv.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, srv.URL, nil)
	_, err := doWithRetry(t.Context(), client, req, 2)
	if err == nil {
		t.Fatal("expected error after retries exhausted")
	}
}

func TestDoWithRetryNonTransient(t *testing.T) {
	var attempts atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, srv.URL, nil)
	resp, err := doWithRetry(t.Context(), client, req, 3)
	if err != nil {
		t.Fatalf("expected no error (404 returned directly), got: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
	if n := attempts.Load(); n != 1 {
		t.Errorf("expected 1 attempt (no retry for 404), got %d", n)
	}
}
