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

var ipfsURI = "/ipfs/"

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
func GetCID(wrt http.ResponseWriter, req *http.Request) (string, error) {
	CIDs, ok := req.URL.Query()["CID"]

	// Query()["keyName"] will return an array of items,
	// we only want the single item.
	if !ok || len(CIDs[0]) < 1 {
		// TODO: Make this its own function
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

// This function generates a key, exports it to "<keyName>.key" file in current dir, then delete from local keystore.
// returns newly generated key file to user
// Correct
func GetKey(w http.ResponseWriter, req *http.Request) {
	keyName := "temp" // user input from API or self-generated non-clashing name
	w.Header().Set("Content-Disposition", "attachment; filename="+string(keyName))
	w.Header().Set("Content-Type", "application/octet-stream")
	fmt.Println("Generating key...")
	key, err := genKey(keyName)
	if err != nil {
		panic(err)
		return
	}

	fmt.Println("Exporting key...")
	err = exportKey(keyName)
	if err != nil {
		panic(err)
		return
	}

	fmt.Println("Deleting key...")
	err = deleteKey(key.Name) // delete temp key from local node keystore
	if err != nil {
		panic(err)
		return
	}

	log.Println("Url Param 'keyName' is: " + string(keyName))
	http.ServeFile(w, req, string(keyName) + ".key") // serve key to user to download

	fmt.Println("Deleting exported file...")
	err = diskDelete(keyName) // delete key from disk
	if err != nil {
		panic(err)
		return
	}
	return
}

// This function takes a CID and file.key and publishes IPNS record to IPFS
// Returns ACK & IPNS path
// Correct
func PostKey(w http.ResponseWriter, r *http.Request) {
	const dir = "KeyStore"
	FileName, err := saveFile(w, r, dir, 32 << 10) // grab uploaded .key file
	if err != nil {
		fmt.Println("Error in saveFile: ", err)
		return
	}
	name := strings.Split(FileName, ".")[0]

	fmt.Println("Getting CID...")
	CID, err := GetCID(w,r) // grab CID from query parameter
	if err != nil {
		fmt.Println("Error in GetCID: ", err)
		return
	}

	fmt.Println("Importing Key...")
	err = importKey(name, dir+"/"+FileName) //import key to local node keystore
	if err != nil {
		fmt.Println("Error in importKey: ", err)
		return
	}

	fmt.Println("Publishing to IPNS...")
	pubResp, err := publishToIPNS(ipfsURI + CID, name) //publish IPNS record to IPFS
	if err != nil {
		fmt.Println("Error in publishtoIPNS: ", err)
		return
	}

	fmt.Printf("\nresponse Name: %s\nresponse Value: %s\n", pubResp.Name, pubResp.Value)

	return
}

// Grab uploaded file and add to ipfs
// returns ipfsPath (ipfsSURI + CID)
// Correct
func addFile(w http.ResponseWriter, r *http.Request) {
	const dir = "Uploads"

	FileName, err := saveFile(w, r, dir, 32 << 20) // grab uploaded content & save to disk
	if err != nil {
		fmt.Println("Error in saveFile: ", err)
		return
	}

	ipfsPath, err := addToIPFS(dir+"/"+FileName) // add content file to IPFS
	if err != nil {
		fmt.Println("Error in addToIPFS: ", err)
		return
	}
	// TODO: Return this over http.response
	fmt.Println("ipfs Path: ", ipfsPath)
	return
}

// Grabs file with specified file size and save to specified dir
// returns name of file
func saveFile(w http.ResponseWriter, r *http.Request, dir string, size int64) (string, error) {
	var FileName string

	fmt.Println("Saving file...")
	r.ParseMultipartForm(size) // limit max input length
	file, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
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

func main() {
	// Used to test if keys need to be passed as objects or ints?
	// testFunctions(ipfsPath)


	// handles api/website routes.
	router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", index)
	router.HandleFunc("/getKey", GetKey).Methods("GET")
	router.HandleFunc("/postKey", PostKey).Methods("POST")
	router.HandleFunc("/addFile", addFile).Methods("POST")

	fmt.Printf("Starting server at port 8082\n")
	log.Fatal(http.ListenAndServe(":8082", router))
}