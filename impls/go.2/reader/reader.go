package reader

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	. "mal/types"
)

type Token string
type Reader struct {
	tokens []Token
	pos    int
}

var tokenPattern = regexp.MustCompile("[\\s,]*(~@|[\\[\\]{}()'`~^@]|\"(?:\\\\.|[^\\\\\"])*\"?|;.*|[^\\s\\[\\]{}('\"`,;)]*)")

func (r *Reader) next() (Token, error) {
	if r.pos >= len(r.tokens) {
		return "", errors.New("no more token")
	}
	r.pos++
	return r.tokens[r.pos-1], nil
}

func (r *Reader) peek() (Token, error) {
	if r.pos >= len(r.tokens) {
		return "", errors.New("no more token")
	}
	return r.tokens[r.pos], nil
}

func ReadStr(s string) (MalType, error) {
	r := Reader{
		tokens: tokenize(s),
		pos:    0,
	}
	return readForm(&r)
}

func tokenize(s string) []Token {
	var tokens []Token
	for _, t := range tokenPattern.FindAllStringSubmatch(s, -1) {
		if len(t[1]) > 0 {
			tokens = append(tokens, Token(t[1]))
		}
	}
	return tokens
}

var singleReaderMacros = map[string]string{
	"'":  "quote",
	"`":  "quasiquote",
	"~":  "unquote",
	"~@": "splice-unquote",
	"@":  "deref",
}

var dualReaderMacros = map[string]string{
	"^": "with-meta",
}

type metComment struct{}

func (metComment) Error() string {
	return "#<comment>"
}

var MetComment = metComment{}

func readForm(r *Reader) (MalType, error) {
	token, err := r.peek()
	if err != nil {
		return nil, err
	}
	if token == "(" {
		return readList(r)
	} else if token == "[" {
		return readVector(r)
	} else if token == "{" {
		return readMap(r)
	} else if sym, ok := singleReaderMacros[string(token)]; ok {
		return readSingleMacro(r, MalSymbol(sym))
	} else if sym, ok := dualReaderMacros[string(token)]; ok {
		return readDualMacro(r, MalSymbol(sym))
	} else if token[0] == ';' {
		r.next() // Consume a comment.
		return nil, MetComment
	} else {
		return readAtom(r)
	}
}

func readSingleMacro(r *Reader, sym MalSymbol) (MalType, error) {
	r.next()
	f, err := readForm(r)
	if err != nil {
		return nil, err
	}
	return NewMalList(sym, f), nil
}

func readDualMacro(r *Reader, sym MalSymbol) (MalType, error) {
	r.next()
	f1, err := readForm(r)
	if err != nil {
		return nil, err
	}
	f2, err := readForm(r)
	if err != nil {
		return nil, err
	}
	return NewMalList(sym, f2, f1), nil
}

func readList(r *Reader) (MalType, error) {
	r.next() // Swallow '('.
	l := NewMalList()
	for true {
		t, err := r.peek()
		if err != nil {
			return nil, fmt.Errorf("unbalanced pair(')'): %w", err)
		}
		if t == ")" {
			r.next()
			break
		}
		m, err := readForm(r)
		if errors.Is(err, MetComment) {
			continue
		} else if err != nil {
			return nil, err
		}
		l.List = append(l.List, m)
	}
	return l, nil
}

func readVector(r *Reader) (MalType, error) {
	r.next() // Swallow '['.
	mv := NewMalVector()
	for true {
		t, err := r.peek()
		if err != nil {
			return nil, fmt.Errorf("unbalanced pair(']'): %w", err)
		}
		if t == "]" {
			r.next()
			break
		}
		m, err := readForm(r)
		if err != nil {
			return nil, err
		}
		mv.Vector = append(mv.Vector, m)
	}
	return mv, nil
}

func readMap(r *Reader) (MalType, error) {
	r.next() // Swallow '['.
	mm := NewMalMap()

	expectKey := true
	var key MalType = nil
	for true {
		t, err := r.peek()
		if err != nil {
			return nil, fmt.Errorf("unbalanced pair('}'): %w", err)
		}
		if t == "}" {
			if !expectKey {
				return nil, fmt.Errorf("key-value pairs are not matched")
			}
			r.next()
			break
		}
		m, err := readForm(r)
		if err != nil {
			return nil, err
		}

		if expectKey {
			key = m
			expectKey = false
		} else {
			mm.Map[key] = m
			expectKey = true
		}
	}
	return mm, nil
}

func readAtom(r *Reader) (MalType, error) {
	t, err := r.next()
	if err != nil {
		return nil, err
	}
	switch t {
	case "nil":
		return MalNil{}, nil
	case "true":
		return MalBool(true), nil
	case "false":
		return MalBool(false), nil
	}
	if i, err := strconv.ParseInt(string(t), 10, 32); err == nil {
		return MalInt(i), nil
	}
	if s := string(t); strings.HasPrefix(s, `"`) {
		escaped := false
		ended := false
		var output string
		for _, c := range s[1:] {
			if ended {
				return nil, fmt.Errorf("unbalanced string: %v", t)
			}
			if escaped {
				switch c {
				case '\\':
					output += string('\\')
				case 'n':
					output += string('\n')
				case '"':
					output += string('"')
				}
				escaped = false
			} else {
				switch c {
				case '\\':
					escaped = true
				case '"':
					ended = true
				default:
					output += string(c)
				}
			}
		}
		if !ended {
			return nil, fmt.Errorf("unbalanced string: %v", t)
		}
		return MalString(output), nil
	}
	if t[0] == ':' {
		return MalKeyword(t[1:]), nil
	}

	return MalSymbol(t), nil
}
