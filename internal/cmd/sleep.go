package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/salmonumbrella/eightsleep-cli/internal/output"
)

var sleepCmd = &cobra.Command{
	Use:   "sleep",
	Short: "Sleep analytics commands",
}

var sleepDayCmd = &cobra.Command{
	Use:   "day",
	Short: "Fetch sleep metrics for a date (YYYY-MM-DD)",
	Args:  cobra.NoArgs,
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
		date := viper.GetString("date")
		if date == "" {
			date = time.Now().Format("2006-01-02")
		}
		tz, err := resolveTimezone(viper.GetString("timezone"))
		if err != nil {
			return err
		}
		day, err := cl.GetSleepDay(ctx, date, tz)
		if err != nil {
			return err
		}
		rows := []map[string]any{
			{
				"date":           day.Date,
				"score":          day.Score,
				"tnt":            day.Tnt,
				"resp_rate":      day.Respiratory,
				"heart_rate":     day.HeartRate,
				"duration":       day.Duration,
				"latency_asleep": day.LatencyAsleep,
				"latency_out":    day.LatencyOut,
				"hrv_score":      day.SleepQuality.HRV.Score,
			},
		}
		fields := viper.GetStringSlice("fields")
		if err := validateFields(fields, []string{"date", "score", "duration", "latency_asleep", "latency_out", "tnt", "resp_rate", "heart_rate", "hrv_score"}); err != nil {
			return err
		}
		rows = output.FilterFields(rows, fields)
		headers := []string{"date", "score", "duration", "latency_asleep", "latency_out", "tnt", "resp_rate", "heart_rate", "hrv_score"}
		if len(fields) > 0 {
			headers = fields
		}
		return output.Print(outputFormat(), headers, rows)
	},
}

func init() {
	sleepCmd.PersistentFlags().String("date", "", "date YYYY-MM-DD (default today)")
	_ = viper.BindPFlag("date", sleepCmd.PersistentFlags().Lookup("date"))
	sleepCmd.AddCommand(sleepDayCmd)
}
