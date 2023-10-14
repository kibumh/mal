_PROMPT = "user> "


def READ(arg: str) -> str:
    return arg


def EVAL(arg: str) -> str:
    return arg


def PRINT(arg: str) -> str:
    return arg


def rep(arg: str) -> None:
    return PRINT(EVAL(READ(arg)))


def main() -> None:
    while True:
        print(_PROMPT, end="")
        print(rep(input()))


if __name__ == "__main__":
    main()
