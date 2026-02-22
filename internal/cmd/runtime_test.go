package cmd

import (
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/salmonumbrella/eightsleep-cli/internal/client"
)

func TestParseTimeout(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		got, err := parseTimeout("30s")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got != 30*time.Second {
			t.Fatalf("expected 30s, got %v", got)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		if _, err := parseTimeout("nope"); err == nil {
			t.Fatalf("expected error for invalid duration")
		}
	})

	t.Run("zero", func(t *testing.T) {
		if _, err := parseTimeout("0s"); err == nil {
			t.Fatalf("expected error for zero duration")
		}
	})
}

func TestApplyClientRuntimeConfig(t *testing.T) {
	viper.Reset()
	viper.Set("timeout", "3s")
	viper.Set("retries", 4)
	c := client.New("email", "pass", "", "", "")
	if err := applyClientRuntimeConfig(c); err != nil {
		t.Fatalf("apply config: %v", err)
	}
	if c.HTTP.Timeout != 3*time.Second {
		t.Fatalf("timeout not applied: %v", c.HTTP.Timeout)
	}
	if c.MaxRetries != 4 {
		t.Fatalf("retries not applied: %d", c.MaxRetries)
	}
}

func TestValidateFields(t *testing.T) {
	if err := validateFields([]string{"a"}, []string{"a", "b"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := validateFields([]string{"c"}, []string{"a", "b"}); err == nil {
		t.Fatalf("expected error for unknown field")
	}
}

func TestRequestContext(t *testing.T) {
	viper.Reset()
	viper.Set("timeout", "1s")
	cmd := &cobra.Command{}
	ctx, cancel, err := requestContext(cmd)
	if err != nil {
		t.Fatalf("requestContext: %v", err)
	}
	defer cancel()
	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatalf("expected deadline")
	}
	if time.Until(deadline) > time.Second {
		t.Fatalf("expected deadline within 1s, got %v", time.Until(deadline))
	}
}
