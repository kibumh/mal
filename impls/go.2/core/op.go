package core

import (
	"fmt"
	"reflect"

	. "mal/types"
)

func add(ps ...MalType) (MalType, error) {
	return ps[0].(MalInt) + ps[1].(MalInt), nil
}

func sub(ps ...MalType) (MalType, error) {
	return ps[0].(MalInt) - ps[1].(MalInt), nil
}

func mul(ps ...MalType) (MalType, error) {
	return ps[0].(MalInt) * ps[1].(MalInt), nil
}

func div(ps ...MalType) (MalType, error) {
	return ps[0].(MalInt) / ps[1].(MalInt), nil
}

func lt(ps ...MalType) (MalType, error) {
	mv1, ok := ps[0].(MalInt)
	if !ok {
		return MalBool(false), nil
	}
	mv2, ok := ps[1].(MalInt)
	if !ok {
		return MalBool(false), nil
	}
	return MalBool(mv1 < mv2), nil
}

func le(ps ...MalType) (MalType, error) {
	mv1, ok := ps[0].(MalInt)
	if !ok {
		return MalBool(false), nil
	}
	mv2, ok := ps[1].(MalInt)
	if !ok {
		return MalBool(false), nil
	}
	return MalBool(mv1 <= mv2), nil
}

func gt(ps ...MalType) (MalType, error) {
	mv1, ok := ps[0].(MalInt)
	if !ok {
		return MalBool(false), nil
	}
	mv2, ok := ps[1].(MalInt)
	if !ok {
		return MalBool(false), nil
	}
	return MalBool(mv1 > mv2), nil
}

func ge(ps ...MalType) (MalType, error) {
	mv1, ok := ps[0].(MalInt)
	if !ok {
		return MalBool(false), nil
	}
	mv2, ok := ps[1].(MalInt)
	if !ok {
		return MalBool(false), nil
	}
	return MalBool(mv1 >= mv2), nil
}
func eq(ps ...MalType) (MalType, error) {
	switch mv1 := ps[0].(type) {
	case MalNil:
		_, ok := ps[1].(MalNil)
		return MalBool(ok), nil
	case MalBool:
		mv2, ok := ps[1].(MalBool)
		if !ok {
			return MalBool(false), nil
		}
		return MalBool(mv1 == mv2), nil
	case MalInt:
		mv2, ok := ps[1].(MalInt)
		if !ok {
			return MalBool(false), nil
		}
		return MalBool(mv1 == mv2), nil
	case MalSymbol:
		mv2, ok := ps[1].(MalSymbol)
		if !ok {
			return MalBool(false), nil
		}
		return MalBool(mv1 == mv2), nil
	case MalKeyword:
		mv2, ok := ps[1].(MalKeyword)
		if !ok {
			return MalBool(false), nil
		}
		return MalBool(mv1 == mv2), nil
	case MalString:
		mv2, ok := ps[1].(MalString)
		if !ok {
			return MalBool(false), nil
		}
		return MalBool(mv1 == mv2), nil
	case MalFunc:
		return MalBool(false), nil // How to comapre two functions?
	case MalList:
		return eqListish([]MalType(mv1.List), ps[1])
	case MalVector:
		return eqListish([]MalType(mv1.Vector), ps[1])
	case MalMap:
		mv2, ok := ps[1].(MalMap)
		if !ok {
			return MalBool(false), nil
		}
		return MalBool(reflect.DeepEqual(mv1, mv2)), nil
	}
	panic(fmt.Sprintf("unreachable: eq(%v)", ps))
}

func eqListish(ms1 []MalType, mv2 MalType) (MalType, error) {
	var ms2 []MalType
	if mlist2, ok := mv2.(MalList); ok {
		ms2 = []MalType(mlist2.List)
	} else if mvec2, ok := mv2.(MalVector); ok {
		ms2 = []MalType(mvec2.Vector)
	} else {
		return MalBool(false), nil
	}

	if len(ms1) != len(ms2) {
		return MalBool(false), nil
	}
	for i := range ms1 {
		same, err := eq(ms1[i], ms2[i])
		if err != nil {
			return nil, err
		}
		if !same.(MalBool) {
			return MalBool(false), nil
		}
	}
	return MalBool(true), nil
}
