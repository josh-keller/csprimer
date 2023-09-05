#include <assert.h>
#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define STARTING_BUCKETS 8

unsigned long hash(char *str) {
  unsigned long hash = 5381;
  int c;

  while ((c = *str++)) {
    hash = ((hash << 5) + hash) + c;
  }

  return hash;
}

typedef struct Node *Node;

struct Node {
  char *key;
  void *val;
  Node next;
};

typedef struct {
  Node *buckets;
  size_t bucket_count;
} Hashmap;

Hashmap *Hashmap_new(void) {
  Hashmap *h = malloc(sizeof(Hashmap));
  h->bucket_count = STARTING_BUCKETS;

  h->buckets = calloc(STARTING_BUCKETS, sizeof(Node));
  // for (int i = 0; i < STARTING_BUCKETS; i++) {
  //   printf("Bucket %d: %p\n", i, h->buckets[i]);
  // }
  return h;
}

void Hashmap_free(Hashmap *h) {
  for (int i = 0; i < h->bucket_count; i++) {
    for (Node curr = h->buckets[i], next = NULL; curr != NULL; curr = next) {
      next = curr->next;
      free(curr->key);
      // Do I need to free val too? I'm thinking no because the caller might
      // retain a reference
      free(curr);
    }
  }
  free(h->buckets);
  free(h);
}

void Hashmap_set(Hashmap *h, char *key, void *val) {
  size_t key_hash = hash(key) % h->bucket_count;
  printf("Setting key %s in bucket %zu to value %p\n", key, key_hash, val);

  if (h->buckets[key_hash] == NULL) {
    printf("No bucket here, creating!\n");
    h->buckets[key_hash] = (Node)malloc(sizeof(Node));
    h->buckets[key_hash]->key = strdup(key);
    h->buckets[key_hash]->val = val;
    h->buckets[key_hash]->next = NULL;
    return;
  }

  Node next = h->buckets[key_hash];
  Node curr = NULL;
  
  while (next != NULL) {
    curr = next;
    next = curr->next;

    // printf("Bucket exists, walking Nodes!\n");
    // printf("Current key: %s.\n", curr->key);
    if (strcmp(curr->key, key) == 0) {
      // printf("Node exists with key");
      curr->val = val;
      return;
    }
  }

  curr->next = (Node)malloc(sizeof(Node));
  curr->next->key = strdup(key);
  curr->next->val = val;
  curr->next->next = NULL;
  return;
}

void *Hashmap_get(Hashmap *h, char *key) {
  size_t key_hash = hash(key) % h->bucket_count;
  // printf("Looking for key %s.\n", key);

  Node curr = h->buckets[key_hash];
  while (curr != NULL) {
    // printf("Current key: %s, val: %p.\n", curr->key, curr->val);
    if (strcmp(curr->key, key) == 0) {
      // printf("Match, returning!\n");
      return curr->val;
    }

    // printf("No match, next!\n");

    curr = curr->next;
  }

  return NULL;
}

int main() {
  Hashmap *h = Hashmap_new();

  // basic get/set functionality
  int a = 5;
  float b = 7.2;
  Hashmap_set(h, "item a", &a);
  Hashmap_set(h, "item b", &b);
  assert(Hashmap_get(h, "item a") == &a);
  assert(Hashmap_get(h, "item b") == &b);

  // using the same key should override the previous value
  int c = 20;
  Hashmap_set(h, "item a", &c);
  assert(Hashmap_get(h, "item a") == &c);

  // // basic delete functionality
  // Hashmap_delete(h, "item a");
  // assert(Hashmap_get(h, "item a") == NULL);
  //
  // // handle collisions correctly
  // // note: this doesn't necessarily test expansion
  // int i, n = STARTING_BUCKETS * 10, ns[n];
  // char key[MAX_KEY_SIZE];
  // for (i = 0; i < n; i++) {
  //   ns[i] = i;
  //   sprintf(key, "item %d", i);
  //   Hashmap_set(h, key, &ns[i]);
  // }
  // for (i = 0; i < n; i++) {
  //   sprintf(key, "item %d", i);
  //   assert(Hashmap_get(h, key) == &ns[i]);
  // }

  Hashmap_free(h);
  /*
     stretch goals:
     - expand the underlying array if we start to get a lot of collisions
     - support non-string keys
     - try different hash functions
     - switch from chaining to open addressing
     - use a sophisticated rehashing scheme to avoid clustered collisions
     - implement some features from Python dicts, such as reducing space use,
     maintaing key ordering etc. see https://www.youtube.com/watch?v=npw4s1QTmPg
     for ideas
     */
  printf("ok\n");
}
