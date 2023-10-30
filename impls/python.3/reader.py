import re
from typing import TypeAlias

import mw

Token: TypeAlias = str

_PATTERN = re.compile(
    r"""[\s,]*(~@|[\[\]{}()'`~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"`,;)]*)"""
)


def tokenize(s: str) -> [Token]:
    return _PATTERN.findall(s)


class SyntaxError(Exception):
    pass


class Reader:
    def __init__(self, tokens: [Token]):
        self._tokens = tokens
        self._pos = 0

    def next(self) -> Token | None:
        t = self.peek()
        if t is None:
            return t
        self._pos += 1
        return t

    def peek(self) -> Token:
        if self._pos >= len(self._tokens):
            return None
        return self._tokens[self._pos]


def read_str(s: str) -> mw.Expr:
    tokens = tokenize(s)
    reader = Reader(tokens)
    return read_form(reader)


def read_form(r: Reader) -> mw.Expr | None:
    match r.peek():
        case "":
            return None
        case "(":
            _ = r.next()
            return read_list(r)
        case "[":
            _ = r.next()
            return read_vector(r)
        case "{":
            _ = r.next()
            return read_map(r)
        case "'" | "`" | "~" | "~@" | "@":
            return read_macro(r)
        case default:
            return read_atom(r)


_MACROS = {
    "'": mw.Symbol("quote"),
    "`": mw.Symbol("quasiquote"),
    "~": mw.Symbol("unquote"),
    "~@": mw.Symbol("splice-unquote"),
    "@": mw.Symbol("deref"),
}


def read_macro(r: Reader) -> mw.Expr:
    m = r.next()
    if m not in _MACROS:
        raise Syntax("Unexpected reader macro")
    return [_MACROS[m], read_form(r)]


def read_map(r: Reader) -> mw.Map:
    m = []
    while True:
        k = read_form(r)
        if k is None:
            raise SyntaxError("unbalanced braces")
        if k == mw.Symbol("}"):
            return mw.Map(m)
        v = read_form(r)
        if v is None:
            raise SyntaxError("unbalanced braces")
        m.append((k, v))


def read_vector(r: Reader) -> mw.Vector:
    v = []
    while True:
        expr = read_form(r)
        if expr is None:
            raise SyntaxError("unbalanced brackets")
        if expr == mw.Symbol("]"):
            return mw.Vector(v)
        v.append(expr)


def read_list(r: Reader) -> mw.Expr:  # List[Expr]?
    l = []
    while True:
        expr = read_form(r)
        if expr is None:
            raise SyntaxError("unbalanced parenthesis")
        if expr == mw.Symbol(")"):
            return l
        l.append(expr)


_ESCAPES = {'"': '"', "n": "\n", "\\": "\\"}


def _read_string(t: Token) -> mw.Expr:
    if t[-1] != '"' or len(t) == 1:
        raise SyntaxError("unbalanced string")

    parsed = ""
    escaped = False
    for c in t[1:-1]:
        if escaped:
            if c not in _ESCAPES:
                raise "unsupported escape"
            parsed += _ESCAPES[c]
            escaped = False
            continue
        if c == "\\":
            escaped = True
            continue
        parsed += c
    if escaped:
        raise SyntaxError("unbalanced string")

    return parsed


def read_atom(r: Reader) -> mw.Expr:
    t = r.next()
    if t is None:
        raise SyntaxError("What?")
    if t == "nil":
        return mw.nil
    if t in ("true", "false"):
        return t == "true"
    try:
        return int(t)
    except ValueError:
        pass
    if isinstance(t, str):
        if t[0] == '"':
            return _read_string(t)
        if t[0] == ":":
            return mw.Keyword(t[1:])
        return mw.Symbol(t)
    raise SyntaxError(f"Unexpected token, {t}")
