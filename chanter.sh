#/bin/bash

touch fullrun.txt
for i in `seq 0 24`;
do
    #echo "./play_icfp2015 -ai chanterai -f qualifiers/problem_${i}.json"
    ./play_icfp2015 -ai chanterai -f qualifiers/problem_${i}.json >> \
        fullrun.txt
done
