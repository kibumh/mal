package core

import (
	. "mal/types"
)

func nilp(ps ...MalType) (MalType, error) {
	_, ok := ps[0].(MalNil)
	return MalBool(ok), nil
}
