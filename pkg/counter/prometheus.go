package counter

import (
	"encoding/json"

	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	"github.com/prometheus/prometheus/promql/parser"
)

type Counter map[string]int

type Output struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func (c Counter) Add(s string) {
	c[s] += 1
}

func (c Counter) MarshalJSON() ([]byte, error) {
	var o []Output
	for k, v := range c {
		o = append(o, Output{k, v})
	}
	return json.Marshal(o)
}

func GetExprs(groups []*rulefmt.RuleGroups) []string {
	var res []string
	for _, ruleGroups := range groups {
		for _, group := range ruleGroups.Groups {
			for _, rule := range group.Rules {
				res = append(res, rule.Expr.Value)
			}
		}
	}
	return res
}

func CountExprs(exprs []string) Counter {
	c := Counter{}
	for _, expr := range exprs {
		c.Add(expr)
	}
	return c
}

func CountSelectors(exprs []string) (Counter, error) {
	c := Counter{}
	for _, expr := range exprs {
		ex, err := parser.ParseExpr(expr)
		if err != nil {
			return nil, err
		}
		parser.Inspect(ex, func(node parser.Node, nodes []parser.Node) error {
			switch n := node.(type) {
			case *parser.VectorSelector:
				c.Add(n.String())
			}
			return nil
		})
	}
	return c, nil
}

func CountMetricNames(exprs []string) (names, undefined Counter, err error) {
	names, undefined = Counter{}, Counter{}
	for _, expr := range exprs {
		ex, err := parser.ParseExpr(expr)
		if err != nil {
			return nil, nil, err
		}
		parser.Inspect(ex, func(node parser.Node, nodes []parser.Node) error {
			switch n := node.(type) {
			case *parser.VectorSelector:
				got := false
				for _, m := range n.LabelMatchers {
					if m.Name == labels.MetricName && m.Type == labels.MatchEqual {
						got = true
						names.Add(m.Value)
					}
				}
				if !got {
					undefined.Add(n.String())
				}
			}
			return nil
		})
	}
	return names, undefined, nil
}

type AnalyzePromqlOutput struct {
	Exprs         Counter `json:"exprs,omitempty"`
	Selectors     Counter `json:"selectors,omitempty"`
	Metrics       Counter `json:"metrics,omitempty"`
	NoMetricNames Counter `json:"no_metric_names,omitempty"`
}

func AnalyzePromql(exprs []string) (AnalyzePromqlOutput, error) {
	out := AnalyzePromqlOutput{
		Exprs: CountExprs(exprs),
	}
	if res, err := CountSelectors(exprs); err != nil {
		return AnalyzePromqlOutput{}, err
	} else {
		out.Selectors = res
	}
	if res1, res2, err := CountMetricNames(exprs); err != nil {
		return AnalyzePromqlOutput{}, err
	} else {
		out.Metrics = res1
		out.NoMetricNames = res2
	}
	return out, nil
}
