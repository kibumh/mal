package core

import (
	. "mal/types"
)

func numberp(ps ...MalType) (MalType, error) {
	_, ok := ps[0].(MalInt)
	// TODO: MalFloat
	return MalBool(ok), nil
}
