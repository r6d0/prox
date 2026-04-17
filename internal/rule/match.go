package rule

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const notName = "!"
const anyName = "any"
const allName = "all"
const startRegexp = "("
const endRegexp = ")"

// The abstraction for comparing some values.
// See RegexMatcher, EqMatcher, NotMatcher.
type Matcher interface {
	Match(string) bool
}

// The matcher for comparing a value with a regular expression.
type RegexMatcher struct {
	Regexp *regexp.Regexp
	Invert bool
}

func (mch *RegexMatcher) Match(value string) bool {
	result := mch.Regexp.Match([]byte(value))
	if mch.Invert {
		return !result
	}
	return result
}

// The matcher for comparing a value with a string.
type EqMatcher struct {
	Value  string
	Invert bool
}

func (mch *EqMatcher) Match(value string) bool {
	if mch.Invert {
		return mch.Value != value
	}
	return mch.Value == value
}

// The logical OR.
type AnyMatcher struct {
	Matchers []Matcher
	Invert   bool
}

func (mch *AnyMatcher) Match(value string) bool {
	result := false
	for _, matcher := range mch.Matchers {
		if matcher.Match(value) {
			result = true
			break
		}
	}

	if mch.Invert {
		return !result
	}
	return result
}

// The logical AND.
type AllMatcher struct {
	Matchers []Matcher
	Invert   bool
}

func (mch *AllMatcher) Match(value string) bool {
	result := true
	for _, matcher := range mch.Matchers {
		if !matcher.Match(value) {
			result = false
			break
		}
	}

	if mch.Invert {
		return !result
	}
	return result
}

// The function creates new matcher by an expression.
func NewMatcher(expression any) (Matcher, error) {
	return parseMatcher(expression)
}

func parseMatcher(expression any) (Matcher, error) {
	var err error
	var matcher Matcher

	if value, ok := expression.(string); ok {
		invert, parsed := parseInvert(value)
		if ok, reg := parseRegexp(parsed); ok {
			var compiled *regexp.Regexp
			if compiled, err = regexp.Compile(reg); err == nil {
				matcher = &RegexMatcher{Regexp: compiled, Invert: invert}
			}
		} else {
			matcher = &EqMatcher{Value: parsed, Invert: invert}
		}
	} else if obj, ok := expression.(map[string]any); ok {
		if len(obj) == 1 {
			for key, value := range obj {
				var matchers []Matcher
				matchers, err = parseArray(value)

				invert, name := parseInvert(key)
				switch name {
				case allName:
					if err == nil {
						matcher = &AllMatcher{Matchers: matchers, Invert: invert}
					}
				case anyName:
					if err == nil {
						matcher = &AnyMatcher{Matchers: matchers, Invert: invert}
					}
				default:
					err = fmt.Errorf("the field [%s] is not supported", name)
				}
			}
		} else {
			err = errors.New("the expression must contain a single field [any] or [all]")
		}
	}
	return matcher, err
}

func parseInvert(value string) (bool, string) {
	if strings.HasPrefix(value, notName) {
		return true, value[1:]
	}
	return false, value
}

func parseRegexp(value string) (bool, string) {
	if strings.HasPrefix(value, startRegexp) && strings.HasSuffix(value, endRegexp) {
		return true, value[1 : len(value)-1]
	}
	return false, value
}

func parseArray(value any) ([]Matcher, error) {
	var err error
	var matchers []Matcher

	if array, ok := value.([]any); ok {
		matchers = make([]Matcher, len(array))

		for index, expression := range array {
			matcher, parseErr := parseMatcher(expression)
			if parseErr == nil {
				matchers[index] = matcher
			} else {
				err = errors.Join(err, parseErr)
			}
		}
	} else {
		err = errors.New("the value must be an array")
	}

	if err == nil {
		return matchers, nil
	}
	return nil, err
}
