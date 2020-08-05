import toml


class Struct:
    def __init__(self, **entries):
        self.__dict__.update(entries)


config = toml.load("pwgo.toml")


chapters = Struct(**config)
a = []
for method in dir(chapters):
    if not method.startswith("__"):
        print(method)
