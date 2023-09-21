package match

import (
	"maps"
	"slices"
)

type Operator int

const (
	Equal Operator = iota
)

type Condition struct {
	Operator Operator
	Value    any
}

type Conditions struct {
	Any []Condition
}

type Rule[T any] struct {
	Conditions map[string]Conditions
	Outcome    T
}

type Matcher[T any] interface {
	Match(object map[string]any) []T
}

type MatcherBuilder[T any] interface {
	AddRules(rules []Rule[T]) MatcherBuilder[T]
	Build() Matcher[T]
}

func NewMatcherBuilder[T any]() MatcherBuilder[T] {
	return &treeBuilder[T]{}
}

type treeBuilder[T any] struct {
	rules []Rule[T]
}

func (b *treeBuilder[T]) AddRules(rules []Rule[T]) MatcherBuilder[T] {
	b.rules = append(b.rules, rules...)
	return b
}

func (b *treeBuilder[T]) Build() Matcher[T] {
	var outcomes []T
	for i := len(b.rules) - 1; i >= 0; i-- {
		if len(b.rules[i].Conditions) == 0 {
			outcomes = append(outcomes, b.rules[i].Outcome)
			b.rules = slices.Delete(b.rules, i, i+1)
		}
	}

	return &node[T]{
		outcomes: outcomes,
		children: b.children(),
	}
}

func (b *treeBuilder[T]) children() map[string]map[any]Matcher[T] {
	var children map[string]map[any]Matcher[T]
	for len(b.rules) > 0 {
		prop := b.groupingProperty()

		groups := make(map[any]treeBuilder[T])
		for i := len(b.rules) - 1; i >= 0; i-- {
			if conditions, ok := b.rules[i].Conditions[prop]; ok {
				delete(b.rules[i].Conditions, prop)
				for _, condition := range conditions.Any {
					rule := Rule[T]{
						Outcome:    b.rules[i].Outcome,
						Conditions: maps.Clone(b.rules[i].Conditions),
					}

					builder := groups[condition.Value]
					builder.rules = append(builder.rules, rule)
					groups[condition.Value] = builder
				}
				b.rules = slices.Delete(b.rules, i, i+1)
			}
		}

		if children == nil {
			children = make(map[string]map[any]Matcher[T])
		}
		children[prop] = make(map[any]Matcher[T])
		for value, rules := range groups {
			children[prop][value] = rules.Build()
		}
	}
	return children
}

func (b *treeBuilder[T]) groupingProperty() string {
	for _, rule := range b.rules {
		for prop := range rule.Conditions {
			return prop
		}
	}
	panic("no rules have properties")
}

type node[T any] struct {
	outcomes []T
	children map[string]map[any]Matcher[T]
}

func (n *node[T]) Match(object map[string]any) []T {
	outcomes := n.outcomes
	for property, children := range n.children {
		if value, ok := object[property]; ok {
			if child, ok := children[value]; ok {
				outcomes = append(outcomes, child.Match(object)...)
			}
		}
	}
	return outcomes
}
