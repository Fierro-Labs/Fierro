package main

import (
	"fmt"
	"os"
    "time"
	"context"


	shell "github.com/ipfs/go-ipfs-api"
    pb "github.com/ipfs/go-ipns/pb"
    ipns "github.com/ipfs/go-ipns"
    ic "github.com/libp2p/go-libp2p-core/crypto"
)

// var ipfsURI = "/ipfs/"

type PublishResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Create an IPNS entry with a 2 day lifespan before needing to revive
// Correct
func createEntry(ipfsPath string, sk ic.PrivKey) (*pb.IpnsEntry, error) {
	ipfsPathByte := []byte(ipfsPath)
	eol := time.Now().Add(time.Hour * 48)
	entry, err := ipns.Create(sk, ipfsPathByte, 1, eol, 0)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

// Adds file to IPFS, given path/filename,
// returns CID
// Correct
func addToIPFS(file string) (string, error) {
	sh := shell.NewShell(localhost)
	fileReader, err := os.Open(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s ", err)
		panic(err)
		return "", err
	}
	cid, err:= sh.Add(fileReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s ", err)
		panic(err)
		return "", err
	}
	fmt.Printf("added %s\n", ipfsURI + cid)
	ipfsPath := ipfsURI + cid
	return ipfsPath, err
}

// Custom publishing function returns the response object and error
// Kept as close as possible to Publish method found at gh.com/ipfs/go-ipns
// Correct
func Publish(contentHash string, key string) (*PublishResponse, error) {
	var pubResp PublishResponse
	sh := shell.NewShell(localhost)
	req := sh.Request("name/publish", contentHash).Option("key", key)
	req.Option("resolve", true)
	err := req.Exec(context.Background(), &pubResp)
	if err != nil {
		panic(err)
		return nil, err
	}

	return &pubResp, nil
}

// This function is needed to let the world know your Record exists.
// Correct
func publishToIPNS(ipfsPath string, KeyName string) (*PublishResponse, error) {
    pubResp, err := Publish(ipfsPath, KeyName)
    if err != nil {
		panic(err)
	}

	if pubResp.Value != ipfsPath {
		fmt.Printf("\nExpected to receive %s but got %s", ipfsPath, pubResp.Value)
	}

	fmt.Printf("\nresponse Name: %s\nresponse Value: %s\n", pubResp.Name, pubResp.Value)
	
	return pubResp, nil
}