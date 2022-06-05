package core

import (
	"time"

	. "mal/types"
)

func timeMs(ps ...MalType) (MalType, error) {
	return MalInt(time.Now().UnixMilli()), nil
}
