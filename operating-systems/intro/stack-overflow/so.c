#include <stdio.h>
#include <stdlib.h>
#include <sys/resource.h>
#include <unistd.h>

void f(int x, long p) {
  char y[4096];
  fprintf(stderr, "frame %d (%p) - %ld bytes\n", x, &x, p - (long)&x);
  f(x+1, p);
}

void start() {
  int x = 0;
  f(x, (long)&x);
}

int main(int argc, char *argv[])
{
  long page_size = sysconf(_SC_PAGESIZE);
  printf("Page size: %ld\n", page_size);
  struct rlimit limit;
  getrlimit(RLIMIT_STACK, &limit);
  printf("Limits: %lu, %luM\n", limit.rlim_cur, limit.rlim_max / (1024 * 1024));
  limit.rlim_cur = 0;
  setrlimit(RLIMIT_STACK, &limit);
  getrlimit(RLIMIT_STACK, &limit);
  printf("Limits: %lu, %luM\n", limit.rlim_cur, limit.rlim_max / (1024 * 1024));
  start();
  
  return EXIT_SUCCESS;
}
