package core

import (
	. "mal/types"
)

func fnp(ps ...MalType) (MalType, error) {
	if _, ok := ps[0].(MalFunc); ok {
		return MalBool(true), nil
	}
	if mtf, ok := ps[0].(MalTCOFunc); ok {
		return MalBool(!mtf.IsMacro), nil
	}
	return MalBool(false), nil
}

func macrop(ps ...MalType) (MalType, error) {
	if mtf, ok := ps[0].(MalTCOFunc); ok {
		return MalBool(mtf.IsMacro), nil
	}
	return MalBool(false), nil
}
