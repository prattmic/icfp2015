import datetime
import json
import sys

data = json.load(sys.stdin)

def to_date(d):
    datestring = d["createdAt"]

    return datetime.datetime.strptime(datestring, "%Y-%m-%dT%H:%M:%S.%fZ")

l = sorted(data, key=to_date, reverse=True)

json.dump(l, sys.stdout, indent=4, sort_keys=True)
