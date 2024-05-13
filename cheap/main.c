#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <stdbool.h>

typedef struct {
    char* category;
    char* size;
    float price;
} item_t;

item_t* create_item(char* category, char* size, float price) {
    item_t* result = (item_t*)malloc(sizeof(item_t));
    result->category = category;
    result->size = size;
    result->price = price;
    return result;
}

typedef struct {
    char* category;
    char* size;
} order_t;

typedef bool (*less_t)(void*, void*);

typedef struct {
    void** arr;
    size_t size;
    size_t count;
    less_t less;
} arraylist_t;

arraylist_t create_arraylist(size_t initial_size, less_t less) {
    arraylist_t result = {
        malloc(sizeof(void*) * initial_size),
        initial_size,
        0,
        less
    };
    return result;
}

void arraylist_append(arraylist_t* arr, void* element) {
    if(arr->count+1 >= arr->size) {
        arr->size <<= 1;
        arr->arr = realloc(arr->arr, sizeof(void*) * arr->size);
    }
    arr->arr[arr->count++] = element;
}

void arraylist_insert(arraylist_t* arr, size_t index, void* element) {
    if(arr->count+1 >= arr->size) {
        arr->size <<= 1;
        arr->arr = realloc(arr->arr, sizeof(void*) * arr->size);
    }
    if(index == arr->count) {
        arraylist_append(arr, element);
        return;
    }
    for(size_t i = arr->count-1; i >= index && i; i--) {
        arr->arr[i+1] = arr->arr[i];
    }
    if(!index) arr->arr[1] = arr->arr[0];
    arr->arr[index] = element;
    arr->count++;
}

void arraylist_sorted_insert(arraylist_t* arr, void* element) {
    for(size_t i = 0; i < arr->count; i++) {
        if(arr->less(element, arr->arr[i])) {
            arraylist_insert(arr, i, element);
            return;
        }
    }
    arraylist_insert(arr, arr->count, element);
}

void* arraylist_delete(arraylist_t* arr, size_t index) {
    if(arr->count <= index) return NULL;
    void* result = arr->arr[index];
    for(size_t i = index+1; i <= arr->count; i++) {
        arr->arr[i-1] = arr->arr[i];
    }
    arr->count--;
    return result;
}

bool item_less(void* a, void* b) {
    return ((item_t*)a)->price < ((item_t*)b)->price;
}

typedef struct {
    char* key;
    arraylist_t value;
} kvp_t;

kvp_t* create_kvp(item_t* item) {
    kvp_t* kvp = (kvp_t*)malloc(sizeof(kvp_t));
    kvp->key = item->category;
    kvp->value = create_arraylist(100, item_less);
    return kvp;
}

typedef struct {
    arraylist_t* bins;
    size_t bcount;
} hashmap_t;

size_t hash_string(const char* s) {
    size_t result = 0;
    const int p = 31;
    const int m = 1e9 + 9;
    size_t p_pow = 1;
    size_t slen = strlen(s);
    for(size_t i = 0; i < slen; i++) {
        result += s[i] * p_pow;
        p_pow = p_pow * p;
    }
    return result % m;
}

hashmap_t create_hashmap(size_t bins) {
    hashmap_t result = {
        malloc(sizeof(arraylist_t) * bins),
        bins
    };
    for(size_t bi = 0; bi < bins; bi++) {
        result.bins[bi] = create_arraylist(100, NULL);
    }
    return result;
}

kvp_t* hashmap_get(hashmap_t hm, char* value) {
    size_t hash = hash_string(value);
    size_t bin = hash % hm.bcount;
    arraylist_t barr = hm.bins[bin];
    if(!barr.count) return NULL;
    for(size_t bi = 0; bi < barr.count; bi++) {
        if(!strcmp(((kvp_t*)barr.arr[bi])->key, value)) {
            return (kvp_t*)barr.arr[bi];
        }
    }
    return NULL;
}

kvp_t* hashmap_get_or_create_kvp(hashmap_t hm, item_t* value) {
    kvp_t* result = hashmap_get(hm, value->category);
    if(result) return result;

    size_t hash = hash_string(value->category);
    size_t bin = hash % hm.bcount;

    fprintf(stderr, "item hashed to %zu and placed in %zu\n", hash, bin);

    kvp_t* new_value = create_kvp(value);

    arraylist_append(&hm.bins[bin], new_value);

    return new_value;
}

void hashmap_put(hashmap_t hm, item_t* value) {
    kvp_t* kvp = hashmap_get_or_create_kvp(hm, value);
    arraylist_sorted_insert(&kvp->value, value);
}

arraylist_t split_string(char* s) {
    arraylist_t result = create_arraylist(3, NULL);
    char* tok = strtok(s, " ");
    while(tok) {
        arraylist_append(&result, tok);
        tok = strtok(NULL, " ");
    }
    return result;
}

int main()
{
    hashmap_t inventory = create_hashmap(1000);

    int c;
    scanf("%d", &c);
    int p;
    scanf("%d", &p); fgetc(stdin);

    fprintf(stderr, "%d\n%d\n", c, p);

    for (int i = 0; i < c; i++) {
        char item[101];
        scanf("%[^\n]", item); fgetc(stdin);

        fprintf(stderr, "%s\n", item);

        item_t* item_s = (item_t*)malloc(sizeof(item_t));
        arraylist_t bits = split_string(item);
        size_t clen = strlen(bits.arr[0]), slen = strlen(bits.arr[1]);

        item_s->category = malloc(sizeof(char) * (clen + 1));
        item_s->category[clen] = 0;
        memcpy(item_s->category, bits.arr[0], sizeof(char) * clen);

        item_s->size = malloc(sizeof(char) * (slen + 1));
        item_s->size[slen] = 0;
        memcpy(item_s->size, bits.arr[1], sizeof(char) * slen);
        
        sscanf(bits.arr[2], "%f", &item_s->price);

        hashmap_put(inventory, item_s);
    }

    for (int i = 0; i < p; i++) {
        char order[101];
        scanf("%[^\n]", order); fgetc(stdin);

        fprintf(stderr, "%s\n", order);

        arraylist_t bits = split_string(order);

        kvp_t* value = hashmap_get(inventory, bits.arr[0]);
        if(!value) {
            printf("NONE\n");
            continue;
        }
        bool found = false;
        for(size_t i = 0; i < value->value.count; i++) {
            item_t* item = value->value.arr[i];
            if(!strcmp(item->category, bits.arr[0]) && !strcmp(item->size, bits.arr[1])) {
                arraylist_delete(&value->value, i);
                printf("%d\n", (int)item->price);
                found = true;
                break;
            }
        }
        if(!found) {
            printf("NONE\n");
        }
    }

    return 0;
}