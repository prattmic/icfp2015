#!/bin/bash
# submit.sh - How to run:
#     ./icfp -f qualifiers/problem_1.json | ./submit.sh

read OUTPUT
# MIRI Fan Club
TEAM_ID=296
# aray's auth token
API_TOKEN=3EJx4cAtwDXCrGc3TpcPUZ4TGAMBXHQNLrt53Y9j17o=
# SEND THE THING
curl --user :$API_TOKEN -X POST -H "Content-Type: application/json" \
    -d $OUTPUT https://davar.icfpcontest.org/teams/$TEAM_ID/solutions
