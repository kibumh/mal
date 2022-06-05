package main

import (
	"fmt"
	"log"

	. "mal/printer"
	. "mal/reader"
	. "mal/types"

	"github.com/chzyer/readline"
)

func READ(line string) (MalType, error) {
	return ReadStr(line)
}

func evalAst(mv MalType, env *Env) (MalType, error) {
	switch v := mv.(type) {
	case MalSymbol:
		sv, err := env.Find(v)
		if err != nil {
			return nil, err
		}
		return sv, nil
	case MalVector:
		var mvec MalVector
		for _, cv := range v.Vector {
			cev, err := EVAL(cv, env)
			if err != nil {
				return nil, err
			}
			mvec.Vector = append(mvec.Vector, cev)
		}
		return mvec, nil
	case MalMap:
		mmap := NewMalMap()
		for key, val := range v.Map {
			keyev, err := EVAL(key, env)
			if err != nil {
				return nil, err
			}
			valev, err := EVAL(val, env)
			if err != nil {
				return nil, err
			}
			mmap.Map[keyev] = valev
		}
		return mmap, nil
	case MalList:
		var ml MalList
		for _, cv := range v.List {
			cev, err := EVAL(cv, env)
			if err != nil {
				return nil, err
			}
			ml.List = append(ml.List, cev)
		}
		return ml, nil
	default:
		return v, nil
	}
}

func EVAL(mv MalType, env *Env) (MalType, error) {
	switch v := mv.(type) {
	case MalInt, MalSymbol, MalKeyword, MalString, MalVector, MalMap:
		return evalAst(v, env)
	case MalList:
		if len(v.List) == 0 {
			return v, nil
		}
		if op, ok := v.List[0].(MalSymbol); ok {
			switch op {
			case MalSymbol("def!"):
				key, ok := v.List[1].(MalSymbol)
				if !ok {
					return nil, fmt.Errorf("bind key is not a symbol, %v", v.List[1])
				}
				value, err := EVAL(v.List[2], env)
				if err != nil {
					return nil, err
				}
				env.Set(key, value)
				return value, nil
			case MalSymbol("let*"):
				newEnv := NewEnv(env, NewMalList(), nil)

				var bindings []MalType
				if ml, ok := v.List[1].(MalList); ok {
					bindings = ml.List
				} else if mvec, ok := v.List[1].(MalVector); ok {
					bindings = mvec.Vector
				} else {
					return nil, fmt.Errorf("binding is not a list nor a vector, %v", v.List[1])
				}

				for i := 0; i < len(bindings); i += 2 {
					key, ok := bindings[i].(MalSymbol)
					if !ok {
						return nil, fmt.Errorf("bind key is not a symbol, %v", bindings[i])
					}
					value, err := EVAL(bindings[i+1], newEnv)
					if err != nil {
						return nil, fmt.Errorf("cannot eval %v, %w", bindings[i+1], err)
					}
					newEnv.Set(key, value)
				}
				return EVAL(v.List[2], newEnv)
			}
		}
		lv, err := evalAst(v, env)
		if err != nil {
			return nil, err
		}
		v = lv.(MalList)
		if err != nil {
			return nil, err
		}
		return v.List[0].(MalFunc).Body(v.List[1:]...)
	}
	panic("Unreachable")
}

func PRINT(mv MalType) string {
	return PrintStr(mv, true)
}

func rep(line string, env *Env) string {
	mv, err := READ(line)
	if err != nil {
		return err.Error()
	}
	mv, err = EVAL(mv, env)
	if err != nil {
		return err.Error()
	}
	return PRINT(mv)
}

func initEnv() *Env {
	env := NewEnv(nil, NewMalList(), nil)
	env.Set(MalSymbol("+"), NewMalFunc(func(mvs ...MalType) (MalType, error) { return mvs[0].(MalInt) + mvs[1].(MalInt), nil }))
	env.Set(MalSymbol("-"), NewMalFunc(func(mvs ...MalType) (MalType, error) { return mvs[0].(MalInt) - mvs[1].(MalInt), nil }))
	env.Set(MalSymbol("*"), NewMalFunc(func(mvs ...MalType) (MalType, error) { return mvs[0].(MalInt) * mvs[1].(MalInt), nil }))
	env.Set(MalSymbol("/"), NewMalFunc(func(mvs ...MalType) (MalType, error) { return mvs[0].(MalInt) / mvs[1].(MalInt), nil }))
	return env
}

func main() {
	rl, err := readline.New("user> ")
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()

	env := initEnv()
	for true {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		ret := rep(string(line), env)
		fmt.Println(ret)
	}
}
