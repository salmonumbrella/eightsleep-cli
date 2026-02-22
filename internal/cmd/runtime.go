package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/salmonumbrella/eightsleep-cli/internal/client"
	"github.com/salmonumbrella/eightsleep-cli/internal/output"
)

// newClient is a variable (not const) to allow tests to inject mock clients.
var newClient = client.New

func newClientFromConfig() (*client.Client, error) {
	c := newClient(
		viper.GetString("email"),
		viper.GetString("password"),
		viper.GetString("user_id"),
		viper.GetString("client_id"),
		viper.GetString("client_secret"),
	)
	if err := applyClientRuntimeConfig(c); err != nil {
		return nil, err
	}
	return c, nil
}

func requireClient() (*client.Client, error) {
	if err := requireAuthFields(); err != nil {
		return nil, err
	}
	return newClientFromConfig()
}

func outputFormat() output.Format {
	return output.Format(viper.GetString("output"))
}

func requestContext(cmd *cobra.Command) (context.Context, context.CancelFunc, error) {
	timeout, err := parseTimeout(viper.GetString("timeout"))
	if err != nil {
		return nil, func() {}, err
	}
	parent := cmd.Context()
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithTimeout(parent, timeout)
	return ctx, cancel, nil
}

func applyClientRuntimeConfig(c *client.Client) error {
	timeout, err := parseTimeout(viper.GetString("timeout"))
	if err != nil {
		return err
	}
	retries := viper.GetInt("retries")
	if retries < 0 {
		return fmt.Errorf("retries must be >= 0")
	}
	c.HTTP.Timeout = timeout
	c.MaxRetries = retries
	return nil
}

func parseTimeout(raw string) (time.Duration, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, fmt.Errorf("timeout cannot be empty")
	}
	timeout, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid timeout %q: %w", raw, err)
	}
	if timeout <= 0 {
		return 0, fmt.Errorf("timeout must be > 0")
	}
	return timeout, nil
}

func validateFields(fields []string, allowed []string) error {
	if len(fields) == 0 {
		return nil
	}
	allowedSet := make(map[string]struct{}, len(allowed))
	for _, a := range allowed {
		allowedSet[a] = struct{}{}
	}
	var invalid []string
	for _, field := range fields {
		if _, ok := allowedSet[field]; !ok {
			invalid = append(invalid, field)
		}
	}
	if len(invalid) > 0 {
		return fmt.Errorf("unknown fields: %s (allowed: %s)", strings.Join(invalid, ", "), strings.Join(allowed, ", "))
	}
	return nil
}
