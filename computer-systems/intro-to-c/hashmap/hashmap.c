#include <assert.h>
#include <stddef.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define STARTING_BUCKETS 8
#define MAX_KEY_SIZE 32

typedef uint32_t Hash;

typedef struct Node {
  char *key;
  void *val;
  Hash hash;
  struct Node *next;
} Node;

typedef struct {
  Node **buckets;
  size_t bucket_count;
} Hashmap;

Hash hash(const char *str) {
  Hash hash = 5381;
  int c;

  while ((c = *str++)) {
    hash = ((hash << 5) + hash) + c;
  }

  return hash;
}

Node *Node_new(char *key, void *val, Hash hash, Node *next) {
  char *new_key = strdup(key);

  if (new_key == NULL) return NULL;

  Node *new_node = (Node *)malloc(sizeof(Node));
  if (new_node == NULL) {
    return NULL;
  }

  new_node->key = new_key;
  new_node->hash = hash;
  new_node->val = val;
  new_node->next = next;
  return new_node;
}

Hashmap *Hashmap_new(void) {
  Hashmap *h = malloc(sizeof(Hashmap));
  if (h == NULL) {
    return NULL;
  }
  h->bucket_count = STARTING_BUCKETS;
  h->buckets = calloc(STARTING_BUCKETS, sizeof(Node));
  if (h->buckets == NULL) {
    free(h);
    return NULL;
  }
  return h;
}

void Hashmap_free(Hashmap *h) {
  for (int i = 0; i < h->bucket_count; i++) {
    for (Node *curr = h->buckets[i], *next = NULL; curr != NULL; curr = next) {
      next = curr->next;
      free(curr->key);
      free(curr);
    }
  }
  free(h->buckets);
  free(h);
}

void Hashmap_set(Hashmap *h, char *key, void *val) {
  Hash key_hash = hash(key);
  Hash bucket = key_hash % h->bucket_count;

  if (h->buckets[bucket] == NULL) {
    h->buckets[bucket] = Node_new(key, val, key_hash, NULL);
    return;
  }

  Node *next = h->buckets[bucket];
  Node *curr = NULL;
  
  while (next != NULL) {
    curr = next;
    next = curr->next;
    if (strcmp(curr->key, key) == 0) {
      curr->val = val;
      return;
    }
  }

  curr->next = Node_new(key, val, key_hash, NULL);
  return;
}

void *Hashmap_get(Hashmap *h, char *key) {
  uint32_t key_hash = hash(key) % h->bucket_count;

  Node *curr = h->buckets[key_hash];
  while (curr != NULL) {
    if (strcmp(curr->key, key) == 0) {
      return curr->val;
    }
    curr = curr->next;
  }

  return NULL;
}

void Hashmap_delete(Hashmap *h, char *key) {
  uint32_t key_hash = hash(key) % h->bucket_count;
  Node *curr = h->buckets[key_hash];
  Node *prev = NULL;
  
  while (curr != NULL) {
    if (strcmp(curr->key, key) == 0) {
      if (prev != NULL) {
        prev->next = curr->next;
      } else {
        h->buckets[key_hash] = curr->next;
      }

      free(curr->key);
      free(curr);
      return;
    }
    prev = curr;
    curr = curr->next;
  }
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

  // basic delete functionality
  Hashmap_delete(h, "item a");
  assert(Hashmap_get(h, "item a") == NULL);

  // handle collisions correctly
  // note: this doesn't necessarily test expansion
  int i, n = STARTING_BUCKETS * 10, ns[n];
  char key[MAX_KEY_SIZE];
  for (i = 0; i < n; i++) {
    ns[i] = i;
    sprintf(key, "item %d", i);
    Hashmap_set(h, key, &ns[i]);
  }
  for (i = 0; i < n; i++) {
    sprintf(key, "item %d", i);
    assert(Hashmap_get(h, key) == &ns[i]);
  }

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
