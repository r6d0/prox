package rule

import (
	"net"
	"net/http"
	"regexp"
	"strings"
)

const NOT_FUNC = "not"
const REGEX_FUNC = "regexp"
const START_FUNC = "("
const END_FUNC = ")"

// The abstraction for comparing some values.
// See RegexMatcher, EqMatcher, NotMatcher.
type Matcher interface {
	Match(string) bool
}

// The matcher for comparing a value with a regular expression.
type RegexMatcher struct {
	Regexp *regexp.Regexp
}

func (mch *RegexMatcher) Match(value string) bool {
	return mch.Regexp.Match([]byte(value))
}

// The matcher for comparing a value with a string.
type EqMatcher struct {
	Value string
}

func (mch *EqMatcher) Match(value string) bool {
	return mch.Value == value
}

// The matcher to invert the result of another matcher.
type NotMatcher struct {
	Matcher Matcher
}

func (mch *NotMatcher) Match(value string) bool {
	return !mch.Matcher.Match(value)
}

// The rule for checking the client's address.
type FromAddrRule struct {
	Matchers []Matcher
}

func (rule *FromAddrRule) CheckHTTP(req *http.Request) (bool, int) {
	addr := req.RemoteAddr
	if host, _, err := net.SplitHostPort(addr); err == nil {
		addr = host
	}

	for _, matcher := range rule.Matchers {
		if matcher.Match(addr) {
			return true, http.StatusOK
		}
	}
	return false, http.StatusForbidden
}

// The function creates new instance of FromAddrRule.
//
// Available expressions: 127.0.0.1, not(127.0.0.1), regexp(127.*), not(regexp(127.*))
func NewFromAddrRule(expressions []string) (Rule, error) {
	matchers := make([]Matcher, len(expressions))
	for index, expression := range expressions {
		hasNot := false
		if hasNot = CheckFunc(expression, NOT_FUNC); hasNot {
			expression = CutFunc(expression, NOT_FUNC)
		}

		hasRegex := false
		if hasRegex = CheckFunc(expression, REGEX_FUNC); hasRegex {
			expression = CutFunc(expression, REGEX_FUNC)
		}

		var matcher Matcher
		if hasRegex {
			regexpValue, err := regexp.Compile(expression)
			if err == nil {
				matcher = &RegexMatcher{Regexp: regexpValue}
			} else {
				return nil, err
			}
		} else {
			matcher = &EqMatcher{Value: expression}
		}

		if hasNot {
			matcher = &NotMatcher{Matcher: matcher}
		}
		matchers[index] = matcher
	}
	return &FromAddrRule{Matchers: matchers}, nil
}

func CheckFunc(value string, name string) bool {
	return strings.HasPrefix(value, name+START_FUNC) && strings.HasSuffix(value, END_FUNC)
}

func CutFunc(value string, name string) string {
	value, _ = strings.CutPrefix(value, name+START_FUNC)
	value, _ = strings.CutSuffix(value, END_FUNC)
	return value
}
