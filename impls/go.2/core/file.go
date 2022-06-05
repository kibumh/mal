package core

import (
	"fmt"
	"os"

	. "mal/types"
)

func slurp(ps ...MalType) (MalType, error) {
	path, ok := ps[0].(MalString)
	if !ok {
		return nil, fmt.Errorf("argument[0] is not a string, %v", ps)
	}
	buf, err := os.ReadFile(string(path))
	if err != nil {
		return nil, err
	}
	return MalString(string(buf)), nil
}
