import sys

import core
import mw
import printer
import reader

_PROMPT = "user> "


def eval_ast(e: mw.Expr, env) -> mw.Expr:
    if isinstance(e, mw.Symbol):
        return env.get(e)
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
    while True:
        if not isinstance(e, list):
            return eval_ast(e, env)
        if len(e) == 0:
            return e
        match e[0]:
            case mw.Symbol("def!"):
                return env.set(e[1], EVAL(e[2], env))
            case mw.Symbol("let*"):
                new_env = mw.Env(env)
                binds = e[1]
                if isinstance(binds, mw.Vector):
                    binds = binds.vector
                for k, v in zip(binds[::2], binds[1::2]):
                    new_env.set(k, EVAL(v, new_env))  # env or new_env?
                e, env = e[2], new_env
            case mw.Symbol("do"):
                [EVAL(c, env) for c in e[1:-1]]  # FIXME: step4: Why eval_ast?
                e = e[-1]
            case mw.Symbol("if"):
                cond = EVAL(e[1], env)
                cond = not (cond is mw.nil or cond is False)
                if cond:
                    e = e[2]
                else:
                    e = e[3] if len(e) == 4 else mw.nil
            case mw.Symbol("fn*"):
                return mw.Fn(
                    body=e[2],
                    params=e[1],
                    env=env,
                    fn=None,  # TODO: lambda *args: EVAL(e[2], mw.Env(env, e[1], args))
                    eval_fn=EVAL,
                )
            case default:
                e = eval_ast(e, env)
                if not isinstance(e[0], mw.Fn):
                    return e[0](*e[1:])
                fn = e[0]
                args = e[1:]
                e = fn.body
                env = mw.Env(fn.env, fn.params, args)


def PRINT(e: mw.Expr) -> str:
    return printer.pr_str(e)


def rep(arg: str, env) -> None:
    return PRINT(EVAL(READ(arg), env))


def main() -> None:
    repl_env = mw.Env()
    for k, v in core.ns.items():
        repl_env.set(k, v)
    repl_env.set(mw.Symbol("eval"), lambda e: EVAL(e, repl_env))
    repl_env.set(mw.Symbol("*ARGV*"), sys.argv[2:])

    rep("(def! not (fn* (a) (if a false true)))", repl_env)
    rep(
        '(def! load-file (fn* (f) (eval (read-string (str "(do " (slurp f) "\nnil)")))))',
        repl_env,
    )

    if len(sys.argv) > 1:
        rep('(load-file "' + sys.argv[1] + '")', repl_env)
        sys.exit(0)

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
