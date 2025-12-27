package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/salmonumbrella/eightsleep-cli/internal/auth"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Eight Sleep via browser",
	Long: `Opens your browser to authenticate with Eight Sleep.

Your credentials are stored securely in your system keychain and used
to generate an OAuth token for API access. The token is cached locally
so you don't need to re-authenticate for subsequent commands.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		fmt.Println("Opening browser for Eight Sleep authentication...")
		fmt.Println("Waiting for authentication to complete...")

		server := auth.NewLoginServer()
		result, err := server.Start(ctx)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		if result.Error != nil {
			return result.Error
		}

		fmt.Println()
		fmt.Printf("Successfully authenticated as %s\n", result.Email)
		fmt.Printf("User ID: %s\n", result.UserID)
		fmt.Println()
		fmt.Println("You can now use eightsleep commands. Try: eightsleep status")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
