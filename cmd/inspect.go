package cmd

import (
	"fmt"
	"os"

	"github.com/sorolens/sorolens-cli/internal/format"
	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect <contract-id>",
	Short: "Show a summary of a tracked contract",
	Args:  cobra.ExactArgs(1),
	RunE:  runInspect,
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}

func runInspect(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	id := args[0]

	contract, err := apiClient.GetContract(ctx, id)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return err
	}

	if flagJSON {
		return format.PrintJSON(nil, contract)
	}

	format.SectionHeader(nil, "Contract: "+id)
	format.KeyValueTable(nil, [][2]string{
		{"Alias", contract.Alias},
		{"Network", contract.Network},
		{"First Seen", contract.FirstSeen},
		{"Event Count", fmt.Sprintf("%d", contract.EventCount)},
		{"Invocation Count", fmt.Sprintf("%d", contract.InvocationCount)},
		{"Success Rate", fmt.Sprintf("%.1f%%", contract.SuccessRate*100)},
		{"Storage Entries", fmt.Sprintf("%d", contract.StorageEntries)},
		{"Nearest TTL Expiry (ledger)", fmt.Sprintf("%d", contract.NearestTTLExpiry)},
	})
	return nil
}
