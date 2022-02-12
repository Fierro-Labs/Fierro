package main

import (
	"fmt"
	"context"
	"os/exec"

	shell "github.com/ipfs/go-ipfs-api"

)

// Delete key from local node keystore
// Correct
func deleteKey(keyName string) error {
	sh := shell.NewShell(localhost)
	_, err := sh.KeyRm(context.Background(), keyName)
	if err != nil {
		return err
	}
	return nil
}

// This function deletes the exported key from  dir
func diskDelete(keyName string) error {
	args := []string{keyName+".key"}
	cmd := exec.Command("rm", args...)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println("Error in disk delete: ", err.Error())
		return err
	}
	fmt.Println(string(stdout))
	return nil
}

// This function takes a key name and searches for it in local node Keystore.
// returns nil if sucessfull & stores key as file in current dir.
func exportKey(keyName string) error {
	args := []string{"key", "export", keyName}
	cmd := exec.Command("ipfs", args...)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println("Error in export: ", err.Error())
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
		fmt.Println("Error in import exec: ", err)
		return err
	}
	fmt.Println(string(stdout))
	return nil
}

// This function generates a key using local node and returns it
// correct
func genKey(keyName string) (*shell.Key, error) {
	sh := shell.NewShell(localhost) //grab local node

	key, err := sh.KeyGen(context.Background(), keyName) //generate temp key to local node

	return key, err	
}