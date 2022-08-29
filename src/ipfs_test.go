package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

func TestMain(m *testing.M) {
	startIPFS()
	sigInt := m.Run()

	defer KillIPFSCtx(context.Background(), "ipfs")
	os.Exit(sigInt)
}

func KillIPFSCtx(ctx context.Context, name string) error {
	processes, err := process.ProcessesWithContext(ctx)
	if err != nil {
		return err
	}
	for _, p := range processes {
		n, err := p.NameWithContext(ctx)
		if err != nil {
			return err
		}
		if n == name {
			err = p.KillWithContext(ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func startIPFS() error {
	cmd := exec.CommandContext(context.Background(), "ipfs", "daemon")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(4 * time.Second)
	return err
}
