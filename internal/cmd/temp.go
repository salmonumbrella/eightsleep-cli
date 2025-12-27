package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/salmonumbrella/eightsleep-cli/internal/client"
	"github.com/salmonumbrella/eightsleep-cli/internal/daemon"
)

var tempCmd = &cobra.Command{
	Use:   "temp <value>",
	Short: "Set pod temperature (e.g., 68F, 20C, or heating level -100..100)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuthFields(); err != nil {
			return err
		}
		lvl, err := daemon.ParseTemp(args[0])
		if err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		if err := cl.SetTemperature(context.Background(), lvl); err != nil {
			return err
		}
		fmt.Printf("temperature set (level %d)\n", lvl)
		return nil
	},
}
