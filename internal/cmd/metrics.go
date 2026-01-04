package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/salmonumbrella/eightsleep-cli/internal/client"
	"github.com/salmonumbrella/eightsleep-cli/internal/output"
)

var metricsCmd = &cobra.Command{Use: "metrics", Short: "Sleep metrics (trends, intervals)"}

var metricsTrendsCmd = &cobra.Command{Use: "trends", RunE: func(cmd *cobra.Command, args []string) error {
	if err := requireAuthFields(); err != nil {
		return err
	}
	from := viper.GetString("from")
	to := viper.GetString("to")
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	var out any
	if err := cl.Metrics().Trends(context.Background(), from, to, &out); err != nil {
		return err
	}
	return output.Print(output.Format(viper.GetString("output")), []string{"trends"}, []map[string]any{{"trends": out}})
}}

var metricsIntervalsCmd = &cobra.Command{Use: "intervals", RunE: func(cmd *cobra.Command, args []string) error {
	if err := requireAuthFields(); err != nil {
		return err
	}
	id := viper.GetString("id")
	cl := client.New(viper.GetString("email"), viper.GetString("password"), viper.GetString("user_id"), viper.GetString("client_id"), viper.GetString("client_secret"))
	var out any
	if err := cl.Metrics().Intervals(context.Background(), id, &out); err != nil {
		return err
	}
	return output.Print(output.Format(viper.GetString("output")), []string{"interval"}, []map[string]any{{"interval": out}})
}}


func init() {
	metricsTrendsCmd.Flags().String("from", "", "from date YYYY-MM-DD")
	metricsTrendsCmd.Flags().String("to", "", "to date YYYY-MM-DD")
	_ = viper.BindPFlag("from", metricsTrendsCmd.Flags().Lookup("from"))
	_ = viper.BindPFlag("to", metricsTrendsCmd.Flags().Lookup("to"))
	metricsIntervalsCmd.Flags().String("id", "", "session id")
	_ = viper.BindPFlag("id", metricsIntervalsCmd.Flags().Lookup("id"))

	metricsCmd.AddCommand(metricsTrendsCmd, metricsIntervalsCmd)
}
