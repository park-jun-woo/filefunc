# ff:func feature=validate type=util control=sequence
# ff:what test: no circular import C
import os


def func_c():
    return os.path.join("a", "b")
