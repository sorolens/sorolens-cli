package cmd

import (
	"fmt"
	"os"

	"github.com/sorolens/sorolens-cli/internal/format"
	"github.com/spf13/cobra"
)

var ttlCmd = &cobra.Command{
	Use:   "ttl <contract-id>",
	Short: "Show storage entry TTLs for a contract",
	Args:  cobra.ExactArgs(1),
	RunE:  runTTL,
}

func init() {
	rootCmd.AddCommand(ttlCmd)
}

func runTTL(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	id := args[0]

	resp, err := apiClient.GetStorage(ctx, id)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return err
	}

	if flagJSON {
		return format.PrintJSON(nil, resp)
	}

	format.SectionHeader(nil, "Storage TTLs: "+id)

	headers := []string{"Key Hash", "Durability", "Live Until Ledger", "Ledgers Left"}
	rows := make([][]string, len(resp.Entries))
	for i, e := range resp.Entries {
		level := format.TTLLevelFor(e.LedgersLeft)
		ledgersLeft := format.ColorTTL(fmt.Sprintf("%d", e.LedgersLeft), level)
		rows[i] = []string{
			format.TruncateHash(e.KeyHash),
			e.Durability,
			fmt.Sprintf("%d", e.LiveUntilLedger),
			ledgersLeft,
		}
	}

	format.RenderTable(nil, headers, rows)
	return nil
}
