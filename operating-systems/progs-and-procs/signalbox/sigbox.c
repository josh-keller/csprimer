#include <errno.h>
#include <stdlib.h>
#include <string.h>
#include <sys/ioctl.h>
#include <signal.h>
#include <stdio.h>
#include <termio.h>
#include <fcntl.h>
#include <string.h>
#include <unistd.h>

#define HORIZ_PAD 5
#define VERT_PAD 2

volatile sig_atomic_t needs_update = 0;

void draw_box(unsigned short width, unsigned short height) {
  unsigned short box_width = width - 2 * HORIZ_PAD;
  unsigned short box_height = height - 2 * VERT_PAD;
  char tb_line[257];
  char mid_line[257];
  int i = 0;
  for (; i < HORIZ_PAD; i++) {
    tb_line[i] = ' ';
    mid_line[i] = ' ';
  }

  tb_line[i] = '*';
  mid_line[i] = '|';
  i++;
  
  for (; i < HORIZ_PAD + box_width - 1; i++) {
    tb_line[i] = '-';
    mid_line[i] = ' ';
  }

  tb_line[i] = '*';
  mid_line[i] = '|';
  i++;

  tb_line[i] = 0;
  mid_line[i] = 0;
  
  // Clear screen
  printf("\e[1;1H\e[2J");
  for (i = 0; i < VERT_PAD; i++) printf("\n");
  printf("%s\n", tb_line);
  i++;
  for (; i < VERT_PAD + box_height - 1; i++) printf("%s\n", mid_line);
  printf("%s\n", tb_line);
}

void sig_handler() {
  needs_update = 1;
}

void update_box() {
  struct winsize ws;
  int fd = open("/dev/tty", O_RDONLY);
  if (fd < 0) {
    fprintf(stderr, "Error opening tty: %s\n", strerror(errno));
    exit(1);
  }
  ioctl(fd, TIOCGWINSZ, &ws);
  
  draw_box(ws.ws_col, ws.ws_row);
  needs_update = 0;
}

int main(int argc, char *argv[]) {
  struct sigaction sa;
  sa.sa_handler = sig_handler;
  sigaction(SIGWINCH, &sa, NULL);

  update_box();

  while (1) {
    pause();
    if (needs_update) {
      update_box();
    }
  }

  return 0;
}
