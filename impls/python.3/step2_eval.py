import mw
import printer
import reader

_PROMPT = "user> "


class EnvError(Exception):
    pass


def eval_ast(e: mw.Expr, env) -> mw.Expr:
    if isinstance(e, mw.Symbol):
        try:
            return env[e]
        except KeyError as exc:
            raise EnvError(f"{e} not found") from exc
    if isinstance(e, mw.Vector):
        return mw.Vector([EVAL(child, env) for child in e.vector])
    if isinstance(e, mw.Map):
        return mw.Map([(k, EVAL(v, env)) for k, v in e.m])
    if isinstance(e, list):
        return [EVAL(child, env) for child in e]
    return e


def READ(s: str) -> str:
    return reader.read_str(s)


def EVAL(e: mw.Expr, env) -> mw.Expr:
    if not isinstance(e, list):
        return eval_ast(e, env)
    if len(e) == 0:
        return e
    e = eval_ast(e, env)
    return e[0](*e[1:])


def PRINT(e: mw.Expr) -> str:
    return printer.pr_str(e)


def rep(arg: str, env) -> None:
    return PRINT(EVAL(READ(arg), env))


def main() -> None:
    env = {
        mw.Symbol("+"): lambda a, b: a + b,
        mw.Symbol("-"): lambda a, b: a - b,
        mw.Symbol("*"): lambda a, b: a * b,
        mw.Symbol("/"): lambda a, b: a // b,
    }
    while True:
        print(_PROMPT, end="")
        try:
            print(rep(input(), env))
        except reader.SyntaxError as e:
            print(e)
        except EnvError as e:
            print(e)


if __name__ == "__main__":
    main()
