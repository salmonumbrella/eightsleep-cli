package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/salmonumbrella/eightsleep-cli/internal/client"
	"github.com/salmonumbrella/eightsleep-cli/internal/output"
)

var alarmCmd = &cobra.Command{
	Use:   "alarm",
	Short: "Manage alarms",
}

var alarmListCmd = &cobra.Command{
	Use:   "list",
	Short: "List alarms",
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
		alarms, err := cl.ListAlarms(ctx)
		if err != nil {
			return err
		}
		rows := make([]map[string]any, 0, len(alarms))
		for _, a := range alarms {
			rows = append(rows, map[string]any{
				"routine_id":   a.RoutineID,
				"routine_name": a.RoutineName,
				"id":           a.ID,
				"time":         a.Time,
				"enabled":      a.Enabled,
				"days":         a.DaysOfWeek,
				"vibration":    a.Vibration,
			})
		}
		format := outputFormat()
		if format != output.FormatJSON {
			for _, row := range rows {
				if days, ok := row["days"].([]int); ok {
					row["days"] = formatDays(days)
				}
			}
		}
		fields := viper.GetStringSlice("fields")
		if err := validateFields(fields, []string{"routine_id", "routine_name", "id", "time", "enabled", "days", "vibration"}); err != nil {
			return err
		}
		rows = output.FilterFields(rows, fields)
		headers := []string{"routine_id", "routine_name", "id", "time", "enabled", "days", "vibration"}
		if len(fields) > 0 {
			headers = fields
		}
		return output.Print(format, headers, rows)
	},
}

var alarmCreateCmd = &cobra.Command{
	Use:        "create",
	Short:      "Create a one-off alarm (alias)",
	Deprecated: "use 'alarm one-off' instead",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runOneOffAlarm(cmd)
	},
}

var alarmOneOffCmd = &cobra.Command{
	Use:   "one-off",
	Short: "Create a one-off alarm",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runOneOffAlarm(cmd)
	},
}

var alarmUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an alarm",
	Args:  cobra.ExactArgs(1),
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
		patch := alarmPatch{}
		if f := viper.GetString("time"); f != "" {
			patch.Time = &f
		}
		if cmd.Flags().Changed("days") {
			patch.Days = viper.GetIntSlice("days")
			patch.DaysSet = true
		}
		if cmd.Flags().Changed("enabled") {
			val := viper.GetBool("enabled")
			patch.Enabled = &val
		}
		if cmd.Flags().Changed("no-vibration") {
			val := !viper.GetBool("no-vibration")
			patch.Vibration = &val
		}
		if patch.Empty() {
			return fmt.Errorf("no fields to update")
		}
		state, err := cl.ListRoutines(ctx)
		if err != nil {
			return err
		}
		routineID, _ := cmd.Flags().GetString("routine")
		routine, alarm, err := findRoutineAlarm(state.Routines, routineID, args[0])
		if err != nil {
			return err
		}
		applyAlarmPatch(routine, alarm, patch)
		if err := cl.UpdateRoutine(ctx, routine.ID, *routine); err != nil {
			return err
		}
		fmt.Println("updated")
		return nil
	},
}

var alarmDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Disable an alarm",
	Args:  cobra.ExactArgs(1),
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
		state, err := cl.ListRoutines(ctx)
		if err != nil {
			return err
		}
		routineID, _ := cmd.Flags().GetString("routine")
		routine, alarm, err := findRoutineAlarm(state.Routines, routineID, args[0])
		if err != nil {
			return err
		}
		alarm.Enabled = false
		alarm.DisabledIndividually = true
		if err := cl.UpdateRoutine(ctx, routine.ID, *routine); err != nil {
			return err
		}
		fmt.Println("disabled")
		return nil
	},
}

func init() {
	addOneOffFlags(alarmCreateCmd)

	addOneOffFlags(alarmOneOffCmd)

	alarmUpdateCmd.Flags().String("time", "", "HH:MM time")
	alarmUpdateCmd.Flags().IntSlice("days", nil, "Comma-separated days 0=Sun..6=Sat")
	alarmUpdateCmd.Flags().Bool("enabled", true, "Set enabled true/false")
	alarmUpdateCmd.Flags().Bool("no-vibration", false, "Disable vibration")
	alarmUpdateCmd.Flags().String("routine", "", "Routine id (optional)")
	alarmDeleteCmd.Flags().String("routine", "", "Routine id (optional)")
	// 9 minutes is the industry-standard snooze duration, dating back to the 1956 GE Snooz-Alarm.
	// Most alarm clocks and phones (including iPhone) use this default, and Eight Sleep likely follows suit.
	alarmSnoozeCmd.Flags().Int("minutes", 9, "Snooze minutes")
	_ = viper.BindPFlag("time", alarmUpdateCmd.Flags().Lookup("time"))
	_ = viper.BindPFlag("days", alarmUpdateCmd.Flags().Lookup("days"))
	_ = viper.BindPFlag("enabled", alarmUpdateCmd.Flags().Lookup("enabled"))
	_ = viper.BindPFlag("no-vibration", alarmUpdateCmd.Flags().Lookup("no-vibration"))

	// add subcommands
	alarmCmd.AddCommand(alarmListCmd, alarmCreateCmd, alarmOneOffCmd, alarmUpdateCmd, alarmDeleteCmd, alarmSnoozeCmd, alarmDismissCmd, alarmDismissAllCmd, alarmVibeCmd)
}

// snooze
var alarmSnoozeCmd = &cobra.Command{Use: "snooze <id>", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
	cl, err := requireClient()
	if err != nil {
		return err
	}
	minutes, _ := cmd.Flags().GetInt("minutes")
	if minutes <= 0 {
		return fmt.Errorf("--minutes must be > 0")
	}
	ctx, cancel, err := requestContext(cmd)
	if err != nil {
		return err
	}
	defer cancel()
	return cl.Alarms().Snooze(ctx, args[0], minutes)
}}

var alarmDismissCmd = &cobra.Command{Use: "dismiss <id>", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
	cl, err := requireClient()
	if err != nil {
		return err
	}
	ctx, cancel, err := requestContext(cmd)
	if err != nil {
		return err
	}
	defer cancel()
	return cl.Alarms().Dismiss(ctx, args[0])
}}

var alarmDismissAllCmd = &cobra.Command{Use: "dismiss-all", Short: "Dismiss next active alarm", RunE: func(cmd *cobra.Command, args []string) error {
	cl, err := requireClient()
	if err != nil {
		return err
	}
	ctx, cancel, err := requestContext(cmd)
	if err != nil {
		return err
	}
	defer cancel()
	return cl.Alarms().DismissAll(ctx)
}}

var alarmVibeCmd = &cobra.Command{Use: "vibration-test", RunE: func(cmd *cobra.Command, args []string) error {
	cl, err := requireClient()
	if err != nil {
		return err
	}
	ctx, cancel, err := requestContext(cmd)
	if err != nil {
		return err
	}
	defer cancel()
	return cl.Alarms().VibrationTest(ctx)
}}

type alarmPatch struct {
	Time      *string
	Days      []int
	DaysSet   bool
	Enabled   *bool
	Vibration *bool
}

func (p alarmPatch) Empty() bool {
	return p.Time == nil && !p.DaysSet && p.Enabled == nil && p.Vibration == nil
}

func applyAlarmPatch(routine *client.Routine, alarm *client.RoutineAlarm, patch alarmPatch) {
	if patch.Time != nil {
		alarm.EnabledSince = time.Now().UTC().Format("2006-01-02T15:04:05Z")
		if alarm.Time != "" || alarm.TimeWithOffset == nil {
			alarm.Time = *patch.Time
			alarm.TimeWithOffset = nil
		} else {
			alarm.TimeWithOffset.Time = *patch.Time
		}
	}
	if patch.DaysSet {
		routine.Days = patch.Days
	}
	if patch.Enabled != nil {
		alarm.Enabled = *patch.Enabled
		alarm.DisabledIndividually = !*patch.Enabled
	}
	if patch.Vibration != nil {
		if alarm.Settings == nil {
			alarm.Settings = &client.RoutineAlarmSettings{}
		}
		if alarm.Settings.Vibration == nil {
			alarm.Settings.Vibration = &client.RoutineAlarmVibration{}
		}
		alarm.Settings.Vibration.Enabled = *patch.Vibration
	}
}

func addOneOffFlags(cmd *cobra.Command) {
	cmd.Flags().String("time", "", "HH:MM time")
	cmd.Flags().Bool("disabled", false, "Create disabled")
	cmd.Flags().Bool("no-vibration", false, "Disable vibration")
	cmd.Flags().Bool("no-thermal", false, "Disable thermal")
	cmd.Flags().Int("thermal-level", 0, "Thermal level (-100..100)")
	cmd.Flags().Int("vibration-level", 50, "Vibration power level (0..100)")
	cmd.Flags().String("vibration-pattern", "RISE", "Vibration pattern")
	_ = viper.BindPFlag("time", cmd.Flags().Lookup("time"))
	_ = viper.BindPFlag("disabled", cmd.Flags().Lookup("disabled"))
	_ = viper.BindPFlag("no-vibration", cmd.Flags().Lookup("no-vibration"))
	_ = viper.BindPFlag("no-thermal", cmd.Flags().Lookup("no-thermal"))
	_ = viper.BindPFlag("thermal-level", cmd.Flags().Lookup("thermal-level"))
	_ = viper.BindPFlag("vibration-level", cmd.Flags().Lookup("vibration-level"))
	_ = viper.BindPFlag("vibration-pattern", cmd.Flags().Lookup("vibration-pattern"))
}

func runOneOffAlarm(cmd *cobra.Command) error {
	cl, err := requireClient()
	if err != nil {
		return err
	}
	ctx, cancel, err := requestContext(cmd)
	if err != nil {
		return err
	}
	defer cancel()
	timeStr := viper.GetString("time")
	if timeStr == "" {
		return fmt.Errorf("--time is required (HH:MM format)")
	}
	if err := validateOneOffAlarmInputs(timeStr, viper.GetInt("vibration-level"), viper.GetInt("thermal-level"), viper.GetString("vibration-pattern")); err != nil {
		return err
	}
	alarm := client.OneOffAlarm{
		Time:             timeStr,
		Enabled:          !viper.GetBool("disabled"),
		VibrationEnabled: !viper.GetBool("no-vibration"),
		VibrationLevel:   viper.GetInt("vibration-level"),
		VibrationPattern: viper.GetString("vibration-pattern"),
		ThermalEnabled:   !viper.GetBool("no-thermal"),
		ThermalLevel:     viper.GetInt("thermal-level"),
	}
	if err := cl.SetOneOffAlarm(ctx, alarm); err != nil {
		return err
	}
	fmt.Println("created one-off alarm")
	return nil
}

func validateOneOffAlarmInputs(timeStr string, vibrationLevel, thermalLevel int, vibrationPattern string) error {
	if !validAlarmTime(timeStr) {
		return fmt.Errorf("--time must be valid time (HH:MM format)")
	}
	if vibrationLevel < 0 || vibrationLevel > 100 {
		return fmt.Errorf("--vibration-level must be between 0 and 100")
	}
	if thermalLevel < -100 || thermalLevel > 100 {
		return fmt.Errorf("--thermal-level must be between -100 and 100")
	}
	if strings.TrimSpace(vibrationPattern) == "" {
		return fmt.Errorf("--vibration-pattern cannot be empty")
	}
	return nil
}

func validAlarmTime(timeStr string) bool {
	if _, err := time.Parse("15:04", timeStr); err == nil {
		return true
	}
	if _, err := time.Parse("15:04:05", timeStr); err == nil {
		return true
	}
	return false
}

func formatDays(days []int) string {
	if len(days) == 0 {
		return ""
	}
	names := []string{"sun", "mon", "tue", "wed", "thu", "fri", "sat"}
	out := make([]string, 0, len(days))
	for _, d := range days {
		if d >= 0 && d < len(names) {
			out = append(out, names[d])
		} else {
			out = append(out, fmt.Sprintf("%d", d))
		}
	}
	return strings.Join(out, ",")
}

func findRoutineAlarm(routines []client.Routine, routineID, alarmID string) (*client.Routine, *client.RoutineAlarm, error) {
	for i := range routines {
		r := &routines[i]
		if routineID != "" && r.ID != routineID {
			continue
		}
		if r.Override != nil {
			for j := range r.Override.Alarms {
				if r.Override.Alarms[j].AlarmID == alarmID {
					return r, &r.Override.Alarms[j], nil
				}
			}
		}
		for j := range r.Alarms {
			if r.Alarms[j].AlarmID == alarmID {
				return r, &r.Alarms[j], nil
			}
		}
	}
	if routineID != "" {
		return nil, nil, fmt.Errorf("alarm %s not found in routine %s", alarmID, routineID)
	}
	return nil, nil, fmt.Errorf("alarm %s not found", alarmID)
}
