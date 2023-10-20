import mw
from typing import Callable


def pr_str(e: mw.Expr) -> str:
    if e is mw.nil:
        return "nil"
    if e is True:
        return "true"
    if e is False:
        return "false"
    if isinstance(e, int):
        return str(e)
    if isinstance(e, mw.Symbol):
        return e
    if isinstance(e, list):
        return "(" + " ".join(pr_str(pr_str(c)) for c in e) + ")"
    if isinstance(e, Callable):
        return "#<function>"
    raise ValueError("Unexpected expression, " + e)
