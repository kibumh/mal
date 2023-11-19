import dataclasses
import itertools
from typing import Any, Callable, List as pyList, Optional, TypeAlias


class MWError(Exception):
    pass


class Comment:
    pass


class Nil:
    pass


nil = Nil()


@dataclasses.dataclass(eq=True, frozen=True, order=True)
class Symbol:
    sym: str


@dataclasses.dataclass(eq=True, frozen=True, order=True)
class Keyword:
    k: str


# TODO: As is_macro is set after creation, we can't set frozen.
@dataclasses.dataclass
class Fn:
    body: Any  # TODO: Expr
    params: pyList[Symbol]
    env: Any  # TODO: Env
    fn: Callable
    eval_fn: Callable  # Hack to pass eval function to core module
    is_macro: bool = False


# NOTE(kibumh): We had used native python list as mal's list.
# But there was a problem. As python list is mutable, it can't be used
# as a key of a dictionary. Hence, mw.List is introduced.
@dataclasses.dataclass(eq=True, frozen=True, order=True)
class List:
    l: Any  # pyList["Expr"]

    def __len__(self):
        return len(self.l)

    def __iter__(self):
        return self.l.__iter__()

    def __getitem__(self, index):
        return self.l[index]


@dataclasses.dataclass(eq=True, frozen=True, order=True)
class Vector:
    v: Any  # pyList["Expr"]

    def __len__(self):
        return len(self.v)

    def __iter__(self):
        return self.v.__iter__()

    def __getitem__(self, index):
        return self.v[index]


@dataclasses.dataclass(eq=True, frozen=True, order=True)
class Map:
    m: Any  # Dict[Expr, Expr]

    def __len__(self):
        return len(self.m)

    def __iter__(self):
        # NOTE(kibumh): eq 구현을 쉽게 하기 위해 정렬된 순서로 돌자.
        for k in sorted(self.m.keys()):
            yield List([k, self.m[k]])

    def keys(self) -> List:
        return List(list(self.m.keys()))

    def values(self) -> List:
        return List(list(self.m.values()))

    def items(self):
        return self.m.items()


@dataclasses.dataclass
class Atom:
    a: Any  # Expr


Expr: TypeAlias = (
    Comment
    | Nil
    | bool
    | int
    | Symbol
    | Keyword
    | str
    | List
    | Vector
    | Map
    | Atom
    | Callable
    | Fn
)


###############
# Environment #
###############


class Env:
    def __init__(
        self,
        outer: Optional["Env"] = None,
        binds: pyList | None = None,
        exprs: pyList[Expr] | None = None,
    ):
        self._outer = outer
        self._data = {}
        if not binds:
            return

        for i, (k, v) in enumerate(itertools.zip_longest(binds, exprs)):
            if k == Symbol("&"):
                self._data[binds[i + 1]] = List(exprs[i:])
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
            raise MWError(f"'{key.sym}' not found")
        return d._data[key]
