package rule

import (
	"errors"
	"math"
	"strconv"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/rulefmt"
)

type GroupsOutput struct {
	Groups []GroupOutput `yaml:"groups"`
}

type GroupOutput struct {
	Name     string         `yaml:"name"`
	Interval model.Duration `yaml:"interval,omitempty"`
	Rules    []rulefmt.Rule `yaml:"rules"`
}

func Format(groups *rulefmt.RuleGroups) *GroupsOutput {
	var res []GroupOutput
	for _, group := range groups.Groups {
		g := GroupOutput{
			Name:     group.Name,
			Interval: group.Interval,
		}
		var r []rulefmt.Rule
		for _, rule := range group.Rules {
			r = append(r, rulefmt.Rule{
				Record:      rule.Record.Value,
				Alert:       rule.Alert.Value,
				Expr:        rule.Expr.Value,
				For:         rule.For,
				Labels:      rule.Labels,
				Annotations: rule.Annotations,
			})
		}
		g.Rules = r
		res = append(res, g)
	}
	return &GroupsOutput{Groups: res}
}

func ruleSize(groups []*rulefmt.RuleGroups) (ruleSize int) {
	for _, ruleGroups := range groups {
		for _, group := range ruleGroups.Groups {
			ruleSize += len(group.Rules)
		}
	}
	return
}

func Compress(groups []*rulefmt.RuleGroups, maxRule int, maxGroup int) (*rulefmt.RuleGroups, error) {
	ruleSize := ruleSize(groups)
	if maxRule == 0 {
		maxRule = ruleSize
	}
	if maxGroup == 0 {
		maxGroup = int(math.Ceil(float64(ruleSize) / float64(maxRule)))
	}
	if ruleSize > maxRule*maxGroup {
		return nil, errors.New("too many rule exists. consider reducing rules")
	}

	rmap := map[model.Duration][]rulefmt.RuleNode{}
	for _, ruleGroups := range groups {
		for _, group := range ruleGroups.Groups {
			if v, ok := rmap[group.Interval]; ok {
				rmap[group.Interval] = append(v, group.Rules...)
			} else {
				rmap[group.Interval] = group.Rules
			}
		}
	}

	res := rulefmt.RuleGroups{}
	for interval, rules := range rmap {
		for i, chunk := range splitRules(rules, maxRule) {
			res.Groups = append(res.Groups, rulefmt.RuleGroup{
				Name:     "group_" + strconv.Itoa(i),
				Interval: interval,
				Rules:    chunk,
			})
		}
	}
	return &res, nil
}

func splitRules(rules []rulefmt.RuleNode, chunkSize int) [][]rulefmt.RuleNode {
	var res [][]rulefmt.RuleNode
	for i := 0; i < len(rules); i += chunkSize {
		last := i + chunkSize
		if last > len(rules) {
			last = len(rules)
		}
		res = append(res, rules[i:last])
	}
	return res
}
