package core

import (
	"errors"
	"fmt"
	. "mal/types"
)

func listf(ps ...MalType) (MalType, error) {
	return NewMalList(ps...), nil
}

func listp(ps ...MalType) (MalType, error) {
	_, ok := ps[0].(MalList)
	return MalBool(ok), nil
}

func emptyp(ps ...MalType) (MalType, error) {
	if ml, ok := ps[0].(MalList); ok {
		return MalBool(len(ml.List) == 0), nil
	}
	if mvec, ok := ps[0].(MalVector); ok {
		return MalBool(len(mvec.Vector) == 0), nil
	}
	return nil, errors.New("not a list")
}

func count(ps ...MalType) (MalType, error) {
	if _, ok := ps[0].(MalNil); ok {
		return MalInt(0), nil
	}
	if ml, ok := ps[0].(MalList); ok {
		return MalInt(len(ml.List)), nil
	}
	if mvec, ok := ps[0].(MalVector); ok {
		return MalInt(len(mvec.Vector)), nil
	}
	return nil, errors.New("not a list")
}

func cons(ps ...MalType) (MalType, error) {
	if len(ps) != 2 {
		return nil, fmt.Errorf("the number of arguments is not 2: %v", ps)
	}

	var listish []MalType
	if mlist, ok := ps[1].(MalList); ok {
		listish = mlist.List
	} else if mvec, ok := ps[1].(MalVector); ok {
		listish = mvec.Vector
	} else if _, ok = ps[1].(MalNil); ok {
		listish = nil
	} else {
		return nil, fmt.Errorf("second argument is not list-ish: %v", ps)
	}
	return NewMalList(append([]MalType{ps[0]}, listish...)...), nil
}

func concat(ps ...MalType) (MalType, error) {
	ml := NewMalList()
	for _, p := range ps {
		if pl, ok := p.(MalList); ok {
			ml.List = append(ml.List, pl.List...)
		} else if pv, ok := p.(MalVector); ok {
			ml.List = append(ml.List, pv.Vector...)
		} else {
			return nil, fmt.Errorf("one of arguments is not a list: %v", ps)
		}
	}
	return ml, nil
}

func nth(ps ...MalType) (MalType, error) {
	seq, err := convertSeq(ps[0])
	if err != nil {
		return nil, err
	}
	i, ok := ps[1].(MalInt)
	if !ok {
		return nil, fmt.Errorf("index is not an integer: list(%v), index(%v)", ps[0], ps[1])
	}
	if int(i) >= len(seq) {
		return nil, fmt.Errorf("index is out-of-range: list(%v), index(%v)", ps[0], ps[1])
	}
	return seq[int(i)], nil
}

func first(ps ...MalType) (MalType, error) {
	seq, err := convertSeq(ps[0])
	if err != nil {
		return nil, err
	}
	if len(seq) == 0 {
		return MalNil{}, nil
	}
	return seq[0], nil
}

func rest(ps ...MalType) (MalType, error) {
	seq, err := convertSeq(ps[0])
	if err != nil {
		return nil, err
	}
	if len(seq) == 0 {
		return NewMalList(), nil
	}
	return NewMalList(seq[1:]...), nil
}

func apply(ps ...MalType) (MalType, error) {
	args, err := convertSeq(ps[len(ps)-1])
	if err != nil {
		return nil, err
	}
	for i := len(ps) - 2; i >= 1; i-- {
		args = append([]MalType{ps[i]}, args...)
	}
	if f, ok := ps[0].(MalFunc); ok {
		return f.Body(args...)
	}
	if f, ok := ps[0].(MalTCOFunc); ok {
		return f.EvalFn(f.Body, NewEnv(f.Env, f.Params, args))
	}
	return nil, fmt.Errorf("the first argument is not a function: %v", ps)
}

func mapf(ps ...MalType) (MalType, error) {
	args, err := convertSeq(ps[1])
	if err != nil {
		return nil, err
	}
	mapped := NewMalList()
	if f, ok := ps[0].(MalFunc); ok {
		for _, a := range args {
			newA, err := f.Body(a)
			if err != nil {
				return nil, err
			}
			mapped.List = append(mapped.List, newA)
		}
		return mapped, nil
	}
	if f, ok := ps[0].(MalTCOFunc); ok {
		for _, a := range args {
			newA, err := f.EvalFn(f.Body, NewEnv(f.Env, f.Params, []MalType{a}))
			if err != nil {
				return nil, err
			}
			mapped.List = append(mapped.List, newA)
		}
		return mapped, nil
	}
	return nil, fmt.Errorf("the first argument is not a function: %v", ps)
}
