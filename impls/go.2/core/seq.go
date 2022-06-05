package core

import (
	"fmt"
	. "mal/types"
)

func conj(ps ...MalType) (MalType, error) {
	if v, ok := ps[0].(MalVector); ok {
		return NewMalVector(append(v.Vector, ps[1:]...)...), nil
	}
	if l, ok := ps[0].(MalList); ok {
		ll := []MalType{}
		for i := len(ps) - 1; i > 0; i-- {
			ll = append(ll, ps[i])
		}
		return NewMalList(append(ll, l.List...)...), nil
	}
	if _, ok := ps[0].(MalNil); ok {
		return NewMalList(ps[1]), nil
	}
	return nil, fmt.Errorf("conj: not list-ish: %v", ps)
}

func convertSeq(mv MalType) ([]MalType, error) {
	if v, ok := mv.(MalVector); ok {
		return v.Vector, nil
	}
	if l, ok := mv.(MalList); ok {
		return l.List, nil
	}
	if _, ok := mv.(MalNil); ok {
		return nil, nil
	}
	return nil, fmt.Errorf("not list-ish: %v", mv)
}

func seq(ps ...MalType) (MalType, error) {
	if s, err := convertSeq(ps[0]); err == nil {
		if len(s) == 0 {
			return MalNil{}, nil
		}
		return NewMalList(s...), nil
	}
	if s, ok := ps[0].(MalString); ok {
		if len(s) == 0 {
			return MalNil{}, nil
		}
		cs := NewMalList()
		for _, c := range s {
			cs.List = append(cs.List, MalString(string(c)))
		}
		return cs, nil
	}
	return nil, fmt.Errorf("seq: not seq-able, %v", ps)
}

func sequentialp(ps ...MalType) (MalType, error) {
	_, err := convertSeq(ps[0])
	return MalBool(err == nil && ps[0] != MalNil{}), nil
}
