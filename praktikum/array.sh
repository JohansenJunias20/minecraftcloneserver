#!/bin/bash

# a=(a a a b c)
# b=${a[@]:1:2}
genap=()
ganjil=()
echo $1
for((i=0;i<$1;i++))
do
    read a
    if [ $((a%2)) -eq 0 ]
    then
        genap=("${genap[@]}" $a)
        # sort genap
    fi
    else
done
echo ${genap[@]}
# echo ${b}

