package rule

import (
	"net/http"
)

// The rule for checking the target address.
type ToAddrRule struct {
	Matcher Matcher
}

func (rule *ToAddrRule) CheckHTTP(req *http.Request) (bool, int) {
	addr := req.RequestURI
	if rule.Matcher.Match(addr) {
		return true, http.StatusOK
	}
	return false, http.StatusForbidden
}

// The function creates new instance of ToAddrRule.
//
// Available expressions: See Matcher.
func NewToAddrRule(expression any) (Rule, error) {
	matcher, err := NewMatcher(expression)
	if err == nil {
		return &ToAddrRule{Matcher: matcher}, nil
	}
	return nil, err
}
