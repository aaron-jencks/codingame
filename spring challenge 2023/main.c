#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <stdbool.h>
#include <stdarg.h>

void debugf(char* fmt, ...) {
    va_list args;
    va_start(args, fmt);
    vfprintf(stderr, fmt, args);
    va_end(args);
}

typedef struct {
    int type;
    int count;
    int my_ants;
    int opp_ants;
} cell_t;

void display_cell(cell_t c) {
    debugf("Cell: Type: %s, Remaining: %d, My Ants: %d, Opponent Ants: %d\n",
        c.type ? "Crystal" : "Empty", c.count, c.my_ants, c.opp_ants);
}

typedef struct {
    int ncells;
    cell_t* cells;
    int bcount;
    int* friendly_bases;
    int* opponent_bases;
} state_t;

void display_state(state_t s) {
    debugf("State: (%d cells, %d bases)\nCells:\n", s.ncells, s.bcount);
    for(size_t i = 0; i < s.ncells; i++) {
        debugf("\t");
        display_cell(s.cells[i]);
    }
    debugf("Bases:\n\tFriendly:");
    for(size_t i = 0; i < s.bcount; i++) {
        debugf(" %d", s.friendly_bases[i]);
    }
    debugf("\n\tOpponent:");
    for(size_t i = 0; i < s.bcount; i++) {
        debugf(" %d", s.opponent_bases[i]);
    }
    debugf("\n");
}

int main()
{
    // amount of hexagonal cells in this map
    state_t current_state;
    
    scanf("%d", &current_state.ncells);

    current_state.cells = (cell_t*)malloc(sizeof(cell_t)*current_state.ncells);
    for (int i = 0; i < current_state.ncells; i++) {
        // the index of the neighbouring cell for each direction
        int neigh_0;
        int neigh_1;
        int neigh_2;
        int neigh_3;
        int neigh_4;
        int neigh_5;
        scanf("%d%d%d%d%d%d%d%d", 
            &current_state.cells[i].type, &current_state.cells[i].count, 
            &neigh_0, &neigh_1, &neigh_2, &neigh_3, &neigh_4, &neigh_5);
        current_state.cells[i].my_ants = 0;
        current_state.cells[i].opp_ants = 0;
    }
    scanf("%d", &current_state.bcount);
    current_state.friendly_bases = (int*)malloc(sizeof(int) * current_state.bcount);
    current_state.opponent_bases = (int*)malloc(sizeof(int) * current_state.bcount);

    for (int i = 0; i < current_state.bcount; i++) {
        scanf("%d", &current_state.friendly_bases[i]);
    }
    for (int i = 0; i < current_state.bcount; i++) {
        scanf("%d", &current_state.opponent_bases[i]);
    }

    // game loop
    while (1) {
        for (int i = 0; i < current_state.ncells; i++) {
            scanf("%d%d%d", 
                &current_state.cells[i].count, 
                &current_state.cells[i].my_ants, 
                &current_state.cells[i].opp_ants);
        }

        display_state(current_state);

        int* cell_throughput = (int*)malloc(sizeof(int) * current_state.ncells);
        for(size_t i = 0; i < current_state.ncells; i++) {
            
        }


        // WAIT | LINE <sourceIdx> <targetIdx> <strength> | BEACON <cellIdx> <strength> | MESSAGE <text>
        printf("WAIT\n");
    }

    return 0;
}