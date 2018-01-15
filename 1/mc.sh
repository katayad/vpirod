#!/bin/bash

function f {
    x=$(echo "$1" | bc -l)
    echo "2*$x*$x" | bc -l
}

function rand {
    echo "$RANDOM/32767.0" | bc -l
}

sum=0
N=$1
for i in `seq 1 $N`;
do
    res=$(rand)
    res=$(f $res)
    sum=$(echo "$sum+$res" | bc -l)
done
echo "$sum/$N" | bc -l
