package main

import (
	"fmt"
	"context"
	"os"
	"os/exec"

	shell "github.com/ipfs/go-ipfs-api"

)

// This function generates a key using local node and returns it
func genKey(keyName string) (*shell.Key, error) {
	sh := shell.NewShell(localhost) //grab local node

	key, err := sh.KeyGen(context.Background(), keyName) //generate temp key to local node
	if err != nil {
		fmt.Println("Error in node keyGen: ", err)
		return nil, err
	}
	return key, err	
}

// Delete key from local node keystore
func deleteKey(keyName string) error {
	sh := shell.NewShell(localhost)
	_, err := sh.KeyRm(context.Background(), keyName)
	if err != nil {
		fmt.Println("Error in node key delete: ", err)
		return err
	}
	return nil
}

// This function deletes the exported key from disk
func diskDelete(keyPath string) error {
	args := []string{keyPath}
	cmd := exec.Command("rm", args...)
	_, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in diskDelete: %s ", err)
		return err
	}
	return nil
}

// This function takes a key name and searches for it in local node Keystore.
// returns nil if sucessfull & stores key as file in current dir.
func exportKey(keyName string) error {
	args := []string{"key", "export", keyName}
	cmd := exec.Command("ipfs", args...)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in exportKey: %s ", err)
		return err
	}
	fmt.Println(string(stdout))
	return nil
}

// This function takes a key name and file name, searches for it in current dir.
// returns nil if sucessfull & stores key in local node.
func importKey(keyName string, fileName string) error {
	args := []string{"key", "import", keyName, fileName}
	cmd := exec.Command("ipfs", args...)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in importKey: %s", err)
		return err
	}
	fmt.Println(string(stdout))
	return nil
}