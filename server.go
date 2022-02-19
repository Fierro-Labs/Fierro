package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

    "github.com/gorilla/mux"
	"github.com/gammazero/deque"
	"github.com/robfig/cron"
	shell "github.com/ipfs/go-ipfs-api"
)

const FILE = "Hello"
const localhost = "localhost:5001"

const ipfsURI = "/ipfs/"
const ipnsURI = "/ipns/"

var q deque.Deque


func index(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: HomePage")
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

// This function generates a key, exports it to "<keyName>.key" file in current dir, then delete from local keystore.
// returns newly generated key file to user
func GetKey(w http.ResponseWriter, r *http.Request) {
	keyName := "temp" // user input from API or self-generated non-clashing name
	fmt.Println("Generating key...")
	key, err := genKey(keyName)
	if err != nil {
		writeJSONError(w, "Error in genKey", err)
		return
	}
	
	fmt.Println("Exporting key...")
	err = exportKey(keyName)
	if err != nil {
		writeJSONError(w, "Error in keyName", err)
		return
	}
	
	fmt.Println("Deleting key...")
	err = deleteKey(key.Name) // delete temp key from local node keystore
	if err != nil {
		writeJSONError(w, "Error in deleteKey", err)
		return
	}


	// log.Println("Url Param 'keyName' is: " + string(keyName))
	w.Header().Set("Content-Disposition", "attachment; filename="+string(keyName))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	http.ServeFile(w, r, string(keyName) + ".key") // serve key to user to download


	fmt.Println("Deleting exported key...")
	err = diskDelete(keyName+".key") // delete key from disk
	if err != nil {
		panic(err)
		return
	}

	// fmt.Println("Deleting exported key a second time...")
	// err = diskDelete(keyName) // delete key from disk again to force error
	// if err != nil {
	// 	writeJSONError(w, "Error in deleteKey", err)
	// 	return
	// }
}

// This function will save a key to node, then delete the uploaded file from disk
// Returns 200 & key name as confirmation
func PostKey(w http.ResponseWriter, r *http.Request) {
	const dir = "KeyStore"
	FileName, err := saveFile(r, dir, 32 << 10) // grab uploaded .key file
	if err != nil {
		writeJSONError(w, "Error in saveFile", err)
		return
	}
	name := strings.Split(FileName, ".")[0]

	fmt.Println("Importing Key...")
	err = importKey(name, dir+"/"+FileName) //import key to local node keystore
	if err != nil {
		writeJSONError(w, "Error in importKey", err)
		return
	}

	fmt.Println("Deleting saved key from disk...")
	err = diskDelete(dir+"/"+FileName) // delete key from disk
	if err != nil {
		writeJSONError(w, "Error in deleteKey", err)
		return
	}
	writeJSONSuccess(w, "Success - Saved key", name)
}

// This function will delete a key from the local node keystore
func DeleteKey(w http.ResponseWriter, r *http.Request) {
	keyName, ok := GetParam(r, "keyName")
	if ok != true {
		writeJSONError(w, keyName, nil)
		return
	}

	fmt.Println("Deleting key...")
	err := deleteKey(keyName) // delete temp key from local node keystore
	if err != nil {
		writeJSONError(w, "Error in deleteKey", err)
		return
	}
	writeJSONSuccess(w, "Success - Deleted key", keyName)
}

// This function will accept a ipns key and resolve it
// returns the ipfs path it resolved.
func GetRecord(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting IPNS Key...")
	ipnsKey, ok := GetParam(r, "ipnskey") // grab ipnskey from query parameter
	if ok != true {
		writeJSONError(w, "Error with getting ipnskey", nil)
		return
	}

	fmt.Println("Resolving IPNS Record...")
	ipfsPath, err := resolve(ipnsKey) // download content and return ipfs path
	if err != nil {
		writeJSONError(w, "Error in resolve", err)
		return
	}

	writeJSONSuccess(w, "Success - GetRecord", ipfsPath)
}

// This function takes a CID and file.key and publishes brand new IPNS records to IPFS
// IPFS Node handles republishing automatically in the background as long as it is up and running
// Returns ACK & IPNS path
func PostRecord(w http.ResponseWriter, r *http.Request) {
	const dir = "KeyStore"
	
	fmt.Println("Getting CID...")
	CID, ok := GetParam(r, "CID") // grab CID from query parameter
	if ok != true {
		writeJSONError(w, "Error with getting CID: "+CID, nil)
		return
	}

	FileName, err := saveFile(r, dir, 32 << 10) // grab uploaded .key file
	if err != nil {
		writeJSONError(w, "Error in saveFile", err)
		return
	}
	name := strings.Split(FileName, ".")[0]

	fmt.Println("Importing Key...")
	err = importKey(name, dir+"/"+FileName) //import key to local node keystore
	if err != nil {
		writeJSONError(w, "Error in importKey", err)
		return
	}

	fmt.Println("Publishing to IPNS...")
	pubResp, err := publishToIPNS(ipfsURI + CID, name) //publish IPNS record to IPFS
	if err != nil {
		writeJSONError(w, "Error in publishToIPNS", err)
		return
	}

	fmt.Println("Deleting exported key...")
	err = diskDelete(dir+"/"+FileName) // delete key from disk
	if err != nil {
		writeJSONError(w, "Error in diskDelete", err)
		return
	}

	fmt.Printf("\nresponse Name: %s\nresponse Value: %s\n", pubResp.Name, pubResp.Value)
	writeJSONSuccess(w, "Success - PostKey", ipnsURI+pubResp.Name)
}

// This function takes an IPNS Key and file.key and resolves IPNS record
// IPFS Node handles republishing automatically in the background as long as it is up and running
// Returns ACK & resolved content
func PutRecord(w http.ResponseWriter, r *http.Request) {
	const dir = "KeyStore"
	
	fmt.Println("Getting IPNS Key...")
	key, ok := GetParam(r, "ipnskey") // grab key from query parameter
	if ok != true {
		writeJSONError(w, "Error with getting key: "+key, nil)
		return
	}

	FileName, err := saveFile(r, dir, 32 << 10) // grab uploaded .key file
	if err != nil {
		writeJSONError(w, "Error in saveFile", err)
		return
	}
	name := strings.Split(FileName, ".")[0]

	fmt.Println("Importing Key...")
	err = importKey(name, dir+"/"+FileName) //import key to local node keystore
	if err != nil {
		writeJSONError(w, "Error in importKey", err)
		return
	}

	fmt.Println("Resolving IPNS Record...")
	ipfsPath, err := resolve(key) // download content and return ipfs path
	if err != nil {
		writeJSONError(w, "Error in resolve", err)
		return
	}

	fmt.Println("Deleting saved key from disk...")
	err = diskDelete(dir+"/"+FileName) // delete key from disk
	if err != nil {
		writeJSONError(w, "Error in deleteKey", err)
		return
	}
	writeJSONSuccess(w, "Success - PutRecord", ipfsPath)
}

// This function will take an ipns key and add it to the queue to be resolved later
func FollowRecord(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting IPNS Key...")
	ipnsKey, ok := GetParam(r, "ipnskey") // grab ipnskey from query parameter
	if ok != true {
		writeJSONError(w, "Error with getting ipnskey", nil)
		return
	}

	q.PushBack(ipnsKey) // add record to queue

	writeJSONSuccess(w, "Success - Will continue to check back later", ipnsKey)
}

// This function will check for if a record is in queue and delete it
// Meant for manual deletion of record 
// func StopFollowing(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Getting IPNS Key...")
// 	ipnsKey, ok := GetParam(r, "ipnskey") // grab ipnskey from query parameter
// 	if ok != true {
// 		writeJSONError(w, "Error with getting ipnskey", nil)
// 		return
// 	}

// 	if q.Index(ipnsKey) == -1 {
// 		writeJSONError(w, "Key not in queue", nil)
// 	} else {
// 		q.Remove(q.Index(func (ipnsKey interface{} bool) {return ipnsKey}))
// 		writeJSONSuccess(w, "Success - Removed key from queue", ipnsKey)
// 	}
// 	return
// }

// Helper function so that we can call internally and externally
// returns the key at front of queue
func follow(internal bool) string {
	defer func() { // <- executes only if there is a panic when using the queue
		if (recover() != nil){
			fmt.Sprintf("Error with queue operations %s")
			return
		}
	}()
	ipnsKeyInt := q.PopFront()
	ipnsKey := fmt.Sprintf("%s", ipnsKeyInt) // Used to convert interface to string
	// fmt.Printf("Success - top key is: %s\n", ipnsKey)
	q.PushBack(ipnsKey)
	if internal {
		ipfsPath, err := resolve(ipnsKey)
		if err != nil {
			fmt.Printf("Error in resolve %s", err)
			return ""
		}
		return ipfsPath
	}
	return ipnsKey
}

// Grab uploaded file and add to ipfs
// returns ipfsPath (ipfsURI + CID)
func AddFile(w http.ResponseWriter, r *http.Request) {
	const dir = "Uploads"

	FileName, err := saveFile(r, dir, 32 << 20) // grab uploaded content & save to disk
	if err != nil {
		writeJSONError(w, "Error in saveFile", err)
		return
	}

	fmt.Println("Adding to IPFS...")
	ipfsPath, err := addToIPFS(dir+"/"+FileName) // add content file to IPFS
	if err != nil {
		writeJSONError(w, "Error in addToIPFS", err)
		return
	}
	// TODO: Return this over http.response
	fmt.Println("ipfs Path: ", ipfsPath)
	writeJSONSuccess(w, "Success - addFile", ipfsPath)
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


func main() {
	// Used to test out individual functions
	// testFunctions(ipfsPath)

	// create, init, & start cron job
	c := cron.New()
	c.AddFunc("@every 2m", func() { 
		ipfsPath := follow(true)
		fmt.Println(ipfsPath)	
	})
	c.Start()
	
	// Watch for ctrl^c, close out cron job
	channel := make(chan os.Signal, 1) 
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func(){
		<-channel
		c.Stop()
		os.Exit(1)
	}()

	// q.PushBack("k51qzi5uqu5dm876hw4kh2mn58rnajofhoohohymt9bui38q6ogsa0rrct6fnh")
	// handles api/website routes
	router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", index)
	router.HandleFunc("/getKey", GetKey).Methods("GET")
	router.HandleFunc("/postKey", PostKey).Methods("POST")
	router.HandleFunc("/deleteKey", DeleteKey).Methods("DELETE")
	router.HandleFunc("/postRecord", PostRecord).Methods("POST")
	router.HandleFunc("/putRecord", PutRecord).Methods("PUT")
	router.HandleFunc("/getRecord", GetRecord).Methods("GET")
	router.HandleFunc("/followRecord", FollowRecord).Methods("GET")
	router.HandleFunc("/addFile", AddFile).Methods("POST")


	fmt.Printf("Starting server at port 8082\n")
	log.Fatal(http.ListenAndServe(":8082", router))
}