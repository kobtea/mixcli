package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/kobtea/mixcli/pkg/counter"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	"github.com/spf13/cobra"
)

// ruleCmd represents the rule command
var ruleCmd = &cobra.Command{
	Use:   "rule",
	Short: "prometheus rule",
}

var ruleAnalyzeCmd = &cobra.Command{
	Use:   "analyze <file>...",
	Short: "analyze rule groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		var groups []*rulefmt.RuleGroups
		for _, file := range args {
			rg, err := rulefmt.ParseFile(file)
			if err != nil {
				return err[0]
			}
			groups = append(groups, rg)
		}
		exprs := counter.GetExprs(groups)
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
	rootCmd.AddCommand(ruleCmd)
	ruleCmd.AddCommand(ruleAnalyzeCmd)
}
