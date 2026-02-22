package cmd

import (
	"encoding/json"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestStatusCommandJSON(t *testing.T) {
	setupTestEnv(t)
	viper.Set("output", "json")
	out := captureStdout(t, func() {
		if err := statusCmd.RunE(statusCmd, []string{}); err != nil {
			t.Fatalf("status: %v", err)
		}
	})
	var rows []map[string]any
	if err := json.Unmarshal([]byte(out), &rows); err != nil {
		t.Fatalf("parse json: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	if rows[0]["mode"] != "on" {
		t.Fatalf("expected mode on, got %v", rows[0]["mode"])
	}
}

func TestStatusCommandInvalidFields(t *testing.T) {
	setupTestEnv(t)
	viper.Set("fields", []string{"nope"})
	if err := statusCmd.RunE(statusCmd, []string{}); err == nil {
		t.Fatalf("expected fields validation error")
	}
}

func TestOnOffTempCommands(t *testing.T) {
	setupTestEnv(t)
	if err := onCmd.RunE(onCmd, []string{}); err != nil {
		t.Fatalf("on: %v", err)
	}
	if err := offCmd.RunE(offCmd, []string{}); err != nil {
		t.Fatalf("off: %v", err)
	}
	if err := tempCmd.RunE(tempCmd, []string{"68F"}); err != nil {
		t.Fatalf("temp: %v", err)
	}
}

func TestSleepDayCommand(t *testing.T) {
	setupTestEnv(t)
	viper.Set("date", "2024-01-01")
	out := captureStdout(t, func() {
		if err := sleepDayCmd.RunE(sleepDayCmd, []string{}); err != nil {
			t.Fatalf("sleep day: %v", err)
		}
	})
	var rows []map[string]any
	if err := json.Unmarshal([]byte(out), &rows); err != nil {
		t.Fatalf("parse json: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	if rows[0]["date"] != "2024-01-01" {
		t.Fatalf("expected date 2024-01-01, got %v", rows[0]["date"])
	}
}

func TestSleepRangeCommand(t *testing.T) {
	setupTestEnv(t)
	resetFlagsOnCleanup(t, sleepRangeCmd)
	if err := sleepRangeCmd.Flags().Set("from", "2024-01-01"); err != nil {
		t.Fatalf("set from: %v", err)
	}
	if err := sleepRangeCmd.Flags().Set("to", "2024-01-02"); err != nil {
		t.Fatalf("set to: %v", err)
	}
	out := captureStdout(t, func() {
		if err := sleepRangeCmd.RunE(sleepRangeCmd, []string{}); err != nil {
			t.Fatalf("sleep range: %v", err)
		}
	})
	var rows []map[string]any
	if err := json.Unmarshal([]byte(out), &rows); err != nil {
		t.Fatalf("parse json: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
}

func TestMetricsCommands(t *testing.T) {
	setupTestEnv(t)
	resetFlagsOnCleanup(t, metricsTrendsCmd)
	resetFlagsOnCleanup(t, metricsIntervalsCmd)
	if err := metricsTrendsCmd.Flags().Set("from", "2024-01-01"); err != nil {
		t.Fatalf("set from: %v", err)
	}
	if err := metricsTrendsCmd.Flags().Set("to", "2024-01-02"); err != nil {
		t.Fatalf("set to: %v", err)
	}
	if err := metricsTrendsCmd.RunE(metricsTrendsCmd, []string{}); err != nil {
		t.Fatalf("metrics trends: %v", err)
	}
	if err := metricsIntervalsCmd.Flags().Set("id", "session-1"); err != nil {
		t.Fatalf("set id: %v", err)
	}
	if err := metricsIntervalsCmd.RunE(metricsIntervalsCmd, []string{}); err != nil {
		t.Fatalf("metrics intervals: %v", err)
	}
}

func TestDeviceCommands(t *testing.T) {
	setupTestEnv(t)
	for _, cmd := range []*cobra.Command{
		deviceCmd.Commands()[0],
		deviceCmd.Commands()[1],
		deviceCmd.Commands()[2],
	} {
		if err := cmd.RunE(cmd, []string{}); err != nil {
			t.Fatalf("device %s: %v", cmd.Use, err)
		}
	}
}

func TestAlarmListCommand(t *testing.T) {
	setupTestEnv(t)
	out := captureStdout(t, func() {
		if err := alarmListCmd.RunE(alarmListCmd, []string{}); err != nil {
			t.Fatalf("alarm list: %v", err)
		}
	})
	var rows []map[string]any
	if err := json.Unmarshal([]byte(out), &rows); err != nil {
		t.Fatalf("parse json: %v", err)
	}
	if len(rows) == 0 {
		t.Fatalf("expected alarms")
	}
}

func TestAlarmOneOffCommand(t *testing.T) {
	setupTestEnv(t)
	resetFlagsOnCleanup(t, alarmOneOffCmd)
	// Set flags directly on command (not viper) since runOneOffAlarm reads from cmd.Flags()
	if err := alarmOneOffCmd.Flags().Set("time", "07:30"); err != nil {
		t.Fatalf("set time: %v", err)
	}
	if err := alarmOneOffCmd.Flags().Set("vibration-level", "50"); err != nil {
		t.Fatalf("set vibration-level: %v", err)
	}
	if err := alarmOneOffCmd.Flags().Set("vibration-pattern", "RISE"); err != nil {
		t.Fatalf("set vibration-pattern: %v", err)
	}
	if err := alarmOneOffCmd.RunE(alarmOneOffCmd, []string{}); err != nil {
		t.Fatalf("alarm one-off: %v", err)
	}
}

func TestAlarmUpdateAndDeleteCommands(t *testing.T) {
	setupTestEnv(t)
	resetFlagsOnCleanup(t, alarmUpdateCmd)
	resetFlagsOnCleanup(t, alarmDeleteCmd)
	if err := alarmUpdateCmd.Flags().Set("enabled", "false"); err != nil {
		t.Fatalf("set enabled: %v", err)
	}
	out := captureStdout(t, func() {
		if err := alarmUpdateCmd.RunE(alarmUpdateCmd, []string{"a1"}); err != nil {
			t.Fatalf("alarm update: %v", err)
		}
	})
	if out == "" {
		t.Fatalf("expected update output")
	}
	if err := alarmDeleteCmd.RunE(alarmDeleteCmd, []string{"a1"}); err != nil {
		t.Fatalf("alarm delete: %v", err)
	}
}

func TestAlarmActionCommands(t *testing.T) {
	setupTestEnv(t)
	resetFlagsOnCleanup(t, alarmSnoozeCmd)
	if err := alarmSnoozeCmd.Flags().Set("minutes", "5"); err != nil {
		t.Fatalf("set minutes: %v", err)
	}
	if err := alarmSnoozeCmd.RunE(alarmSnoozeCmd, []string{"a1"}); err != nil {
		t.Fatalf("alarm snooze: %v", err)
	}
	if err := alarmDismissCmd.RunE(alarmDismissCmd, []string{"a1"}); err != nil {
		t.Fatalf("alarm dismiss: %v", err)
	}
	if err := alarmDismissAllCmd.RunE(alarmDismissAllCmd, []string{}); err != nil {
		t.Fatalf("alarm dismiss-all: %v", err)
	}
	if err := alarmVibeCmd.RunE(alarmVibeCmd, []string{}); err != nil {
		t.Fatalf("alarm vibration-test: %v", err)
	}
}

func TestWhoamiCommand(t *testing.T) {
	setupTestEnv(t)
	out := captureStdout(t, func() {
		if err := whoamiCmd.RunE(whoamiCmd, []string{}); err != nil {
			t.Fatalf("whoami: %v", err)
		}
	})
	if out == "" {
		t.Fatalf("expected whoami output")
	}
}

func TestVersionCommand(t *testing.T) {
	out := captureStdout(t, func() {
		versionCmd.Run(versionCmd, []string{})
	})
	if out != Version {
		t.Fatalf("expected %q, got %q", Version, out)
	}
}

func TestCompletionCommand(t *testing.T) {
	out := captureStdout(t, func() {
		if err := completionCmd.RunE(completionCmd, []string{"bash"}); err != nil {
			t.Fatalf("completion: %v", err)
		}
	})
	if out == "" {
		t.Fatalf("expected completion output")
	}
}
