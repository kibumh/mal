from typing import List, TypeAlias

Symbol: TypeAlias = str

Expr: TypeAlias = int | Symbol | List["Expr"]
