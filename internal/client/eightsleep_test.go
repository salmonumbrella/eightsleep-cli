package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// mockServer builds a test server that can serve a handful of endpoints the client expects.
func mockServer(t *testing.T) (*httptest.Server, *Client) {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"user":{"userId":"uid-123","currentDevice":{"id":"dev-1"}}}`))
	})

	mux.HandleFunc("/users/uid-123/temperature", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"currentLevel":5,"currentState":{"type":"on"}}`))
			return
		}
		if r.Method == http.MethodPut {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.NotFound(w, r)
	})

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		// first call rate limits, second succeeds
		if r.Header.Get("X-Test-Retry") == "done" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"ok":true}`))
			return
		}
		w.WriteHeader(http.StatusTooManyRequests)
	})

	srv := httptest.NewServer(mux)

	// client with pre-set token to skip auth
	c := New("email", "pass", "", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	return srv, c
}

func TestRequireUserFilledAutomatically(t *testing.T) {
	srv, c := mockServer(t)
	defer srv.Close()

	// UserID empty; GetStatus should fetch it from /users/me
	st, err := c.GetStatus(context.Background())
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	if c.UserID != "uid-123" {
		t.Fatalf("expected user id populated, got %s", c.UserID)
	}
	if st.CurrentLevel != 5 || st.CurrentState.Type != "on" {
		t.Fatalf("unexpected status %+v", st)
	}
}

func Test429Retry(t *testing.T) {
	count := 0
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		count++
		if count == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()

	start := time.Now()
	if err := c.do(context.Background(), http.MethodGet, "/ping", nil, nil, nil); err != nil {
		t.Fatalf("do retry: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 attempts, got %d", count)
	}
	if elapsed := time.Since(start); elapsed < 2*time.Second {
		t.Fatalf("expected backoff, got %v", elapsed)
	}
}

func TestContextTimeout(t *testing.T) {
	// Server that delays response beyond context timeout
	mux := http.NewServeMux()
	mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		// Sleep longer than the context timeout
		time.Sleep(500 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c := New("email", "pass", "uid", "", "")
	c.BaseURL = srv.URL
	c.token = "t"
	c.tokenExp = time.Now().Add(time.Hour)
	c.HTTP = srv.Client()
	c.HTTP.Timeout = 100 * time.Millisecond // Short HTTP timeout

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := c.do(ctx, http.MethodGet, "/slow", nil, nil, nil)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	// Should be context deadline exceeded
	if ctx.Err() != context.DeadlineExceeded {
		t.Fatalf("expected context.DeadlineExceeded, got ctx.Err()=%v, err=%v", ctx.Err(), err)
	}
}

func TestRetryExhaustion(t *testing.T) {
	t.Run("MaxRetries=0", func(t *testing.T) {
		count := 0
		mux := http.NewServeMux()
		mux.HandleFunc("/rate-limited", func(w http.ResponseWriter, r *http.Request) {
			count++
			w.WriteHeader(http.StatusTooManyRequests)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		c := New("email", "pass", "uid", "", "")
		c.BaseURL = srv.URL
		c.token = "t"
		c.tokenExp = time.Now().Add(time.Hour)
		c.HTTP = srv.Client()
		c.MaxRetries = 0 // No retries

		err := c.do(context.Background(), http.MethodGet, "/rate-limited", nil, nil, nil)
		if err == nil {
			t.Fatal("expected error after exhausting retries, got nil")
		}
		if count != 1 {
			t.Fatalf("expected 1 attempt with MaxRetries=0, got %d", count)
		}
	})

	t.Run("MaxRetries=1", func(t *testing.T) {
		count := 0
		mux := http.NewServeMux()
		mux.HandleFunc("/rate-limited", func(w http.ResponseWriter, r *http.Request) {
			count++
			w.WriteHeader(http.StatusTooManyRequests)
		})
		srv := httptest.NewServer(mux)
		defer srv.Close()

		c := New("email", "pass", "uid", "", "")
		c.BaseURL = srv.URL
		c.token = "t"
		c.tokenExp = time.Now().Add(time.Hour)
		c.HTTP = srv.Client()
		c.MaxRetries = 1 // One retry (2 total attempts)

		start := time.Now()
		err := c.do(context.Background(), http.MethodGet, "/rate-limited", nil, nil, nil)
		if err == nil {
			t.Fatal("expected error after exhausting retries, got nil")
		}
		if count != 2 {
			t.Fatalf("expected 2 attempts with MaxRetries=1, got %d", count)
		}
		// Should have waited for backoff between attempts
		if elapsed := time.Since(start); elapsed < time.Second {
			t.Fatalf("expected backoff delay, elapsed=%v", elapsed)
		}
	})
}
