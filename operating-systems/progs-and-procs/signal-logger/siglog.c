#include <errno.h>
#include <fcntl.h>
#include <float.h>
#include <netinet/in.h>
#include <setjmp.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>

pid_t mypid;

volatile uint64_t handled = 0;
volatile int sigint_handled = 0, quit = 0;
sigjmp_buf jmp_env_seg, jmp_env_fpe, jmp_env_xfsz, jmp_env_bus, jmp_env_ill,
    jmp_env_sys;

// Source:
// https://www.man7.org/tlpi/code/online/dist/pgsjc/handling_SIGTSTP.c.html
static void tstpHandler(int sig) {
  sigset_t tstpMask, prevMask;
  int savedErrno;
  struct sigaction sa;

  // Reset errno if it is changed here
  savedErrno = errno;

  if (signal(SIGTSTP, SIG_DFL) == SIG_ERR) {
    quit = 1;
    return;
  }

  printf("Raising sigstp\n");
  raise(SIGTSTP);

  sigemptyset(&tstpMask);
  sigaddset(&tstpMask, SIGTSTP);
  if (sigprocmask(SIG_UNBLOCK, &tstpMask, &prevMask) == -1) {
    exit(1);
  }

  // This will trigger SIG_TTIN if process is put in background
  printf("Press any key ");
  char c = getc(stdin);
  printf("\n");

  if (sigprocmask(SIG_SETMASK, &prevMask, NULL) == -1) {
    exit(1);
  }

  sigemptyset(&sa.sa_mask);
  sa.sa_flags = SA_RESTART;
  sa.sa_handler = tstpHandler;
  if (sigaction(SIGTSTP, &sa, NULL) == -1) {
    exit(1);
  }

  printf("Exiting SIGSTP handler\n");
  errno = savedErrno;
}

void handle(int sig) {
  handled |= (1 << sig);
  printf("Caught %d: %s (%d total)\n", sig, sys_signame[sig],
         __builtin_popcount(handled));

  switch (sig) {
  case SIGINT:
    if (sigint_handled) {
      exit(0);
    }
    sigint_handled = 1;
    printf("Next SIGINT will stop the program.\n");
    break;
  case SIGTSTP:
    tstpHandler(sig);
    break;
  case SIGSEGV:
    siglongjmp(jmp_env_seg, 1);
  case SIGFPE:
    siglongjmp(jmp_env_fpe, 1);
  case SIGXFSZ:
    siglongjmp(jmp_env_xfsz, 1);
  case SIGBUS:
    siglongjmp(jmp_env_bus, 1);
  case SIGILL:
    siglongjmp(jmp_env_ill, 1);
  case SIGSYS:
    siglongjmp(jmp_env_sys, 1);
  }
}

void segfault() {
  int *segfault = NULL;
  if (sigsetjmp(jmp_env_seg, 1) == 0) {
    int foo = *segfault;
  }
}

void bus_err() {
  int *bus;
  if (sigsetjmp(jmp_env_bus, 1) == 0) {
    *bus = 7;
  }
}

void file_size_limit() {
  struct rlimit curr_size;
  getrlimit(RLIMIT_FSIZE, &curr_size);
  printf("Initial filesize limits: %llu, %llu\n", curr_size.rlim_cur,
         curr_size.rlim_max);
  struct rlimit new_size = {0, curr_size.rlim_max};
  setrlimit(RLIMIT_FSIZE, &new_size);
  if (sigsetjmp(jmp_env_xfsz, 1) == 0) {
    // Set rlimit of no files, try to create a file
    getrlimit(RLIMIT_FSIZE, &curr_size);
    printf("Changed file size limits: %llu, %llu\n", curr_size.rlim_cur,
           curr_size.rlim_max);

    FILE *filePointer;
    filePointer = fopen("/tmp/testsignals", "w");
    if (filePointer == NULL) {
      printf("Could not open file");
      exit(1);
    }
    fputs("Testing this file.\n", filePointer);
    fclose(filePointer);
  } else {
    new_size.rlim_cur = curr_size.rlim_max;
    new_size.rlim_max = curr_size.rlim_max;
    setrlimit(RLIMIT_FSIZE, &new_size);
    getrlimit(RLIMIT_FSIZE, &curr_size);
    printf("Restored file size limits: %llu, %llu\n", curr_size.rlim_cur,
           curr_size.rlim_max);
  }
}

int set_up_socket() {
  int sockfd;
  struct sockaddr_in server_addr;
  int flags;

  sockfd = socket(AF_INET, SOCK_STREAM, 0);
  if (sockfd < 0) {
    perror("socket");
    exit(EXIT_FAILURE);
  }

  memset(&server_addr, 0, sizeof(server_addr));
  server_addr.sin_family = AF_INET;
  server_addr.sin_addr.s_addr = htonl(INADDR_ANY);
  server_addr.sin_port = htons(12345);

  if (bind(sockfd, (struct sockaddr *)&server_addr, sizeof(server_addr)) < 0) {
    perror("bind");
    close(sockfd);
    exit(EXIT_FAILURE);
  }

  flags = fcntl(sockfd, F_GETFL, 0);
  fcntl(sockfd, F_SETFL, flags | O_ASYNC | O_NONBLOCK);

  fcntl(sockfd, F_SETOWN, getpid());

  // fcntl(sockfd, F_SETSIG, SIGURG);
  if (listen(sockfd, 5) < 0) {
    perror("listen");
    close(sockfd);
    exit(EXIT_FAILURE);
  }

  printf("Waiting for connections on fd %d...\n", sockfd);
  return sockfd;
}

int main(int argc, char *argv[]) {
  mypid = getpid();
  // Register all valid signals
  for (int i = 0; i < NSIG; i++) {
    signal(i, handle);
  }

  printf("Signal logger PID: %d\n", mypid);

  // Send user signals
  kill(mypid, SIGUSR1);
  kill(mypid, SIGUSR2);

  segfault();
  bus_err();
  file_size_limit();

  if (sigsetjmp(jmp_env_fpe, 1) == 0) {
    printf("Float norm max^2: %f\n", FLT_NORM_MAX * FLT_NORM_MAX);
  }

  if (sigsetjmp(jmp_env_ill, 1) == 0) {
    __asm__ __volatile__(".word 0xffffffff");
  }

  if (sigsetjmp(jmp_env_sys, 1) == 0) {
    __asm__ __volatile__("mov x16, #0xabcd;"
                         "svc #0x80;");
  }

  int sockfd = set_up_socket();

  // Child
  if (0 == fork()) {
    exit(0);
  }

  // Alarm
  alarm(1);

  // spin
  while (!quit)
    sleep(1);

  close(sockfd);
}
