# ff:func feature=validate type=util control=sequence
# ff:what test: circular import A
from .b import func_b


def func_a():
    return func_b()
