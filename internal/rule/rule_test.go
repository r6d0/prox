package rule_test

import (
	"fmt"
	"prox/internal"
	"prox/internal/rule"
	"testing"
)

func TestNewRulesEmpty(t *testing.T) {
	config, _ := internal.NewJsonConfig("D:\\andrey\\workspace\\prox\\default_prox_config.json")

	fmt.Println(config)

	rules, err := rule.NewRules(&rule.RequestRulesConfig{})
	if err != nil {
		t.Fatalf("error should be nil, but error is %v", err)
	}

	if len(rules) > 0 {
		t.Fatalf("rules should be empty, but it is %v", rules)
	}
}
