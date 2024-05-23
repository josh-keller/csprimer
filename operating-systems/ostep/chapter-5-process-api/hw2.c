// Write a program that opens a file (with the open() system call)
// and then calls fork() to create a new process. Can both the child
// and parent access the file descriptor returned by open()? What
// happens when they are writing to the file concurrently, i.e., at the
// same time?

#include <fcntl.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

int main(int argc, char *argv[]) {
  printf("Start of process %d\n", (int) getpid());
  int fd = open("/tmp/ostep-temp-file.txt", O_CREAT|O_WRONLY|O_TRUNC, S_IRWXU);
  printf("Opened fd %d. Calling fork...\n", fd);
  int rc = fork();
  if (rc < 0) {
    // fork failed
    fprintf(stderr, "fork failed\n");
    exit(1);
  } else if (rc == 0) {
    printf("Writing to fd %d...\n", fd);
    write(fd, "Writing from the child\n", strlen("Writing from the child\n"));
  } else {
    printf("Writing to fd %d...\n", fd);
    write(fd, "Writing from the parent\n", strlen("Writing from the parent\n"));
    printf("Child finished. Closing file...\n");
    close(fd);
  }
  
  return 0;
}

/* My take:
 * It seems like 
 */
