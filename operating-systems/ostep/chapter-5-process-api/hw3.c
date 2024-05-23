// Write another program using fork(). The child process should
// print “hello”; the parent process should print “goodbye”. You should
// try to ensure that the child process always prints first; can you do
// this without calling wait() in the parent?
#include <string.h>
#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <fcntl.h>

int main(int argc, char *argv[]) {
  int rc = fork();

  if (rc < 0) {
    fprintf(stderr, "fork failed\n");
    exit(1);
  } else if (rc == 0) {
    printf("hello\n");
  } else {
    int fd = open("/dev/null/", O_APPEND);
    write(fd, "waiting", strlen("waiting"));
    close(fd);
    sleep(1);
    printf("goodbye\n");
  }

  return 0;
}
