import functools
import operator
from typing import Callable

import mw
import printer
import reader


def _singularp(x: mw.Expr):
    return not isinstance(x, (mw.List, mw.Vector, mw.Map))


def _eq(x: mw.Expr, y: mw.Expr) -> bool:
    if _singularp(x):
        if _singularp(y):
            return x == y
        return False
    if _singularp(y):
        return False
    if len(x) != len(y):
        return False
    if isinstance(x, mw.Map) != isinstance(y, mw.Map):
        return False
    return all(_eq(p, q) for p, q in zip(x, y))


def _throw(e: mw.Expr):
    raise mw.MWError(e)


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
    return len(e) if isinstance(e, mw.List) or isinstance(e, mw.Vector) else 0


def _slurp(e: mw.Expr) -> str:
    if not isinstance(e, str):
        raise mw.MWError("path is not string, %s", e)
    with open(e, "r") as fp:
        return fp.read()


def _reset(e1: mw.Expr, e2: mw.Expr) -> mw.Expr:
    if not isinstance(e1, mw.Atom):
        raise mw.MWError("first element is not atom")
    e1.a = e2
    return e2


def _swap(e: mw.Expr, f: mw.Expr, *es) -> mw.Expr:
    if not isinstance(e, mw.Atom):
        raise mw.MWError("first element is not atom")
    # TODO(kibum): Use apply?
    if isinstance(f, Callable):
        e.a = f(e.a, *es)
    elif isinstance(f, mw.Fn):
        e.a = f.eval_fn(f.body, mw.Env(f.env, f.params, mw.List([e.a] + list(es))))
    else:
        raise mw.MWError("second element is not function")
    return e.a


def _nth(l: mw.Expr, n: mw.Expr) -> mw.Expr:
    if not isinstance(l, (mw.List, mw.Vector)):
        raise mw.MWError("nth: argument is not `sequential`")
    if not isinstance(n, int):
        raise mw.MWError("nth: index is not integer")
    try:
        return l[n]
    except IndexError as e:
        raise mw.MWError("nth: index out of range") from e


def _first(l: mw.Expr) -> mw.Expr:
    if isinstance(l, mw.Nil):
        return mw.nil
    if not isinstance(l, (mw.List, mw.Vector)):
        raise mw.MWError("first: argument is not `sequential`")
    if not l:
        return mw.nil
    return l[0]


def _rest(l: mw.Expr) -> mw.Expr:
    if isinstance(l, mw.Nil):
        return mw.List([])
    if not isinstance(l, (mw.List, mw.Vector)):
        raise mw.MWError("rest: argument is not `sequential`")
    return mw.List(l[1:] if l else [])


def _apply(f: mw.Expr, *es) -> mw.Expr:
    args = []
    for e in es[:-1]:
        args.append(e)
    args += es[-1]

    if isinstance(f, Callable):
        return f(*args)
    return f.eval_fn(f.body, mw.Env(f.env, f.params, mw.List(args)))


def _map(f: mw.Expr, l: mw.Expr) -> mw.Expr:
    if isinstance(f, Callable):
        return mw.List([f(c) for c in l])
    return mw.List(
        [f.eval_fn(f.body, mw.Env(f.env, f.params, mw.List([c]))) for c in l]
    )


def _hash_map(*es) -> mw.Map:
    if len(es) % 2 != 0:
        raise mw.MWError("Odd number of arguments are given for 'hash-map'")
    m = dict()
    for k, v in zip(es[::2], es[1::2]):
        m[k] = v
    print(m)
    return mw.Map(m)


def _assoc(mwm: mw.Map, *es) -> mw.Map:
    if len(es) % 2 != 0:
        raise mw.MWError("Odd number of arguments are given for 'hash-map'")
    m = mwm.m.copy()
    for k, v in zip(es[::2], es[1::2]):
        m[k] = v
    return mw.Map(m)


def _dissoc(mwm: mw.Map, *es) -> mw.Map:
    m = mwm.m.copy()
    for k in es:
        m.pop(k, None)
    return mw.Map(m)


ns = {
    mw.Symbol("="): _eq,
    mw.Symbol("throw"): _throw,
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
    mw.Symbol("nil?"): lambda e: isinstance(e, mw.Nil),
    mw.Symbol("true?"): lambda e: isinstance(e, bool) and e,
    mw.Symbol("false?"): lambda e: isinstance(e, bool) and not e,
    mw.Symbol("symbol"): lambda e: e if isinstance(e, mw.Symbol) else mw.Symbol(e),
    mw.Symbol("symbol?"): lambda e: isinstance(e, mw.Symbol),
    mw.Symbol("keyword"): lambda e: e if isinstance(e, mw.Keyword) else mw.Keyword(e),
    mw.Symbol("keyword?"): lambda e: isinstance(e, mw.Keyword),
    mw.Symbol("list"): lambda *es: mw.List(list(es)),
    mw.Symbol("list?"): lambda e: isinstance(e, mw.List),
    mw.Symbol("vec"): lambda e: e if isinstance(e, mw.Vector) else mw.Vector(e),
    mw.Symbol("vector"): lambda *es: mw.Vector(es),
    mw.Symbol("vector?"): lambda e: isinstance(e, mw.Vector),
    mw.Symbol("empty?"): lambda e: len(e) == 0,
    mw.Symbol("count"): _count,
    mw.Symbol("sequential?"): lambda e: isinstance(e, (mw.List, mw.Vector)),
    mw.Symbol("read-string"): reader.read_str,
    mw.Symbol("slurp"): _slurp,
    mw.Symbol("hash-map"): _hash_map,
    mw.Symbol("map?"): lambda e: isinstance(e, mw.Map),
    mw.Symbol("assoc"): _assoc,
    mw.Symbol("dissoc"): _dissoc,
    mw.Symbol("get"): lambda mwm, k: (
        mwm.m.get(k, mw.nil) if isinstance(mwm, mw.Map) else mw.nil
    ),
    mw.Symbol("contains?"): lambda mwm, k: k in mwm.m,
    mw.Symbol("keys"): lambda mwm: mwm.keys(),
    mw.Symbol("vals"): lambda mwm: mwm.values(),
    mw.Symbol("atom"): lambda e: mw.Atom(e),
    mw.Symbol("atom?"): lambda e: isinstance(e, mw.Atom),
    # TODO(kibumh): atom이 아니면 어떡하지?
    mw.Symbol("deref"): lambda e: e.a if isinstance(e, mw.Atom) else mw.Nil,
    mw.Symbol("reset!"): _reset,
    mw.Symbol("swap!"): _swap,
    mw.Symbol("cons"): lambda e, l: mw.List(
        [e] + (l.v if isinstance(l, mw.Vector) else l.l)
    ),
    # TODO(kibumh): 모두 list인지 체크할 것
    mw.Symbol("concat"): lambda *es: mw.List(
        functools.reduce(
            operator.add, [e.v if isinstance(e, mw.Vector) else e.l for e in es], []
        )
    ),
    mw.Symbol("nth"): _nth,
    mw.Symbol("first"): _first,
    mw.Symbol("rest"): _rest,
    mw.Symbol("apply"): _apply,
    mw.Symbol("map"): _map,
}
