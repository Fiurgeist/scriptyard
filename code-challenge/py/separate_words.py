import collections


# task: separate a string of unseparated words given a dictionary of valid words
# constrain: input is a valid sentence
# time: 20min
def parse(data, dictionary):
    arr = list(data)

    return flatten(_parse(arr, dictionary, 0))


def _parse(arr, dictionary, i):
    if i >= len(arr):
        return []
    cur = ""
    for c in arr[i:]:
        cur += c
        if cur in dictionary:
            ret = _parse(arr, dictionary, i + len(cur))
            if ret is not None:
                return [cur, ret]


def flatten(x):
    result = []
    if x is None:
        return result
    for el in x:
        if isinstance(x, collections.Iterable) and not isinstance(el, str):
            result.extend(flatten(el))
        else:
            result.append(el)
    return result


print(parse("thebesttest", ["best", "test", "the"]))
print(parse("therecyclebike", ["cycle", "there", "the", "cyclebike", "recycle", "bike"]))
print(parse("therecycledbike", ["there", "the", "recycle", "cycled", "recycled", "bike"]))
