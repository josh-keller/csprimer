# Signal Logger Notes

|     | Name      | Default Action    | Description                                     | Triggered     |
| --- | --------- | ----------------- | ----------------------------------------------- | ------------- |
| 1   | SIGHUP    | terminate process | terminal line hangup                            |               |
| 2   | SIGINT    | terminate process | interrupt program                               | x             |
| 3   | SIGQUIT   | create core image | quit program                                    | x             |
| 4   | SIGILL    | create core image | illegal                                         |               |
| 5   | SIGTRAP   | create core image | trace trap                                      |               |
| 6   | SIGABRT   | create core image | abort program (formerly SIGIOT)                 |               |
| 7   | SIGEMT    | create core image | emulate instruction executed                    |               |
| 8   | SIGFPE    | create core image | floating-point exception                        |               |
| 9   | SIGKILL   | terminate process | kill program                                    | x             |
| 10  | SIGBUS    | create core image | bus error                                       |               |
| 11  | SIGSEGV   | create core image | segmentation violation                          |               |
| 12  | SIGSYS    | create core image | non-existent system call invoked                |               |
| 13  | SIGPIPE   | terminate process | write on a pipe with no reader                  |               |
| 14  | SIGALRM   | terminate process | real-time timer expired                         | x             |
| 15  | SIGTERM   | terminate process | software termination signal                     | x             |
| 16  | SIGURG    | discard signal    | urgent condition present on socket              |               |
| 17  | SIGSTOP   | stop process      | stop (cannot be caught or ignored)              |               |
| 18  | SIGTSTP   | stop process      | stop signal generated from keyboard             | CTRL-Z        |
| 19  | SIGCONT   | discard signal    | continue after stop                             | fg            |
| 20  | SIGCHLD   | discard signal    | child status has changed                        | child exit    |
| 21  | SIGTTIN   | stop process      | background read attempted from control terminal |               |
| 22  | SIGTTOU   | stop process      | background write attempted to control terminal  |               |
| 23  | SIGIO     | discard signal    | I/O is possible on a descriptor (see fcntl(2))  |               |
| 24  | SIGXCPU   | terminate process | cpu time limit exceeded (see setrlimit(2))      |               |
| 25  | SIGXFSZ   | terminate process | file size limit exceeded (see setrlimit(2))     |               |
| 26  | SIGVTALRM | terminate process | virtual time alarm (see setitimer(2))           |               |
| 27  | SIGPROF   | terminate process | profiling timer alarm (see setitimer(2))        |               |
| 28  | SIGWINCH  | discard signal    | Window size change                              | Change window |
| 29  | SIGINFO   | discard signal    | status request from keyboard                    | CTRL-T        |
| 30  | SIGUSR1   | terminate process | User defined signal 1                           | kill()        |
| 31  | SIGUSR2   | terminate process | User defined signal 2                           | kill()        |

## Terminal input

```sh
$ stty -a

speed 9600 baud; 57 rows; 233 columns;
lflags: icanon isig iexten echo echoe -echok echoke -echonl echoctl
        -echoprt -altwerase -noflsh -tostop -flusho -pendin -nokerninfo
        -extproc
iflags: -istrip icrnl -inlcr -igncr ixon -ixoff ixany imaxbel iutf8
        -ignbrk brkint -inpck -ignpar -parmrk
oflags: opost onlcr -oxtabs -onocr -onlret
cflags: cread cs8 -parenb -parodd hupcl -clocal -cstopb -crtscts -dsrflow
        -dtrflow -mdmbuf
cchars: discard = ^O; dsusp = ^Y; eof = ^D; eol = <undef>;
        eol2 = <undef>; erase = ^?; intr = ^C; kill = ^U; lnext = ^V;
        min = 1; quit = ^\; reprint = ^R; start = ^Q; status = ^T;
        stop = ^S; susp = ^Z; time = 0; werase = ^W;
```

cchars of interest:

- intr = `^C`
- kill = `^U`
- quit = `^\`
- status = `^T`
- susp = `^Z`

## Resources

- <https://www.cs.kent.edu/~ruttan/sysprog/lectures/signals.html#:~:text=The%20most%20common%20way%20of,SIGINT%20)%20to%20the%20running%20process>.
