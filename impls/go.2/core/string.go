package core

import (
	. "mal/types"
)

func stringp(ps ...MalType) (MalType, error) {
	_, ok := ps[0].(MalString)
	return MalBool(ok), nil
}
