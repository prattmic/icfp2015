import datetime
import argparse
import json
import requests

# MIRI Fan Club
TEAM_ID = 296
# aray's auth token
API_TOKEN ='3EJx4cAtwDXCrGc3TpcPUZ4TGAMBXHQNLrt53Y9j17o='
address = 'https://davar.icfpcontest.org/teams/{}/solutions'.format(TEAM_ID)


class Entry(object):
    def __init__(self, json_dict):
        self._raw = json_dict
        self.score = json_dict['score']
        self.tag = json_dict['tag']
        self.team_id = json_dict['teamId']
        self.solution = json_dict['solution']
        self.power_score = json_dict['powerScore']
        self.author_id = json_dict['authorId']
        self.seed = json_dict['seed']
        self.problem_id = json_dict['problemId']
        self.created_at = datetime.datetime.strptime(json_dict['createdAt'][:-1], "%Y-%m-%dT%H:%M:%S.%f")

    def __cmp__(self, other):
        return cmp(self.score, other.score)

    def __str__(self):
        return 'score: {:6} tag: {} solution: {}'.format(
            self.score,
            self.tag,
            self.solution,
        )


def max_score(solutions):
    solutions.sort(reverse=True, key=lambda obj: obj.score)
    return solutions[0]


def latest_score(solutions):
    solutions.sort(reverse=True, key=lambda obj: obj.created_at)
    return solutions[0]


def print_max_scores(scoring):
    key_list = scoring.keys()
    key_list.sort()
    for key in key_list:
        a = max_score(scoring[key])
        print '{:2}:{:7} = {:5}'.format(key[0], key[1], a.score)


def print_latest_scores(scoring):
    key_list = scoring.keys()
    key_list.sort()
    for key in key_list:
        a = latest_score(scoring[key])
        print '{:2}:{:7} = {} @ {}'.format(key[0], key[1], a, a.created_at)


def print_suboptimal_scores(scoring):
    key_list = scoring.keys()
    key_list.sort()
    for key in key_list:
        a = latest_score(scoring[key])
        b = max_score(scoring[key])
        if a.score < b.score and a.score is not None:
            print '{:2}:{:7} = {:5} < {:5} {}'.format(key[0], key[1], a.score, b.score, b.tag)
            print '                          {}'.format(a)
            print '                          {}'.format(b)

def post_optimal_scores(scoring):
    for problem in scoring.keys():
        a = latest_score(scoring[problem])
        b = max_score(scoring[problem])
        if a.score < b.score:
            post_score(b)


def post_score(entry):
    assert isinstance(entry, Entry)
    formatted_ans = [
        {
            "problemId": entry.problem_id,
            "seed": entry.seed,
            "tag": entry.tag,
            "solution": entry.solution,
        }
    ]
    resp = requests.post(address, auth=('', API_TOKEN), headers={'Content-type': 'application/json'}, data=json.dumps(formatted_ans))


def read_server():
    resp = requests.get(address, auth=('', API_TOKEN))
    raw_scores = json.loads(resp.content)
    scores = {}
    for entry in raw_scores:
        newEntry = Entry(entry)
        if (newEntry.problem_id, newEntry.seed) not in scores:
            scores[newEntry.problem_id, newEntry.seed] = []
        scores[newEntry.problem_id, newEntry.seed].append(newEntry)
    return scores


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Solution Manager')
    parser.add_argument('-max_scores', action="store_true", default=False)
    parser.add_argument('-latest_scores', action="store_true", default=False)
    parser.add_argument('-sub_optimal_scores', action="store_true", default=False)
    parser.add_argument('-post_optimal_scores', action="store_true", default=False)
    args = parser.parse_args()

    scores = read_server()

    if args.max_scores:
        print_max_scores(scores)

    if args.latest_scores:
        print_latest_scores(scores)

    if args.sub_optimal_scores:
        print_suboptimal_scores(scores)

    if args.post_optimal_scores:
        post_optimal_scores(scores)


