package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/salmonumbrella/eightsleep-cli/internal/daemon"
)

var tempCmd = &cobra.Command{
	Use:   "temp <value>",
	Short: "Set pod temperature (e.g., 68F, 20C, or heating level -100..100)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cl, err := requireClient()
		if err != nil {
			return err
		}
		lvl, err := daemon.ParseTemp(args[0])
		if err != nil {
			return err
		}
		ctx, cancel, err := requestContext(cmd)
		if err != nil {
			return err
		}
		defer cancel()
		if err := cl.SetTemperature(ctx, lvl); err != nil {
			return err
		}
		fmt.Printf("temperature set (level %d)\n", lvl)
		return nil
	},
}
