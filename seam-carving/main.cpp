#include <iostream>
#include <iomanip>
#include <string>
#include <vector>
#include <algorithm>
#include <stdint.h>

using namespace std;

struct image_t {
    uint8_t h;
    uint8_t w;
    uint8_t** image;

    uint8_t remove_pixel(uint8_t x, uint8_t y) {
        uint8_t v = image[y][x];
        for(uint8_t cx = x+1; cx < w; cx++) {
            image[y][cx-1] = image[y][cx];
        }
        return v;
    }

    inline uint8_t* operator [](int row) { return image[row]; }
};

int16_t find_coord_energy(uint8_t x, uint8_t y, image_t image) {
    int16_t dx = (0 < x && x < image.w-1) ? ((int16_t)image.image[y][x+1]) - ((int16_t)image.image[y][x-1]) : 0;
    int16_t dy = (0 < y && y < image.h-1) ? ((int16_t)image.image[y+1][x]) - ((int16_t)image.image[y-1][x]) : 0;
    return abs(dx) + abs(dy);
}

struct path_t {
    vector<uint8_t> path;
    int16_t total_energy;

    inline bool valid(uint8_t h) {
        return path.size() == h;
    }

    inline bool operator < (const path_t& p) const {
        return total_energy < p.total_energy || (total_energy == p.total_energy && path[0] < p.path[0]);
    }

    inline uint8_t operator [] (int i) { return path[i]; }
};

ostream& operator << (ostream& out, const path_t& p) {
    out << "path(" << p.total_energy << "): ";
    for(auto pv: p.path) out << (int)pv << ", ";
    return out;
}

struct heatmap_t {
    image_t image;
    int16_t** energies;
    int16_t** path_map;
};

heatmap_t generate_heatmap(image_t image) {
    heatmap_t hm = {image, new int16_t*[image.h], new int16_t*[image.h]};
#ifdef DEBUG
    cerr << "energies:" << endl;
#endif
    for(uint8_t i = 0; i < image.h; i++) {
        hm.energies[i] = new int16_t[image.w];
        hm.path_map[i] = new int16_t[image.w];
        for(uint8_t j = 0; j < image.w; j++) {
            hm.energies[i][j] = find_coord_energy(j, i, image);
#ifdef DEBUG
            cerr << setw(3) << (int)hm.energies[i][j] << " ";
#endif
        }
#ifdef DEBUG
        cerr << endl;
#endif
    }
#ifdef DEBUG
    cerr << "path costs:" << endl;
#endif
    for(uint8_t col = 0; col < image.w; col++) {
        hm.path_map[image.h-1][col] = hm.energies[image.h-1][col];
#ifdef DEBUG
        cerr << setw(4) << hm.path_map[image.h-1][col] << ' ';
#endif
    }
#ifdef DEBUG
    cerr << endl;
#endif
    // memoize the energy map for use in pathing
    for(int16_t row = image.h-2; row >= 0; row--) {
        for(uint8_t col = 0; col < image.w; col++) {
            int16_t pv = 32767; // int16 positive max
            if(col > 0) pv = hm.path_map[row+1][col-1];
            int16_t tpv = hm.path_map[row+1][col];
            if(tpv < pv) pv = tpv;
            if(col < image.w-1) {
                tpv = hm.path_map[row+1][col+1];
                if(tpv < pv) pv = tpv;
            }
            hm.path_map[row][col] = hm.energies[row][col] + pv;
#ifdef DEBUG
            cerr << setw(4) << hm.path_map[row][col] << ' ';
#endif
        }
#ifdef DEBUG
        cerr << endl;
#endif
    }
    return hm;
}

path_t find_path(heatmap_t hm) {
    path_t result = {vector<uint8_t>(), 0};
    uint8_t current_col = 255;
    int16_t current_val = 32767;
    for(uint8_t c = 0; c < hm.image.w; c++) {
        if(hm.path_map[0][c] < current_val) {
            current_col = c;
            current_val = hm.path_map[0][c];
        }
    }
    result.path.push_back(current_col);
    result.total_energy = current_val;
    for(uint8_t h = 1; h < hm.image.h; h++) {
        uint8_t next_col;
        current_val = 32767;
        if(current_col > 0) {
            next_col = current_col-1;
            current_val = hm.path_map[h][current_col-1];
        }
        if(hm.path_map[h][current_col] < current_val) {
            next_col = current_col;
            current_val = hm.path_map[h][current_col];
        }
        if(current_col < hm.image.w-1 && hm.path_map[h][current_col+1] < current_val) {
            next_col = current_col+1;
            current_val = hm.path_map[h][current_col+1];
        }
        result.path.push_back(next_col);
        current_col = next_col;
    }
    return result;
}

struct seam_cut_return_t {
    image_t img;
    heatmap_t hm;
    int16_t path_energy;
};

seam_cut_return_t perform_seam_cut(image_t img, heatmap_t hm) {
    if(!img.w) return {img, hm, 0};
    path_t best_path = find_path(hm);
    // TODO add 2d array of booleans that toggle a pixel as having been cropped or not
    // that way we don't have to actually remove them from the image
    for(uint8_t row = 0; row < img.h; row++) {
        img.remove_pixel(best_path[row], row);
    }
    img.w--;
    // TODO perform optimal regeneration of heatmap
    // based on width of chosen path and width of current paths
    return {img, hm, best_path.total_energy};
}

int main()
{
    string magic;
    getline(cin, magic);
    int w;
    int h;
    cin >> w >> h; cin.ignore();
    string comment;
    int v;
    cin >> comment >> v; cin.ignore();
    int maxintensity;
    cin >> maxintensity; cin.ignore();

#ifdef DEBUG
    cerr << w << ' ' << h << endl;
#endif
#ifdef DEBUG
    cerr << v << endl;
#endif
#ifdef DEBUG
    cerr << maxintensity << endl;
#endif

    uint8_t** image = new uint8_t*[h];
    for (int i = 0; i < h; i++) {
        image[i] = new uint8_t[w];
        for (int j = 0; j < w; j++) {
            int value;
            cin >> value; cin.ignore();
            image[i][j] = (uint8_t)value;
#ifdef DEBUG
            cerr << setw(3) << (int)image[i][j] << ' ';
#endif
        }
#ifdef DEBUG
        cerr << endl;
#endif
    }

    image_t img_state = {(uint8_t)h, (uint8_t)w, image};
    heatmap_t hm = generate_heatmap(img_state);

    while(img_state.w > v) {
        seam_cut_return_t cut = perform_seam_cut(img_state, hm);
        cout << cut.path_energy << endl;
        if(cut.img.w == img_state.w) break;
        img_state = cut.img;
        hm = generate_heatmap(img_state);
    }
}