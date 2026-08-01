package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/sorolens/sorolens-cli/internal/format"
	"github.com/spf13/cobra"
)

var (
	eventsFollow bool
	eventsType   string
	eventsLimit  int
	eventsCSV    bool
)

var eventsCmd = &cobra.Command{
	Use:   "events <contract-id>",
	Short: "List events for a contract",
	Args:  cobra.ExactArgs(1),
	RunE:  runEvents,
}

func init() {
	rootCmd.AddCommand(eventsCmd)
	eventsCmd.Flags().BoolVarP(&eventsFollow, "follow", "f", false,
		"Poll for new events every 3 seconds")
	eventsCmd.Flags().StringVar(&eventsType, "type", "",
		"Filter by event type")
	// default 50 matches the Sorolens API default page size
	eventsCmd.Flags().IntVar(&eventsLimit, "limit", 50,
		"Max events to return (default matches API page size)")
	eventsCmd.Flags().BoolVar(&eventsCSV, "csv", false,
		"Output RFC 4180 CSV with header row")
}

func runEvents(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	id := args[0]

	var csvWriter *csv.Writer
	if eventsCSV {
		csvWriter = csv.NewWriter(os.Stdout)
		csvWriter.Write([]string{"time", "type", "tx_hash", "summary"})
	}

	seenCursor := ""
	for {
		resp, err := apiClient.GetEvents(ctx, id, eventsType, seenCursor, eventsLimit)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return err
		}

		if flagJSON && !eventsFollow {
			return format.PrintJSON(nil, resp)
		}

		if eventsCSV {
			for _, e := range resp.Events {
				csvWriter.Write([]string{e.Time, e.Type, e.TxHash, e.Summary})
			}
			csvWriter.Flush()
		} else {
			rows := make([][]string, len(resp.Events))
			for i, e := range resp.Events {
				rows[i] = []string{e.Time, e.Type, format.TruncateHash(e.TxHash), e.Summary}
			}
			if len(rows) > 0 {
				format.RenderTable(nil,
					[]string{"Time", "Type", "Tx Hash", "Summary"},
					rows,
				)
			}
		}

		if !eventsFollow {
			break
		}
		seenCursor = resp.NextCursor
		time.Sleep(3 * time.Second)
	}
	return nil
}
