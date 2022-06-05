package core

import (
	"fmt"
	. "mal/types"
)

func symbol(ps ...MalType) (MalType, error) {
	s, ok := ps[0].(MalString)
	if !ok {
		return nil, fmt.Errorf("symbol: not a string, %v", ps)
	}
	return MalSymbol(s), nil
}

func symbolp(ps ...MalType) (MalType, error) {
	_, ok := ps[0].(MalSymbol)
	return MalBool(ok), nil
}
