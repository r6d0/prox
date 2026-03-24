package rule_test

import (
	"net/http"
	"prox/internal/rule"
	"testing"
)

func TestCheckFuncTrue(t *testing.T) {
	if !rule.CheckFunc("not()", rule.NOT_FUNC) {
		t.Fatal("result should be true")
	}

	if !rule.CheckFunc("not(127.0.0.1)", rule.NOT_FUNC) {
		t.Fatal("result should be true")
	}
}

func TestCheckFuncFalse(t *testing.T) {
	if rule.CheckFunc("not(", rule.NOT_FUNC) {
		t.Fatal("result should be false")
	}

	if rule.CheckFunc("not)", rule.NOT_FUNC) {
		t.Fatal("result should be false")
	}

	if rule.CheckFunc("not 127.0.0.1", rule.NOT_FUNC) {
		t.Fatal("result should be false")
	}

	if rule.CheckFunc("not", rule.NOT_FUNC) {
		t.Fatal("result should be false")
	}
}

func TestCutFunc(t *testing.T) {
	result := rule.CutFunc("not(127.0.0.1)", rule.NOT_FUNC)
	if result != "127.0.0.1" {
		t.Fatalf("result should be [127.0.0.1], but it is [%v]", result)
	}
}

func TestNewFromAddrRuleSuccess(t *testing.T) {
	newRule, err := rule.NewFromAddrRule(
		[]string{"127.0.0.1", "not(127.0.0.1)", "regexp(127.*)", "not(regexp(127.*))"},
	)

	if err != nil {
		t.Fatalf("error should be nil, but error is %v", err)
	}

	if len(newRule.(*rule.FromAddrRule).Matchers) != 4 {
		t.Fatalf("length should be 4, but it is %v", len(newRule.(*rule.FromAddrRule).Matchers))
	}
}

func TestCheckHTTPTrue(t *testing.T) {
	newRule, _ := rule.NewFromAddrRule(
		[]string{"127.0.0.1", "regexp(128.*)"},
	)

	result, _ := newRule.CheckHTTP(&http.Request{RemoteAddr: "127.0.0.1"})
	if !result {
		t.Fatal("result should be true, but it is false")
	}

	result, _ = newRule.CheckHTTP(&http.Request{RemoteAddr: "127.0.0.1:9594"})
	if !result {
		t.Fatal("result should be true, but it is false")
	}

	result, _ = newRule.CheckHTTP(&http.Request{RemoteAddr: "128.0.0.1"})
	if !result {
		t.Fatal("result should be true, but it is false")
	}

	result, _ = newRule.CheckHTTP(&http.Request{RemoteAddr: "128.0.0.1:9594"})
	if !result {
		t.Fatal("result should be true, but it is false")
	}
}

func TestCheckHTTPFalse(t *testing.T) {
	newRule, _ := rule.NewFromAddrRule(
		[]string{"not(127.0.0.1)"},
	)

	result, _ := newRule.CheckHTTP(&http.Request{RemoteAddr: "127.0.0.1"})
	if result {
		t.Fatal("result should be false, but it is true")
	}

	result, _ = newRule.CheckHTTP(&http.Request{RemoteAddr: "127.0.0.1:9594"})
	if result {
		t.Fatal("result should be false, but it is true")
	}

	newRule, _ = rule.NewFromAddrRule(
		[]string{"not(regexp(128.*))"},
	)

	result, _ = newRule.CheckHTTP(&http.Request{RemoteAddr: "128.0.0.1"})
	if result {
		t.Fatal("result should be false, but it is true")
	}

	result, _ = newRule.CheckHTTP(&http.Request{RemoteAddr: "128.0.0.1:9594"})
	if result {
		t.Fatal("result should be false, but it is true")
	}
}
