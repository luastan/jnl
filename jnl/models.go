package jnl

import (
	"context"
	"fmt"
	"hash/fnv"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type CommandInfo struct {
	Command  []string  `json:"command"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end,omitempty"`
	ExitCode int       `json:"exitCode,omitempty"`
	Dir      string    `json:"dir,omitempty"`
	// 	Cancel   time.Time `json:"cancel"`
}

func (c *CommandInfo) hasEnded() bool {
	return !c.End.Equal(time.Time{})
}

type Execution struct {
	Info *CommandInfo
	wg   *sync.WaitGroup
	cmd  *exec.Cmd
	ctx  context.Context
}

func (e *Execution) Start() error {
	e.Info.Start = time.Now()
	cmdId := fmt.Sprintf(
		"%s:%d",
		strings.Join(e.Info.Command, " "),
		e.Info.Start.UnixNano(),
	)
	h := fnv.New64a()
	_, err := h.Write([]byte(cmdId))
	if err != nil {
		return err
	}
	e.Info.Dir = fmt.Sprintf("%d", h.Sum64())
	stdoutReader, err := e.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderrReader, err := e.cmd.StderrPipe()
	if err != nil {
		return err
	}
	e.wg.Add(2)
	go OutErrRoutine(e.wg, stdoutReader, os.Stdout)
	go OutErrRoutine(e.wg, stderrReader, os.Stderr)

	_ = SaveState(e)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		sig := <-sigs
		_ = e.cmd.Process.Signal(sig)
	}()
	return e.cmd.Start()
}

func (e *Execution) Wait() error {
	err := e.cmd.Wait()
	e.Info.End = time.Now()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				e.Info.ExitCode = status.ExitStatus()
			}
		} else {
			return err
		}
	} else {
		e.Info.ExitCode = 0
	}

	_ = SaveState(e)

	e.wg.Wait()
	return nil
}

func NewExecution(ctx context.Context, command []string) *Execution {
	c := exec.CommandContext(ctx, command[0], command[1:]...)
	c.Stdin = os.Stdin
	return &Execution{
		Info: &CommandInfo{
			Command: command,
		},
		ctx: ctx,
		cmd: c,
		wg:  &sync.WaitGroup{},
	}
}
