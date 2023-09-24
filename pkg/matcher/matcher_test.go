package match_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	match "github.com/tenstad/property-rules-matcher/pkg/matcher"
)

func TestNil(t *testing.T) {
	t.Parallel()

	matcher, err := match.NewMatcherBuilder[*struct{}]().
		AddRules(
			[]match.Rule[*struct{}]{
				{
					Outcome: nil,
					Conditions: map[string]match.Conditions{
						"a": {Any: []match.Condition{
							{Value: "x"},
							{Value: nil},
						}},
						"b": {Any: []match.Condition{
							{Value: nil},
						}},
					},
				},
				{
					Outcome: nil,
					Conditions: map[string]match.Conditions{
						"c": {Any: nil},
					},
				},
			},
		).Build()
	if err != nil {
		t.Fatal(err.Error())
	}

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
		{
			object: map[string]interface{}{
				"c": "x",
			},
			outcomes: nil,
		},
		{outcomes: nil, object: map[string]interface{}{"a": "x", "b": "x"}},
		{outcomes: nil, object: map[string]interface{}{"a": "x"}},
		{outcomes: nil, object: map[string]interface{}{"a": nil}},
		{outcomes: nil, object: map[string]interface{}{"b": nil}},
	} {
		outcomes, err := matcher.Match(tt.object)
		if err != nil {
			t.Fatal(err.Error())
		}
		sort := cmpopts.SortSlices(func(a, b string) bool { return a < b })
		if diff := cmp.Diff(tt.outcomes, outcomes, sort); diff != "" {
			t.Fatalf("test case %v failed with diff:\n%s", i, diff)
		}
	}
}

func TestValueTypes(t *testing.T) {
	t.Parallel()

	matcher, err := match.NewMatcherBuilder[string]().
		AddRules(
			[]match.Rule[string]{
				{
					Outcome: "alpha", Conditions: map[string]match.Conditions{
						"a": {Any: []match.Condition{{Value: nil}}},
					},
				},
				{
					Outcome: "bravo", Conditions: map[string]match.Conditions{
						"a": {Any: []match.Condition{{Value: false}}},
					},
				},
				{
					Outcome: "charlie", Conditions: map[string]match.Conditions{
						"a": {Any: []match.Condition{{Value: 0}}},
					},
				},
				{
					Outcome: "delta", Conditions: map[string]match.Conditions{
						"a": {Any: []match.Condition{{Value: 0.0}}},
					},
				},
				{
					Outcome: "foxtrot", Conditions: map[string]match.Conditions{
						"a": {Any: []match.Condition{{Value: ""}}},
					},
				},
			},
		).Build()
	if err != nil {
		t.Fatal(err.Error())
	}

	for i, tt := range []struct {
		object   map[string]interface{}
		outcomes []string
	}{
		{
			object:   map[string]interface{}{"a": nil},
			outcomes: []string{"alpha"},
		},
		{
			object:   map[string]interface{}{"a": false},
			outcomes: []string{"bravo"},
		},
		{
			object:   map[string]interface{}{"a": true},
			outcomes: nil,
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
			object:   map[string]interface{}{"a": 0.0},
			outcomes: []string{"delta"},
		},
		{
			object:   map[string]interface{}{"a": 0.1},
			outcomes: nil,
		},
		{
			object:   map[string]interface{}{"a": ""},
			outcomes: []string{"foxtrot"},
		},
		{
			object:   map[string]interface{}{"a": "x"},
			outcomes: nil,
		},
	} {
		outcomes, err := matcher.Match(tt.object)
		if err != nil {
			t.Fatal(err.Error())
		}
		sort := cmpopts.SortSlices(func(a, b string) bool { return a < b })
		if diff := cmp.Diff(tt.outcomes, outcomes, sort); diff != "" {
			t.Fatalf("test case %v failed with diff:\n%s", i, diff)
		}
	}
}

func TestInvalidRuleTypes(t *testing.T) {
	t.Parallel()

	str := ""

	for i, tt := range []struct {
		rules []match.Rule[string]
		err   string
	}{
		{
			rules: []match.Rule[string]{{
				Outcome: "alpha", Conditions: map[string]match.Conditions{
					"a": {Any: []match.Condition{{Value: int32(0)}}},
				},
			}},
			err: "invalid value type: int32",
		},
		{
			rules: []match.Rule[string]{{
				Outcome: "alpha", Conditions: map[string]match.Conditions{
					"a": {Any: []match.Condition{{Value: &str}}},
				},
			}},
			err: "invalid value type: *string",
		},
		{
			rules: []match.Rule[string]{{
				Outcome: "alpha", Conditions: map[string]match.Conditions{
					"a": {Any: []match.Condition{{Value: []string{}}}},
				},
			}},
			err: "invalid value type: []string",
		},
		{
			rules: []match.Rule[string]{{
				Outcome: "alpha", Conditions: map[string]match.Conditions{
					"a": {Any: []match.Condition{{Value: struct{}{}}}},
				},
			}},
			err: "invalid value type: struct {}",
		},
		{
			rules: []match.Rule[string]{{
				Outcome: "alpha", Conditions: map[string]match.Conditions{
					"a": {Any: []match.Condition{{Value: &struct{}{}}}},
				},
			}},
			err: "invalid value type: *struct {}",
		},
		{
			rules: []match.Rule[string]{{
				Outcome: "alpha", Conditions: map[string]match.Conditions{
					"a": {Any: []match.Condition{{Value: map[string]string{"a": "b"}}}},
				},
			}},
			err: "invalid value type: map[string]string",
		},
	} {
		_, err := match.NewMatcherBuilder[string]().
			AddRules(tt.rules).Build()
		if diff := cmp.Diff(tt.err, err.Error()); diff != "" {
			t.Fatalf("test case %v failed with diff:\n%s", i, diff)
		}
	}
}

func TestInvalidObjectTypes(t *testing.T) {
	t.Parallel()

	matcher, err := match.NewMatcherBuilder[string]().
		AddRules(
			[]match.Rule[string]{
				{
					Outcome: "alpha", Conditions: map[string]match.Conditions{
						"a": {Any: []match.Condition{{Value: nil}}},
					},
				},
			},
		).Build()
	if err != nil {
		t.Fatal(err.Error())
	}

	for i, tt := range []struct {
		object map[string]interface{}
		err    string
	}{
		{
			object: map[string]interface{}{"a": []string{}},
			err:    "invalid value type: []string",
		},
		{
			object: map[string]interface{}{"a": struct{}{}},
			err:    "invalid value type: struct {}",
		},
		{
			object: map[string]interface{}{"a": int32(0)},
			err:    "invalid value type: int32",
		},
	} {
		_, err := matcher.Match(tt.object)
		if diff := cmp.Diff(tt.err, err.Error()); diff != "" {
			t.Fatalf("test case %v failed with diff:\n%s", i, diff)
		}
	}
}

func TestManyAnyValues(t *testing.T) {
	t.Parallel()

	matcher, err := match.NewMatcherBuilder[string]().
		AddRules(
			[]match.Rule[string]{
				{
					Outcome: "alpha",
					Conditions: map[string]match.Conditions{
						"a": {Any: []match.Condition{
							{Value: "x"},
							{Value: "y"},
							{Value: "z"},
						}},
						"b": {Any: []match.Condition{
							{Value: "x"},
							{Value: "y"},
							{Value: "z"},
						}},
					},
				},
				{
					Outcome: "bravo",
					Conditions: map[string]match.Conditions{
						"a": {Any: []match.Condition{
							{Value: "x"},
							{Value: "y"},
						}},
						"b": {Any: []match.Condition{
							{Value: "x"},
							{Value: "y"},
						}},
					},
				},
			},
		).Build()
	if err != nil {
		t.Fatal(err.Error())
	}

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
		{outcomes: nil, object: map[string]interface{}{"a": "x"}},
		{outcomes: nil, object: map[string]interface{}{"a": "y"}},
		{outcomes: nil, object: map[string]interface{}{"a": "z"}},
		{outcomes: nil, object: map[string]interface{}{"b": "x"}},
		{outcomes: nil, object: map[string]interface{}{"b": "y"}},
		{outcomes: nil, object: map[string]interface{}{"b": "z"}},
	} {
		outcomes, err := matcher.Match(tt.object)
		if err != nil {
			t.Fatal(err.Error())
		}
		sort := cmpopts.SortSlices(func(a, b string) bool { return a < b })
		if diff := cmp.Diff(tt.outcomes, outcomes, sort); diff != "" {
			t.Fatalf("test case %v failed with diff:\n%s", i, diff)
		}
	}
}

func TestListOutcome(t *testing.T) {
	t.Parallel()

	matcher, err := match.NewMatcherBuilder[[]string]().
		AddRules(
			[]match.Rule[[]string]{
				{
					Outcome: []string{"alpha", "charlie"},
					Conditions: map[string]match.Conditions{
						"a": {Any: []match.Condition{
							{Value: "x"},
						}},
					},
				},
				{
					Outcome: []string{"bravo", "alpha"},
					Conditions: map[string]match.Conditions{
						"b": {Any: []match.Condition{
							{Value: "x"},
						}},
					},
				},
			},
		).Build()
	if err != nil {
		t.Fatal(err.Error())
	}

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
		outcomes, err := matcher.Match(tt.object)
		if err != nil {
			t.Fatal(err.Error())
		}
		sort := cmpopts.SortSlices(func(a, b []string) bool { return a[0] < b[0] })
		if diff := cmp.Diff(tt.outcomes, outcomes, sort); diff != "" {
			t.Fatalf("test case %v failed with diff:\n%s", i, diff)
		}
	}
}
