package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"

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
	_, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in exportKey: %s ", err)
		return err
	}
	return nil
}

// This function takes a key name and file name, searches for it in current dir.
// returns nil if sucessfull & stores key in local node.
func importKey(keyName string, fileName string) error {
	args := []string{"key", "import", keyName, fileName}
	cmd := exec.Command("ipfs", args...)
	_, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in importKey: %s", err)
		return err
	}
	return nil
}

// This function generates a key, exports it to "<keyName>.key" file in current dir, then delete from local keystore.
// returns newly generated key file to user
func GetKey(w http.ResponseWriter, r *http.Request) {
	keyName, ok := HasParam(r, "keyName")
	var key *shell.Key
	var err error

	if ok != true {
		keyName := addRdmSuffix("temp")
		key, err = genKey(keyName)
	} else {
		key, err = genKey(keyName)
	}

	if err != nil {
		writeJSONError(w, "Error in genKey", err)
		return
	}

	err = exportKey(key.Name)
	if err != nil {
		writeJSONError(w, "Error in keyName", err)
		return
	}

	err = deleteKey(key.Name) // delete temp key from local node keystore
	if err != nil {
		writeJSONError(w, "Error in deleteKey", err)
		return
	}

	// log.Println("Url Param 'keyName' is: " + string(keyName))
	w.Header().Set("Content-Disposition", "attachment; filename="+key.Name+".key")
	w.Header().Set("Content-Type", "application/octet-stream")
	// w.WriteHeader(http.StatusOK)
	http.ServeFile(w, r, key.Name+".key") // serve key to user to download

	err = diskDelete(key.Name + ".key") // delete key from disk
	if err != nil {
		panic(err)
	}
}

// This function will save a key to node, then delete the uploaded file from disk
// Returns 200 & key name as confirmation
func PostKey(w http.ResponseWriter, r *http.Request) {
	var dir = abs + "/KeyStore"

	FilePath, err := saveFile(r, dir, 32<<10) // grab uploaded .key file
	if err != nil {
		writeJSONError(w, "Error in saveFile", err)
		return
	}
	name := removeExtenstion(path.Base(FilePath))

	err = importKey(name, FilePath) //import key to local node keystore
	if err != nil {
		writeJSONError(w, "Error in importKey", err)
		return
	}

	err = diskDelete(FilePath) // delete key from disk
	if err != nil {
		writeJSONError(w, "Error in deleteKey", err)
		return
	}
	writeJSONSuccess(w, "Success - Saved key", name)
}

// This function will delete a key from the local node keystore
func DeleteKey(w http.ResponseWriter, r *http.Request) {
	keyName, ok := HasParam(r, "keyName")
	if ok != true {
		writeJSONError(w, keyName, nil)
		return
	}

	err := deleteKey(keyName) // delete temp key from local node keystore
	if err != nil {
		writeJSONError(w, "Error in deleteKey", err)
		return
	}
	writeJSONSuccess(w, "Success - Deleted key", keyName)
}
