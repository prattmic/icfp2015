import json
import requests
import time

# CHANGE THIS
TAG_BASE = "fwkjefekwe"

# CHANGE THIS
commands = r"lllllllllllllalalalalalaakalppppppppllllllllllllllalalalalalaakalplllllllllllllalalalalalalaldpppppadddlllllllllllllalalalalalaladpabdallllllllllllalalalalalalaappdllllllllllllllalalalaalpdpllllllllllllllalalalalpddplldpalllllllllllalalalalakpppdadaapaappklpppddbbbdpppddbllllllllllllalalalalakpppdadaapaakblllllllllllllalalalalkpkpkbblllllllllllllalalaappapaapapakkakallllllllllllalalaapappaaldpdbdpdbdpdlllllllllllllalaaaaakkkkpppappapaapdplllllllllllllalalaaapppappappallllllllllllllalaapaakppapppkkbkppddllllllllllllllalaappaddallllllllllllllalaakaapppplllllllllllllalalappkbbkapkpkalllllllllllllalplkblllllllllllalaapalpppapappppkalllllllllllalappapappakpkkbkaplllllllllllllaapdalllllllllllladaallllllllaaaaaaalpkkalllllllllllppaaapappapalllllllllaaaaaadlllllllllaapklpappakalllllllllllappkappppppdaapkkallllllllllaapppkpkkpppppddbbbbdldlllllllllaapdadlllllllllllblalpldddllllllllllblddddlllllllaaakpkkbakllllllllldpdbdppddbbllllllalpapapapappalllllaaalpppapadlldpllllllaaldddllllllpapaaapddddllllaaappkbbkakpkppppdpalllllaappkpalllllaapalllkllbblbblblldadppaplllbldadbdllllbkppppapppddbbdapdalblpkakkbddlapaaddddllpdbl"

API_TOKEN = "3EJx4cAtwDXCrGc3TpcPUZ4TGAMBXHQNLrt53Y9j17o="

URL = "https://davar.icfpcontest.org/teams/296/solutions"

def send_command(command, tag):
    d = json.dumps([{
        "problemId": 2,
        "seed": 0,
        "tag": tag,
        "solution": command,
    }])

    headers = {'Content-Type': 'application/json'}

    r = requests.post(URL, data=d, headers=headers,
                      auth=('', API_TOKEN))
    print "Sent %s: %s %d %s" % (d, r, r.status_code, r.content)

def get_tag(tag):
    r = requests.get(URL, auth=('', API_TOKEN))
    d = r.json()

    for s in d:
        if s["tag"] == tag:
            return s
    return None

def wait_for_score(tag):
    score = None
    while True:
        t = get_tag(tag)
        print "Got tag: %s" % t
        score = t["score"]
        if score is not None:
            break

        time.sleep(60)

    return score

if __name__ == "__main__":
    tag_count = 0

    first = 0
    last = len(commands) - 1

    i = len(commands) / 2

    score = None
    while first <= last:
        i = (last + first) / 2
        print "first: %d, last: %d" % (first, last)
        print "i: %s" % i

        tag = TAG_BASE + str(tag_count)
        tag_count += 1
        print "tag: %s" % tag

        com = commands[:i]
        print "truncated commands: %s" % com

        send_command(com, tag)
        score = wait_for_score(tag)

        if score == 0:
            # Down
            last = i - 1
        else:
            # Up
            first = i + 1

    print "Done!"
    print "first: %d last: %d" % (first, last)
    print "score: %d" % score
