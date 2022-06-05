package core

import (
	"fmt"
	"strings"

	. "mal/printer"
	. "mal/types"
)

func prn(ps ...MalType) (MalType, error) {
	s, err := prStr(ps...)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(s.(MalString)))
	return MalNil{}, nil
}

func println(ps ...MalType) (MalType, error) {
	var ss []string
	for _, mv := range ps {
		ss = append(ss, PrintStr(mv, false))
	}
	fmt.Println(strings.Join(ss, " "))
	return MalNil{}, nil
}

func prStr(ps ...MalType) (MalType, error) {
	var ss []string
	for _, mv := range ps {
		ss = append(ss, PrintStr(mv, true))
	}
	return MalString(strings.Join(ss, " ")), nil
}

func str(ps ...MalType) (MalType, error) {
	var ss []string
	for _, mv := range ps {
		ss = append(ss, PrintStr(mv, false))
	}
	return MalString(strings.Join(ss, "")), nil
}
