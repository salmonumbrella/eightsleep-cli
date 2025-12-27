package secrets

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/99designs/keyring"
)

const serviceName = "eightsleep-cli"

type Store interface {
	Keys() ([]string, error)
	Set(name string, creds Credentials) error
	Get(name string) (Credentials, error)
	Delete(name string) error
	List() ([]Credentials, error)
}

type KeyringStore struct {
	ring keyring.Keyring
}

type Credentials struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type storedCredentials struct {
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	cachedStore   Store
	cachedStoreMu sync.Mutex
)

func OpenDefault() (Store, error) {
	cachedStoreMu.Lock()
	defer cachedStoreMu.Unlock()

	if cachedStore != nil {
		return cachedStore, nil
	}

	home, _ := os.UserHomeDir()
	ring, err := keyring.Open(keyring.Config{
		ServiceName: serviceName,
		AllowedBackends: []keyring.BackendType{
			keyring.FileBackend,
		},
		FileDir:          filepath.Join(home, ".config", "eightsleep-cli", "keyring"),
		FilePasswordFunc: filePassword,
	})
	if err != nil {
		return nil, err
	}
	cachedStore = &KeyringStore{ring: ring}
	return cachedStore, nil
}

func filePassword(_ string) (string, error) {
	return serviceName + "-fallback", nil
}

func (s *KeyringStore) Keys() ([]string, error) {
	return s.ring.Keys()
}

func (s *KeyringStore) Set(name string, creds Credentials) error {
	name = normalize(name)
	if name == "" {
		return fmt.Errorf("missing account name")
	}
	if creds.Email == "" {
		return fmt.Errorf("missing email")
	}
	if creds.Password == "" {
		return fmt.Errorf("missing password")
	}
	if creds.CreatedAt.IsZero() {
		creds.CreatedAt = time.Now().UTC()
	}

	payload, err := json.Marshal(storedCredentials{
		Email:     creds.Email,
		Password:  creds.Password,
		CreatedAt: creds.CreatedAt,
	})
	if err != nil {
		return err
	}

	return s.ring.Set(keyring.Item{
		Key:  credentialKey(name),
		Data: payload,
	})
}

func (s *KeyringStore) Get(name string) (Credentials, error) {
	name = normalize(name)
	if name == "" {
		return Credentials{}, fmt.Errorf("missing account name")
	}
	item, err := s.ring.Get(credentialKey(name))
	if err != nil {
		return Credentials{}, err
	}
	var stored storedCredentials
	if err := json.Unmarshal(item.Data, &stored); err != nil {
		return Credentials{}, err
	}

	return Credentials{
		Name:      name,
		Email:     stored.Email,
		Password:  stored.Password,
		CreatedAt: stored.CreatedAt,
	}, nil
}

func (s *KeyringStore) Delete(name string) error {
	name = normalize(name)
	if name == "" {
		return fmt.Errorf("missing account name")
	}
	return s.ring.Remove(credentialKey(name))
}

func (s *KeyringStore) List() ([]Credentials, error) {
	keys, err := s.Keys()
	if err != nil {
		return nil, err
	}
	var out []Credentials
	for _, k := range keys {
		name, ok := ParseCredentialKey(k)
		if !ok {
			continue
		}
		creds, err := s.Get(name)
		if err != nil {
			return nil, err
		}
		out = append(out, creds)
	}
	return out, nil
}

func ParseCredentialKey(k string) (name string, ok bool) {
	const prefix = "account:"
	if !strings.HasPrefix(k, prefix) {
		return "", false
	}
	rest := strings.TrimPrefix(k, prefix)
	if strings.TrimSpace(rest) == "" {
		return "", false
	}
	return rest, true
}

func credentialKey(name string) string {
	return fmt.Sprintf("account:%s", name)
}

func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func ValidateAccountName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("account name cannot be empty")
	}
	if len(name) > 64 {
		return fmt.Errorf("account name too long (max 64 characters)")
	}
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_') {
			return fmt.Errorf("account name can only contain letters, numbers, hyphens, and underscores")
		}
	}
	return nil
}

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// SetOpenDefaultForTest allows tests to override the store opener
var openDefaultFunc = OpenDefault

func SetOpenDefaultForTest(fn func() (Store, error)) (restore func()) {
	cachedStoreMu.Lock()
	defer cachedStoreMu.Unlock()
	prev := openDefaultFunc
	prevStore := cachedStore
	openDefaultFunc = fn
	cachedStore = nil
	return func() {
		cachedStoreMu.Lock()
		defer cachedStoreMu.Unlock()
		openDefaultFunc = prev
		cachedStore = prevStore
	}
}
