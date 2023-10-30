import mw
from typing import Callable, Optional


# NOTE: \ should be the first one to be replaced.
_REPLACES = [(r"\\", "\\"), (r"\"", '"'), (r"\n", "\n")]


def _print_string(s: str, readably: bool) -> str:
    if not readably:
        return s
    for new, org in _REPLACES:
        s = s.replace(org, new)
    return '"' + s + '"'


def pr_str(e: mw.Expr, readably: Optional[bool] = True) -> str:
    if e is mw.nil:
        return "nil"
    if e is True:
        return "true"
    if e is False:
        return "false"
    if isinstance(e, int):
        return str(e)
    if isinstance(e, mw.Symbol):
        return e.sym
    if isinstance(e, mw.Keyword):
        return ":" + e.keyword
    if isinstance(e, str):
        return _print_string(e, readably)
    if isinstance(e, mw.Vector):
        return "[" + " ".join(pr_str(c, readably) for c in e.vector) + "]"
    if isinstance(e, list):
        return "(" + " ".join(pr_str(c, readably) for c in e) + ")"
    if isinstance(e, mw.Map):
        return (
            "{"
            + " ".join(pr_str(k, readably) + " " + pr_str(v, readably) for k, v, in e.m)
            + "}"
        )
    if isinstance(e, Callable):
        return "#<function>"
    if isinstance(e, mw.Fn):
        return "#<tcofunction>"
    raise ValueError("Unexpected expression, " + e)
