package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/salmonumbrella/eightsleep-cli/internal/client"
	"github.com/salmonumbrella/eightsleep-cli/internal/tokencache"
)

func setupTestEnv(t *testing.T) *client.Client {
	t.Helper()
	useTempKeyring(t)
	viper.Reset()
	viper.Set("email", "user@example.com")
	viper.Set("password", "pass")
	viper.Set("user_id", "uid-123")
	viper.Set("timezone", "UTC")
	viper.Set("output", "json")
	viper.Set("fields", []string{})
	viper.Set("timeout", "2s")
	viper.Set("retries", 0)

	srv := newTestServer(t)
	t.Cleanup(srv.Close)

	c := client.New("user@example.com", "pass", "uid-123", "", "")
	c.BaseURL = srv.URL
	c.AppBaseURL = srv.URL
	c.HTTP = srv.Client()
	if err := tokencache.Save(c.Identity(), "tok", time.Now().Add(time.Hour), c.UserID); err != nil {
		t.Fatalf("save token: %v", err)
	}

	prev := newClient
	newClient = func(email, password, userID, clientID, clientSecret string) *client.Client {
		return c
	}
	t.Cleanup(func() { newClient = prev })
	return c
}

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{
			"user": map[string]any{
				"userId": "uid-123",
				"currentDevice": map[string]any{
					"id": "dev-1",
				},
			},
		})
	})

	mux.HandleFunc("/users/uid-123/temperature", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			writeJSON(w, map[string]any{
				"currentLevel": 10,
				"currentState": map[string]any{"type": "on"},
			})
		case http.MethodPut:
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/users/uid-123/trends", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{
			"days": []map[string]any{
				{
					"day":                  "2024-01-01",
					"score":                90,
					"tnt":                  1,
					"respiratoryRate":      12.5,
					"heartRate":            60.0,
					"latencyAsleepSeconds": 300,
					"latencyOutSeconds":    120,
					"sleepDurationSeconds": 28800,
					"sleepQualityScore": map[string]any{
						"hrv":             map[string]any{"score": 85},
						"respiratoryRate": map[string]any{"score": 80},
					},
				},
			},
		})
	})

	mux.HandleFunc("/users/uid-123/intervals/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{
			"intervals": []map[string]any{{"id": "session-1"}},
		})
	})

	mux.HandleFunc("/devices/dev-1", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{"id": "dev-1"})
	})
	mux.HandleFunc("/devices/dev-1/peripherals", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{"peripherals": []string{"fan"}})
	})
	mux.HandleFunc("/devices/dev-1/online", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{"online": true})
	})

	mux.HandleFunc("/v2/users/uid-123/routines", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			writeJSON(w, map[string]any{
				"settings": map[string]any{
					"routines": []map[string]any{
						{
							"id":   "r1",
							"name": "Weekdays",
							"days": []int{1, 2, 3, 4, 5},
							"alarms": []map[string]any{
								{
									"alarmId": "a1",
									"enabled": true,
									"time":    "07:00",
									"settings": map[string]any{
										"vibration": map[string]any{"enabled": true},
									},
								},
							},
						},
					},
				},
				"state": map[string]any{
					"nextAlarm": map[string]any{
						"alarmId":       "a1",
						"nextTimestamp": "2024-01-02T07:00:00Z",
					},
				},
			})
		case http.MethodPut:
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/v2/users/uid-123/routines/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/v1/users/uid-123/routines", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/users/uid-123/vibration-test", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	return httptest.NewServer(mux)
}

func writeJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payload)
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w
	fn()
	_ = w.Close()
	os.Stdout = old
	data, _ := io.ReadAll(r)
	return strings.TrimSpace(string(data))
}

// flagState captures the value and Changed state of a flag.
type flagState struct {
	value   string
	changed bool
}

// resetFlagsOnCleanup saves the current flag values and Changed states for a
// command, and uses t.Cleanup() to restore them after the test completes.
// This prevents flag state leakage between tests.
func resetFlagsOnCleanup(t *testing.T, cmd *cobra.Command) {
	t.Helper()

	saved := make(map[string]flagState)
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		saved[f.Name] = flagState{
			value:   f.Value.String(),
			changed: f.Changed,
		}
	})

	t.Cleanup(func() {
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if state, ok := saved[f.Name]; ok {
				_ = f.Value.Set(state.value)
				f.Changed = state.changed
			}
		})
	})
}
