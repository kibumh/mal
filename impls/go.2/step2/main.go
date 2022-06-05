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

var replEnv = NewEnv(nil,
	NewMalList(MalSymbol("+"),
		MalSymbol("-"),
		MalSymbol("*"),
		MalSymbol("/")),
	[]MalType{
		NewMalFunc(func(mvs ...MalType) (MalType, error) { return mvs[0].(MalInt) + mvs[1].(MalInt), nil }),
		NewMalFunc(func(mvs ...MalType) (MalType, error) { return mvs[0].(MalInt) - mvs[1].(MalInt), nil }),
		NewMalFunc(func(mvs ...MalType) (MalType, error) { return mvs[0].(MalInt) * mvs[1].(MalInt), nil }),
		NewMalFunc(func(mvs ...MalType) (MalType, error) { return mvs[0].(MalInt) / mvs[1].(MalInt), nil }),
	},
)

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
	panic(fmt.Sprintf("Unreachable: mv(%v)", mv))
}

func PRINT(mv MalType) string {
	return PrintStr(mv, true)
}

func rep(line string) string {
	mv, err := READ(line)
	if err != nil {
		return err.Error()
	}
	mv, err = EVAL(mv, replEnv)
	if err != nil {
		return err.Error()
	}
	return PRINT(mv)
}

func main() {
	rl, err := readline.New("user> ")
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()

	for true {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		ret := rep(string(line))
		fmt.Println(ret)
	}
}
