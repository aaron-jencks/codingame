import sys
import math
from typing import Tuple, List

# TODO do this again, but 
# generate a distance map for each tower and then overlay the towers and min them
# could store and iterate over the coordinates instead of doing the whole grid

w = int(input())
print(w, file=sys.stderr)
h = int(input())
print(h, file=sys.stderr)
grid = []
distances = []
owner = []
towers = []
for i in range(h):
    line = input()
    print(line, file=sys.stderr)
    grid.append([c for c in line])
    distances.append([50]*w)
    owner.append([-1]*w)
    for ci, c in enumerate(line):
        if c != '.' and c != '#':
            tid = len(towers)
            towers.append((tid, c, ci, i))
            distances[i][ci] = 0
            owner[i][ci] = tid


def find_coord_neighbors(c: Tuple[int, int]) -> List[Tuple[int, int]]:
    result = []
    cx, cy = c
    if cx > 0:
        result.append((cx-1, cy))

    if cy > 0:
        result.append((cx, cy-1))

    if cy < h-1:
        result.append((cx, cy+1))

    if cx < w-1:
        result.append((cx+1, cy))

    return result


for ti, t in enumerate(towers):
    tid, tch, tx, ty = t
    q = [(0, tx, ty)]
    while len(q) > 0:
        dist, ex, ey = q.pop(0)

        neighbors = find_coord_neighbors((ex, ey))
        for nx, ny in neighbors:
            if grid[ny][nx] == '#':
                continue

            if distances[ny][nx] > dist+1:
                distances[ny][nx] = dist+1
                grid[ny][nx] = tch
                owner[ny][nx] = tid
                q.append((dist+1, nx, ny))
            elif distances[ny][nx] == dist+1 and owner[ny][nx] != tid:
                grid[ny][nx] = '+'
                owner[ny][nx] = len(towers)
                q.append((dist+1, nx, ny))


print('\n'.join(''.join(row) for row in grid))
