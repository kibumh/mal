import dataclasses
from typing import Any, Callable, List, Optional, TypeAlias

Symbol: TypeAlias = str


class Nil:
    pass


nil = Nil()


@dataclasses.dataclass
class Fn:
    body: Any  # TODO: Expr
    params: List[Symbol]
    env: Any  # TODO: Env
    fn: Callable


Expr: TypeAlias = Nil | bool | int | Symbol | List["Expr"] | Callable | Fn


###############
# Environment #
###############


class EnvError(Exception):
    pass


class Env:
    def __init__(
        self,
        outer: Optional["Env"] = None,
        binds: List | None = None,
        exprs: List[Expr] | None = None,
    ):
        self._outer = outer
        self._data = {}
        if binds:
            for i, (k, v) in enumerate(zip(binds, exprs)):
                if k == "&":
                    self._data[binds[i + 1]] = exprs[i:]
                    break
                self._data[k] = v

    def set(self, key: Symbol, value: Expr) -> Expr:
        self._data[key] = value
        return value

    def find(self, key: Symbol) -> Optional["Env"]:
        if key in self._data:
            return self
        if self._outer is None:
            return None
        return self._outer.find(key)

    def get(self, key: Symbol) -> Expr:
        d = self.find(key)
        if d is None:
            raise EnvError(f"'{key}' not found")
        return d._data[key]
