#include <iostream>
#include <string>
#include <vector>
#include <algorithm>

using namespace std;

const char EMPTY = '.';
const string VALID = "ABCDEFGHIJKLMNOP";
const int SIZE = 16;
const int BLOCK_SIZE = SIZE >> 2;

typedef struct {
    int row;
    int col;
} coord_t;

class openings {
    public:
    int board_size;
    int* rows;
    int* cols;
    int** quad;
    vector<coord_t> coords;

    openings(int size) {
        board_size = size;
        rows = new int[size];
        cols = new int[size];
        int qsize = size >> 2;
        quad = new int*[qsize];
        for(int i = 0; i < qsize; i++) {
            quad[i] = new int[qsize];
        }
        coords = vector<coord_t>();
    }

    ~openings() {
        delete[] rows;
        delete[] cols;
        int qsize = board_size >> 2;
        for(int i = 0; i < qsize; i++) {
            delete[] quad[i];
        }
        delete[] quad;
    }
};

int main()
{
    for (int i = 0; i < 16; i++) {
        string row;
        getline(cin, row);
    }

    // Write an answer using cout. DON'T FORGET THE "<< endl"
    // To debug: cerr << "Debug messages..." << endl;

    cout << "answer" << endl;
}