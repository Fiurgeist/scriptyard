import random

# quick and dirty 3d20 rolls with skill value and difficulty modifier (DSA rule set)
def test_3d20(attrs, skill=0, difficulty=0, iterations=10000, seed=42):
    random.seed(seed)
    attrs_copy = attrs.copy()
    attrs_copy.sort()

    e_skill = skill - difficulty

    succ_row = 0
    succ_sort = 0
    for _ in range(iterations):
        results = [
            random.randint(1, 20),
            random.randint(1, 20),
            random.randint(1, 20),
        ]
        r_copy = results.copy()
        r_copy.sort()
        row = True
        sort = True
        if e_skill < 0:
            for idx, a in enumerate(attrs):
                if a + e_skill < results[idx]:
                    row = False
                    break

            for idx, a in enumerate(attrs_copy):
                if a + e_skill < r_copy[idx]:
                    sort = False
                    break
        else:
            e_skill_copy = e_skill
            for idx, a in enumerate(attrs):
                diff = results[idx] - a
                if diff > 0:
                    e_skill_copy -= diff
                    if e_skill_copy < 0:
                        row = False
                        break
            e_skill_copy = e_skill
            for idx, a in enumerate(attrs_copy):
                diff = r_copy[idx] - a
                if diff > 0:
                    e_skill_copy -= diff
                    if e_skill_copy < 0:
                        sort = False
                        break
        if row:
            succ_row += 1
        if sort:
            succ_sort += 1

    print("row: %s" % (succ_row/iterations))
    print("sort: %s" % (succ_sort/iterations))
