import mw
import printer
import reader

_PROMPT = "user> "


def eval_ast(e: mw.Expr, env) -> mw.Expr:
    if isinstance(e, mw.Symbol):
        return env.get(e)
    if isinstance(e, mw.Vector):
        return mw.Vector([EVAL(child, env) for child in e])
    if isinstance(e, mw.Map):
        return mw.Map({k: EVAL(v, env) for k, v in e.items()})
    if isinstance(e, mw.List):
        return mw.List([EVAL(child, env) for child in e])
    return e


def READ(s: str) -> str:
    return reader.read_str(s)


def EVAL(e: mw.Expr, env) -> mw.Expr:
    if not isinstance(e, mw.List):
        return eval_ast(e, env)
    if len(e) == 0:
        return e
    match e[0]:
        case mw.Symbol("def!"):
            return env.set(e[1], EVAL(e[2], env))
        case mw.Symbol("let*"):
            new_env = mw.Env(env)
            binds = e[1]
            for k, v in zip(binds[::2], binds[1::2]):
                new_env.set(k, EVAL(v, new_env))  # env or new_env?
            return EVAL(e[2], new_env)
        case default:
            e = eval_ast(e, env)
            return e[0](*e[1:])


def PRINT(e: mw.Expr) -> str:
    return printer.pr_str(e)


def rep(arg: str, env) -> None:
    return PRINT(EVAL(READ(arg), env))


def main() -> None:
    env = mw.Env()
    env.set(mw.Symbol("+"), lambda a, b: a + b)
    env.set(mw.Symbol("-"), lambda a, b: a - b)
    env.set(mw.Symbol("*"), lambda a, b: a * b)
    env.set(mw.Symbol("/"), lambda a, b: a // b)
    while True:
        print(_PROMPT, end="")
        try:
            print(rep(input(), env))
        except reader.SyntaxError as e:
            print(e)
        except mw.MWError as e:
            print(e)


if __name__ == "__main__":
    main()
