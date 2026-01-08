package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/salmonumbrella/eightsleep-cli/internal/output"
)

var sleepRangeCmd = &cobra.Command{
	Use:   "range",
	Short: "Fetch sleep metrics for a date range",
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
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		if from == "" || to == "" {
			return fmt.Errorf("--from and --to are required")
		}
		layout := "2006-01-02"
		start, err := time.Parse(layout, from)
		if err != nil {
			return err
		}
		end, err := time.Parse(layout, to)
		if err != nil {
			return err
		}
		if end.Before(start) {
			return fmt.Errorf("to must be >= from")
		}
		tz, err := resolveTimezone(viper.GetString("timezone"))
		if err != nil {
			return err
		}
		rows := []map[string]any{}
		for d := start; !d.After(end); d = d.Add(24 * time.Hour) {
			day, err := cl.GetSleepDay(ctx, d.Format(layout), tz)
			if err != nil {
				return err
			}
			rows = append(rows, map[string]any{
				"date":       day.Date,
				"score":      day.Score,
				"duration":   day.Duration,
				"tnt":        day.Tnt,
				"resp_rate":  day.Respiratory,
				"heart_rate": day.HeartRate,
				"hrv_score":  day.SleepQuality.HRV.Score,
			})
		}
		fields := viper.GetStringSlice("fields")
		if err := validateFields(fields, []string{"date", "score", "duration", "tnt", "resp_rate", "heart_rate", "hrv_score"}); err != nil {
			return err
		}
		rows = output.FilterFields(rows, fields)
		headers := []string{"date", "score", "duration", "tnt", "resp_rate", "heart_rate", "hrv_score"}
		if len(fields) > 0 {
			headers = fields
		}
		return output.Print(outputFormat(), headers, rows)
	},
}

func init() {
	sleepRangeCmd.Flags().String("from", "", "start date YYYY-MM-DD")
	sleepRangeCmd.Flags().String("to", "", "end date YYYY-MM-DD")
	if sleepCmd != nil {
		sleepCmd.AddCommand(sleepRangeCmd)
	}
}
