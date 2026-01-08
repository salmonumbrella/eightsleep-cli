package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/salmonumbrella/eightsleep-cli/internal/client"
	"github.com/salmonumbrella/eightsleep-cli/internal/output"
)

var deviceCmd = &cobra.Command{Use: "device", Short: "Device info and priming"}

func deviceSimple(name string, fn func(ctx context.Context, cl *client.Client) (any, error)) *cobra.Command {
	return &cobra.Command{Use: name, RunE: func(cmd *cobra.Command, args []string) error {
		cl, err := requireClient()
		if err != nil {
			return err
		}
		ctx, cancel, err := requestContext(cmd)
		if err != nil {
			return err
		}
		defer cancel()
		res, err := fn(ctx, cl)
		if err != nil {
			return err
		}
		return output.Print(outputFormat(), []string{name}, []map[string]any{{name: res}})
	}}
}

func init() {
	deviceCmd.AddCommand(
		deviceSimple("info", func(ctx context.Context, cl *client.Client) (any, error) {
			return cl.Device().Info(ctx)
		}),
		deviceSimple("peripherals", func(ctx context.Context, cl *client.Client) (any, error) {
			return cl.Device().Peripherals(ctx)
		}),
		deviceSimple("online", func(ctx context.Context, cl *client.Client) (any, error) {
			return cl.Device().Online(ctx)
		}),
	)
}
