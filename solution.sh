#!/bin/bash
# solution.sh - Download our solutions

# MIRI Fan Club
TEAM_ID=296
# aray's auth token
API_TOKEN=3EJx4cAtwDXCrGc3TpcPUZ4TGAMBXHQNLrt53Y9j17o=
curl --user :$API_TOKEN https://davar.icfpcontest.org/teams/$TEAM_ID/solutions \
    | python -m json.tool | less
