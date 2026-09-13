package cmd

import (
	"fmt"
	"os"

	"github.com/sorolens/sorolens-cli/internal/client"
	"github.com/sorolens/sorolens-cli/internal/format"
	"github.com/spf13/cobra"
)

var watchdogCmd = &cobra.Command{
	Use:   "watchdog",
	Short: "Query the on-chain contract health watchdog",
	Long: `The watchdog command group inspects the Sorolens on-chain watchdog:
monitored contracts, aggregate stats, health-check history, and alerts.`,
}

var watchdogStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show aggregate watchdog stats",
	Args:  cobra.NoArgs,
	RunE:  runWatchdogStats,
}

var watchdogListCmd = &cobra.Command{
	Use:   "list",
	Short: "List monitored contracts",
	Args:  cobra.NoArgs,
	RunE:  runWatchdogList,
}

var watchdogStatusCmd = &cobra.Command{
	Use:   "status <contract-id>",
	Short: "Show status detail for one monitored contract",
	Args:  cobra.ExactArgs(1),
	RunE:  runWatchdogStatus,
}

var (
	watchdogAlertsSeverity string
	watchdogAlertsLimit    int
)

var watchdogAlertsCmd = &cobra.Command{
	Use:   "alerts [contract-id]",
	Short: "List watchdog alerts (all, or for one contract)",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runWatchdogAlerts,
}

var watchdogHistoryLimit int

var watchdogHistoryCmd = &cobra.Command{
	Use:   "history <contract-id>",
	Short: "Show health-check history for a contract",
	Args:  cobra.ExactArgs(1),
	RunE:  runWatchdogHistory,
}

func init() {
	rootCmd.AddCommand(watchdogCmd)
	watchdogCmd.AddCommand(watchdogStatsCmd)
	watchdogCmd.AddCommand(watchdogListCmd)
	watchdogCmd.AddCommand(watchdogStatusCmd)
	watchdogCmd.AddCommand(watchdogAlertsCmd)
	watchdogCmd.AddCommand(watchdogHistoryCmd)

	watchdogAlertsCmd.Flags().StringVar(&watchdogAlertsSeverity, "severity", "",
		"Filter by severity (Info, Warning, Critical)")
	watchdogAlertsCmd.Flags().IntVar(&watchdogAlertsLimit, "limit", 20,
		"Max alerts to return")

	watchdogHistoryCmd.Flags().IntVar(&watchdogHistoryLimit, "limit", 20,
		"Max health-check entries to return")
}

// truncateContractID shortens a Soroban contract/account ID for table display:
// first 4 + "..." + last 3.
func truncateContractID(id string) string {
	if len(id) <= 10 {
		return id
	}
	return id[:4] + "..." + id[len(id)-3:]
}

func runWatchdogStats(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	stats, err := apiClient.GetWatchdogStats(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return err
	}
	if flagJSON {
		return format.PrintJSON(nil, stats)
	}
	fmt.Printf("Monitored Contracts: %d\n", stats.TotalMonitored)
	fmt.Printf("Healthy:             %d\n", stats.Healthy)
	fmt.Printf("Degraded:            %d\n", stats.Degraded)
	fmt.Printf("Unresponsive:        %d\n", stats.Unresponsive)
	fmt.Printf("Total Alerts:        %d\n", stats.TotalAlerts)
	fmt.Printf("Critical Alerts:     %d\n", stats.CriticalAlerts)
	return nil
}

func runWatchdogList(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	resp, err := apiClient.GetWatchdogContracts(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return err
	}
	if flagJSON {
		return format.PrintJSON(nil, resp)
	}
	headers := []string{"Contract ID", "Name", "Status", "Last Check", "Interval"}
	rows := make([][]string, len(resp.Contracts))
	for i, c := range resp.Contracts {
		rows[i] = []string{
			truncateContractID(c.ContractID),
			c.Name,
			c.Status,
			c.LastCheck,
			fmt.Sprintf("%ds", c.CheckInterval),
		}
	}
	format.RenderTable(nil, headers, rows)
	return nil
}

func runWatchdogStatus(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	id := args[0]
	c, err := apiClient.GetWatchdogContract(ctx, id)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return err
	}
	if flagJSON {
		return format.PrintJSON(nil, c)
	}
	format.SectionHeader(nil, "Watchdog Contract: "+c.ContractID)
	format.KeyValueTable(nil, [][2]string{
		{"Contract ID", c.ContractID},
		{"Name", c.Name},
		{"Owner", c.Owner},
		{"Status", c.Status},
		{"Last Check", c.LastCheck},
		{"Check Interval", fmt.Sprintf("%ds", c.CheckInterval)},
		{"Registered", c.RegisteredAt},
		{"Updated", c.UpdatedAt},
	})
	return nil
}

func runWatchdogAlerts(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	var resp *client.WatchdogAlertsResponse
	var err error
	if len(args) == 1 {
		resp, err = apiClient.GetWatchdogContractAlerts(ctx, args[0], watchdogAlertsSeverity, watchdogAlertsLimit)
	} else {
		resp, err = apiClient.GetWatchdogAlerts(ctx, watchdogAlertsSeverity)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return err
	}
	if flagJSON {
		return format.PrintJSON(nil, resp)
	}
	alerts := resp.Alerts
	if len(args) != 1 && watchdogAlertsLimit > 0 && watchdogAlertsLimit < len(alerts) {
		alerts = alerts[:watchdogAlertsLimit]
	}
	headers := []string{"Severity", "Contract", "Message", "Timestamp"}
	rows := make([][]string, len(alerts))
	for i, a := range alerts {
		rows[i] = []string{a.Severity, truncateContractID(a.ContractID), a.Message, a.Timestamp}
	}
	format.RenderTable(nil, headers, rows)
	return nil
}

func runWatchdogHistory(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	id := args[0]
	resp, err := apiClient.GetWatchdogHealth(ctx, id, watchdogHistoryLimit)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return err
	}
	if flagJSON {
		return format.PrintJSON(nil, resp)
	}
	headers := []string{"Status", "Metadata", "Ledger", "Tx Hash", "Timestamp"}
	rows := make([][]string, len(resp.HealthChecks))
	for i, h := range resp.HealthChecks {
		rows[i] = []string{
			h.Status,
			h.Metadata,
			fmt.Sprintf("%d", h.Ledger),
			format.TruncateHash(h.TxHash),
			h.Timestamp,
		}
	}
	format.RenderTable(nil, headers, rows)
	return nil
}
