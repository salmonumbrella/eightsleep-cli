package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/salmonumbrella/eightsleep-cli/internal/config"
	"github.com/salmonumbrella/eightsleep-cli/internal/secrets"
	"github.com/salmonumbrella/eightsleep-cli/internal/tokencache"
)

var (
	rootCmd = &cobra.Command{
		Use:   "eightsleep",
		Short: "Control your Eight Sleep Pod from the terminal",
	}
	logger = log.New(os.Stderr)
)

// Execute is the entry point for main.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("config", "", "config file (default ~/.config/eightsleep-cli/config.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose logging")
	rootCmd.PersistentFlags().String("account", "", "account name from keyring (see: eightsleep auth list)")
	rootCmd.PersistentFlags().String("email", "", "Eight Sleep account email")
	rootCmd.PersistentFlags().String("password", "", "Eight Sleep account password")
	rootCmd.PersistentFlags().String("client-id", "", "Eight Sleep client ID (optional; defaults to public app client)")
	rootCmd.PersistentFlags().String("client-secret", "", "Eight Sleep client secret (optional; defaults to public app client)")
	rootCmd.PersistentFlags().String("user-id", "", "Eight Sleep user ID")
	rootCmd.PersistentFlags().String("timezone", "local", "IANA timezone (e.g., America/New_York) or 'local'")
	rootCmd.PersistentFlags().String("output", "table", "output format: table|json|csv")
	rootCmd.PersistentFlags().StringSlice("fields", []string{}, "output fields filter")
	rootCmd.PersistentFlags().Bool("quiet", false, "suppress config load message")

	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("account", rootCmd.PersistentFlags().Lookup("account"))
	_ = viper.BindPFlag("email", rootCmd.PersistentFlags().Lookup("email"))
	_ = viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	_ = viper.BindPFlag("client_id", rootCmd.PersistentFlags().Lookup("client-id"))
	_ = viper.BindPFlag("client_secret", rootCmd.PersistentFlags().Lookup("client-secret"))
	_ = viper.BindPFlag("user_id", rootCmd.PersistentFlags().Lookup("user-id"))
	_ = viper.BindPFlag("timezone", rootCmd.PersistentFlags().Lookup("timezone"))
	_ = viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	_ = viper.BindPFlag("fields", rootCmd.PersistentFlags().Lookup("fields"))
	_ = viper.BindPFlag("config-quiet", rootCmd.PersistentFlags().Lookup("quiet"))

	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(onCmd)
	rootCmd.AddCommand(offCmd)
	rootCmd.AddCommand(tempCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(tracksCmd)
	rootCmd.AddCommand(featsCmd)
	rootCmd.AddCommand(sleepCmd)
	rootCmd.AddCommand(daemonCmd)
	rootCmd.AddCommand(alarmCmd)
	rootCmd.AddCommand(scheduleCmd)
	rootCmd.AddCommand(presenceCmd)
	rootCmd.AddCommand(tempModeCmd)
	rootCmd.AddCommand(audioCmd)
	rootCmd.AddCommand(baseCmd)
	rootCmd.AddCommand(deviceCmd)
	rootCmd.AddCommand(metricsCmd)
	rootCmd.AddCommand(autopilotCmd)
	rootCmd.AddCommand(travelCmd)
	rootCmd.AddCommand(householdCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(whoamiCmd)
	rootCmd.AddCommand(completionCmd)
}

func initConfig() {
	cfg, err := config.Load(viper.GetString("config"), viper.GetBool("config-quiet"))
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	// ensure env works on the main viper, too
	viper.SetEnvPrefix("EIGHTSLEEP")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv()
	// merge into viper defaults
	viper.SetDefault("email", cfg.Email)
	viper.SetDefault("password", cfg.Password)
	viper.SetDefault("user_id", cfg.UserID)
	viper.SetDefault("client_id", cfg.ClientID)
	viper.SetDefault("client_secret", cfg.ClientSecret)
	viper.SetDefault("timezone", cfg.Timezone)
	viper.SetDefault("output", cfg.Output)
	viper.SetDefault("fields", cfg.Fields)
	viper.SetDefault("verbose", cfg.Verbose)

	// Load credentials from keyring if account is specified
	if account := viper.GetString("account"); account != "" {
		store, err := secrets.OpenDefault()
		if err == nil {
			creds, err := store.Get(account)
			if err == nil {
				viper.Set("email", creds.Email)
				viper.Set("password", creds.Password)
			}
		}
	}

	if err := config.WarnInsecurePerms(viper.ConfigFileUsed()); err != nil {
		logger.Warn(err.Error())
	}

	if viper.GetBool("verbose") {
		log.SetLevel(log.DebugLevel)
	}
}

func requireAuthFields() error {
	// If email and password are set, we have credentials
	if viper.GetString("email") != "" && viper.GetString("password") != "" {
		return nil
	}

	// Otherwise, check if we have a cached token that can be used
	// Use the same identity construction as Client to ensure key match
	id := tokencache.Identity{
		BaseURL:  "https://client-api.8slp.net/v1", // defaultBaseURL
		ClientID: viper.GetString("client_id"),
		Email:    viper.GetString("email"),
	}
	if id.ClientID == "" {
		id.ClientID = "0894c7f33bb94800a03f1f4df13a4f38" // defaultClientID
	}

	if cached, err := tokencache.Load(id, viper.GetString("user_id")); err == nil {
		if cached.UserID != "" && viper.GetString("user_id") == "" {
			viper.Set("user_id", cached.UserID)
		}
		return nil
	}

	// No credentials and no cached token
	missing := []string{}
	if viper.GetString("email") == "" {
		missing = append(missing, "email")
	}
	if viper.GetString("password") == "" {
		missing = append(missing, "password")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required auth fields: %s", strings.Join(missing, ", "))
	}
	return nil
}
