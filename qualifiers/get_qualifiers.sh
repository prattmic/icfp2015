#!/bin/bash

for i in {0..23}; do
	wget http://icfpcontest.org/problems/problem_$i.json
done
