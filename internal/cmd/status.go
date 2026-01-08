package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/salmonumbrella/eightsleep-cli/internal/output"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show device status",
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
		st, err := cl.GetStatus(ctx)
		if err != nil {
			return err
		}
		row := map[string]any{"mode": st.CurrentState.Type, "level": st.CurrentLevel}
		fields := viper.GetStringSlice("fields")
		if err := validateFields(fields, []string{"mode", "level"}); err != nil {
			return err
		}
		rows := output.FilterFields([]map[string]any{row}, fields)
		headers := fields
		if len(headers) == 0 {
			headers = []string{"mode", "level"}
		}
		return output.Print(outputFormat(), headers, rows)
	},
}
