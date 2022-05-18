package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
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

func TestAddFile(t *testing.T) {
	resp := make(map[string]string)

	// create file
	path := "/tmp/dat"
	rr := addTmpFileToIPFS(path)

	// check response
	if rr.Code != http.StatusOK {
		t.Errorf("Error: %v\n Status Code: %v",
			rr.Body, rr.Code)
	}
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	fmt.Println("response value:", resp["value"])

	// delete tmp file
	err = os.Remove(path)
	check(err)
}

func TestForceIPFSErrors(t *testing.T) {
	//  // d.n.e = does not exist
	tests := []struct {
		desc    string
		fx      func() error
		testnum int
	}{
		{
			testnum: 1,
			desc:    "dir d.n.e",
			fx:      func() error { _, err := addToIPFS(abs+"stati", "r"); return err },
		},
		{
			testnum: 2,
			desc:    "file d.n.e",
			fx:      func() error { _, err := addToIPFS(abs+"/src/Tests/Helloooo", ""); return err },
		},
		{
			testnum: 3,
			desc:    "Rslv bad input",
			fx:      func() error { _, err := resolve(""); return err },
		},
		{
			testnum: 4,
			desc:    "dir daemon down",
			fx:      func() error { _, err := addToIPFS(abs+"/static", "r"); return err },
		},
		{
			testnum: 5,
			desc:    "file daemon down",
			fx:      func() error { _, err := addToIPFS(abs+"/src/Tests/Hello", ""); return err },
		},
		{
			testnum: 6,
			desc:    "Pub daemon down",
			fx:      func() error { _, err := publishToIPNS("", ""); return err },
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.testnum == 4 {
				KillIPFSCtx(context.Background(), "ipfs")
			}
			err := tt.fx()
			if err == nil {
				t.Errorf("function %s did not error out.", tt.desc)
			}
			fmt.Printf("\n%s failed: %s\n", tt.desc, err)
		})
	}

	if err := startIPFS(); err != nil {
		fmt.Println("Could not start daemon", err)
		t.Errorf("Couldn't start ipfs daemon %s", err)
	}
}

func addTmpFileToIPFS(path string) *httptest.ResponseRecorder {
	pl := []byte("Hello World")
	file, err := os.Create(path)
	check(err)
	defer file.Close()
	_, err = file.Write(pl)
	file.Seek(0, 0) //reset pointer to start of file

	// create writer to send to API
	w, body := createWriter(file)

	// Create request
	req, err := http.NewRequest("POST", "/addFile", body)
	check(err)
	req.Header.Add("Content-Type", w.FormDataContentType())

	// execute request
	return executeRequest(AddFile, req)
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
