package core

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	. "mal/reader"
	. "mal/types"
)

func readLine(ps ...MalType) (MalType, error) {
	fmt.Print(string(ps[0].(MalString)))
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	text = strings.TrimRight(text, "\n")
	if err != nil {
		return MalNil{}, nil
	}
	return MalString(text), nil
}

func readString(ps ...MalType) (MalType, error) {
	ms, ok := ps[0].(MalString)
	if !ok {
		return nil, fmt.Errorf("argument[0] is not a string, %v", ps)
	}
	return ReadStr(string(ms))
}
