package main

import (
	"fmt"
	"os"
	"context"

	shell "github.com/ipfs/go-ipfs-api"
)

// var ipfsURI = "/ipfs/"

type PublishResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ResolvedPath struct {
	Path string
}


// Adds file to IPFS, given path/filename,
// returns CID
func addToIPFS(file string) (string, error) {
	sh := shell.NewShell(localhost)
	fileReader, err := os.Open(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: %s ", err)
		return "", err
	}
	cid, err:= sh.Add(fileReader)
	if err != nil {
		fmt.Println("Error in adding file to ipfs: ", err)
		return "", err
	}
	ipfsPath := ipfsURI + cid
	fmt.Printf("Added %s\n", ipfsPath)
	return ipfsPath, nil
}

// Custom publishing function returns the response object and error
// Kept as close as possible to Publish method found at gh.com/go-ipfs-api
func PublishToIPFS(contentHash string, key string) (*PublishResponse, error) {
	var pubResp PublishResponse
	sh := shell.NewShell(localhost)
	req := sh.Request("name/publish", contentHash).Option("key", key)
	req.Option("resolve", true)
	err := req.Exec(context.Background(), &pubResp)
	if err != nil {
		return nil, err
	}

	return &pubResp, nil
}

// This function is needed to let the world know your Record exists.
func publishToIPNS(ipfsPath string, KeyName string) (*PublishResponse, error) {
    pubResp, err := PublishToIPFS(ipfsPath, KeyName)
    if err != nil {
		fmt.Fprintf(os.Stderr, "Error in Publish: %s ", err)
		return nil, err
	}

	if pubResp.Value != ipfsPath {
		fmt.Printf("\nExpected to receive %s but got %s", ipfsPath, pubResp.Value)
		return nil, err
	}

	fmt.Printf("\nresponse Name: %s\nresponse Value: %s\n", pubResp.Name, pubResp.Value)
	
	return pubResp, nil
}

// This function will resolve/download the content pointed to by the record. 
func resolve(ipnsKey string) (string, error) {
	var path ResolvedPath
	sh := shell.NewShell(localhost)
	req := sh.Request("name/resolve", ipnsKey).Option("dht-timeout", "180s") // timeout after 3 minutes
	err := req.Exec(context.Background(), &path)
	if err != nil {
		return "", err
	}
	return path.Path, nil
}
