package rule_test

import (
	"prox/internal/rule"
	"testing"
)

func TestNewRulesEmpty(t *testing.T) {
	rules, err := rule.NewRules(&rule.RequestRulesConfig{})
	if err != nil {
		t.Fatalf("error should be nil, but error is %v", err)
	}

	if len(rules) > 0 {
		t.Fatalf("rules should be empty, but it is %v", rules)
	}
}
