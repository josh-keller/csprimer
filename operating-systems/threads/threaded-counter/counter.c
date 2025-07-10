#include <pthread.h>
#include <signal.h>
#include <stdio.h>

#define NUMBER_OF_THREADS 2

const u_int64_t EACH_COUNT = 1000000000;
volatile u_int64_t counter = 0;
pthread_mutex_t lock;

void handle_int() { printf("Counter: %llu\n", counter); }

void *thread_entry(void *arg) {
  u_int64_t local_counter = 0;
  for (u_int64_t i = 0; i < EACH_COUNT; i++) {
    local_counter++;
  }
  pthread_mutex_lock(&lock);
  counter += local_counter;
  pthread_mutex_unlock(&lock);
  return NULL;
}

int main() {
  signal(SIGINT, handle_int);
  pthread_t threads[NUMBER_OF_THREADS];
  pthread_mutex_init(&lock, NULL);
  for (int i = 0; i < NUMBER_OF_THREADS; i++) {
    pthread_create(&threads[i], NULL, thread_entry, NULL);
  }

  for (int i = 0; i < NUMBER_OF_THREADS; i++) {
    pthread_join(threads[i], NULL);
  }
  printf("Final count: %llu (expected %llu)\n", counter,
         NUMBER_OF_THREADS * EACH_COUNT);
}
