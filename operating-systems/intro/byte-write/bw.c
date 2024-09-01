#include <errno.h>
#include <fcntl.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <unistd.h>
#include <stdio.h>

int main() {
  const char *filename = "./output.txt";
  printf("Opening...\n");
  int fd = open(filename, O_CREAT | O_TRUNC | O_WRONLY, 0644);
  if (fd == -1) {
    printf("Error: %s\n", strerror(errno));
    exit(1);
  }
  struct stat info;
  stat(filename, &info);
  off_t blocks = info.st_blocks;

  for (u_int i = 1; i <= (1024 * 1024); i++) {
    int ret = write(fd, "h", 1);
    if (ret == -1) {
      printf("Error: %s\n", strerror(errno));
      exit(1);
    }
    stat(filename, &info);
    if (info.st_blocks > blocks) {
      printf("Size on disk changed at byte %d, from %ld to %ld blocks\n", i, blocks, info.st_blocks);
      blocks = info.st_blocks;
    }
  }

  close(fd);
}
