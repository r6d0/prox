package rule

import (
	"net"
	"net/http"
)

// The rule for checking the client's address.
type FromAddrRule struct {
	Matcher Matcher
}

func (rule *FromAddrRule) CheckHTTP(req *http.Request) (bool, int) {
	addr := req.RemoteAddr
	if host, _, err := net.SplitHostPort(addr); err == nil {
		addr = host
	}

	if rule.Matcher.Match(addr) {
		return true, http.StatusOK
	}
	return false, http.StatusForbidden
}

// The function creates new instance of FromAddrRule.
//
// Available expressions: See Matcher.
func NewFromAddrRule(expression any) (Rule, error) {
	matcher, err := NewMatcher(expression)
	if err == nil {
		return &FromAddrRule{Matcher: matcher}, nil
	}
	return nil, err
}
