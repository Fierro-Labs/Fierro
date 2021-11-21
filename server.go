package main

// nodemon --exec go run server.go --signal SIGTERM

import (
	"fmt"
	"log"
	"net/http"
    "time"
    "encoding/json"

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
func CreateEntryWithEmbed(ipfsPath string, publicKey crypto.PubKey, privateKey crypto.PrivKey) (*pb.IpnsEntry, error) {
	ipfsPathByte := []byte(ipfsPath)
	eol := time.Now().Add(time.Hour * 48)
	entry, err := ipns.Create(privateKey, ipfsPathByte, 1, eol, 0)
	if err != nil {
		return nil, err
	}
	err = ipns.EmbedPublicKey(publicKey, entry)
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



    ipnsRecord, err := CreateEntryWithEmbed(ipfsPath, privkey.GetPublic(), privkey)
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

    // might have to use Fprintf to write to w, to return as api
	fmt.Printf("Welcome to the IPNSKeyServer!\n")
	fmt.Printf("Provate key: %d %d", privateKey, publicKey)
}

func main() {
	myRouter := mux.NewRouter().StrictSlash(true)
    mux.HandleFunc("/", index.html) // New code
	mux.HandleFunc("/postKey", postKey)
	mux.HandleFunc("/getKey", getKey)

	fmt.Printf("Starting server at port 8082\n")
	log.Fatal(http.ListenAndServe(":8082", myRouter))
        
}

