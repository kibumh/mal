package printer

import (
	"fmt"
	"strconv"
	"strings"

	. "mal/types"
)

func PrintStr(mv MalType, readably bool) string {
	switch v := mv.(type) {
	case MalNil:
		return "nil"
	case *MalAtom:
		return fmt.Sprintf("(atom %v)", v.Value)
	case MalBool:
		if v {
			return "true"
		} else {
			return "false"
		}
	case MalInt:
		return strconv.FormatInt(int64(v), 10)
	case MalKeyword:
		return ":" + string(v)
	case MalString:
		s := string(v)
		if readably {
			s = strings.ReplaceAll(s, `\`, `\\`)
			s = strings.ReplaceAll(s, "\n", `\n`)
			s = strings.ReplaceAll(s, `"`, `\"`)
			s = `"` + s + `"`
		}
		return s
	case MalSymbol:
		return string(v)
	case MalList:
		var ss []string
		for _, cmv := range v.List {
			ss = append(ss, PrintStr(cmv, readably))
		}
		return "(" + strings.Join(ss, " ") + ")"
	case MalVector:
		var ss []string
		for _, cmv := range v.Vector {
			ss = append(ss, PrintStr(cmv, readably))
		}
		return "[" + strings.Join(ss, " ") + "]"
	case MalMap:
		var ss []string
		for key, value := range v.Map {
			ss = append(ss, PrintStr(key, readably), PrintStr(value, readably))
		}
		return "{" + strings.Join(ss, " ") + "}"
	case MalFunc, MalTCOFunc:
		return "#<function>"
	}
	panic(fmt.Sprintf("unreachable. can't print %v of type %T", mv, mv))
}
