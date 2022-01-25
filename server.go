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
    "time"
	"io"
	// "encoding/hex"
	"os"
    // "encoding/json"

    // "github.com/gorilla/mux"
	
	shell "github.com/ipfs/go-ipfs-api"
    pb "github.com/ipfs/go-ipns/pb"
    ipns "github.com/ipfs/go-ipns"
    crypto "github.com/libp2p/go-libp2p-core/crypto"
)

const filepath = "Hello"
const ipfsPath = "/ipfs/QmSVjCYjy4jYZynyC2i5GeFgjhq1bLCK2vrkRz5ffnssqo"
const localhost = "localhost:5001"

func index(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: homePage")
}

// Correct! Adds file to IPFS, given path/filename,
// returns CID
func addToIPFS(file string) {
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
	fmt.Printf("added %s\n", cid)
}

// This function is needed to let the world know your Record exists.\
// May not be correct procedure
func PublishToIPNS(ipfsPath string, privateKey string) {
	sh := shell.NewShell(localhost)
    resp, err := sh.PublishWithDetails(ipfsPath, privateKey, time.Second, time.Second, false)
    if err != nil {
		log.Fatal(err)
	}

	if resp.Value != ipfsPath {
		fmt.Sprintf("Expected to receive %s but got %s", ipfsPath, resp.Value)
	}
	fmt.Println("response value: %s\n", resp.Value)
}

// CreateEntryWithEmbed shows how you can create an IPNS entry
// and embed it with a public key. For ed25519 keys this is not needed
// so attempting to embed with an ed25519 key, will not actually embed the key
func CreateEntryWithEmbed(ipfsPath string, privateKey crypto.PrivKey) (*pb.IpnsEntry, error) {
	ipfsPathByte := []byte(ipfsPath)
	eol := time.Now().Add(time.Hour * 48)
	entry, err := ipns.Create(privateKey, ipfsPathByte, 1, eol, 0)
	if err != nil {
		return nil, err
	}
	err = ipns.EmbedPublicKey(privateKey.GetPublic(), entry)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

// Create IPNS entry & embed Public key into entry, & upload to IPFS return entry to enter in contract. 
// ipfsPath string, privkey crypto.PrivKey
// TODO: verify parsing works correctly
// TODO: Verify ipnsRecord gets created properly 
func postKey(w http.ResponseWriter, r *http.Request) {
    // ******** parse data here *********
	var ipfsPath string
	var privkey crypto.PrivKey

	bodyBytes, err := io.ReadAll(r.Body)
	// verify there was no error
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	fmt.Printf("body string %s\n", bodyString)
	// ************************************

	
	// This line creates an IPNS record & embeds users' public key
    ipnsRecord, err := CreateEntryWithEmbed(ipfsPath, privkey)
    // Verify there was no error
	if err != nil {
	    panic(err)
    }

	// print to console.
    fmt.Printf("POST request successful %s\n", ipnsRecord)
    fmt.Printf("entry = %s\n")
}

// Generate private and public key to return to user, will need to use to post to IPFS.
// Function is correct for phase 1
// TODO: properly display contents of private and public keys
func getKey(w http.ResponseWriter, r *http.Request) {
    privateKey, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
    // verify there was no error
	if err != nil {
        panic(err)
    }

	// print to console. 
	fmt.Printf("Welcome to the IPNSKeyServer!\n")
	// spew.Printf("PrivateKey getKey(): %#+v", privateKey)
	fmt.Printf("Private key: %d \n", privateKey) //privateKey.GetPublic() returns the public key as well.
}


// Generate keys and embed records. Meant to test how keys are needed to be passed.
func testFunctions() {
	fmt.Println("Creating key:pair...")
	privKey, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	// rawKey, err := privKey.Raw()
	// hexkey := hex.EncodeToString(rawKey)
	// fmt.Printf("hex pub key: %s\n", hexkey)
    // verify there was no error
	if err != nil {
        panic(err)
    }

	fmt.Println("Creating IPNS record...", privKey)
	ipnsRecord, err := CreateEntryWithEmbed(ipfsPath, privKey)
    // // Verify there was no error
	// if err != nil {
	//     panic(err)
    // }

	fmt.Sprint(ipnsRecord)
	// PublishToIPNS(ipfsPath, hexkey)

	fmt.Printf("functions individually work.\n\n")
}

func main() {
	// handles api/website routes.
	// Used to test if keys need to be passed as objects or ints?
	addToIPFS(filepath)
	testFunctions()

	// router := mux.NewRouter().StrictSlash(true)
    // router.HandleFunc("/", index)
	// router.HandleFunc("/getKey", getKey).Methods("GET")
	// router.HandleFunc("/postKey", postKey).Methods("POST")

	// fmt.Printf("Starting server at port 8082\n")
	// log.Fatal(http.ListenAndServe(":8082", router))
}