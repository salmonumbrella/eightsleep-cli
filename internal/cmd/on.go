package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var onCmd = &cobra.Command{
	Use:   "on",
	Short: "Turn pod on",
	RunE: func(cmd *cobra.Command, args []string) error {
		cl, err := requireClient()
		if err != nil {
			return err
		}
		ctx, cancel, err := requestContext(cmd)
		if err != nil {
			return err
		}
		defer cancel()
		if err := cl.TurnOn(ctx); err != nil {
			return err
		}
		fmt.Println("pod turned on")
		return nil
	},
}
