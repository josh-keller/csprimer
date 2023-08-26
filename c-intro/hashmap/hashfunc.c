#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[])
{
  if (argc != 2) {
    fprintf(stderr, "Usage: ./hashfunc <string to be hashed>\n");
    return EXIT_FAILURE;
  }

  char *str = argv[1];

  uint32_t hash = 0;

  for (int i = 0; str[i]; i++) {
    uint32_t highorder = hash & 0xf8000000;
    hash <<= 5;
    hash = hash ^ ()
    int 
  }


  return EXIT_SUCCESS;
}
