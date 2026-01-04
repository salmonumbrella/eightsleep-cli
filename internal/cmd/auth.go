package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"

	"github.com/salmonumbrella/eightsleep-cli/internal/auth"
	"github.com/salmonumbrella/eightsleep-cli/internal/client"
	"github.com/salmonumbrella/eightsleep-cli/internal/secrets"
	"github.com/salmonumbrella/eightsleep-cli/internal/tokencache"
)

var openSecretsStore = secrets.OpenDefault

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication and account management",
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate via browser",
	Long: `Opens a browser window to authenticate with Eight Sleep.

This provides a guided setup experience with:
  - Secure credential entry in browser
  - Connection testing before saving
  - Secure credential storage in keychain

Examples:
  eightsleep auth login`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		// Handle interrupt
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-sigChan
			cancel()
		}()

		// Open keyring store for the login server
		store, err := openSecretsStore()
		if err != nil {
			return fmt.Errorf("failed to open keyring: %w", err)
		}

		fmt.Println("Opening browser for Eight Sleep authentication...")
		fmt.Println("Complete the setup in your browser, then return here.")

		server := auth.NewLoginServer(store)
		result, err := server.Start(ctx)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		if result.Error != nil {
			return result.Error
		}

		fmt.Println()
		fmt.Printf("Successfully authenticated as %s\n", result.Email)
		if result.AccountName != "" {
			fmt.Printf("Account saved as: %s\n", result.AccountName)
		}
		fmt.Printf("User ID: %s\n", result.UserID)
		fmt.Println()
		fmt.Println("You can now use eightsleep commands. Try: eightsleep status")

		return nil
	},
}

var authAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add account credentials",
	Long: `Add account credentials for API authentication.

Credentials are stored securely in your system's keychain.

Examples:
  # Add a new account (prompts for password)
  eightsleep auth add my-account --email user@example.com

  # You'll be prompted securely for your password`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := strings.TrimSpace(args[0])

		if err := secrets.ValidateAccountName(name); err != nil {
			return fmt.Errorf("invalid account name: %w", err)
		}

		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")

		if email == "" {
			fmt.Fprint(os.Stderr, "Email: ")
			reader := bufio.NewReader(os.Stdin)
			line, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read email: %w", err)
			}
			email = strings.TrimSpace(line)
		}

		email = strings.TrimSpace(email)
		if err := secrets.ValidateEmail(email); err != nil {
			return fmt.Errorf("invalid email: %w", err)
		}

		if password == "" {
			fmt.Fprint(os.Stderr, "Password: ")
			key, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				// Fallback for non-terminal
				reader := bufio.NewReader(os.Stdin)
				line, _ := reader.ReadString('\n')
				password = strings.TrimSpace(line)
			} else {
				password = string(key)
				fmt.Fprintln(os.Stderr)
			}
		}

		password = strings.TrimSpace(password)
		if password == "" {
			return fmt.Errorf("password cannot be empty")
		}

		// Test credentials before saving
		fmt.Println("Testing credentials...")
		c := client.New(email, password, "", "", "")
		if err := c.Authenticate(cmd.Context()); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		store, err := openSecretsStore()
		if err != nil {
			return fmt.Errorf("failed to open keyring: %w", err)
		}

		err = store.Set(name, secrets.Credentials{
			Email:    email,
			Password: password,
		})
		if err != nil {
			return fmt.Errorf("failed to store credentials: %w", err)
		}

		fmt.Printf("Added account: %s\n", name)
		return nil
	},
}

var authListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured accounts",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := openSecretsStore()
		if err != nil {
			return fmt.Errorf("failed to open keyring: %w", err)
		}

		creds, err := store.List()
		if err != nil {
			return fmt.Errorf("failed to list accounts: %w", err)
		}

		output := viper.GetString("output")

		if output == "json" {
			return printJSON(map[string]any{"accounts": creds})
		}

		if len(creds) == 0 {
			fmt.Println("No accounts configured")
			fmt.Println()
			fmt.Println("Add an account with:")
			fmt.Println("  eightsleep auth add <name> --email <email>")
			fmt.Println()
			fmt.Println("Or authenticate via browser:")
			fmt.Println("  eightsleep auth login")
			return nil
		}

		fmt.Printf("%-20s %-30s %s\n", "NAME", "EMAIL", "CREATED")
		for _, c := range creds {
			fmt.Printf("%-20s %-30s %s\n", c.Name, c.Email, c.CreatedAt.Format("2006-01-02"))
		}
		return nil
	},
}

var authRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove account credentials",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		store, err := openSecretsStore()
		if err != nil {
			return fmt.Errorf("failed to open keyring: %w", err)
		}

		if err := store.Delete(name); err != nil {
			return fmt.Errorf("failed to remove account: %w", err)
		}

		fmt.Printf("Removed account: %s\n", name)
		return nil
	},
}

var authTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test account credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		account, _ := cmd.Flags().GetString("account")

		// If no account specified, try to use email/password from config
		if account == "" {
			email := viper.GetString("email")
			password := viper.GetString("password")

			if email != "" && password != "" {
				fmt.Printf("Testing credentials for: %s\n", email)
				c := client.New(email, password, "", "", "")
				if err := c.Authenticate(cmd.Context()); err != nil {
					return fmt.Errorf("authentication failed: %w", err)
				}
				fmt.Println("Credentials valid")
				return nil
			}

			return fmt.Errorf("no account specified and no credentials in config")
		}

		store, err := openSecretsStore()
		if err != nil {
			return fmt.Errorf("failed to open keyring: %w", err)
		}

		creds, err := store.Get(account)
		if err != nil {
			return fmt.Errorf("account not found: %s", account)
		}

		fmt.Printf("Testing account: %s (email: %s)\n", account, creds.Email)

		c := client.New(creds.Email, creds.Password, "", "", "")
		if err := c.Authenticate(cmd.Context()); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		fmt.Println("Credentials valid")
		return nil
	},
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear cached authentication token",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.New(
			viper.GetString("email"),
			viper.GetString("password"),
			viper.GetString("user_id"),
			viper.GetString("client_id"),
			viper.GetString("client_secret"),
		)
		if err := tokencache.Clear(c.Identity()); err != nil {
			return fmt.Errorf("clear token: %w", err)
		}
		fmt.Println("Logged out (token cache cleared)")
		return nil
	},
}

func init() {
	authAddCmd.Flags().String("email", "", "Eight Sleep account email")
	authAddCmd.Flags().String("password", "", "Eight Sleep account password (omit to prompt securely)")

	authTestCmd.Flags().String("account", "", "Account name to test")

	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authAddCmd)
	authCmd.AddCommand(authListCmd)
	authCmd.AddCommand(authRemoveCmd)
	authCmd.AddCommand(authTestCmd)
	authCmd.AddCommand(authLogoutCmd)
}
