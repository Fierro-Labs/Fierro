package main

// nodemon --exec go run server.go --signal SIGTERM

// curl -s --request   GET
//         --header    "Content-Type: application/json"
//         --write-out "\n%{http_code}\n"
//         http://localhost:13802/getKey

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
    "time"
	"context"
	// "io"
	// "io/ioutil"

    // "github.com/gorilla/mux"
	
	shell "github.com/ipfs/go-ipfs-api"
    pb "github.com/ipfs/go-ipns/pb"
    ipns "github.com/ipfs/go-ipns"
    ic "github.com/libp2p/go-libp2p-core/crypto"
	keystore "github.com/ipfs/go-ipfs-keystore"
	// peer "github.com/libp2p/go-libp2p-core/peer"
	// ke "github.com/ipfs/go-ipfs/core/commands/keyencode"
	// fsrepo "github.com/ipfs/go-ipfs/repo/fsrepo"
)

const file = "Hello"
const localhost = "localhost:5001"

var ipfsURI = "/ipfs/"
var ks *keystore.FSKeystore

func index(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: homePage")
}

type PublishResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}



// Adds file to IPFS, given path/filename,
// returns CID
// Correct
func addToIPFS(file string) (string, error) {
	sh := shell.NewShell(localhost)
	fileReader, err := os.Open(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s ", err)
		os.Exit(1)
	}
	cid, err:= sh.Add(fileReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s ", err)
		os.Exit(1)
	}
	fmt.Printf("added %s\n", ipfsURI + cid)
	ipfsPath := ipfsURI + cid
	return ipfsPath, err
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

// Custom publishing function returns the response object and error
// Correct?
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
// Correct?
func publishToIPNS(ipfsPath string, KeyName string) {
    pubResp, err := Publish(ipfsPath, KeyName)
    if err != nil {
		panic(err)
	}

	if pubResp.Value != ipfsPath {
		fmt.Printf("\nExpected to receive %s but got %s", ipfsPath, pubResp.Value)
	}

	fmt.Printf("\nresponse Name: %s\nresponse Value: %s\n", pubResp.Name, pubResp.Value)
}

// This function takes a key name and searches for it in local node Keystore.
// returns nil if sucessfull & stores key as file in current dir.
func exportKey(keyName string) error {
	args := []string{"key", "export", keyName}
	cmd := exec.Command("ipfs", args...)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(string(stdout))
	return nil
}

// This function takes a key name and file name, searches for it in current dir.
// returns nil if sucessfull & stores key in local node.
func importKey(keyName string, fileName string) error{
	args := []string{"key", "import", keyName, fileName}
	cmd := exec.Command("ipfs", args...)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(string(stdout))
	return nil
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


func main() {
	ipfsPath, err := addToIPFS(file)
	if err != nil{
		log.Panic(err)
	}
	// Used to test if keys need to be passed as objects or ints?
	testFunctions(ipfsPath)


	// handles api/website routes.
	// router := mux.NewRouter().StrictSlash(true)
    // router.HandleFunc("/", index)
	// router.HandleFunc("/getKey", getKey).Methods("GET")
	// router.HandleFunc("/postKey", postKey).Methods("POST")

	// fmt.Printf("Starting server at port 8082\n")
	// log.Fatal(http.ListenAndServe(":8082", router))
}