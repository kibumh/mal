package types

import "fmt"

type MalType interface {
	IsMalType() bool
}

type MalBool bool

func (_ MalBool) IsMalType() bool {
	return true
}

type MalFunc struct {
	Body func(...MalType) (MalType, error)
	Meta MalType
}

func (_ MalFunc) IsMalType() bool {
	return true
}

func NewMalFunc(f func(...MalType) (MalType, error)) MalFunc {
	return MalFunc{f, MalNil{}}
}

type MalInt int64

func (_ MalInt) IsMalType() bool {
	return true
}

type MalKeyword string // TODO(kibumh): Intern it.
func (_ MalKeyword) IsMalType() bool {
	return true
}

type MalNil struct{}

func (_ MalNil) IsMalType() bool {
	return true
}

type MalString string

func (_ MalString) IsMalType() bool {
	return true
}

type MalSymbol string

func (_ MalSymbol) IsMalType() bool {
	return true
}

type MalTCOFunc struct {
	Body    MalType
	Params  MalList
	Env     *Env
	IsMacro bool
	EvalFn  func(mv MalType, env *Env) (MalType, error)
	Meta    MalType
	// Fn     interface{} // FIXME
}

func (_ MalTCOFunc) IsMalType() bool {
	return true
}

type MalException struct {
	Value MalType
}

func (_ MalException) IsMalType() bool {
	return true
}
func (_ MalException) Error() string {
	return "MalException"
}

type MalAtom struct {
	Value MalType
}

func (_ *MalAtom) IsMalType() bool {
	return true
}

func NewMalAtom(v MalType) *MalAtom {
	return &MalAtom{v}
}

type MalList struct {
	List []MalType
	Meta MalType
}

func (_ MalList) IsMalType() bool {
	return true
}

func NewMalList(ps ...MalType) MalList {
	return MalList{
		List: ps,
		Meta: MalNil{},
	}
}

type MalVector struct {
	Vector []MalType
	Meta   MalType
}

func (_ MalVector) IsMalType() bool {
	return true
}
func NewMalVector(ps ...MalType) MalVector {
	return MalVector{
		Vector: ps,
		Meta:   MalNil{},
	}

}

type MalMap struct {
	Map  map[MalType]MalType
	Meta MalType
}

func (_ MalMap) IsMalType() bool {
	return true
}
func NewMalMap() MalMap {
	return MalMap{
		Map:  make(map[MalType]MalType),
		Meta: MalNil{},
	}
}

// env.go
type Env struct {
	outer *Env
	data  (map[MalSymbol]MalType)
}

func NewEnv(outer *Env, binds MalList, exprs []MalType) *Env {
	data := make(map[MalSymbol]MalType)

	for i := range binds.List {
		param := binds.List[i].(MalSymbol)
		if param == MalSymbol("&") {
			mlist := NewMalList()
			mlist.List = append(mlist.List, exprs[i:]...)
			data[binds.List[i+1].(MalSymbol)] = mlist
			break
		} else {
			data[param] = exprs[i]
		}
	}
	env := &Env{
		outer: outer,
		data:  data,
	}
	return env
}

func (env *Env) Set(k MalSymbol, mv MalType) {
	env.data[k] = mv
}

func (env *Env) Find(k MalSymbol) (MalType, error) {
	if mv, ok := env.data[k]; ok {
		return mv, nil
	}
	if env.outer != nil {
		return env.outer.Find(k)
	}
	return nil, fmt.Errorf("'%v' not found", k)
}
