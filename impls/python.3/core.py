import operator

import mw
import printer


def _atomp(x: mw.Expr):
    return not isinstance(x, (list, mw.Vector, mw.Map))


def _eq(x: mw.Expr, y: mw.Expr) -> bool:
    if _atomp(x):
        if _atomp(y):
            return x == y
        return False
    if _atomp(y):
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
}
