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
	"io"
	"os"
	"io/ioutil"
    "time"
	"context"

    // "github.com/gorilla/mux"
	
	shell "github.com/ipfs/go-ipfs-api"
    pb "github.com/ipfs/go-ipns/pb"
    ipns "github.com/ipfs/go-ipns"
    ic "github.com/libp2p/go-libp2p-core/crypto"
	keystore "github.com/ipfs/go-ipfs-keystore"
	peer "github.com/libp2p/go-libp2p-core/peer"
)

const filepath = "Hello"
const ipfsPath = "/ipfs/QmSVjCYjy4jYZynyC2i5GeFgjhq1bLCK2vrkRz5ffnssqo"
const localhost = "localhost:5001"
var ks *keystore.FSKeystore

func index(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: homePage")
}

type KeyOutput struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type KeyOutputList struct {
	Keys []KeyOutput
}

var out struct {
	Path string
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

// CreateEntryWithEmbed shows how you can create an IPNS entry
// and embed it with a public key. For ed25519 keys this is not needed
// so attempting to embed with an ed25519 key, will not actually embed the key
func createEntry(ipfsPath string, sk ic.PrivKey) (*pb.IpnsEntry, error) {
	ipfsPathByte := []byte(ipfsPath)
	eol := time.Now().Add(time.Hour * 48)
	entry, err := ipns.Create(sk, ipfsPathByte, 1, eol, 0)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

// This function is needed to let the world know your Record exists.
// May not be correct procedure
func publishToIPNS(ipfsPath string, KeyName string) {
	sh := shell.NewShell(localhost)
    resp, err := sh.PublishWithDetails(ipfsPath, KeyName, time.Second, time.Second, false)
    if err != nil {
		panic(err)
	}

	if resp.Value != ipfsPath {
		fmt.Printf("\nExpected to receive %s but got %s", ipfsPath, resp.Value)
	}
	fmt.Println("response value: %s\n", resp.Value)
}


// Generate keys and embed records. Meant to test how keys are needed to be passed.
func testFunctions(ks keystore.Keystore) {
	const KeyName = "temp"
	var KeyGenResp KeyOutput
	// var KeyOutResp KeyOutput

	sh := shell.NewShell(localhost)

	// req := sh.api.repo.Keystore.Get("self")
	// err := req.Exec(context.Background(), &KeyGenResp)
	fmt.Println("Creating key:pair...")
	sk, _, err := ic.GenerateKeyPair(ic.Ed25519, 256)
	if err != nil {
		panic(err)
	}
	fmt.Println(KeyGenResp.Name, KeyGenResp.Id)
	pid, err := peer.IDFromPublicKey(sk.GetPublic()) //"create" peerID from the private key
	peerID := pid.Pretty() //convert type peer.ID to string
	fmt.Printf("PeerID from key: %s\n", peerID)

	err = ks.Put(KeyName, sk)
	hasKey, err := ks.Has(KeyName)
	fmt.Printf("Ks has key: %t\n", hasKey)

	// At this point the newly generated key should be in my local repo
	// This should enable me to publish to ipfs using this key
	
	fmt.Println("Creating IPNS record...")
	ipnsRecord, err := createEntry(ipfsPath, sk)
	if err != nil {
	    panic(err)
    }

	fmt.Printf("IPNS value: %s\n", ipnsRecord.Value)

	publishToIPNS(string(ipnsRecord.Value), KeyName)

	err = sh.Request("resolve", peerID).Exec(context.Background(), &out)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Resolve ipns key: %s\n", out.Path)
	// rm temp key 
	err = ks.Delete(KeyName)
	if err != nil {
		panic(err)
	}

	fmt.Printf("functions individually work.\n\n")
}

// Create IPNS entry & embed Public key into entry, & upload to IPFS return entry to enter in contract. 
// ipfsPath string, sk ic.PrivKey
// TODO: verify parsing works correctly
// TODO: Verify ipnsRecord gets created properly 
func postKey(w http.ResponseWriter, r *http.Request) {
    // ******** parse data here *********
	// var ipfsPath string
	// var sk ic.PrivKey

	bodyBytes, err := io.ReadAll(r.Body)
	// verify there was no error
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	fmt.Printf("body string %s\n", bodyString)
	// ************************************

	
	// This line creates an IPNS record & embeds users' public key
    // ipnsRecord, err := createEntry(ipfsPath, sk)
    // Verify there was no error
	if err != nil {
	    panic(err)
    }

	// print to console.
    // fmt.Printf("POST request successful %s\n", ipnsRecord)
    fmt.Printf("entry = %s\n")
}

// Generate private and public key to return to user, will need to use to post to IPFS.
// Function is correct for phase 1
// TODO: properly display contents of private and public keys
func getKey(w http.ResponseWriter, r *http.Request) {
    sk, _, err := ic.GenerateKeyPair(ic.RSA, 2048)
    // verify there was no error
	if err != nil {
        panic(err)
    }

	// print to console. 
	fmt.Printf("Welcome to the IPNSKeyServer!\n")
	// spew.Printf("PrivateKey getKey(): %#+v", sk)
	fmt.Printf("Private key: %d \n", sk) //sk.GetPublic() returns the public key as well.
}

func main() {
	tdir, err := ioutil.TempDir("", "keystore-test")
	if err != nil {
		log.Fatal(err)
	}
	ks, err := keystore.NewFSKeystore(tdir)
	if err != nil {
		log.Fatal(err)
	}
	// handles api/website routes.
	addToIPFS(filepath)
	// Used to test if keys need to be passed as objects or ints?
	testFunctions(ks)

	// router := mux.NewRouter().StrictSlash(true)
    // router.HandleFunc("/", index)
	// router.HandleFunc("/getKey", getKey).Methods("GET")
	// router.HandleFunc("/postKey", postKey).Methods("POST")

	// fmt.Printf("Starting server at port 8082\n")
	// log.Fatal(http.ListenAndServe(":8082", router))
}