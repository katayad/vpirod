import sys
import pandas
from nltk.tokenize import RegexpTokenizer
import mincemeat

chars_lines = pandas.read_csv("./southpark/All-seasons.csv")[['Character', 'Line']].values
token = RegexpTokenizer(r'\w+')
for line in chars_lines:
    line[1] = token.tokenize(line[1].lower())

def mapfn(key, value):
    yield value[0], set(value[1])

def reducefn(key, value):
    return len(set().union(*value))

server = mincemeat.Server()
server.mapfn = mapfn
server.reducefn = reducefn
server.datasource = dict(enumerate(chars_lines))

results = server.run_server(password='changeme')
with open('res.csv', 'w') as f:
    f.write('Character, UniqueWords\n')
    for name, words in results.iteritems():
        f.write(name + ',' + str(words) + '\n')
