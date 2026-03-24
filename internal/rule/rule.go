package rule

import (
	"net/http"
)

// The rules configuration.
type RequestRulesConfig struct {
	FromAddr []string `json:"fromAddr"`
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

	fromAddr := config.FromAddr
	if len(fromAddr) > 0 {
		var rule Rule
		if rule, err = NewFromAddrRule(fromAddr); err == nil {
			rules = append(rules, rule)
		}
	}
	return rules, err
}
