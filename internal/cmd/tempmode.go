package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/salmonumbrella/eightsleep-cli/internal/output"
)

var tempModeCmd = &cobra.Command{Use: "tempmode", Short: "Temperature modes (nap, hot flash, events)"}

// --- nap ---

var napCmd = &cobra.Command{Use: "nap", Short: "Nap mode controls"}

var napOnCmd = &cobra.Command{Use: "on", Short: "Activate nap mode", RunE: func(cmd *cobra.Command, args []string) error {
	cl, err := requireClient()
	if err != nil {
		return err
	}
	ctx, cancel, err := requestContext(cmd)
	if err != nil {
		return err
	}
	defer cancel()
	if err := cl.TempModes().NapActivate(ctx); err != nil {
		return err
	}
	fmt.Println("nap mode activated")
	return nil
}}

var napOffCmd = &cobra.Command{Use: "off", Short: "Deactivate nap mode", RunE: func(cmd *cobra.Command, args []string) error {
	cl, err := requireClient()
	if err != nil {
		return err
	}
	ctx, cancel, err := requestContext(cmd)
	if err != nil {
		return err
	}
	defer cancel()
	if err := cl.TempModes().NapDeactivate(ctx); err != nil {
		return err
	}
	fmt.Println("nap mode deactivated")
	return nil
}}

var napExtendCmd = &cobra.Command{Use: "extend", Short: "Extend nap mode", RunE: func(cmd *cobra.Command, args []string) error {
	cl, err := requireClient()
	if err != nil {
		return err
	}
	ctx, cancel, err := requestContext(cmd)
	if err != nil {
		return err
	}
	defer cancel()
	if err := cl.TempModes().NapExtend(ctx); err != nil {
		return err
	}
	fmt.Println("nap mode extended")
	return nil
}}

var napStatusCmd = &cobra.Command{Use: "status", Short: "Show nap mode status", RunE: func(cmd *cobra.Command, args []string) error {
	cl, err := requireClient()
	if err != nil {
		return err
	}
	ctx, cancel, err := requestContext(cmd)
	if err != nil {
		return err
	}
	defer cancel()
	var out any
	if err := cl.TempModes().NapStatus(ctx, &out); err != nil {
		return err
	}
	return printDynamic("nap_status", out)
}}

// --- hot flash ---

var hotflashCmd = &cobra.Command{Use: "hotflash", Short: "Hot flash mode controls"}

var hotflashOnCmd = &cobra.Command{Use: "on", Short: "Activate hot flash mode", RunE: func(cmd *cobra.Command, args []string) error {
	cl, err := requireClient()
	if err != nil {
		return err
	}
	ctx, cancel, err := requestContext(cmd)
	if err != nil {
		return err
	}
	defer cancel()
	if err := cl.TempModes().HotFlashActivate(ctx); err != nil {
		return err
	}
	fmt.Println("hot flash mode activated")
	return nil
}}

var hotflashOffCmd = &cobra.Command{Use: "off", Short: "Deactivate hot flash mode", RunE: func(cmd *cobra.Command, args []string) error {
	cl, err := requireClient()
	if err != nil {
		return err
	}
	ctx, cancel, err := requestContext(cmd)
	if err != nil {
		return err
	}
	defer cancel()
	if err := cl.TempModes().HotFlashDeactivate(ctx); err != nil {
		return err
	}
	fmt.Println("hot flash mode deactivated")
	return nil
}}

var hotflashStatusCmd = &cobra.Command{Use: "status", Short: "Show hot flash mode status", RunE: func(cmd *cobra.Command, args []string) error {
	cl, err := requireClient()
	if err != nil {
		return err
	}
	ctx, cancel, err := requestContext(cmd)
	if err != nil {
		return err
	}
	defer cancel()
	var out any
	if err := cl.TempModes().HotFlashStatus(ctx, &out); err != nil {
		return err
	}
	return printDynamic("hotflash_status", out)
}}

// --- events ---

var tempEventsCmd = &cobra.Command{Use: "events", Short: "Show temperature events", RunE: func(cmd *cobra.Command, args []string) error {
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
	var out any
	if err := cl.TempModes().TempEvents(ctx, from, to, &out); err != nil {
		return err
	}
	return printDynamic("events", out)
}}

// printDynamic renders an arbitrary API response using output.Print.
// For maps it extracts sorted keys as headers; for slices of maps it uses
// keys from the first element; otherwise it wraps the value under fallback.
func printDynamic(fallback string, v any) error {
	format := outputFormat()
	fields := viper.GetStringSlice("fields")

	switch val := v.(type) {
	case map[string]any:
		headers := mapKeys(val)
		rows := []map[string]any{val}
		rows = output.FilterFields(rows, fields)
		if len(fields) > 0 {
			headers = fields
		}
		return output.Print(format, headers, rows)
	case []any:
		if len(val) == 0 {
			return output.Print(format, []string{fallback}, nil)
		}
		if first, ok := val[0].(map[string]any); ok {
			headers := mapKeys(first)
			rows := make([]map[string]any, 0, len(val))
			for _, item := range val {
				if m, ok := item.(map[string]any); ok {
					rows = append(rows, m)
				}
			}
			rows = output.FilterFields(rows, fields)
			if len(fields) > 0 {
				headers = fields
			}
			return output.Print(format, headers, rows)
		}
		return output.Print(format, []string{fallback}, []map[string]any{{fallback: val}})
	default:
		return output.Print(format, []string{fallback}, []map[string]any{{fallback: val}})
	}
}

// mapKeys returns sorted keys from a map.
func mapKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func init() {
	tempEventsCmd.Flags().String("from", "", "from date (ISO 8601)")
	tempEventsCmd.Flags().String("to", "", "to date (ISO 8601)")

	napCmd.AddCommand(napOnCmd, napOffCmd, napExtendCmd, napStatusCmd)
	hotflashCmd.AddCommand(hotflashOnCmd, hotflashOffCmd, hotflashStatusCmd)
	tempModeCmd.AddCommand(napCmd, hotflashCmd, tempEventsCmd)
}
