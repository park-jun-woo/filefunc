# ff:func feature=validate type=util control=sequence
# ff:what test: circular import B
from .c import func_c


def func_b():
    return func_c()
