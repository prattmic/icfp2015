import argparse

__author__ = 'Joe'

E = ('b', 'c', 'e', 'f', 'y', '2')
W = ('p', "'", '!', '.', '0', '3')
SE = ('a', 'g', 'h', 'i', 'j', '4')
SW = ('l', 'm', 'n', 'o', ' ', '5')
CW = ('d', 'q', 'r', 'v', 'z', '1')
CCW = ('k', 's', 't', 'u', 'w', 'x')


class West(object):
    pass


class East(object):
    pass


class Southeast(object):
    pass


class Southwest(object):
    pass


class Clockwise(object):
    pass


class Counterclockwise(object):
    pass

moves_2_letters = {
    West: W,
    East: E,
    Southeast: SE,
    Southwest: SW,
    Clockwise: CW,
    Counterclockwise: CCW,
}


def valid_power(moves):
    if len(moves) > 51:
        return False

    for i in xrange(len(moves)-1):
        if type(moves[i]) is West and type(moves[i+1]) is East:
            return False
        if type(moves[i]) is East and type(moves[i+1]) is West:
            return False
        if type(moves[i]) is Clockwise and type(moves[i+1]) is Counterclockwise:
            return False
        if type(moves[i]) is Counterclockwise and type(moves[i+1]) is Clockwise:
            return False

    return True


def read_word(raw_word):
    moves = []
    for letter in raw_word:
        new_move = None
        for key, value in moves_2_letters.iteritems():
            if letter in value:
                new_move = key()
        if new_move is None:
            raise Exception
        moves.append(new_move)
    return moves


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Power word analyzer')
    parser.add_argument('-test_proposed', action="store_true", default=False)
    parser.add_argument('-list_moves_confirmed', action="store_true", default=False)
    args = parser.parse_args()

    if args.test_proposed:
        with open('proposed') as f:
            for line in f:
                line = line.strip().lower()
                print '{} : {}'.format(line, valid_power(line))

    if args.list_moves_confirmed:
        with open('confirmed') as f:
            for line in f:
                line = line.strip().lower()
                print line
                for move in read_word(line):
                    print type(move)

