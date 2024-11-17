// Plan:
//
// 1. REPL (getc/puts)
// 2. Parse, print argv/argc (strtok)
// 3. Execute (fork/exec)
// 4. Refeactor/details signal handling, etc

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>

#define PROMPT "$ "
const int MAX_INPUT_LIMIT = 4096;

int execute(int s_argc, char *s_argv[]) {
  // Execute
  printf("argc: %d\n", s_argc);
  for (int j = 0; s_argv[j]; j++) {
    if (j != 0) {
      printf("|");
    }

    printf("%s", s_argv[j]);
  }
  printf("\n");
  printf("Executing %s\n", s_argv[0]);

  __pid_t pid = fork();
  if (pid < 0) {
    fprintf(stderr, "Fork error: %s\n", strerror(errno));
    return 1;
  } else if (pid == 0) {
    int ret = execvp(s_argv[0], s_argv);
  } else {
    int status;
    int child_pid = waitpid(pid, &status, 0);
    if (child_pid < 0) {
      fprintf(stderr, "Wait error: %s\n", strerror(errno));
      return 1;
    } else {
      fprintf(stderr, "Child returned.\n");
      return status;
    }
  }
}

int main(int argc, char *argv[]) {
  char input_buffer[MAX_INPUT_LIMIT];
  int s_argc = 0;
  char *s_argv[1024];
  char *str, *token;

  while (1) {
    // Prompt
    printf(PROMPT);

    // Get intput
    fgets(input_buffer, MAX_INPUT_LIMIT, stdin);

    str = input_buffer;

    // Strip trailing newline
    str[strcspn(str, "\n")] = '\0';

    // Parse
    for (int i = 0;; i++, str = NULL) {
      token = strtok(str, " ");
      if (token == NULL) {
        s_argc = i;
        s_argv[i] = NULL;
        break;
      }
      s_argv[i] = token;
    }
    if (s_argc == 0) {
      continue;
    }

    int result = execute(s_argc, s_argv);
    if (result == 0) {
      printf("Success...\n");
    } else {
      fprintf(stderr, "Exec returned error code %d: %s\n", result,
              strerror(errno));
    }

    // Cleanup
    s_argc = 0;
  }

  return EXIT_SUCCESS;
}
