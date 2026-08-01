package cmd

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/sorolens/sorolens-cli/internal/client"
	"github.com/sorolens/sorolens-cli/internal/format"
	"github.com/spf13/cobra"
)

var (
	flagAPIURL  string
	flagJSON    bool
	flagNoColor bool
	flagTimeout time.Duration
)

// apiClient is initialised in PersistentPreRunE before any command runs.
var apiClient *client.Client

var rootCmd = &cobra.Command{
	Use:   "sorolens",
	Short: "Terminal access to Soroban contract observability data",
	Long: `sorolens is a CLI for inspecting Soroban smart contracts via the
Sorolens API. It lets you view contract summaries, events, storage TTLs,
and register contracts for tracking.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if flagNoColor || os.Getenv("NO_COLOR") != "" {
			format.SetColor(false)
		}
		httpCl := &http.Client{Timeout: flagTimeout}
		apiClient = client.New(flagAPIURL, httpCl)
		return nil
	},
}

// Execute is the entry-point called from main.
func Execute() {
	ctx := context.Background()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

func init() {
	defaultURL := os.Getenv("SOROLENS_API_URL")
	if defaultURL == "" {
		defaultURL = "http://localhost:8080"
	}

	rootCmd.PersistentFlags().StringVar(&flagAPIURL, "api-url", defaultURL,
		"Sorolens API base URL (env: SOROLENS_API_URL)")
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false,
		"Output machine-readable JSON")
	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false,
		"Disable color output (also respects NO_COLOR env var)")
	rootCmd.PersistentFlags().DurationVar(&flagTimeout, "timeout", 10*time.Second,
		"API request timeout")
}
