package core

import (
	"fmt"
	. "mal/types"
)

func cloneMap(m MalMap) MalMap {
	newM := NewMalMap()
	for k, v := range m.Map {
		newM.Map[k] = v
	}
	return newM
}

func asMap(mv MalType) (MalMap, error) {
	if _, ok := mv.(MalNil); ok {
		return MalMap{}, nil
	}
	mm, ok := mv.(MalMap)
	if !ok {
		return MalMap{}, fmt.Errorf("map is expected but, %v is given", mv)
	}
	return mm, nil
}

func assoc(ps ...MalType) (MalType, error) {
	m, err := asMap(ps[0])
	if err != nil {
		return nil, err
	}
	newM := cloneMap(m)
	for i := 1; i < len(ps); i += 2 {
		newM.Map[ps[i]] = ps[i+1]
	}
	return newM, nil
}

func containsp(ps ...MalType) (MalType, error) {
	m, err := asMap(ps[0])
	if err != nil {
		return nil, err
	}
	_, ok := m.Map[ps[1]]
	return MalBool(ok), nil
}

func dissoc(ps ...MalType) (MalType, error) {
	m, err := asMap(ps[0])
	if err != nil {
		return nil, err
	}
	newM := cloneMap(m)
	for _, k := range ps[1:] {
		delete(newM.Map, k)
	}
	return newM, nil
}

func get(ps ...MalType) (MalType, error) {
	m, err := asMap(ps[0])
	if err != nil {
		return nil, err
	}
	v, ok := m.Map[ps[1]]
	if !ok {
		return MalNil{}, nil
	}
	return v, nil
}

func hash_map(ps ...MalType) (MalType, error) {
	if len(ps)%2 != 0 {
		return nil, fmt.Errorf("hash-map: the number of arguments is odd, %v", ps)
	}
	mm := NewMalMap()
	for i := 0; i < len(ps); i += 2 {
		mm.Map[ps[i]] = ps[i+1]
	}
	return mm, nil
}

func mapp(ps ...MalType) (MalType, error) {
	_, ok := ps[0].(MalMap)
	return MalBool(ok), nil
}

func keys(ps ...MalType) (MalType, error) {
	m, err := asMap(ps[0])
	if err != nil {
		return nil, err
	}
	ks := NewMalList()
	for k, _ := range m.Map {
		ks.List = append(ks.List, k)
	}
	return ks, nil
}

func vals(ps ...MalType) (MalType, error) {
	m, err := asMap(ps[0])
	if err != nil {
		return nil, err
	}
	vs := NewMalList()
	for _, v := range m.Map {
		vs.List = append(vs.List, v)
	}
	return vs, nil
}
