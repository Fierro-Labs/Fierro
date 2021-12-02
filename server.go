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
    // "encoding/json"

    "github.com/gorilla/mux"
    pb "github.com/ipfs/go-ipns/pb"
    ipns "github.com/ipfs/go-ipns"
    crypto "github.com/libp2p/go-libp2p-core/crypto"
)

// func PublishToIPNS(ipfsPath string, privateKey string) (string, error) {
//     shell := NewShell("localhost:5001")

//     resp, err := shell.PublishWithDetails(ipfsPath, privateKey, time.Second, time.Second, false)
//     if err != nil {
// 		return nil, err
// 	}

// 	if resp.Value != examplesHashForIPNS {
// 		fmt.Sprintf("Expected to receive %s but got %s", examplesHash, resp.Value)
// 	}
// }



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
	err = ipns.EmbedPublicKey(privateKey.GenerateKeyPair, entry)
	if err != nil {
		return nil, err
	}

    err = PublishToIPNS(ipfsPath, privateKey)
	return entry, nil
}


// Create IPNS entry & embed Public key into entry, & upload to IPFS return entry to enter in contract. 
// ipfsPath string, privkey crypto.PrivKey
func postKey(w http.ResponseWriter, r *http.Request) {
    // parse data here
	var ipfsPath string
	var privkey crypto.PrivKey


	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	fmt.Printf(bodyString)

    ipnsRecord, err := CreateEntryWithEmbed(ipfsPath, privkey)
    if err != nil {
	    panic(err)
    }

    fmt.Printf("POST request successful\n")
    fmt.Printf("entry = %s\n")
}

// Generate private and public key to return to user, will need to use to post to IPFS.
func getKey(w http.ResponseWriter, r *http.Request) {
    privateKey, publicKey, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
    if err != nil {
        panic(err)
    }

	fmt.Printf("Welcome to the IPNSKeyServer!\n")
	fmt.Printf("Provate key: %d %d", privateKey, publicKey)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
    router.Handle("/", index)
	router.HandleFunc("/getKey", getKey).Methods("GET")
	router.HandleFunc("/postKey", postKey).Methods("POST")

	fmt.Printf("Starting server at port 8082\n")
	log.Fatal(http.ListenAndServe(":8082", router))
        
}

