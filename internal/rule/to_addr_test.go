package rule_test

import (
	"net/http"
	"prox/internal/rule"
	"testing"
)

func TestNewToAddrRuleSuccess(t *testing.T) {
	expression := make(map[string]any)
	expression["any"] = []any{
		"127.0.0.1",
		"(128.*)",
	}

	_, err := rule.NewToAddrRule(expression)
	if err != nil {
		t.Fatalf("error should be nil, but error is %v", err)
	}
}

func TestNewToAddrRuleFailure(t *testing.T) {
	expression := make(map[string]any)
	expression["any"] = nil

	_, err := rule.NewToAddrRule(expression)
	if err == nil {
		t.Fatal("error should be not nil, but error is nil")
	}
}

func TestNewToAddrRuleCheckHTTPSuccess(t *testing.T) {
	matcher, _ := rule.NewToAddrRule("127.0.0.1")

	ok, code := matcher.CheckHTTP(&http.Request{RequestURI: "127.0.0.1"})
	if !ok {
		t.Fatal("result should be true")
	}

	if code != http.StatusOK {
		t.Fatal("code should be 200")
	}
}

func TestNewToAddrRuleCheckHTTPFailure(t *testing.T) {
	matcher, _ := rule.NewToAddrRule("127.0.0.1")

	ok, code := matcher.CheckHTTP(&http.Request{RequestURI: "128.0.0.1"})
	if ok {
		t.Fatal("result should be false")
	}

	if code != http.StatusForbidden {
		t.Fatal("code should be 403")
	}
}
