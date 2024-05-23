#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <unistd.h>
#include <sys/times.h>

#define SLEEP_SEC 3
#define NUM_MULS 100000000
#define NUM_MALLOCS 100000
#define MALLOC_SIZE 1000

// TODO define this struct
struct profile_times {
  int pid;
  char desc[81];
  struct tms start;
};

// TODO populate the given struct with starting information
void profile_start(struct profile_times *t) {
  t->pid = getpid();
  times(&t->start);
}

// TODO given starting information, compute and log differences to now
void profile_log(struct profile_times *t) {
  struct tms end;
  times(&end);
  printf("%s:\n", t->desc);
  printf("[pid %d] %ld %ld\n------------\n", t->pid, end.tms_stime - t->start.tms_stime, end.tms_utime - t->start.tms_utime);
}

int main(int argc, char *argv[]) {
  printf("Clock ticks per second: %ld\n", sysconf(_SC_CLK_TCK));
  struct profile_times t;

  // TODO profile doing a bunch of floating point muls
  float x = 1.0;
  sprintf(t.desc, "%d fmuls", NUM_MULS);
  profile_start(&t);
  for (int i = 0; i < NUM_MULS; i++)
    x *= 1.1;
  profile_log(&t);

  // TODO profile doing a bunch of mallocs
  sprintf(t.desc, "%d mallocs of size %d", NUM_MALLOCS, MALLOC_SIZE);
  profile_start(&t);
  void *p;
  for (int i = 0; i < NUM_MALLOCS; i++)
    p = malloc(MALLOC_SIZE);
  profile_log(&t);

  // TODO profile sleeping
  sprintf(t.desc, "sleeping for %d seconds", SLEEP_SEC);
  profile_start(&t);
  sleep(SLEEP_SEC);
  profile_log(&t);
}
