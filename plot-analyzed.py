import matplotlib.pyplot as plt
import os

import sys
import json
from operator import itemgetter

patterns = {}


target = sys.argv[1]
with open(target) as inp:
    data = inp.readlines()[1::]
    for line in data:
        row = line.split()
        k, max_len, text_len, alphabet_size, patts, avg_active_states, data = row
        if int(alphabet_size) != 4 or int(text_len) != 100:
            continue

        if int(k) != 1:
            continue

        if int(patts) not in [100, 1000]:
            continue

        sprobs = json.loads(data)
        probs = {}
        for key, value in sprobs.items():
            probs[int(key)] = [int(value[0]), int(value[1])]

        x = []
        y = []
        for depth,(hit, miss) in sorted(probs.items(), key=itemgetter(0)):

            x.append(int(depth))
            y.append(float(hit)/(hit + miss))

        plt.plot(x, y, linewidth=1,  label='k=%s, |P|=%s' % (k, patts))

    plt.ylim(0, 1)
    plt.legend()
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
