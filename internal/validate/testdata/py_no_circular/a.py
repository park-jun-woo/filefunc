# ff:func feature=validate type=util control=sequence
# ff:what test: no circular import A
from .b import func_b
from .c import func_c


def func_a():
    return func_b() + func_c()
