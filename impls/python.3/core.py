import functools
import operator
from typing import Callable

import mw
import printer
import reader


def _singularp(x: mw.Expr):
    return not isinstance(x, (list, mw.Vector, mw.Map))


def _eq(x: mw.Expr, y: mw.Expr) -> bool:
    if _singularp(x):
        if _singularp(y):
            return x == y
        return False
    if _singularp(y):
        return False
    if len(x) != len(y):
        return False
    # To make (= (list []) [(list)])  => true.
    # TODO(kibumh): Maps?
    return all(_eq(p, q) for p, q in zip(x, y))


def _pr_str(*es) -> str:
    return " ".join(printer.pr_str(e, True) for e in es)


def _str(*es) -> str:
    return "".join(printer.pr_str(e, False) for e in es)


def _prn(*es) -> mw.Nil:
    print(" ".join(printer.pr_str(e, True) for e in es))
    return mw.nil


def _println(*es) -> mw.Nil:
    print(" ".join(printer.pr_str(e, False) for e in es))
    return mw.nil


def _count(e: mw.Expr) -> int:
    return len(e) if isinstance(e, list) or isinstance(e, mw.Vector) else 0


def _slurp(e: mw.Expr) -> str:
    if not isinstance(e, str):
        raise mw.RuntimeError("path is not string, %s", e)
    with open(e, "r") as fp:
        return fp.read()


def _reset(e1: mw.Expr, e2: mw.Expr) -> mw.Expr:
    if not isinstance(e1, mw.Atom):
        raise mw.RuntimeError("first element is not atom")
    e1.v = e2
    return e2


def _swap(e: mw.Expr, f: mw.Expr, *es) -> mw.Expr:
    if not isinstance(e, mw.Atom):
        raise mw.RuntimeError("first element is not atom")
    # TODO(kibum): Use apply?
    if isinstance(f, Callable):
        e.v = f(e.v, *es)
    elif isinstance(f, mw.Fn):
        print(f, es)
        e.v = f.eval_fn(f.body, mw.Env(f.env, f.params, [e.v] + list(es)))
    else:
        raise mw.RuntimeError("second element is not function")
    return e.v


def _nth(l: mw.Expr, n: mw.Expr) -> mw.Expr:
    if isinstance(l, mw.Vector):
        l = l.vector
    if not isinstance(l, list):
        raise mw.RuntimeError("nth: argument is not a list")
    if not isinstance(n, int):
        raise mw.RuntimeError("nth: index is not integer")
    try:
        return l[n]
    except IndexError as e:
        raise mw.RuntimeError("nth: index out of range") from e


def _first(l: mw.Expr) -> mw.Expr:
    if isinstance(l, mw.Nil):
        return mw.nil
    if isinstance(l, mw.Vector):
        l = l.vector
    if not isinstance(l, list):
        raise mw.RuntimeError("first: argument is not a list")
    if not l:
        return mw.nil
    return l[0]


def _rest(l: mw.Expr) -> mw.Expr:
    if isinstance(l, mw.Nil):
        return []
    if isinstance(l, mw.Vector):
        l = l.vector
    if not isinstance(l, list):
        raise mw.RuntimeError("first: argument is not a list")
    if not l:
        return []
    return l[1:]


ns = {
    mw.Symbol("="): _eq,
    mw.Symbol("pr-str"): _pr_str,
    mw.Symbol("str"): _str,
    mw.Symbol("prn"): _prn,
    mw.Symbol("println"): _println,
    mw.Symbol("<"): operator.lt,
    mw.Symbol("<="): operator.le,
    mw.Symbol(">"): operator.gt,
    mw.Symbol(">="): operator.ge,
    mw.Symbol("+"): operator.add,
    mw.Symbol("-"): operator.sub,
    mw.Symbol("*"): operator.mul,
    mw.Symbol("/"): operator.floordiv,
    mw.Symbol("list"): lambda *es: list(es),
    mw.Symbol("list?"): lambda e: isinstance(e, list),
    mw.Symbol("empty?"): lambda e: len(e) == 0,
    mw.Symbol("count"): _count,
    mw.Symbol("read-string"): reader.read_str,
    mw.Symbol("slurp"): _slurp,
    mw.Symbol("atom"): lambda e: mw.Atom(e),
    mw.Symbol("atom?"): lambda e: isinstance(e, mw.Atom),
    # TODO(kibumh): atom이 아니면 어떡하지?
    mw.Symbol("deref"): lambda e: e.v if isinstance(e, mw.Atom) else mw.Nil,
    mw.Symbol("reset!"): _reset,
    mw.Symbol("swap!"): _swap,
    mw.Symbol("cons"): lambda e, l: [e] + (l.vector if isinstance(l, mw.Vector) else l),
    # TODO(kibumh): 모두 list인지 체크할 것
    mw.Symbol("concat"): lambda *es: functools.reduce(
        operator.add, [e.vector if isinstance(e, mw.Vector) else e for e in es], []
    ),
    mw.Symbol("vec"): lambda e: e if isinstance(e, mw.Vector) else mw.Vector(e),
    mw.Symbol("nth"): _nth,
    mw.Symbol("first"): _first,
    mw.Symbol("rest"): _rest,
}
