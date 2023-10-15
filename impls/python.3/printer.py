import mw


def pr_str(e: mw.Expr) -> str:
    if isinstance(e, mw.Symbol):
        return e
    if isinstance(e, int):
        return str(e)
    if isinstance(e, list):
        return "(" + " ".join(pr_str(pr_str(c)) for c in e) + ")"
