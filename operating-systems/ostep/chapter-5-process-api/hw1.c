// Write a program that calls fork(). Before calling fork(), have the
// main process access a variable (e.g., x) and set its value to some-
// thing (e.g., 100). What value is the variable in the child process?
// What happens to the variable when both the child and parent change
// the value of x?

#include <stdio.h>
#include <stdlib.h>
#include <sys/wait.h>
#include <unistd.h>

int main(int argc, char *argv[]) {
  int x = 42;
  int *y = &x;
  int *z = malloc(sizeof(int));
  *z = 1000;

  int rc = fork();
  if (rc < 0) {
    fprintf(stderr, "Fork failed\n");
  } else if (rc == 0) {
    printf("Child (pid: %d):\n", (int) getpid());
    printf("Initial value of x: %d\n", x);
    x = 51;
    printf("Changed value: %d\n", x);
    printf("Address of x is %p\n", y);
    *y = 0;
    printf("Changed by pointer: %d\n", x);
    printf("Value at z: %d\n", *z);
    *z = 2000;
    printf("Changed value at z: %d\n", *z);
    printf("------------------------------\n");
  } else {
    int rc_pid = wait(NULL);
    printf("Parent (pid: %d):\n", (int) getpid());
    printf("Initial value of x: %d\n", x);
    x = 33;
    printf("Changed value: %d\n", x);
    printf("Address of x is %p\n", y);
    *y = 0;
    printf("Changed through *y: %d\n", x);
    printf("Value at z: %d\n", *z);
    *z = 3000;
    printf("Changed value at z: %d\n", *z);
  }

  return 0;
}
