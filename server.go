package main

import (
	"fmt"
	"log"
	"os"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

    "github.com/gorilla/mux"
	shell "github.com/ipfs/go-ipfs-api"
)

const FILE = "Hello"
const localhost = "localhost:5001"

const ipfsURI = "/ipfs/"
const ipnsURI = "/ipns/"


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

// This function parses the parameters in URL to grab the CID.
// Correct
func GetCID(r *http.Request) (string, bool) {
	CIDs, ok := r.URL.Query()["CID"]
	// Query()["keyName"] will return an array of items,
	// we only want the single item.
	if !ok || len(CIDs[0]) < 1 {
		fmt.Println("Error: Missing CID")
		return "Missing CID parameter", !ok
	} 
	return CIDs[0], ok
}

// This function generates a key, exports it to "<keyName>.key" file in current dir, then delete from local keystore.
// returns newly generated key file to user
// Correct
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
	err = diskDelete(keyName) // delete key from disk
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

// This function takes a CID and file.key and publishes IPNS record to IPFS
// Returns ACK & IPNS path
// Need to test error cases and responses.
func PostKey(w http.ResponseWriter, r *http.Request) {
	const dir = "KeyStore"
	FileName, err := saveFile(r, dir, 32 << 10) // grab uploaded .key file
	if err != nil {
		writeJSONError(w, "Error in saveFile", err)
		return
	}
	name := strings.Split(FileName, ".")[0]

	fmt.Println("Getting CID...")
	CID, ok := GetCID(r) // grab CID from query parameter
	if ok != true {
		writeJSONError(w, CID, nil)
		return
	}

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

	fmt.Printf("\nresponse Name: %s\nresponse Value: %s\n", pubResp.Name, pubResp.Value)
	writeJSONSuccess(w, "Success - PostKey", ipnsURI+pubResp.Name)
}

// Grab uploaded file and add to ipfs
// returns ipfsPath (ipfsURI + CID)
// Correct
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
	// Used to test if keys need to be passed as objects or ints?
	// testFunctions(ipfsPath)


	// handles api/website routes.
	router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", index)
	router.HandleFunc("/getKey", GetKey).Methods("GET")
	router.HandleFunc("/postKey", PostKey).Methods("POST")
	router.HandleFunc("/addFile", AddFile).Methods("POST")

	fmt.Printf("Starting server at port 8082\n")
	log.Fatal(http.ListenAndServe(":8082", router))
}