# ff:func feature=validate type=util control=sequence
# ff:what test: no circular import B
import os


def func_b():
    return os.getcwd()
