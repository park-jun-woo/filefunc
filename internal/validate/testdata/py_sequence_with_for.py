# ff:func feature=validate type=rule control=sequence
# ff:what sequence file with for loop — should be A12 violation
def py_sequence_with_for(items):
    for item in items:
        print(item)
