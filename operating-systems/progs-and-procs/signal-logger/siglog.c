#include <errno.h>
#include <float.h>
#include <setjmp.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

pid_t mypid;

volatile uint64_t handled = 0;
volatile int sigint_handled = 0;
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
    exit(1);
  }

  raise(SIGTSTP);

  sigemptyset(&tstpMask);
  sigaddset(&tstpMask, SIGTSTP);
  if (sigprocmask(SIG_UNBLOCK, &tstpMask, &prevMask) == -1) {
    exit(1);
  }

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

  // Child
  if (0 == fork()) {
    exit(0);
  }

  // Alarm
  alarm(1);

  // spin
  for (;;)
    sleep(1);
}
