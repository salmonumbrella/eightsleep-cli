package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Turn pod off",
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
		if err := cl.TurnOff(ctx); err != nil {
			return err
		}
		fmt.Println("pod turned off")
		return nil
	},
}
