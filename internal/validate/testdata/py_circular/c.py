# ff:func feature=validate type=util control=sequence
# ff:what test: circular import C
from .a import func_a


def func_c():
    return func_a()
