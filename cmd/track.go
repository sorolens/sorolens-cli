package cmd

import (
	"fmt"
	"os"

	"github.com/sorolens/sorolens-cli/internal/format"
	"github.com/spf13/cobra"
)

var trackAlias string

var trackCmd = &cobra.Command{
	Use:   "track <contract-id>",
	Short: "Register a contract for tracking",
	Args:  cobra.ExactArgs(1),
	RunE:  runTrack,
}

func init() {
	rootCmd.AddCommand(trackCmd)
	trackCmd.Flags().StringVar(&trackAlias, "alias", "", "Human-readable alias for the contract")
}

func runTrack(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	id := args[0]

	resp, err := apiClient.Track(ctx, id, trackAlias)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return err
	}

	if flagJSON {
		return format.PrintJSON(nil, resp)
	}

	fmt.Printf("Tracking %s", resp.ContractID)
	if resp.Alias != "" {
		fmt.Printf(" (alias: %s)", resp.Alias)
	}
	fmt.Println()
	if resp.Message != "" {
		fmt.Println(resp.Message)
	}
	return nil
}
