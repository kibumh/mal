package core

import (
	"fmt"
	. "mal/types"
)

func meta(ps ...MalType) (MalType, error) {
	if v, ok := ps[0].(MalVector); ok {
		return v.Meta, nil
	}
	if l, ok := ps[0].(MalList); ok {
		return l.Meta, nil
	}
	if m, ok := ps[0].(MalMap); ok {
		return m.Meta, nil
	}
	if f, ok := ps[0].(MalTCOFunc); ok {
		return f.Meta, nil
	}
	if f, ok := ps[0].(MalFunc); ok {
		return f.Meta, nil
	}
	return nil, fmt.Errorf("meta: meta is not applicable to '%v'", ps)
}

func withMeta(ps ...MalType) (MalType, error) {
	if v, ok := ps[0].(MalVector); ok {
		newv := v
		newv.Meta = ps[1]
		return newv, nil
	}
	if l, ok := ps[0].(MalList); ok {
		newl := l
		newl.Meta = ps[1]
		return newl, nil
	}
	if m, ok := ps[0].(MalMap); ok {
		newm := m
		newm.Meta = ps[1]
		return newm, nil
	}
	if f, ok := ps[0].(MalTCOFunc); ok {
		newf := f
		newf.Meta = ps[1]
		return newf, nil
	}
	if f, ok := ps[0].(MalFunc); ok {
		newf := f
		newf.Meta = ps[1]
		return newf, nil
	}
	return nil, fmt.Errorf("meta: with-meta is not applicable to '%v'", ps)
}
