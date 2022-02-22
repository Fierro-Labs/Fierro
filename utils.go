package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"context"
	"encoding/json"
	"io"


	shell "github.com/ipfs/go-ipfs-api"

)

// This function will set the appropiate headers for when there is an error.
func writeJSONError(w http.ResponseWriter, msg string, err error) {
	resp := make(map[string]interface{})
	resp["message"] = msg
	resp["error"] = err
	jsonResp, err:= json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error in JSON marshal. Err: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

// This function will set the appropiate headers for when there is an error.
func writeJSONSuccess(w http.ResponseWriter, msg string, val string) {
	resp := make(map[string]string)
	resp["message"] = msg
	resp["value"] = val
	jsonResp, err:= json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error in JSON marshal. Err: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

// This function grabs the specified parameter value out the URL
func GetParam(r *http.Request, parameter string) (string, bool) {
	params, ok := r.URL.Query()[parameter]
	// Query()[parameter] will return an array of items,
	// we only want the single item.
	if !ok || len(params[0]) < 1 {
		fmt.Println("Error: Missing " + parameter)
		return "Missing "+parameter+" parameter", !ok
	} 
	return params[0], ok
}

// Grabs file with specified file size and save to specified dir
// returns name of file
func saveFile(r *http.Request, dir string, size int64) (string, error) {
	var FileName string

	fmt.Println("Saving file...")
	r.ParseMultipartForm(size) // limit max input length
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s ", err)
		return "", err
	}
	defer file.Close()
	FileName = header.Filename
	
	f, err := os.OpenFile(dir+"/"+FileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Error in OpenFile: ", err)
		return "", err
	}
	io.Copy(f, file)

	return FileName, nil
}

// Generate keys and embed records. Meant to test how keys are needed to be passed.
func testFunctions(ipfsPath string) {
	const KeyName = "temp"
	sh := shell.NewShell(localhost)

	key, err := sh.KeyGen(context.Background(), KeyName) //generate temp key to local node
	if err != nil {
		panic(err)
	}

	// publishToIPNS(ipfsPath, key.Name) // publish ipfsPath to ipfs under temp key
	fmt.Println("Exporting key...")
	err = exportKey(KeyName)
	if err != nil {
		panic(err)
	}

	fmt.Println("Deleting key...")
	err = deleteKey(key.Name) // delete temp key
	if err != nil {
		panic(err)
	}

	err = importKey(key.Name, key.Name + ".key")
	if err != nil {
		panic(err)
	}
	fmt.Printf("functions individually work.\n\n")
}