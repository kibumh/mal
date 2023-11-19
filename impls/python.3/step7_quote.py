import sys

import core
import mw
import printer
import reader

_PROMPT = "user> "


def _quasiquote(e: mw.Expr, do_unquote=True) -> mw.Expr:
    if isinstance(e, mw.List):
        if do_unquote and e and e[0] == mw.Symbol("unquote"):
            return e[1]
        else:
            ret = mw.List([])
            for c in reversed(e):
                if isinstance(c, mw.List) and c and c[0] == mw.Symbol("splice-unquote"):
                    ret = mw.List([mw.Symbol("concat"), c[1], ret])
                else:
                    ret = mw.List([mw.Symbol("cons"), _quasiquote(c), ret])
            return ret
    elif isinstance(e, (mw.Map, mw.Symbol)):
        return mw.List([mw.Symbol("quote"), e])
    elif isinstance(e, mw.Vector):
        return mw.List([mw.Symbol("vec"), _quasiquote(mw.List(e.v), False)])
    else:
        return e


def _eval_ast(e: mw.Expr, env) -> mw.Expr:
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
    while True:
        if not isinstance(e, mw.List):
            return _eval_ast(e, env)
        if len(e) == 0:
            return e
        match e[0]:
            case mw.Symbol("quote"):
                return e[1]
            case mw.Symbol("quasiquote"):
                e = _quasiquote(e[1])  # TCO.
            case mw.Symbol("quasiquoteexpand"):
                return _quasiquote(e[1])
            case mw.Symbol("def!"):
                return env.set(e[1], EVAL(e[2], env))
            case mw.Symbol("let*"):
                new_env = mw.Env(env)
                binds = e[1]
                for k, v in zip(binds[::2], binds[1::2]):
                    new_env.set(k, EVAL(v, new_env))  # env or new_env?
                e, env = e[2], new_env  # TCO.
            case mw.Symbol("do"):
                [EVAL(c, env) for c in e[1:-1]]  # FIXME: step4: Why _eval_ast?
                e = e[-1]  # TCO.
            case mw.Symbol("if"):
                cond = EVAL(e[1], env)
                cond = not (cond is mw.nil or cond is False)
                if cond:
                    e = e[2]  # TCO.
                else:
                    e = e[3] if len(e) == 4 else mw.nil  # TCO.
            case mw.Symbol("fn*"):
                return mw.Fn(
                    body=e[2],
                    params=e[1],
                    env=env,
                    fn=None,  # TODO: lambda *args: EVAL(e[2], mw.Env(env, e[1], args))
                    eval_fn=EVAL,
                )
            case default:
                e = _eval_ast(e, env)
                if not isinstance(e[0], mw.Fn):
                    return e[0](*e[1:])
                fn = e[0]
                args = mw.List(e[1:])
                env = mw.Env(fn.env, fn.params, args)
                e = fn.body  # TCO.


def PRINT(e: mw.Expr) -> str:
    return printer.pr_str(e)


def rep(arg: str, env) -> None:
    return PRINT(EVAL(READ(arg), env))


def main() -> None:
    repl_env = mw.Env()
    for k, v in core.ns.items():
        repl_env.set(k, v)
    repl_env.set(mw.Symbol("eval"), lambda e: EVAL(e, repl_env))
    repl_env.set(mw.Symbol("*ARGV*"), mw.List(sys.argv[2:]))

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
        except mw.MWError as e:
            print(e)


if __name__ == "__main__":
    main()
