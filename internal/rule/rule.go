package rule

import (
	"net/http"
)

// The rules configuration.
type RequestRulesConfig struct {
	From any `json:"from"`
	To   any `json:"to"`
}

// The rule for the request checking.
type Rule interface {
	// The function checks HTTP-request.
	CheckHTTP(req *http.Request) (bool, int)
}

// The function creates rules by configuration.
func NewRules(config *RequestRulesConfig) ([]Rule, error) {
	var err error
	rules := []Rule{}

	if config.From != nil {
		var rule Rule
		if rule, err = NewFromAddrRule(config.From); err == nil {
			rules = append(rules, rule)
		}
	}

	if config.To != nil {
		var rule Rule
		if rule, err = NewToAddrRule(config.To); err == nil {
			rules = append(rules, rule)
		}
	}
	return rules, err
}
