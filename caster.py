__author__ = 'Joe'
import requests
import json
import subprocess
from solution_mgr import read_server
from solution_mgr import Entry
from solution_mgr import TEAM_ID, address, API_TOKEN
from solution_mgr import latest_score

scores = read_server()
ai = 'chanterai'
program = './play_icfp2015'
# program = 'C:\Users\Joe\AppData\Local\Temp\server0go.exe'

def send_score(problemId, seed, tag, solution):
    ans = [{'problemId': problemId, 'seed': seed, 'tag': tag, 'solution': solution}]
    requests.post(address, auth=('', API_TOKEN), headers={'Content-type': 'application/json'}, data=json.dumps(ans))


cmd = '{} -ai {} -f qualifiers/problem_{}.json'
for problem_id in xrange(2):
    proc = subprocess.Popen(cmd.format(program, ai, problem_id), shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)

    raw_data = proc.communicate()
    solutions = json.loads(raw_data[0])
    logs = raw_data[1]
    for sub_sol in solutions:

        score = int(sub_sol['tag'].split(':')[1].strip())
        seed = int(sub_sol['seed'])
        problem_id = int(sub_sol['problemId'])
        tag = sub_sol['tag']

        latest = latest_score(scores[problem_id, seed])
        print 'latest', latest.score, 'this', score
        if score > latest.score:
            print 'Sending improved score of {}, for {}:{}'.format(score, problem_id, seed)
            send_score(problem_id, seed, tag, sub_sol['solution'])


