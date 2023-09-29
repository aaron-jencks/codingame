#!/bin/bash

make referee
tdir="./tests/"
tc=1
while [ -f "$tdir""tc$tc.in.txt" ]; do
    tinf="$tdir""tc$tc.in.txt"
    tof="$tdir""tc$tc.out.txt"
    echo "testing $tdir""tc$tc"
    cat $tinf | ./referee > ./validator.out.txt
    if diff -u "$tof" ./validator.out.txt; then
        echo "pass"
    else
        :
        break
    fi
    ((tc+=1))
done

echo "all tests run successfully"
