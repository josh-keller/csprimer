class MyRange(object):
    def __init__(self, max):
        self.curr = 0
        self.max = max

    def __iter__(self):
        return self

    def __next__(self):
        if self.curr >= self.max:
            raise StopIteration

        curr = self.curr
        self.curr += 1
        return curr



for x in MyRange(10):
    print(x)
