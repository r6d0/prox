package rule_test

import (
	"prox/internal/rule"
	"testing"
)

func TestStringMatcher(t *testing.T) {
	matcher, err := rule.NewMatcher("127.0.0.1")
	if err != nil {
		t.Fatalf("error should be nil, but error is [%v]", err)
	}

	if !matcher.Match("127.0.0.1") {
		t.Fatalf("result should be true, but it is false")
	}

	if matcher.Match("128.0.0.1") {
		t.Fatalf("result should be false, but it is true")
	}
}

func TestInvertStringMatcher(t *testing.T) {
	matcher, err := rule.NewMatcher("!127.0.0.1")
	if err != nil {
		t.Fatalf("error should be nil, but error is [%v]", err)
	}

	if matcher.Match("127.0.0.1") {
		t.Fatalf("result should be false, but it is true")
	}

	if !matcher.Match("128.0.0.1") {
		t.Fatalf("result should be true, but it is false")
	}
}

func TestRegexpMatcher(t *testing.T) {
	matcher, err := rule.NewMatcher("(127.*)")
	if err != nil {
		t.Fatalf("error should be nil, but error is [%v]", err)
	}

	if !matcher.Match("127.0.0.1") {
		t.Fatalf("result should be true, but it is false")
	}

	if matcher.Match("128.0.0.1") {
		t.Fatalf("result should be false, but it is true")
	}
}

func TestInvertRegexpMatcher(t *testing.T) {
	matcher, err := rule.NewMatcher("!(127.*)")
	if err != nil {
		t.Fatalf("error should be nil, but error is [%v]", err)
	}

	if matcher.Match("127.0.0.1") {
		t.Fatalf("result should be false, but it is true")
	}

	if !matcher.Match("128.0.0.1") {
		t.Fatalf("result should be true, but it is false")
	}
}

func TestAnyMatcher(t *testing.T) {
	childExpression := make(map[string]any)
	childExpression["any"] = []any{
		"127.0.0.1",
		"(128.*)",
	}

	expression := make(map[string]any)
	expression["any"] = []any{
		"127.0.0.1",
		"(128.*)",
		childExpression,
	}

	matcher, err := rule.NewMatcher(expression)
	if err != nil {
		t.Fatalf("error should be nil, but error is [%v]", err)
	}

	if !matcher.Match("127.0.0.1") {
		t.Fatalf("result should be true, but it is false")
	}

	if !matcher.Match("128.0.0.1") {
		t.Fatalf("result should be true, but it is false")
	}
}

func TestInvertAnyMatcher(t *testing.T) {
	childExpression := make(map[string]any)
	childExpression["any"] = []any{
		"127.0.0.1",
		"(128.*)",
	}

	expression := make(map[string]any)
	expression["!any"] = []any{
		"127.0.0.1",
		"(128.*)",
		childExpression,
	}

	matcher, err := rule.NewMatcher(expression)
	if err != nil {
		t.Fatalf("error should be nil, but error is [%v]", err)
	}

	if matcher.Match("127.0.0.1") {
		t.Fatalf("result should be false, but it is true")
	}

	if matcher.Match("128.0.0.1") {
		t.Fatalf("result should be false, but it is true")
	}
}

func TestAllMatcher(t *testing.T) {
	childExpression := make(map[string]any)
	childExpression["all"] = []any{
		"127.0.0.1",
		"(127.*)",
	}

	expression := make(map[string]any)
	expression["all"] = []any{
		"127.0.0.1",
		"(127.*)",
		childExpression,
	}

	matcher, err := rule.NewMatcher(expression)
	if err != nil {
		t.Fatalf("error should be nil, but error is [%v]", err)
	}

	if !matcher.Match("127.0.0.1") {
		t.Fatalf("result should be true, but it is false")
	}
}

func TestInvertAllMatcher(t *testing.T) {
	expression := make(map[string]any)
	expression["!all"] = []any{
		"127.0.0.1",
		"(127.*)",
	}

	matcher, err := rule.NewMatcher(expression)
	if err != nil {
		t.Fatalf("error should be nil, but error is [%v]", err)
	}

	if matcher.Match("127.0.0.1") {
		t.Fatalf("result should be false, but it is true")
	}

	if !matcher.Match("129.0.0.1") {
		t.Fatalf("result should be true, but it is false")
	}
}

func TestNewMatcherErrTooManyFields(t *testing.T) {
	expression := make(map[string]any)
	expression["all"] = []any{"127.0.0.1"}
	expression["any"] = []any{"127.0.0.1"}

	_, err := rule.NewMatcher(expression)
	if err == nil {
		t.Fatalf("error should be not nil, but error is nil")
	}
}

func TestNewMatcherErrUnsupportedField(t *testing.T) {
	expression := make(map[string]any)
	expression["unsupported"] = []any{"127.0.0.1"}

	_, err := rule.NewMatcher(expression)
	if err == nil {
		t.Fatalf("error should be not nil, but error is nil")
	}
}

func TestNewMatcherErrValueMustBeArray(t *testing.T) {
	expression := make(map[string]any)
	expression["any"] = "127.0.0.1"

	_, err := rule.NewMatcher(expression)
	if err == nil {
		t.Fatalf("error should be not nil, but error is nil")
	}
}
