import matplotlib.pyplot as plt
import os

import sys
import json
import math
from operator import itemgetter

actual = {}
predicted = {}

fig, ax1 = plt.subplots()
ax2 = ax1.twinx()

target = sys.argv[1]
with open(target) as inp:
    data = inp.readlines()[1::]

    pk = None
    for i, line in enumerate(data):
        if i % 3 != 0:
            continue

        row = line.split()
        k, max_len, text_len, alphabet_size, patts, avg_active_states, data, errs = row
        k = int(k)
        patts = int(patts)
        alphabet_size = int(alphabet_size)

        if int(alphabet_size) != int(sys.argv[2]):
            continue

        if k not in actual:
            actual[k] = {'x': [], 'y': []}
            predicted[k] = {'x': [], 'y': []}

        errs = json.loads(errs)
        avg_active_states = [int(x) for x in json.loads(avg_active_states)][-50::]
        avg = sum(avg_active_states)/len(avg_active_states)
        pred = math.log(patts, alphabet_size)**k

        actual[k]['x'].append(patts)
        predicted[k]['x'].append(patts)

        actual[k]['y'].append(avg)
        predicted[k]['y'].append(pred)

    for k in actual:
        ax1.plot(actual[k]['x'], actual[k]['y'], linewidth=1,  label='k=%s' % k)
        ax2.plot(predicted[k]['x'], predicted[k]['y'], linewidth=1, color='#666666', linestyle='--')

    ax1.legend(fontsize=16)
    plt.savefig('../engrep-paper/paper/figures/e-%s.pdf' % sys.argv[2])
