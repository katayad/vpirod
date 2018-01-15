#!/bin/bash
c_on_p=$(echo "$1/$2" | bc)

echo "" > params.txt

for i in $(seq $2); do
	echo $c_on_p >> params.txt
done

if [ -z $3 ]
	then
		res=$(parallel -a params.txt ./mc.sh)
	else
		scp mc.sh $3:~/
		res=$(parallel -a params.txt -S $3 ~/mc.sh)
fi

sum=$(echo "scale=10; (" $(echo $res | tr ' ' '+') ") / $2" | bc)

echo $sum
