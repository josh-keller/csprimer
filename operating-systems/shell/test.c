#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main(int argc, char *argv[]) {
  char *str1, *str2, *token, *subtoken;
  char *saveptr1, *saveptr2;

  if (argc != 4) {
    fprintf(stderr, "Usage: %s string delim subdelim\n", argv[0]);
    exit(EXIT_FAILURE);
  }

  str1 = argv[1];
  for (int j = 1;; j++, str1 = NULL) {
    token = strtok_r(str1, argv[2], &saveptr1);
    if (token == NULL)
      break;
    printf("%d: %s\n", j, token);

    for (str2 = token;; str2 = NULL) {
      subtoken = strtok_r(str2, argv[3], &saveptr2);
      if (subtoken == NULL)
        break;
      printf(" --> %s\n", subtoken);
    }
  }

  exit(EXIT_SUCCESS);
}
