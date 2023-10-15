import mw
import printer
import reader

_PROMPT = "user> "


def READ(s: str) -> str:
    return reader.read_str(s)


def EVAL(e: mw.Expr) -> mw.Expr:
    return e


def PRINT(e: mw.Expr) -> str:
    return printer.pr_str(e)


def rep(arg: str) -> None:
    return PRINT(EVAL(READ(arg)))


def main() -> None:
    while True:
        print(_PROMPT, end="")
        try:
            print(rep(input()))
        except reader.SyntaxError as e:
            print(e)


if __name__ == "__main__":
    main()
