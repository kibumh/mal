package core

import (
	"fmt"
	. "mal/types"
)

func atom(ps ...MalType) (MalType, error) {
	return NewMalAtom(ps[0]), nil
}

func atomp(ps ...MalType) (MalType, error) {
	_, ok := ps[0].(*MalAtom)
	return MalBool(ok), nil
}

func deref(ps ...MalType) (MalType, error) {
	a, ok := ps[0].(*MalAtom)
	if !ok {
		return nil, fmt.Errorf("deref on non-atom, %v", ps)
	}
	return a.Value, nil
}

func reset(ps ...MalType) (MalType, error) {
	a, ok := ps[0].(*MalAtom)
	if !ok {
		return nil, fmt.Errorf("reset on non-atom, %v", ps)
	}
	a.Value = ps[1]
	return a.Value, nil
}

func swap(ps ...MalType) (MalType, error) {
	a, ok := ps[0].(*MalAtom)
	if !ok {
		return nil, fmt.Errorf("swap on non-atom, %v", ps)
	}
	var ret MalType
	var err error
	args := append([]MalType{a.Value}, ps[2:]...)
	if f, ok := ps[1].(MalFunc); ok {
		ret, err = f.Body(args...)
	} else if f2, ok := ps[1].(MalTCOFunc); ok {
		ret, err = f2.EvalFn(f2.Body, NewEnv(f2.Env, f2.Params, args))
	} else {
		return nil, fmt.Errorf("swap using non-function, %v", ps)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to apply function, %v", ps)
	}
	a.Value = ret
	return a.Value, nil
}
