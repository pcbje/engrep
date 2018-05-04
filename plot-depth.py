import matplotlib.pyplot as plt
import os

import sys
import json
from operator import itemgetter

patterns = {}


target = sys.argv[1]
with open(target) as inp:
    data = inp.readlines()[1::]

    pk = None
    for i, line in enumerate(data):
        if i % 3 != 0:
            continue

        row = line.split()
        k, max_len, text_len, alphabet_size, patts, avg_active_states, data, errs = row

        if int(alphabet_size) != 2:
            continue

        if int(k) != int(sys.argv[2]):
            continue

        xdepths = json.loads(data)
        depths = {}
        for key, value in xdepths.items():
            depths[int(key)] = int(value[0])+  int(value[1])

        x = []
        y = []
        for depth, count in sorted(depths.items(), key=itemgetter(0)):
            x.append(depth)
            y.append(count)

        plt.plot(x, y, linewidth=1,  label='k=%s, |P|=%s' % (k, patts))

    plt.axvline(x=10, color='#cccccc', linestyle='--')
    plt.axvline(x=12.5, color='#cccccc', linestyle='--')
    plt.axvline(x=15, color='#cccccc', linestyle='--')
    plt.legend()
    plt.xlim(0, 20)
    plt.show()

"""
styles = {1: 'solid', 2: 'dashed', 3: 'dotted'}
for k in all_speed:
    x1 = []
    y1 = []

    for patterns, speed in sorted(all_speed[k].items()):
        x1.append(patterns)
        y1.append(speed)

    plt.plot(x1, y1, label='k=%s' % k, linestyle=styles[k], linewidth=1)

plt.ylim(0, 8.5)
plt.grid(axis='y', color='#cccccc')
plt.legend()
plt.savefig('../../paper/figures/speed-%s.pdf' % os.path.basename(target))
plt.close()

for k in all_memory:
    x1 = []
    y1 = []

    for patterns, memory in sorted(all_memory[k].items()):
        x1.append(patterns)
        y1.append(memory)

    plt.plot(x1, y1, label='k=%s' % k, linestyle=styles[k], linewidth=1)

plt.ylim(0, 2)
plt.grid(axis='y', color='#cccccc')
plt.legend()
plt.savefig('../../paper/figures/memory-%s.pdf' % os.path.basename(target))
"""
