package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/salmonumbrella/eightsleep-cli/internal/client"
	"github.com/salmonumbrella/eightsleep-cli/internal/output"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show device status",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireAuthFields(); err != nil {
			return err
		}
		cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
		st, err := cl.GetStatus(context.Background())
		if err != nil {
			return err
		}
		row := map[string]any{"mode": st.CurrentState.Type, "level": st.CurrentLevel}
		fields := viper.GetStringSlice("fields")
		rows := output.FilterFields([]map[string]any{row}, fields)
		headers := fields
		if len(headers) == 0 {
			headers = []string{"mode", "level"}
		}
		return output.Print(output.Format(viper.GetString("output")), headers, rows)
	},
}
