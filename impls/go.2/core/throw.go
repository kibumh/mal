package core

import (
	. "mal/types"
)

func throw(ps ...MalType) (MalType, error) {
	// panic(ps[0])
	return nil, MalException{ps[0]}
}
