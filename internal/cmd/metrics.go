package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/salmonumbrella/eightsleep-cli/internal/output"
)

var metricsCmd = &cobra.Command{Use: "metrics", Short: "Sleep metrics (trends, intervals)"}

var metricsTrendsCmd = &cobra.Command{Use: "trends", RunE: func(cmd *cobra.Command, args []string) error {
	cl, err := requireClient()
	if err != nil {
		return err
	}
	ctx, cancel, err := requestContext(cmd)
	if err != nil {
		return err
	}
	defer cancel()
	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")
	tz, err := resolveTimezone(viper.GetString("timezone"))
	if err != nil {
		return err
	}
	var out any
	if err := cl.Metrics().Trends(ctx, from, to, tz, &out); err != nil {
		return err
	}
	return output.Print(outputFormat(), []string{"trends"}, []map[string]any{{"trends": out}})
}}

var metricsIntervalsCmd = &cobra.Command{Use: "intervals", RunE: func(cmd *cobra.Command, args []string) error {
	cl, err := requireClient()
	if err != nil {
		return err
	}
	ctx, cancel, err := requestContext(cmd)
	if err != nil {
		return err
	}
	defer cancel()
	id, _ := cmd.Flags().GetString("id")
	var out any
	if err := cl.Metrics().Intervals(ctx, id, &out); err != nil {
		return err
	}
	return output.Print(outputFormat(), []string{"interval"}, []map[string]any{{"interval": out}})
}}

func init() {
	metricsTrendsCmd.Flags().String("from", "", "from date YYYY-MM-DD")
	metricsTrendsCmd.Flags().String("to", "", "to date YYYY-MM-DD")
	metricsIntervalsCmd.Flags().String("id", "", "session id")

	metricsCmd.AddCommand(metricsTrendsCmd, metricsIntervalsCmd)
}
