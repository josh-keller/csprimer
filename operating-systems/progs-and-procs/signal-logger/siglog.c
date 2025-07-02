#include <errno.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

pid_t mypid;

volatile uint64_t handled = 0;
volatile int sigint_handled = 0;

static void tstpHandler(int sig) {
  sigset_t tstpMask, prevMask;
  int savedErrno;
  struct sigaction sa;

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

  if (sig == SIGINT) {
    if (sigint_handled) {
      exit(0);
    }

    sigint_handled = 1;
    printf("Next time SIGINT will stop the program.\n");
  }

  if (sig == SIGTSTP) {
    tstpHandler(sig);
  }
}

int main(int argc, char *argv[]) {
  mypid = getpid();
  // Register all valid signals
  for (int i = 0; i < NSIG; i++) {
    signal(i, handle);
  }

  printf("Signal logger PID: %d\n", mypid);

  alarm(1);

  if (0 == fork()) {
    exit(0);
  }

  // Send user signals
  kill(mypid, SIGUSR1);
  kill(mypid, SIGUSR2);

  // spin
  for (;;)
    sleep(1);
}
