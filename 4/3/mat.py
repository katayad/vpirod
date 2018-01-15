import sys
import pandas
from nltk.tokenize import RegexpTokenizer
import mincemeat
import os
import numpy as np

n = 3
inpt = pandas.read_csv("mat_a.csv")[['matrix', 'row', 'col', 'val']].values
inpt = np.concatenate((inpt, pandas.read_csv("mat_b.csv")[['matrix', 'row', 'col', 'val']].values), axis=0)

mats = []
for row in inpt:
    mats.append(list(row.tolist()) + [n])
print(mats)
def mapfn(_, row):
    if row[0] == "a":
        for i in range(row[-1]):
            yield (row[1], i), (row[2], row[3])
    else:
        for i in range(row[-1]):
            yield (i, row[2]), (row[1], row[3])

def reducefn(key, vals):
    sm = 0
    used = {}
    for val in vals:
        if val[0] in used:
            sm += used[val[0]] * val[1]
        else:
            used[val[0]] = val[1]
    return sm % 97

server = mincemeat.Server()
server.mapfn = mapfn
server.reducefn = reducefn
server.datasource = dict(enumerate(mats))

results = server.run_server(password='changeme')

with open("res.csv","w") as f:
    f.write('matrix,row,col,val\n')
    for key, value in results.iteritems():
        f.write("c,%d,%d,%d\n" % (key[0], key[1], value))
