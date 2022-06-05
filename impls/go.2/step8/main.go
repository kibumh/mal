package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	. "mal/core"
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
	for true {
		switch v := mv.(type) {
		case MalNil:
			return v, nil
		case MalBool, MalInt, MalSymbol, MalKeyword, MalFunc, MalString, MalVector, MalMap:
			return evalAst(v, env)
		case MalTCOFunc:
			return v, nil
		case MalList:
			if len(v.List) == 0 {
				return v, nil
			}
			expanded, err := macroExpand(v, env)
			if err != nil {
				return nil, err
			}
			v2, ok := expanded.(MalList)
			if !ok {
				return evalAst(expanded, env)
			}
			v = v2
			if op, ok := v.List[0].(MalSymbol); ok {
				switch op {
				case MalSymbol("def!"), MalSymbol("defmacro!"):
					key, ok := v.List[1].(MalSymbol)
					if !ok {
						return nil, fmt.Errorf("bind key is not a symbol, %v", v.List[1])
					}
					value, err := EVAL(v.List[2], env)
					if err != nil {
						return nil, err
					}
					if mtf, ok := value.(MalTCOFunc); ok {
						mtf.IsMacro = op == MalSymbol("defmacro!")
						value = mtf
					}
					env.Set(key, value)
					return value, nil
				case MalSymbol("macroexpand"):
					return macroExpand(v.List[1], env)
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
					mv = v.List[2]
					env = newEnv
					continue
				case MalSymbol("do"):
					_, err := evalAst(NewMalList(v.List[1:len(v.List)-1]...), env)
					if err != nil {
						return nil, err
					}
					mv = v.List[len(v.List)-1]
					continue // TCO
				case MalSymbol("if"):
					cond, err := EVAL(v.List[1], env)
					if err != nil {
						return nil, err
					}
					if cond != MalNil(struct{}{}) && cond != MalBool(false) {
						mv = v.List[2]
					} else {
						if len(v.List) == 3 {
							return MalNil{}, nil
						}
						mv = v.List[3]
					}
					continue // TCO
				case MalSymbol("fn*"):
					var syms MalList
					if syms, ok = v.List[1].(MalList); !ok {
						syms = NewMalList(v.List[1].(MalVector).Vector...)
					}
					return MalTCOFunc{
						Body:    v.List[2],
						Params:  syms,
						Env:     env,
						IsMacro: false,
						EvalFn:  EVAL,
					}, nil
				case MalSymbol("quote"):
					return v.List[1], nil
				case MalSymbol("quasiquoteexpand"):
					mv, err := quasiquote(v.List[1], env, 1)
					if err != nil {
						return nil, err
					}
					return mv, nil
				case MalSymbol("quasiquote"):
					new_mv, err := quasiquote(v.List[1], env, 1)
					if err != nil {
						return nil, err
					}
					mv = new_mv
					continue // TCO
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
			if tcoFunc, ok := v.List[0].(MalTCOFunc); ok {
				env = NewEnv(tcoFunc.Env, tcoFunc.Params, v.List[1:])
				mv = tcoFunc.Body
			} else {
				return v.List[0].(MalFunc).Body(v.List[1:]...)
			}
		}
	}
	panic(fmt.Sprintf("Unreachable: EVAL %v", mv))
}

func PRINT(mv MalType) string {
	return PrintStr(mv, true)
}

func isMacroCall(mv MalType, env *Env) bool {
	mlist, ok := mv.(MalList)
	if !ok {
		return false
	}
	if len(mlist.List) == 0 {
		return false
	}
	msym, ok := mlist.List[0].(MalSymbol)
	if !ok {
		return false
	}
	mvalue, err := env.Find(msym)
	if err != nil {
		return false
	}
	mtf, ok := mvalue.(MalTCOFunc)
	return ok && mtf.IsMacro
}

func macroExpand(mv MalType, env *Env) (MalType, error) {
	for isMacroCall(mv, env) {
		mlist := mv.(MalList)
		mv2, _ := env.Find(mlist.List[0].(MalSymbol))
		mtf := mv2.(MalTCOFunc)
		mv3, err := EVAL(mtf.Body, NewEnv(mtf.Env, mtf.Params, mlist.List[1:]))
		if err != nil {
			return nil, err
		}
		mv = mv3
	}
	return mv, nil
}

func quasiquote(mv MalType, env *Env, depth int) (MalType, error) {
	var ml []MalType
	var isListish bool
	if mlist, ok := mv.(MalList); ok {
		ml = mlist.List
		isListish = true
	} else if mvec, ok := mv.(MalVector); ok {
		ml = mvec.Vector
		isListish = true
	}
	if isListish {
		if len(ml) == 2 && ml[0] == MalSymbol("unquote") {
			if _, ok := mv.(MalList); ok {
				return ml[1], nil // Hmm... What a smell!
			}
		}
		var mlist MalList
		for i := len(ml) - 1; i >= 0; i-- {
			if cml, ok := ml[i].(MalList); ok && len(cml.List) == 2 && cml.List[0] == MalSymbol("splice-unquote") {
				mlist = NewMalList(MalSymbol("concat"), cml.List[1], mlist)
			} else {
				cml, err := quasiquote(ml[i], env, depth+1)
				if err != nil {
					return nil, err
				}
				mlist = NewMalList(MalSymbol("cons"), cml, mlist)
			}
		}
		if _, ok := mv.(MalList); ok {
			return mlist, nil
		} else {
			return NewMalList(MalSymbol("vec"), mlist), nil
		}
	}
	if mmap, ok := mv.(MalMap); ok {
		return NewMalList(MalSymbol("quote"), mmap), nil
	}
	if msym, ok := mv.(MalSymbol); ok {
		return NewMalList(MalSymbol("quote"), msym), nil
	}
	return mv, nil
}

func rep(line string, env *Env) string {
	mv, err := READ(line)
	if errors.Is(err, MetComment) {
		return ""
	} else if err != nil {
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
	for k, v := range CoreNS {
		env.Set(k, v)
	}

	env.Set(MalSymbol("eval"), NewMalFunc(func(ps ...MalType) (MalType, error) {
		return EVAL(ps[0], env)
	}))

	var args MalList
	if len(os.Args) > 2 {
		for _, a := range os.Args[2:] {
			args.List = append(args.List, MalString(a))
		}
	}
	env.Set(MalSymbol("*ARGV*"), args)

	rep("(def! not (fn* (a) (if a false true)))", env)
	rep("(def! load-file (fn* (f) (eval (read-string (str \"(do \" (slurp f) \"\nnil)\")))))", env)
	rep("(defmacro! cond (fn* (& xs) (if (> (count xs) 0) (list 'if (first xs) (if (> (count xs) 1) (nth xs 1) (throw \"odd number of forms to cond\")) (cons 'cond (rest (rest xs)))))))", env)

	return env
}

func main() {
	env := initEnv()

	if len(os.Args) > 1 {
		rep(fmt.Sprintf("(load-file %q)", os.Args[1]), env)
		return
	}

	var args MalList
	for _, a := range os.Args[1:] {
		args.List = append(args.List, MalString(a))
	}
	env.Set(MalSymbol("*ARGV*"), args)

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
		ret := rep(string(line), env)
		if ret != "" { // TODO: Handle comment better.
			fmt.Println(ret)
		}
	}
}
