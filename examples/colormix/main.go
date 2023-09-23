package main

import (
	"fmt"

	match "github.com/tenstad/property-rules-matcher/pkg/matcher"
)

func main() {
	matcher, err := match.NewMatcherBuilder[string]().
		AddRules(
			[]match.Rule[string]{
				{
					Outcome: "orange",
					Conditions: map[string]match.Conditions{
						"color-a": {Any: []match.Condition{
							{Value: "red"},
						}},
						"color-b": {Any: []match.Condition{
							{Value: "yellow"},
						}},
					},
				},
				{
					Outcome: "orange",
					Conditions: map[string]match.Conditions{
						"color-a": {Any: []match.Condition{
							{Value: "yellow"},
						}},
						"color-b": {Any: []match.Condition{
							{Value: "red"},
							{Value: "orange"},
						}},
					},
				},
				{
					Outcome: "dark",
					Conditions: map[string]match.Conditions{
						"color-a": {Any: []match.Condition{
							{Value: "black"},
							{Value: "eternal darkness"},
						}},
					},
				},
				{
					Outcome: "gray",
					Conditions: map[string]match.Conditions{
						"color-a": {Any: []match.Condition{
							{Value: "black"},
						}},
						"color-b": {Any: []match.Condition{
							{Value: "white"},
						}},
					},
				},
				{
					Outcome: "rainbow",
					Conditions: map[string]match.Conditions{
						"special-sauce": {Any: []match.Condition{
							{Value: "unicorn sparkles"},
							{Value: "magic"},
						}},
					},
				},
			},
		).Build()
	if err != nil {
		panic(err)
	}

	objects := []map[string]any{
		{
			"color-a": "red",
			"color-b": "yellow",
		},
		{
			"color-a": "yellow",
			"color-b": "red",
		},
		{
			"color-a": "black",
			"color-b": "red",
		},
		{
			"color-a":       "black",
			"color-b":       "red",
			"special-sauce": "unicorn sparkles",
		},
		{
			"color-a": "black",
			"color-b": "white",
		},
		{
			"color-a": "eternal darkness",
		},
	}

	for _, object := range objects {
		outcomes, err := matcher.Match(object)
		if err != nil {
			panic(err)
		}
		fmt.Println(object, "=>", outcomes)
	}
}
