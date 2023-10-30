import dataclasses
import itertools
from typing import Any, Callable, List, Optional, TypeAlias


class Nil:
    pass


nil = Nil()


@dataclasses.dataclass(frozen=True)
class Symbol:
    sym: str


@dataclasses.dataclass(frozen=True)
class Keyword:
    keyword: str


@dataclasses.dataclass(frozen=True)
class Fn:
    body: Any  # TODO: Expr
    params: List[Symbol]
    env: Any  # TODO: Env
    fn: Callable


@dataclasses.dataclass
class Vector:
    vector: Any  # list["Expr"]

    def __len__(self):
        return len(self.vector)

    def __iter__(self):
        return self.vector.__iter__()


@dataclasses.dataclass
class Map:
    # The key type of python dict should be immutable.
    # Let's use list as an underlying data structure for a while.
    m: Any  # List[Expr, Expr]


Expr: TypeAlias = (
    Nil | bool | int | Symbol | Keyword | Vector | List["Expr"] | Callable | Fn
)


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
            for i, (k, v) in enumerate(itertools.zip_longest(binds, exprs)):
                if k == Symbol("&"):
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
            raise EnvError(f"'{key.sym}' not found")
        return d._data[key]
