package core

import (
	. "mal/types"
)

func truep(ps ...MalType) (MalType, error) {
	b, ok := ps[0].(MalBool)
	return MalBool(ok && bool(b)), nil
}

func falsep(ps ...MalType) (MalType, error) {
	b, ok := ps[0].(MalBool)
	return MalBool(ok && !bool(b)), nil
}
