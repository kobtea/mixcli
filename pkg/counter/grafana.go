package counter

import (
	"regexp"
	"strings"

	"github.com/grafana-tools/sdk"
	"github.com/prometheus/prometheus/promql/parser"
)

var (
	reLabelValues = regexp.MustCompile(`label_values\((.+),\s*[a-zA-Z_:][a-zA-Z0-9_:]*\)`)
	reQueryResult = regexp.MustCompile(`query_result\((.+)\)`)
)

func GetPromqls(boards []sdk.Board) []string {
	var queries []string
	for _, board := range boards {
		// template
		for _, template := range board.Templating.List {
			if template.Type != "query" {
				continue
			}
			q := template.Query.(string)
			if strings.HasPrefix(q, "label_values") {
				res := reLabelValues.FindStringSubmatch(q)
				if len(res) > 1 {
					queries = append(queries, res[1])
				}
			}
			if strings.HasPrefix(q, "query_result") {
				res := reQueryResult.FindStringSubmatch(q)
				if len(res) > 1 {
					queries = append(queries, res[1])
				}
			}
		}

		// panel
		for _, panel := range board.Panels {
			targets := panel.GetTargets()
			if targets == nil {
				continue
			}
			for _, target := range *targets {
				if _, err := parser.ParseExpr(target.Expr); err == nil {
					queries = append(queries, target.Expr)
				}
			}
		}
		// panel in row
		for _, row := range board.Rows {
			for _, panel := range row.Panels {
				targets := panel.GetTargets()
				if targets == nil {
					continue
				}
				for _, target := range *targets {
					if _, err := parser.ParseExpr(target.Expr); err == nil {
						queries = append(queries, target.Expr)
					}
				}
			}
		}
	}
	return queries
}
