package core

import (
	"fmt"
	. "mal/types"
)

func vec(ps ...MalType) (MalType, error) {
	if pl, ok := ps[0].(MalList); ok {
		return NewMalVector(pl.List...), nil
	} else if _, ok := ps[0].(MalVector); ok {
		return ps[0], nil
	}
	return nil, fmt.Errorf("an argument is not list-ish: %v", ps)
}

func vector(ps ...MalType) (MalType, error) {
	return NewMalVector(ps...), nil
}

func vectorp(ps ...MalType) (MalType, error) {
	_, ok := ps[0].(MalVector)
	return MalBool(ok), nil
}
