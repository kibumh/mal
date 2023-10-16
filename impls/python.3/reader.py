import re
from typing import TypeAlias

import mw

Token: TypeAlias = str

_PATTERN = re.compile(
    r"""[\s,]*(~@|[\[\]{}()'`~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"`,;)]*)"""
)


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


def tokenize(s: str) -> [Token]:
    return _PATTERN.findall(s)


def read_form(r: Reader) -> mw.Expr | None:
    match r.peek():
        case "":
            return None
        case "(":
            _ = r.next()
            return read_list(r)
        case default:
            return read_atom(r)


def read_list(r: Reader) -> mw.Expr:
    l = []
    while True:
        expr = read_form(r)
        if expr is None:
            raise SyntaxError("unbalanced parenthesis")
        if expr == ")":
            return l
        l.append(expr)


def read_atom(r: Reader) -> mw.Expr:
    t = r.next()
    if t is None:
        raise SyntaxError("What?")
    try:
        return int(t)
    except ValueError:
        pass
    return mw.Symbol(t)
