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


def valid_power(word):
    if len(word) > 51:
        return False

    moves = []
    for letter in word:
        new_move = None
        for key, value in moves_2_letters.iteritems():
            if letter in value:
                new_move = key()
        if new_move is None:
            return False
        moves.append(new_move)

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


if __name__ == '__main__':
    with open('proposed') as f:
        for line in f:
            line = line.strip().lower()
            print '{} : {}'.format(line, valid_power(line))

