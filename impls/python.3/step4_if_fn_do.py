import core
import mw
import printer
import reader

_PROMPT = "user> "


def eval_ast(e: mw.Expr, env) -> mw.Expr:
    if isinstance(e, mw.Symbol):
        return env.get(e)
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
    match e[0]:
        case "def!":
            return env.set(e[1], EVAL(e[2], env))
        case "let*":
            new_env = mw.Env(env)
            for k, v in zip(e[1][::2], e[1][1::2]):
                new_env.set(k, EVAL(v, new_env))  # env or new_env?
            return EVAL(e[2], new_env)
        case "do":
            return [EVAL(c, env) for c in e[1:]][-1]  # FIXME: step4: Why eval_ast?
        case "if":
            cond = EVAL(e[1], env)
            cond = not (cond is mw.nil or cond is False)
            then_clause = e[2]
            else_clause = e[3] if len(e) == 4 else mw.nil
            return EVAL(then_clause if cond else else_clause, env)
        case "fn*":
            return lambda *args: EVAL(e[2], mw.Env(env, e[1], args))
        case default:
            e = eval_ast(e, env)
            return e[0](*e[1:])


def PRINT(e: mw.Expr) -> str:
    return printer.pr_str(e)


def rep(arg: str, env) -> None:
    return PRINT(EVAL(READ(arg), env))


def main() -> None:
    repl_env = mw.Env()
    for k, v in core.ns.items():
        repl_env.set(k, v)
    while True:
        print(_PROMPT, end="")
        try:
            print(rep(input(), repl_env))
        except reader.SyntaxError as e:
            print(e)
        except mw.EnvError as e:
            print(e)


if __name__ == "__main__":
    main()
