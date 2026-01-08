package cmd

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/99designs/keyring"
	"github.com/spf13/viper"

	"github.com/salmonumbrella/eightsleep-cli/internal/client"
	"github.com/salmonumbrella/eightsleep-cli/internal/secrets"
	"github.com/salmonumbrella/eightsleep-cli/internal/tokencache"
)

type fakeStore struct {
	creds map[string]secrets.Credentials
}

func (f *fakeStore) Keys() ([]string, error) { return nil, nil }
func (f *fakeStore) Set(name string, creds secrets.Credentials) error {
	if f.creds == nil {
		f.creds = map[string]secrets.Credentials{}
	}
	creds.Name = name
	f.creds[name] = creds
	return nil
}

func (f *fakeStore) Get(name string) (secrets.Credentials, error) {
	creds, ok := f.creds[name]
	if !ok {
		return secrets.Credentials{}, keyring.ErrKeyNotFound
	}
	return creds, nil
}

func (f *fakeStore) Delete(name string) error {
	delete(f.creds, name)
	return nil
}

func (f *fakeStore) List() ([]secrets.Credentials, error) {
	out := make([]secrets.Credentials, 0, len(f.creds))
	for _, c := range f.creds {
		out = append(out, c)
	}
	return out, nil
}
func (f *fakeStore) SetPrimary(name string) error { return nil }
func (f *fakeStore) GetPrimary() (string, error)  { return "", keyring.ErrKeyNotFound }

func setupAuthEnv(t *testing.T) {
	t.Helper()
	useTempKeyring(t)
	viper.Reset()
	viper.Set("timeout", "2s")
	viper.Set("output", "json")
}

func TestAuthListJSON(t *testing.T) {
	setupAuthEnv(t)
	store := &fakeStore{
		creds: map[string]secrets.Credentials{
			"acct": {Name: "acct", Email: "user@example.com", CreatedAt: time.Now().UTC()},
		},
	}
	prev := openSecretsStore
	openSecretsStore = func() (secrets.Store, error) { return store, nil }
	t.Cleanup(func() { openSecretsStore = prev })

	out := captureStdout(t, func() {
		if err := authListCmd.RunE(authListCmd, []string{}); err != nil {
			t.Fatalf("auth list: %v", err)
		}
	})
	var payload map[string]any
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("parse json: %v", err)
	}
	if _, ok := payload["accounts"]; !ok {
		t.Fatalf("expected accounts key")
	}
}

func TestAuthRemove(t *testing.T) {
	setupAuthEnv(t)
	store := &fakeStore{
		creds: map[string]secrets.Credentials{
			"acct": {Name: "acct", Email: "user@example.com", CreatedAt: time.Now().UTC()},
		},
	}
	prev := openSecretsStore
	openSecretsStore = func() (secrets.Store, error) { return store, nil }
	t.Cleanup(func() { openSecretsStore = prev })

	if err := authRemoveCmd.RunE(authRemoveCmd, []string{"acct"}); err != nil {
		t.Fatalf("auth remove: %v", err)
	}
	if _, ok := store.creds["acct"]; ok {
		t.Fatalf("expected account removed")
	}
}

func TestAuthTestWithAccount(t *testing.T) {
	setupAuthEnv(t)
	store := &fakeStore{
		creds: map[string]secrets.Credentials{
			"acct": {Name: "acct", Email: "user@example.com", Password: "pass", CreatedAt: time.Now().UTC()},
		},
	}
	prev := openSecretsStore
	openSecretsStore = func() (secrets.Store, error) { return store, nil }
	t.Cleanup(func() { openSecretsStore = prev })

	cl := client.New("user@example.com", "pass", "", "", "")
	if err := tokencache.Save(cl.Identity(), "tok", time.Now().Add(time.Hour), "uid-123"); err != nil {
		t.Fatalf("save token: %v", err)
	}
	if err := authTestCmd.Flags().Set("account", "acct"); err != nil {
		t.Fatalf("set account: %v", err)
	}
	if err := authTestCmd.RunE(authTestCmd, []string{}); err != nil {
		t.Fatalf("auth test: %v", err)
	}
}

func TestAuthLogout(t *testing.T) {
	setupAuthEnv(t)
	viper.Set("email", "user@example.com")
	viper.Set("password", "pass")
	cl := client.New("user@example.com", "pass", "", "", "")
	if err := tokencache.Save(cl.Identity(), "tok", time.Now().Add(time.Hour), "uid-123"); err != nil {
		t.Fatalf("save token: %v", err)
	}
	if err := authLogoutCmd.RunE(authLogoutCmd, []string{}); err != nil {
		t.Fatalf("auth logout: %v", err)
	}
}
