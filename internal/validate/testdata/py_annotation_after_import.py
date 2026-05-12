import os

# ff:func feature=validate type=rule control=sequence
# ff:what this annotation is after import — should be A6 violation
def py_annotation_after_import():
    return os.getcwd()
