package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/grafana-tools/sdk"
	"github.com/kobtea/mixcli/pkg/counter"
	"github.com/spf13/cobra"
)

// dashboardCmd represents the dashboard command
var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "grafana dashboard",
}

var dashboardAnalyzeCmd = &cobra.Command{
	Use:   "analyze <file>...",
	Short: "analyze dashboard",
	RunE: func(cmd *cobra.Command, args []string) error {
		var boards []sdk.Board
		for _, file := range args {
			b, err := ioutil.ReadFile(file)
			if err != nil {
				return err
			}
			var board sdk.Board
			if err = json.Unmarshal(b, &board); err != nil {
				return err
			}
			boards = append(boards, board)
		}
		exprs := counter.GetPromqls(boards)
		res, err := counter.AnalyzePromql(exprs)
		if err != nil {
			return err
		}
		b, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
	dashboardCmd.AddCommand(dashboardAnalyzeCmd)
}
