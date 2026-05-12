# ff:func feature=validate type=rule control=selection
# ff:what selection file with no match statement
def py_selection_no_match(x):
    if x > 0:
        return "positive"
    return "non-positive"
