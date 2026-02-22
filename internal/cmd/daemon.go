package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/salmonumbrella/eightsleep-cli/internal/daemon"
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run schedule daemon from config file",
	RunE: func(cmd *cobra.Command, args []string) error {
		cl, err := requireClient()
		if err != nil {
			return err
		}
		cfgData, err := readConfigSchedule()
		if err != nil {
			return err
		}
		items, err := parseSchedule(cfgData)
		if err != nil {
			return err
		}
		tzName, err := resolveTimezone(viper.GetString("timezone"))
		if err != nil {
			return err
		}
		loc, err := time.LoadLocation(tzName)
		if err != nil {
			return fmt.Errorf("load timezone: %w", err)
		}
		r := daemon.Runner{
			Items:    items,
			Client:   cl,
			Timezone: loc,
			DryRun:   viper.GetBool("dry-run"),
			Sync:     viper.GetBool("sync-state"),
			PIDFile:  defaultPIDFile(viper.GetString("pid-file")),
		}
		fmt.Printf("daemon started with %d items\n", len(items))
		// Use cmd.Context() directly instead of requestContext() because daemons
		// run indefinitely and should not have a timeout applied.
		return r.Run(cmd.Context())
	},
}

func init() {
	daemonCmd.Flags().Bool("dry-run", false, "log actions without executing")
	daemonCmd.Flags().Bool("sync-state", false, "(reserved) sync device state")
	daemonCmd.Flags().String("pid-file", "", "pid file path (default ~/.config/eightsleep-cli/daemon.pid)")
	_ = viper.BindPFlag("dry-run", daemonCmd.Flags().Lookup("dry-run"))
	_ = viper.BindPFlag("sync-state", daemonCmd.Flags().Lookup("sync-state"))
	_ = viper.BindPFlag("pid-file", daemonCmd.Flags().Lookup("pid-file"))
}

func readConfigSchedule() ([]byte, error) {
	cfg := viper.ConfigFileUsed()
	if cfg == "" {
		return nil, fmt.Errorf("no config file loaded; specify --config")
	}
	return os.ReadFile(cfg)
}

func parseSchedule(data []byte) ([]daemon.ScheduleItem, error) {
	var raw struct {
		Schedule []daemon.ScheduleItem `yaml:"schedule"`
	}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	if len(raw.Schedule) == 0 {
		return nil, fmt.Errorf("no schedule entries found")
	}
	return raw.Schedule, nil
}

func defaultPIDFile(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "eightsleep", "daemon.pid")
}
