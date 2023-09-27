#include <iostream>
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

struct path_return_t {
    path_t best_path;
    bool found;
};

path_return_t find_path_rec(path_t path_in, int16_t** energy_map, uint8_t col, uint8_t row, uint8_t h, uint8_t w) {
    path_t base_path = path_in;
    base_path.path.push_back(col);
    int16_t prev_energy = energy_map[row-1][col];
    base_path.total_energy += prev_energy;

    if(row == h) {
        // base case
        return {base_path, true};
    }

    path_return_t cpath{path_t{}, false};

    // TODO memoization of later paths
    if(col > 0) {
        int16_t lenergy = energy_map[row][col-1];
        if(abs(lenergy-prev_energy) <= 1) {
            cpath = find_path_rec(base_path, energy_map, col-1, row+1, h, w);
        }
    }
    int16_t cenergy = energy_map[row][col];
    if(abs(cenergy-prev_energy) <= 1) {
        path_return_t npath = find_path_rec(base_path, energy_map, col-1, row+1, h, w);
        if(cpath.found) {
            if(npath.found && npath.best_path < cpath.best_path) {
            }
        } else {
            cpath = npath;
        }
    }
    int16_t renergy = energy_map[row][col+1];
    if(abs(renergy-prev_energy) <= 1) {
        path_return_t npath = find_path_rec(base_path, energy_map, col+1, row+1, h, w);
        if(cpath.found) {
            if(npath.found && npath.best_path < cpath.best_path) {
                cpath = npath;
            }
        } else {
            cpath = npath;
        }
    }

    return cpath;
}

path_return_t find_path(int16_t** energy_map, uint8_t col, uint8_t h, uint8_t w) {
    return find_path_rec(path_t{vector<uint8_t>(), 0}, energy_map, col, 1, h, w);
}

struct heatmap_t {
    image_t image;
    int16_t** energies;
    vector<path_t> paths;
};

heatmap_t generate_heatmap(image_t image) {
    heatmap_t hm = {image, new int16_t*[image.h]};
    for(uint8_t i = 0; i < image.h; i++) {
        hm.energies[i] = new int16_t[image.w];
        for(uint8_t j = 0; j < image.w; j++) {
            hm.energies[i][j] = find_coord_energy(j, i, image);
        }
    }
    for(uint8_t c = 0; c < image.w; c++) {
        path_return_t pp = find_path(hm.energies, c, image.h, image.w);
        if(pp.found && pp.best_path.valid(image.h)) hm.paths.push_back(pp.best_path);
    }
    return hm;
}

struct seam_cut_return_t {
    image_t img;
    heatmap_t hm;
    int16_t path_energy;
    bool found;
};

seam_cut_return_t perform_seam_cut(image_t img, heatmap_t hm) {
    if(!hm.paths.size() || !img.w) return {img, hm, 0, false};
    path_t best_path;
    bool first = true;
    for(auto p: hm.paths) {
        if(first) {
            best_path = p;
            first = false;
        } else if(p < best_path) best_path = p;
    }
    // TODO add 2d array of booleans that toggle a pixel as having been cropped or not
    // that way we don't have to actually remove them from the image
    for(uint8_t row = 0; row < img.h; row++) {
        img.remove_pixel(best_path[row], row);
    }
    img.w--;
    // TODO perform optimal regeneration of heatmap
    // based on width of chosen path and width of current paths
    return {img, hm, best_path.total_energy, true};
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

    cerr << w << ' ' << h << endl;
    cerr << v << endl;
    cerr << maxintensity << endl;

    uint8_t** image = new uint8_t*[h];
    for (int i = 0; i < h; i++) {
        image[i] = new uint8_t[w];
        for (int j = 0; j < w; j++) {
            cin >> image[i][j]; cin.ignore();
        }
    }

    image_t img_state = {(uint8_t)h, (uint8_t)w, image};
    heatmap_t hm = generate_heatmap(img_state);

    while(img_state.w > v) {
        seam_cut_return_t cut = perform_seam_cut(img_state, hm);
        cout << cut.path_energy << endl;
        img_state = cut.img;
        hm = generate_heatmap(img_state);
    }
}