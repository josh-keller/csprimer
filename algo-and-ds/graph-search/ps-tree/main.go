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

const (
	lastChildPrefix = "└──"
	childPrefix     = "├──"
	belowLastChild  = "   "
	belowChild      = "│  "
)

var re = regexp.MustCompile(`\s+`)

type Proc struct {
	Pid      int
	Ppid     int
	Command  string
	Children []int
}

func NewProcFromString(s string) (*Proc, error) {
	fields := re.Split(strings.TrimSpace(s), 3)

	pid, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, err
	}

	ppid, err := strconv.Atoi(fields[1])
	if err != nil {
		return nil, err
	}

	command := fields[2]

	return &Proc{
		Pid:     pid,
		Ppid:    ppid,
		Command: command,
	}, nil
}

func (p *Proc) AddChild(childPid int) {
	p.Children = append(p.Children, childPid)
}

type ProcTree map[int]*Proc

func (pt ProcTree) Print() {
	pt.printProc(1, "")
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

func main() {
	cmd := exec.Command("/bin/ps", "-Ac", "-o", "pid=,ppid=,command=")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal("Output error: ", err)
	}

	procs := ProcTree{}

	scanner := bufio.NewScanner(bytes.NewReader(output))

	for scanner.Scan() {
		proc, err := NewProcFromString(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		if existingProc, exists := procs[proc.Pid]; exists {
			proc.Children = existingProc.Children
		}

		procs[proc.Pid] = proc

		_, exists := procs[proc.Ppid]
		if !exists {
			procs[proc.Ppid] = &Proc{
				Pid: proc.Ppid,
			}
		}
		procs[proc.Ppid].AddChild(proc.Pid)
	}

	procs.Print()
}
