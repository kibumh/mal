package core

import (
	. "mal/types"
)

type namespace map[MalSymbol]MalFunc

var CoreNS = namespace{
	// atom.go
	MalSymbol("atom"):   NewMalFunc(atom),
	MalSymbol("atom?"):  NewMalFunc(atomp),
	MalSymbol("deref"):  NewMalFunc(deref),
	MalSymbol("reset!"): NewMalFunc(reset),
	MalSymbol("swap!"):  NewMalFunc(swap),

	// bool.go
	MalSymbol("false?"): NewMalFunc(falsep),
	MalSymbol("true?"):  NewMalFunc(truep),

	// file.go
	MalSymbol("slurp"): NewMalFunc(slurp),

	// func.go
	MalSymbol("fn?"):    NewMalFunc(fnp),
	MalSymbol("macro?"): NewMalFunc(macrop),

	// keyword.go
	MalSymbol("keyword"):  NewMalFunc(keyword),
	MalSymbol("keyword?"): NewMalFunc(keywordp),

	// map.go
	MalSymbol("assoc"):     NewMalFunc(assoc),
	MalSymbol("contains?"): NewMalFunc(containsp),
	MalSymbol("dissoc"):    NewMalFunc(dissoc),
	MalSymbol("get"):       NewMalFunc(get),
	MalSymbol("hash-map"):  NewMalFunc(hash_map),
	MalSymbol("keys"):      NewMalFunc(keys),
	MalSymbol("map?"):      NewMalFunc(mapp),
	MalSymbol("vals"):      NewMalFunc(vals),

	// meta.go
	MalSymbol("meta"):      NewMalFunc(meta),
	MalSymbol("with-meta"): NewMalFunc(withMeta),

	// nil.go
	MalSymbol("nil?"): NewMalFunc(nilp),

	// list.go
	MalSymbol("apply"):  NewMalFunc(apply),
	MalSymbol("concat"): NewMalFunc(concat),
	MalSymbol("cons"):   NewMalFunc(cons),
	MalSymbol("count"):  NewMalFunc(count),
	MalSymbol("empty?"): NewMalFunc(emptyp),
	MalSymbol("first"):  NewMalFunc(first),
	MalSymbol("list"):   NewMalFunc(listf),
	MalSymbol("list?"):  NewMalFunc(listp),
	MalSymbol("map"):    NewMalFunc(mapf),
	MalSymbol("nth"):    NewMalFunc(nth),
	MalSymbol("rest"):   NewMalFunc(rest),

	// number.go
	MalSymbol("number?"): NewMalFunc(numberp),

	// op.go
	MalSymbol("+"): NewMalFunc(add),
	MalSymbol("-"): NewMalFunc(sub),
	MalSymbol("*"): NewMalFunc(mul),
	MalSymbol("/"): NewMalFunc(div),

	MalSymbol("<"):  NewMalFunc(lt),
	MalSymbol("<="): NewMalFunc(le),
	MalSymbol("="):  NewMalFunc(eq),
	MalSymbol(">"):  NewMalFunc(gt),
	MalSymbol(">="): NewMalFunc(ge),

	// print.go
	MalSymbol("pr-str"):  NewMalFunc(prStr),
	MalSymbol("println"): NewMalFunc(println),
	MalSymbol("prn"):     NewMalFunc(prn),
	MalSymbol("str"):     NewMalFunc(str),

	// read.go
	MalSymbol("read-string"): NewMalFunc(readString),
	MalSymbol("readline"):    NewMalFunc(readLine),

	// seq.go
	MalSymbol("conj"):        NewMalFunc(conj),
	MalSymbol("seq"):         NewMalFunc(seq),
	MalSymbol("sequential?"): NewMalFunc(sequentialp),

	// string.go
	MalSymbol("string?"): NewMalFunc(stringp),

	// symbol.go
	MalSymbol("symbol"):  NewMalFunc(symbol),
	MalSymbol("symbol?"): NewMalFunc(symbolp),

	// time.go
	MalSymbol("time-ms"): NewMalFunc(timeMs),

	// throw.go
	MalSymbol("throw"): NewMalFunc(throw),

	// vector.go
	MalSymbol("vec"):     NewMalFunc(vec),
	MalSymbol("vector"):  NewMalFunc(vector),
	MalSymbol("vector?"): NewMalFunc(vectorp),
}

func notImplemented(ps ...MalType) (MalType, error) {
	// panic(fmt.Sprintf("Not implemented, %v", ps))
	return MalNil{}, nil
}
