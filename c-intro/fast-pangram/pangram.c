#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>


bool ispangram(char *s) {
  char c;
  u_int32_t counts = 0;

  for (int i = 0; s[i] != 0; i++) {
    c = s[i];
    if (c >= 65 && c <= 90) {
      c += 32;
    }

    if (c >= 97 && c <= 122) {
      int bit = (1 << (c - 97));
      counts |= bit;
    }
    
    if (counts == 0x03ffffff) {
      return true;
    }
  }
  return false;
}

int main() {
  size_t len;
  ssize_t read;
  char *line = NULL;
  while ((read = getline(&line, &len, stdin)) != -1) {
    if (ispangram(line))
      printf("%s", line);
  }

  if (ferror(stdin))
    fprintf(stderr, "Error reading from stdin");

  free(line);
  fprintf(stderr, "ok\n");
}
