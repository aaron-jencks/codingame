#include <iostream>
#include <string>
#include <vector>
#include <algorithm>

using namespace std;

#define MAX(a,b) ((a) > (b) ? (a) : (b))
#define CRASHED(a,b,r) (MAX(abs(a),abs(b)) <= (r))
#define GRAVITY(v,p) ((v) + (((p)>0) ? -1 : ((p) < 0) ? 1 : 0))

typedef struct {
    int x;
    int y;
} coord_t;

typedef struct {
    coord_t pos;
    coord_t vel;
    bool crashed;
} state_t;

inline state_t step_velocity(state_t s) {
    s.pos.x += s.vel.x;
    s.pos.y += s.vel.y;
    s.vel.x = GRAVITY(s.vel.x, s.pos.x);
    s.vel.y = GRAVITY(s.vel.y, s.pos.y);
    return s;
}

int main()
{
    int radius;
    int x;
    int y;
    int vx;
    int vy;
    int time;
    cin >> radius >> x >> y >> vx >> vy >> time; cin.ignore();

    bool crashed = false;
    state_t current = state_t{coord_t{x, y}, coord_t{vx, vy}, false};
    for(int t = 0; t < time; t++) {
        current.crashed = CRASHED(current.pos.x, current.pos.y, radius);
        if(current.crashed) break;
        current = step_velocity(current);
    }

    cout << current.pos.x << ' ' << current.pos.y << ' ' << (int)current.crashed << endl;
}