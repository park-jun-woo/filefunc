# ff:func feature=validate type=rule control=sequence
# ff:what deeply nested function for Q1 test
def py_deep_nesting():
    if True:
        if True:
            if True:
                return 1
    return 0
