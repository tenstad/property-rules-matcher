package match_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	match "github.com/tenstad/property-rules-matcher/pkg/matcher"
)

func TestNil(t *testing.T) {
	t.Parallel()

	matcher := match.NewMatcherBuilder[*struct{}]().
		AddRules(
			[]match.Rule[*struct{}]{
				{
					Outcome: nil,
					Conditions: map[string]match.Conditions{
						"a": {
							Any: []match.Condition{
								{Value: "x"},
								{Value: nil},
							},
						},
						"b": {
							Any: []match.Condition{
								{Value: nil},
							},
						},
					},
				},
			},
		).Build()

	for i, tt := range []struct {
		object   map[string]interface{}
		outcomes []*struct{}
	}{
		{
			object: map[string]interface{}{
				"a": "x",
				"b": nil,
			},
			outcomes: []*struct{}{nil},
		},
		{
			object: map[string]interface{}{
				"a": nil,
				"b": nil,
			},
			outcomes: []*struct{}{nil},
		},
		{object: map[string]interface{}{"a": "x", "b": "x"}, outcomes: nil},
		{object: map[string]interface{}{"a": "x"}, outcomes: nil},
		{object: map[string]interface{}{"a": nil}, outcomes: nil},
		{object: map[string]interface{}{"b": nil}, outcomes: nil},
	} {
		outcomes := matcher.Match(tt.object)
		sort := cmpopts.SortSlices(func(a, b string) bool { return a < b })
		if diff := cmp.Diff(tt.outcomes, outcomes, sort); diff != "" {
			t.Fatalf("test case %v failed with diff:\n%s", i, diff)
		}
	}
}

func TestValueTypes(t *testing.T) {
	t.Parallel()

	matcher := match.NewMatcherBuilder[string]().
		AddRules(
			[]match.Rule[string]{
				{
					Outcome: "alpha", Conditions: map[string]match.Conditions{
						"a": {Any: []match.Condition{{Value: ""}}},
					},
				},
				{
					Outcome: "bravo", Conditions: map[string]match.Conditions{
						"a": {Any: []match.Condition{{Value: nil}}},
					},
				},
				{
					Outcome: "charlie", Conditions: map[string]match.Conditions{
						"a": {Any: []match.Condition{{Value: 0}}},
					},
				},
				{
					Outcome: "delta", Conditions: map[string]match.Conditions{
						"a": {Any: []match.Condition{{Value: false}}},
					},
				},
				{
					Outcome: "echo", Conditions: map[string]match.Conditions{
						"a": {Any: []match.Condition{{Value: struct{}{}}}},
					},
				},
				{
					Outcome: "foxtrot", Conditions: map[string]match.Conditions{
						"a": {Any: []match.Condition{{Value: &struct{}{}}}},
					},
				},
				{
					Outcome: "golf", Conditions: map[string]match.Conditions{
						"a": {Any: []match.Condition{{Value: struct{ x string }{x: "x"}}}},
					},
				},
				{
					Outcome: "hotel", Conditions: map[string]match.Conditions{
						"a": {Any: []match.Condition{{Value: &struct{ x string }{x: "x"}}}},
					},
				},
			},
		).Build()

	for i, tt := range []struct {
		object   map[string]interface{}
		outcomes []string
	}{
		{
			object:   map[string]interface{}{"a": ""},
			outcomes: []string{"alpha"},
		},
		{
			object:   map[string]interface{}{"a": "x"},
			outcomes: nil,
		},
		{
			object:   map[string]interface{}{"a": nil},
			outcomes: []string{"bravo"},
		},
		{
			object:   map[string]interface{}{"a": 0},
			outcomes: []string{"charlie"},
		},
		{
			object:   map[string]interface{}{"a": -1},
			outcomes: nil,
		},
		{
			object:   map[string]interface{}{"a": false},
			outcomes: []string{"delta"},
		},
		{
			object:   map[string]interface{}{"a": true},
			outcomes: nil,
		},
		{
			object:   map[string]interface{}{"a": struct{}{}},
			outcomes: []string{"echo"},
		},
		{
			object:   map[string]interface{}{"a": &struct{}{}},
			outcomes: []string{"foxtrot"},
		},
		{
			object:   map[string]interface{}{"a": struct{ x string }{x: "x"}},
			outcomes: []string{"golf"},
		},
		{
			object:   map[string]interface{}{"a": struct{ x string }{}},
			outcomes: nil,
		},
		{
			object:   map[string]interface{}{"a": &struct{ x string }{x: "x"}},
			outcomes: nil,
		},
	} {
		outcomes := matcher.Match(tt.object)
		sort := cmpopts.SortSlices(func(a, b string) bool { return a < b })
		if diff := cmp.Diff(tt.outcomes, outcomes, sort); diff != "" {
			t.Fatalf("test case %v failed with diff:\n%s", i, diff)
		}
	}
}

func TestManyAnyValues(t *testing.T) {
	t.Parallel()

	matcher := match.NewMatcherBuilder[string]().
		AddRules(
			[]match.Rule[string]{
				{
					Outcome: "alpha",
					Conditions: map[string]match.Conditions{
						"a": {
							Any: []match.Condition{
								{Value: "x"},
								{Value: "y"},
								{Value: "z"},
							},
						},
						"b": {
							Any: []match.Condition{
								{Value: "x"},
								{Value: "y"},
								{Value: "z"},
							},
						},
					},
				},
				{
					Outcome: "bravo",
					Conditions: map[string]match.Conditions{
						"a": {
							Any: []match.Condition{
								{Value: "x"},
								{Value: "y"},
							},
						},
						"b": {
							Any: []match.Condition{
								{Value: "x"},
								{Value: "y"},
							},
						},
					},
				},
			},
		).Build()

	for i, tt := range []struct {
		object   map[string]interface{}
		outcomes []string
	}{
		{
			object: map[string]interface{}{
				"a": "x",
				"b": "x",
			},
			outcomes: []string{"alpha", "bravo"},
		},
		{
			object: map[string]interface{}{
				"a": "z",
				"b": "z",
			},
			outcomes: []string{"alpha"},
		},
		{object: map[string]interface{}{"a": "x"}, outcomes: nil},
		{object: map[string]interface{}{"a": "y"}, outcomes: nil},
		{object: map[string]interface{}{"a": "z"}, outcomes: nil},
		{object: map[string]interface{}{"b": "x"}, outcomes: nil},
		{object: map[string]interface{}{"b": "y"}, outcomes: nil},
		{object: map[string]interface{}{"b": "z"}, outcomes: nil},
	} {
		outcomes := matcher.Match(tt.object)
		sort := cmpopts.SortSlices(func(a, b string) bool { return a < b })
		if diff := cmp.Diff(tt.outcomes, outcomes, sort); diff != "" {
			t.Fatalf("test case %v failed with diff:\n%s", i, diff)
		}
	}
}

func TestListOutcome(t *testing.T) {
	t.Parallel()

	matcher := match.NewMatcherBuilder[[]string]().
		AddRules(
			[]match.Rule[[]string]{
				{
					Outcome: []string{"alpha", "charlie"},
					Conditions: map[string]match.Conditions{
						"a": {
							Any: []match.Condition{
								{Value: "x"},
							},
						},
					},
				},
				{
					Outcome: []string{"bravo", "alpha"},
					Conditions: map[string]match.Conditions{
						"b": {
							Any: []match.Condition{
								{Value: "x"},
							},
						},
					},
				},
			},
		).Build()

	for i, tt := range []struct {
		object   map[string]interface{}
		outcomes [][]string
	}{
		{
			object: map[string]interface{}{
				"a": "x",
			},
			outcomes: [][]string{{"alpha", "charlie"}},
		},
		{
			object: map[string]interface{}{
				"a": "x",
				"b": "x",
			},
			outcomes: [][]string{{"alpha", "charlie"}, {"bravo", "alpha"}},
		},
	} {
		outcomes := matcher.Match(tt.object)
		sort := cmpopts.SortSlices(func(a, b []string) bool { return a[0] < b[0] })
		if diff := cmp.Diff(tt.outcomes, outcomes, sort); diff != "" {
			t.Fatalf("test case %v failed with diff:\n%s", i, diff)
		}
	}
}
