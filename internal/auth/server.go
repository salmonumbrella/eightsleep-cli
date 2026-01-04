package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/salmonumbrella/eightsleep-cli/internal/client"
	"github.com/salmonumbrella/eightsleep-cli/internal/secrets"
)

// LoginResult contains the result of a browser-based login
type LoginResult struct {
	AccountName string
	Email       string
	// Password is stored temporarily to enable keyring storage in cmd/auth.go.
	// Once stored in keyring via secrets.Store.Set(), this field is no longer accessed.
	// The password is needed because Eight Sleep uses password-grant OAuth and tokens expire,
	// requiring re-authentication with the stored password from keyring.
	Password string
	UserID   string
	Error    error
}

// LoginServer handles the browser-based authentication flow
type LoginServer struct {
	result        chan LoginResult
	shutdown      chan struct{}
	pendingResult *LoginResult
	csrfToken     string
	store         secrets.Store
}

// NewLoginServer creates a new login server
func NewLoginServer(store secrets.Store) *LoginServer {
	// Generate CSRF token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		// Fall back to less random but still unique token
		tokenBytes = []byte(fmt.Sprintf("%d", time.Now().UnixNano()))
	}

	return &LoginServer{
		result:    make(chan LoginResult, 1),
		shutdown:  make(chan struct{}),
		csrfToken: hex.EncodeToString(tokenBytes),
		store:     store,
	}
}

// Start starts the login server and opens the browser
func (s *LoginServer) Start(ctx context.Context) (*LoginResult, error) {
	// Find an available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("failed to start server: %w", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	baseURL := fmt.Sprintf("http://127.0.0.1:%d", port)

	// Create HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleSetup)
	mux.HandleFunc("/submit", s.handleSubmit)
	mux.HandleFunc("/success", s.handleSuccess)
	mux.HandleFunc("/complete", s.handleComplete)
	mux.HandleFunc("/accounts", s.handleListAccounts)
	mux.HandleFunc("/set-primary", s.handleSetPrimary)
	mux.HandleFunc("/remove-account", s.handleRemoveAccount)

	server := &http.Server{
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start server in background
	go func() {
		_ = server.Serve(listener)
	}()

	// Open browser
	go func() {
		_ = openBrowser(baseURL)
	}()

	// Wait for result or context cancellation
	select {
	case result := <-s.result:
		_ = server.Shutdown(context.Background())
		return &result, nil
	case <-ctx.Done():
		_ = server.Shutdown(context.Background())
		return nil, ctx.Err()
	case <-s.shutdown:
		_ = server.Shutdown(context.Background())
		if s.pendingResult != nil {
			return s.pendingResult, nil
		}
		return nil, fmt.Errorf("login cancelled")
	}
}

// handleSetup serves the main login page
func (s *LoginServer) handleSetup(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.New("setup").Parse(setupTemplate)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"CSRFToken": s.csrfToken,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}

// handleSubmit handles the login form submission
func (s *LoginServer) handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Verify CSRF token
	if r.Header.Get("X-CSRF-Token") != s.csrfToken {
		http.Error(w, "Invalid CSRF token", http.StatusForbidden)
		return
	}

	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	// Validate account name
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": false,
			"error":   "Account name is required",
		})
		return
	}
	if err := secrets.ValidateAccountName(req.Name); err != nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Validate credentials with Eight Sleep
	c := client.New(req.Email, req.Password, "", "", "")
	if err := c.Authenticate(r.Context()); err != nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": false,
			"error":   fmt.Sprintf("Authentication failed: %v", err),
		})
		return
	}

	// Get user ID
	if err := c.EnsureUserID(r.Context()); err != nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": false,
			"error":   fmt.Sprintf("Failed to get user info: %v", err),
		})
		return
	}

	// Store credentials directly to keyring
	if s.store != nil {
		creds := secrets.Credentials{
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
		}
		if err := s.store.Set(req.Name, creds); err != nil {
			writeJSON(w, http.StatusOK, map[string]any{
				"success": false,
				"error":   fmt.Sprintf("Failed to store credentials: %v", err),
			})
			return
		}

		// Set as primary if it's the first account
		accounts, _ := s.store.List()
		if len(accounts) == 1 {
			_ = s.store.SetPrimary(req.Name)
		}
	}

	// Store pending result
	s.pendingResult = &LoginResult{
		AccountName: req.Name,
		Email:       req.Email,
		Password:    req.Password,
		UserID:      c.UserID,
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":      true,
		"account_name": req.Name,
		"email":        req.Email,
		"user_id":      c.UserID,
	})
}

// handleSuccess serves the success page
func (s *LoginServer) handleSuccess(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("success").Parse(successTemplate)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"UserEmail": r.URL.Query().Get("email"),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}

// handleComplete signals that login is done
func (s *LoginServer) handleComplete(w http.ResponseWriter, r *http.Request) {
	if s.pendingResult != nil {
		s.result <- *s.pendingResult
	}
	close(s.shutdown)
	writeJSON(w, http.StatusOK, map[string]any{"success": true})
}

// handleListAccounts returns the list of stored accounts
func (s *LoginServer) handleListAccounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.store == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"accounts": []any{},
			"primary":  "",
		})
		return
	}

	accounts, err := s.store.List()
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"accounts": []any{},
			"primary":  "",
		})
		return
	}

	primary, _ := s.store.GetPrimary()

	accountList := make([]map[string]any, 0, len(accounts))
	for _, acc := range accounts {
		accountList = append(accountList, map[string]any{
			"name":       acc.Name,
			"email":      acc.Email,
			"created_at": acc.CreatedAt.Format("2006-01-02"),
		})
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"accounts": accountList,
		"primary":  primary,
	})
}

// handleSetPrimary sets an account as the primary account
func (s *LoginServer) handleSetPrimary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Verify CSRF token
	if r.Header.Get("X-CSRF-Token") != s.csrfToken {
		http.Error(w, "Invalid CSRF token", http.StatusForbidden)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	if s.store == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": false,
			"error":   "Store not available",
		})
		return
	}

	if err := s.store.SetPrimary(req.Name); err != nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
	})
}

// handleRemoveAccount removes an account from the keyring
func (s *LoginServer) handleRemoveAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Verify CSRF token
	if r.Header.Get("X-CSRF-Token") != s.csrfToken {
		http.Error(w, "Invalid CSRF token", http.StatusForbidden)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	if s.store == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": false,
			"error":   "Store not available",
		})
		return
	}

	if err := s.store.Delete(req.Name); err != nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// If we deleted the primary account, clear it or set a new one
	primary, _ := s.store.GetPrimary()
	if primary == req.Name {
		accounts, _ := s.store.List()
		if len(accounts) > 0 {
			_ = s.store.SetPrimary(accounts[0].Name)
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
	})
}

// writeJSON writes a JSON response
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// openBrowser opens the URL in the default browser
func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}
