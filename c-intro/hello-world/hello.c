#include <stdio.h>

int main(int argc, char **argv) {
  char *name;
  
  if (argc == 1) {
    name = "world";
  } else {
    name = argv[1];
  }

  printf("Hello %s!\n", name);
}
