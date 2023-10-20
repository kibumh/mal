import operator

import mw
import printer


def _prn(*es) -> mw.Expr:
    for e in es:
        print(printer.pr_str(e))
    return mw.nil


ns = {
    mw.Symbol("="): operator.eq,
    mw.Symbol("prn"): _prn,
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
    mw.Symbol("count"): lambda e: len(e) if isinstance(e, list) else 0,
}
