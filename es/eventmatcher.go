package es

import (
	"github.com/contextgg/pkg/events"
	"github.com/contextgg/pkg/types"
)

// EventMatcher is a func that can match event to a criteria.
type EventMatcher func(events.Event) bool

// MatchAny matches any event.
func MatchAny() EventMatcher {
	return func(events.Event) bool {
		return true
	}
}

// MatchAnyInRegistry matches any event found in the registry.
func MatchAnyInRegistry() EventMatcher {
	return func(e events.Event) bool {
		_, ok := types.GetTypeData(e.Type)
		return ok
	}
}

// MatchNotLocal only events that aren't local
func MatchNotLocal() EventMatcher {
	return func(e events.Event) bool {
		if d, ok := types.GetTypeData(e.Type); ok {
			return !d.InternalType
		}
		return true
	}
}

// MatchEvent matches a specific event type, nil events never match.
func MatchEvent(t string) EventMatcher {
	return func(e events.Event) bool {
		return e.Type == t
	}
}

// MatchAnyEventOf matches if any of several matchers matches.
func MatchAnyEventOf(allTypes ...interface{}) EventMatcher {
	all := make(map[string]interface{})
	for _, t := range allTypes {
		_, name := types.GetTypeName(t)
		all[name] = t
	}

	return func(e events.Event) bool {
		if _, ok := all[e.Type]; ok {
			return true
		}
		return false
	}
}
