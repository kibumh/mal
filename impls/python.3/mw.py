from typing import List, Optional, TypeAlias

Symbol: TypeAlias = str

Expr: TypeAlias = int | Symbol | List["Expr"]


###############
# Environment #
###############


class EnvError(Exception):
    pass


class Env:
    def __init__(self, outer: Optional["Env"] = None):
        self._outer = outer
        self._data = {}

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
