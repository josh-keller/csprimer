#include <assert.h>
#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define STARTING_CAPACITY 8

typedef struct DA {
  size_t size;
  size_t cap;
  void** arr;
} DA;

DA *DA_new(void) {
  DA *da = (DA *)malloc(sizeof(DA));
  da->size = 0;
  da->arr = malloc(STARTING_CAPACITY * sizeof(void *));
  da->cap = STARTING_CAPACITY;
  return da;
}

int DA_size(DA *da) {
  return da->size;
}

// void DA_grow_cap(DA *da) {
//   da->cap *= 2;
//   void **new_arr = malloc(da->cap * sizeof(void *));
//   memcpy(new_arr, da->arr, sizeof(void *) * da->size);
//   free(da->arr);
//   da->arr = new_arr;
// }

void DA_grow_cap(DA *da) {
  da->cap *= 2;
  da->arr = realloc(da->arr, da->cap * sizeof(void *));
}

void DA_push(DA *da, void *x) {
  if (da->size == da->cap) {
    DA_grow_cap(da);
  }
  da->arr[da->size++] = x;
}

void *DA_pop(DA *da) {
  if (da->size == 0) {
    return NULL;
  }

  da->size--;
  void *ret = da->arr[da->size];
  da->arr[da->size] = NULL;
  return ret;
}

void DA_set(DA *da, void *x, int i) {
  if (i >= da->size || i < 0) {
    return;
  }

  da->arr[i] = x;
}

void *DA_get(DA *da, int i) {
  if (i >= da->size || i < 0) {
    return NULL;
  }
  return da->arr[i];
}

void DA_free(DA *da) {
  free(da->arr);
  free(da);
}

int main() {
  DA *da = DA_new();

  assert(DA_size(da) == 0);

  // basic push and pop test
  int x = 5;
  float y = 12.4;
  DA_push(da, &x);
  DA_push(da, &y);
  assert(DA_size(da) == 2);

  assert(DA_pop(da) == &y);
  assert(DA_size(da) == 1);

  assert(DA_pop(da) == &x);
  assert(DA_size(da) == 0);
  assert(DA_pop(da) == NULL);

  // basic set/get test
  DA_push(da, &x);
  DA_set(da, &y, 0);
  assert(DA_get(da, 0) == &y);
  DA_pop(da);
  assert(DA_size(da) == 0);

  // expansion test
  DA *da2 = DA_new(); // use another DA to show it doesn't get overriden
  DA_push(da2, &x);
  int i, n = 100 * STARTING_CAPACITY, arr[n];
  for (i = 0; i < n; i++) {
    arr[i] = i;
    DA_push(da, &arr[i]);
  }
  assert(DA_size(da) == n);
  for (i = 0; i < n; i++) {
    assert(DA_get(da, i) == &arr[i]);
  }
  for (; n; n--)
    DA_pop(da);
  assert(DA_size(da) == 0);
  assert(DA_pop(da2) == &x); // this will fail if da doesn't expand

  DA_free(da);
  DA_free(da2);
  printf("OK\n");
}
