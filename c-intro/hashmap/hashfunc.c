#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

unsigned long hash(char *str) {
  unsigned long hash = 5381;
  int c;

  while ((c = *str++)) {
    hash = ((hash << 5) + hash) + c;
  }

  return hash;
}

int main(int argc, char *argv[])
{
  if (argc != 2) {
    fprintf(stderr, "Usage: ./hashfunc <string to be hashed>\n");
    return EXIT_FAILURE;
  }

  char *str = argv[1];

  printf("%lu\n", hash(str) % 1000);

  return EXIT_SUCCESS;
}

