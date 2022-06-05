package core

import (
	"fmt"
	. "mal/types"
)

func keyword(ps ...MalType) (MalType, error) {
	if k, ok := ps[0].(MalKeyword); ok {
		return k, nil
	} else if s, ok := ps[0].(MalString); ok {
		return MalKeyword(s), nil
	}
	return nil, fmt.Errorf("keyword: invalid argument, %v", ps)
}

func keywordp(ps ...MalType) (MalType, error) {
	_, ok := ps[0].(MalKeyword)
	return MalBool(ok), nil
}
