package main

import (
	"errors"
	"log"
	"os/exec"
	"strings"
	"syscall"
)

type ChildProgram struct {
	cmd    *exec.Cmd
	stdIn  *NonBlockingWriteCloser
	stdOut *NonBlockingReader
	stdErr *NonBlockingReader

	exitErr *exec.ExitError
}

func NewChildProgram(name string, arg ...string) (*ChildProgram, error) {
	cmd := exec.Command(name, arg...)

	// cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdIn, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	ret := &ChildProgram{
		cmd:    cmd,
		stdIn:  NewNonBlockingWriter(stdIn),
		stdOut: NewNonBlockingReader(stdOut),
		stdErr: NewNonBlockingReader(stdErr),
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (s *ChildProgram) Wait() error {
	err := s.cmd.Wait()
	if err == nil {
		return nil
	}

	if exitErr := (&exec.ExitError{}); errors.As(err, &exitErr) {
		s.exitErr = exitErr
		return nil
	}

	return err
}

func (s *ChildProgram) Stop() error {
	err1 := s.stdIn.Close()
	if err1 != nil {
		log.Println(err1)
	}
	err2 := s.cmd.Process.Kill()
	if err2 != nil {
		log.Println(err2)
	}

	return err2
}

func (s *ChildProgram) IsRunning() bool {
	if s.cmd.ProcessState.Success() {
		return true
	}

	if s.cmd.ProcessState.Exited() {
		return true
	}

	// p, _ := os.FindProcess(123)
	err := s.cmd.Process.Signal(syscall.Signal(0))
	if err != nil {
		return true
	}

	return false
}

func (s *ChildProgram) WriteStdIn(in string) {
	s.stdIn.WriteString(in)
}

func (s *ChildProgram) ReadLineStdOut() string {
	return s.stdOut.GetLine()
}

func (s *ChildProgram) ReadAllStdOut() string {
	b := strings.Builder{}
	for {
		line := s.stdOut.GetLine()
		if line == "" {
			break
		}
		_, err := b.WriteString(line)
		if err != nil { // we shouldn't get an error
			panic(err)
		}
	}
	return b.String()
}

func (s *ChildProgram) ReadLineStdErr() string {
	return s.stdErr.GetLine()
}

func (s *ChildProgram) ReadAllStdErr() string {
	b := strings.Builder{}
	for {
		line := s.stdErr.GetLine()
		if line == "" {
			break
		}
		_, err := b.WriteString(line)
		if err != nil { // we shouldn't get an error
			panic(err)
		}
	}
	return b.String()
}
