import sys
import pandas
from nltk.tokenize import RegexpTokenizer
import mincemeat
import os

books_words = {}
token = RegexpTokenizer(r'\w+')
for book in os.listdir('./sherlock'):
    with open(os.path.join('./sherlock', book), 'r') as f:
        books_words[book] = token.tokenize(f.read().lower())

def mapfn(book, words):
    for word in words:
        yield word, book

def reducefn(word, books):
    res = {}
    for book in books:
        if book not in res:
            res[book] = 0
        res[book] += 1
    return res

server = mincemeat.Server()
server.mapfn = mapfn
server.reducefn = reducefn
server.datasource = books_words

results = server.run_server(password='changeme')
books_list = os.listdir('./sherlock')
with open('res.csv', 'w') as f:
    f.write('word,' + ','.join(books_list) + '\n')
    for word, books_dict in results.iteritems():
        f.write(word + ',' + ','.join([str(books_dict.get(book, 0)) for book in books_list]) + '\n')
