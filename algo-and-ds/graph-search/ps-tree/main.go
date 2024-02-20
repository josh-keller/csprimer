package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Proc struct {
	Pid      int
	Ppid     int
	Command  string
	Children []int
}

type ProcTree map[int]*Proc

const (
	lastChildPrefix = "└──"
	childPrefix     = "├──"
	belowLastChild  = "   "
	belowChild      = "│  "
)

func (pt ProcTree) Print() {
	pt.printProc(1, "")
}

func main() {
	cmd := exec.Command("/bin/ps", "-axc", "-o", "pid=", "-o", "ppid=", "-o", "command=")

	output, err := cmd.Output()
	if err != nil {
		log.Fatal("Output error: ", err)
	}
	scanner := bufio.NewScanner(bytes.NewReader(output))
	procs := ProcTree{}

	/*
		 Pid should be unique
		 Ppid will not be
		 Ppid may or may not have been seen as a pid yet

		 Cases:
		 - Neither Pid or Ppid exist
		   - Create Pid with no children
			 - Create Ppid (no Ppid, no command, add pid to children)
		 - Pid exists but Ppid doesn't
		   - Pid is okay
			 - Add ppid and command
			 - Don't touch children
			 - Create Ppid as above

		 - Ppid exists but Pid doesn't
		   - Create Pid with no children
			 - Add Pid to ppid children

		 - Both exist
		   - Pid is okay
			 - Add ppid and command
			 - Don't touch children
			 - Add Pid to ppid chilren

		 Algorithm:
		   If pid doesn't exist:
			 - Create empty proc at pid (only have pid)
			 If ppid doesn't exist:
			 - Create empty ppid at ppid (only have ppid)
			 Add ppid and command to proc at pid
			 Add pid as a child of ppid
	*/

	for scanner.Scan() {
		re := regexp.MustCompile(`\s+`)
		line := strings.TrimSpace(scanner.Text())
		fields := re.Split(line, 3)

		pid, err := strconv.Atoi(fields[0])
		if err != nil {
			log.Fatal(err)
		}
		ppid, err := strconv.Atoi(fields[1])
		if err != nil {
			log.Fatal(err)
		}

		_, exists := procs[pid]
		if !exists {
			procs[pid] = &Proc{
				Pid: pid,
			}
		}
		procs[pid].Ppid = ppid
		procs[pid].Command = fields[2]

		_, exists = procs[ppid]
		if !exists {
			procs[ppid] = &Proc{
				Pid: ppid,
			}
		}
		procs[ppid].Children = append(procs[ppid].Children, pid)
	}

	procs.Print()
}

func (pt ProcTree) printProc(pid int, prefix string) {
	fmt.Printf("%s%d (%s)\n", prefix, pid, pt[pid].Command)

	for i, childPid := range pt[pid].Children {
		nextPrefix := strings.Replace(prefix, lastChildPrefix, belowLastChild, -1)
		nextPrefix = strings.Replace(nextPrefix, childPrefix, belowChild, -1)

		if i == len(pt[pid].Children)-1 {
			nextPrefix = nextPrefix + lastChildPrefix
		} else {
			nextPrefix = nextPrefix + childPrefix
		}
		pt.printProc(childPid, nextPrefix)
	}
}
