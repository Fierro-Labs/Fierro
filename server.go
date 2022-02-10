package main

// nodemon --exec go run server.go --signal SIGTERM

// curl -s --request   GET
//         --header    "Content-Type: application/json"
//         --write-out "\n%{http_code}\n"
//         http://localhost:13802/getKey

import (
	"fmt"
	"log"
	"os"
	"os/exec"
    "time"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	// "bytes"
	// "io/ioutil"

    "github.com/gorilla/mux"
	
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

// This function deletes the exported key from current (default) dir
func localDelete(keyName string) error {
	args := []string{keyName+".key"}
	cmd := exec.Command("rm", args...)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println("Error in local delete: ", err.Error())
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
func importKey(keyName string, fileName string) error{
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

// This function parses the parameters in URL to grab the CID.
func GetCID(wrt http.ResponseWriter, req *http.Request) (string, error) {
	CIDs, ok := req.URL.Query()["CID"]

	// Query()["keyName"] will return an array of items,
	// we only want the single item.
	if !ok || len(CIDs[0]) < 1 {
		log.Println("Url Param 'CID' is missing")
		wrt.WriteHeader(http.StatusBadRequest)
		wrt.Header().Set("Content-Type", "application/octet-stream")
		resp := make(map[string]string)
		resp["message"] = "Status Bad Request - Missing CID"
		jsonResp, err:= json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error in JSON marshal. Err: %s", err)
			return "", err
		}
		wrt.Write(jsonResp)
	} 
	return CIDs[0], nil
		
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

// This function generates a key, exports it to "<keyName>.key" file in current dir, then delete from local keystore.
// returns newly generated key file to user
func GetKey(w http.ResponseWriter, req *http.Request) {
	keyName := "temp" // user input from API or self-generated non-clashing name
	w.Header().Set("Content-Disposition", "attachment; filename="+string(keyName))
	w.Header().Set("Content-Type", "application/octet-stream")
	fmt.Println("Generating key...")
	key, err := genKey(keyName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/octet-stream")
		resp := make(map[string]string)
		resp["message"] = "Status Bad Request"
		jsonResp, err:= json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
	}

	fmt.Println("Exporting key...")
	err = exportKey(keyName)
	if err != nil {
		panic(err)
		return
	}

	fmt.Println("Deleting key...")
	err = deleteKey(key.Name) // delete temp key
	if err != nil {
		panic(err)
		return
	}

	log.Println("Url Param 'keyName' is: " + string(keyName))
	http.ServeFile(w, req, string(keyName) + ".key") // serve key to user to download

	fmt.Println("Deleting exported file...")
	err = localDelete(keyName)
	if err != nil {
		panic(err)
		return
	}
	return
}


// This function take a file and path and publish to ipfs
// Returns ACK & IPNS path
func PostKey(w http.ResponseWriter, r *http.Request) {
	var FileName string
	fmt.Println("PostKey...")
	r.ParseMultipartForm(32 << 10) // limit max input length
	// var buf bytes.Buffer
	// in your case file would be fileupload
	file, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	FileName = header.Filename
	nameArr := strings.Split(header.Filename, ".")
	name := nameArr[0]
	// fmt.Printf("File: %s\nFile name: %s\nFile extension: %s\n", header.Filename, name[0], name[1])
	// Copy the file data to my buffer
	f, err := os.OpenFile("KeyStore/"+FileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Error in OpenFile: ", err)
		return
	}
	io.Copy(f, file)
	
	CID, err := GetCID(w,r)
	if err != nil {
		fmt.Println("Error in GetCID: ", err)
		return
	}

	err = importKey(name, "KeyStore/"+FileName)
	if err != nil {
		fmt.Println("Error in importKey: ", err)
		return
	}



	publishToIPNS(ipfsURI + CID, FileName)
	// do something with the contents...
	// I normally have a struct defined and unmarshal into a struct, but this will
	// work as an example
	// contents := buf.String()
	// fmt.Println(contents)
	// I reset the buffer in case I want to use it again
	// reduces memory allocations in more intense projects
	// buf.Reset()
	// do something else
	// etc write header
	return
}



func main() {
	// ipfsPath, err := addToIPFS(file)
	// if err != nil{
	// 	log.Panic(err)
	// }
	// Used to test if keys need to be passed as objects or ints?
	// testFunctions(ipfsPath)


	// handles api/website routes.
	router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", index)
	router.HandleFunc("/getKey", GetKey).Methods("GET")
	router.HandleFunc("/postKey", PostKey).Methods("POST")

	fmt.Printf("Starting server at port 8082\n")
	log.Fatal(http.ListenAndServe(":8082", router))
}