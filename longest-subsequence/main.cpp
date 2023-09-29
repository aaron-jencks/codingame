#include <iostream>
#include <string>
#include <vector>
#include <algorithm>

using namespace std;

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

int main()
{
    int n;
    cin >> n; cin.ignore();
    cerr << n << endl;
    int* sequence = new int[n];
    int* P = new int[n];
    int* M = new int[n+1];
    for (int i = 0; i < n; i++) {
        cin >> sequence[i]; cin.ignore();
        M[i] = 0;
        cerr << sequence[i] << ' ';
    }
    cerr << endl;

    int L = 0;
    for(int i = 0; i < n; i++) {
        int lo = 1, hi = L + 1;
        while(lo < hi) {
            int mid = lo + ((hi-lo) >> 1);
            if(sequence[M[mid]] >= sequence[i]) hi = mid;
            else lo = mid + 1;
        }
        int newL = lo;
        P[i] = M[newL-1];
        M[newL] = i;

        if(newL > L) L = newL;
    }

    cout << L << endl;
}