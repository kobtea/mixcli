package cmd

import (
	"fmt"

	"github.com/kobtea/mixcli/pkg/rule"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// alertCmd represents the alert command
var alertCmd = &cobra.Command{
	Use:   "alert",
	Short: "prometheus alert",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("alert called")
	},
}

var alertCompressCmd = &cobra.Command{
	Use:   "compress <file>...",
	Short: "compress alert groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		maxRule, err := cmd.Flags().GetInt("max-rule")
		if err != nil {
			return err
		}
		maxGroup, err := cmd.Flags().GetInt("max-group")
		if err != nil {
			return err
		}

		var groups []*rulefmt.RuleGroups
		for _, file := range args {
			rg, err := rulefmt.ParseFile(file)
			if err != nil {
				return err[0]
			}
			groups = append(groups, rg)
		}
		res, err := rule.Compress(groups, maxRule, maxGroup)
		if err != nil {
			return err
		}
		b, err := yaml.Marshal(rule.Format(res))
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(alertCmd)

	alertCmd.AddCommand(alertCompressCmd)
	alertCompressCmd.Flags().Int("max-rule", 0, "max rules per group")
	alertCompressCmd.Flags().Int("max-group", 0, "max groups")
}
