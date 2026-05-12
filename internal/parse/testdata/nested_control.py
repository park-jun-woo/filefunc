def process(items):
    for item in items:
        if item > 0:
            for sub in range(item):
                print(sub)
