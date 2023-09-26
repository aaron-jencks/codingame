#include <iostream>
#include <string>
#include <vector>
#include <algorithm>
#include <unordered_map>
#include <stdint.h>

using namespace std;

#define MAX(a,b) ((a) > (b) ? (a) : (b))
#define CRASHED(a,b,r) (MAX(abs(a),abs(b)) <= (r))
#define GRAVITY(v,p) ((v) + (((p)>0) ? -1 : ((p) < 0) ? 1 : 0))

// largest known planet orbit distance is 2048, 12 bits

struct coord_t {
    int16_t x;
    int16_t y;

    bool operator==(const coord_t& c) const {
        return x == c.x && y == c.y;
    }
    bool operator!=(const coord_t& c) const {
        return x != c.x || y != c.y;
    }

    uint32_t hash() const {
        return (((uint32_t)x) << 16) + ((uint32_t)y);
    }
};

ostream& operator << (ostream& out, const coord_t& c) {
    return out << "(" << c.x << ", " << c.y << ")";
}

struct state_t {
    coord_t pos;
    coord_t vel;
    bool crashed;
    bool operator==(const state_t& s) const {
        return s.pos == pos && s.vel == vel && s.crashed == crashed;
    }
    bool operator!=(const state_t& s) const {
        return s.pos != pos || s.vel != vel || s.crashed != crashed;
    }
};

ostream& operator << (ostream& out, const state_t& c) {
    return out << c.pos << " going " << c.vel << " and has " << (c.crashed ? "" : "not ") << "crashed";
}

class StateHashFunction {
    public:
    uint64_t operator()(const state_t& s) const {
        return (((uint64_t)s.pos.hash()) << 32) + ((uint64_t)s.vel.hash());
    }
};

inline state_t step_velocity(state_t s) {
    s.pos.x += s.vel.x;
    s.pos.y += s.vel.y;
    s.vel.x = GRAVITY(s.vel.x, s.pos.x);
    s.vel.y = GRAVITY(s.vel.y, s.pos.y);
    return s;
}

vector<state_t> find_loop(state_t start, unordered_map<state_t, state_t, StateHashFunction> map) {
    vector<state_t> result;
    state_t current = start;
    do {
        result.push_back(current);
        current = map[current];
    } while(current != start);
    return result;
}

int main()
{
    int radius;
    int16_t x;
    int16_t y;
    int16_t vx;
    int16_t vy;
    int time;
    cin >> radius >> x >> y >> vx >> vy >> time; cin.ignore();

    unordered_map<state_t, state_t, StateHashFunction> visited;
    state_t previous;
    state_t current = state_t{coord_t{x, y}, coord_t{vx, vy}, false};
    for(int t = 0; t < time; t++) {
        // cerr << current << endl;
        previous = current;
        current = step_velocity(current);
        visited.insert({previous, current});
        if(visited.find(current) != visited.end()) {
            // we're in a loop
            cerr << "we found a loop" << endl;
            vector<state_t> loop = find_loop(current, visited);
            cerr << "the loop has " << loop.size() << " nodes" << endl;
            int tleft = time - (t + 1);
            cerr << "there are " << tleft << " seconds left in the orbit" << endl;
            int final_index = tleft % loop.size();
            state_t final_state = loop[final_index];
            cerr << "final state would be " << final_state << endl;
            current = final_state;
            break;
        }
        current.crashed = CRASHED(current.pos.x, current.pos.y, radius);
        if(current.crashed) break;
    }

    cout << current.pos.x << ' ' << current.pos.y << ' ' << (int)current.crashed << endl;
}