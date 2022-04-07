import (
	"io/ioutil"
	"net/http"
	"encoding/json"

	ic "github.com/libp2p/go-libp2p-core/crypto"
	keystore "github.com/ipfs/go-ipfs-keystore"
	fsrepo "github.com/ipfs/go-ipfs/repo/fsrepo"
	pb "github.com/ipfs/go-ipns/pb"
)

type IpnsEntry struct {
	Name  string
	Value string
}
var ks *keystore.FSKeystore

// Put a key in custom keystore
// Will run into problems when publishing, because the publish function won't know where to find keys.
func putKeyv0(ks keystore.Keystore, keyName string, sk ic.PrivKey) {
	err = ks.Put(KeyName, sk) //insert key to new keystore
	if err != nil {
		panic(err)
	}
}

// Example function to show how to create a custom keystore
func makeKeystorev0() {
	tdir, err := ioutil.TempDir("", "keystore-test")
	if err != nil {
		log.Fatal(err)
	}
	ks, err := keystore.NewFSKeystore(tdir)
	if err != nil {
		log.Fatal(err)
	}
}

// Put key in local node keystore programatically
// Won't work bc daemon will have lock on repo.
func putKeyV1(r fsrepo, keyName string, sk ic.PrivKey){
	err = r.Keystore().Put(KeyName, sk) //insert key to new keystore
	if err != nil {
		panic(err)
	}
}

// Example function to show how to grab local keystore.
func makeKeystorev1() {
	// cfgRoot, err := cmdenv.GetConfigRoot(env)
	// if err != nil {
	// 	panic(err)
	// }
	repo, err := fsrepo.Open("~/.ipfs") // or fsrepo.Open(cfgRoot)
	if err != nil {
		panic(err)
	}
	defer repo.Close()
}

// Generates Ed25519 key pair and returns the Private key, PeerID in k51 format, and error.
// Correct
func generateKey(keyName string) (*ic.PrivKey, string, error) {
	keyEnc, err := ke.KeyEncoderFromString("k") // Creating obj to encode after I have peer ID
	if err != nil {
		log.Panic(err)
		return nil, nil, err
	}

	fmt.Println("Creating key:pair...")
	sk, pk, err := ic.GenerateKeyPair(ic.Ed25519, 256) //generate default standard key
	if err != nil {
		log.Panic(err)
		return nil, nil, err
	}

	pidpk, err := peer.IDFromPublicKey(pk) //"create" peerID from the public key
	if err != nil {
		log.Panic(err)
		return nil, nil, err
	}

	peerID := keyEnc.FormatID(pidpk) //convert peerID into k51 identifier
	fmt.Printf("PeerID from pk: %s\nk51 FormattedID: %s\n", pidpk.String(), peerID)
	return sk, peerID, nil
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


// Just an Example function to show how to call createEntry()
func makeEntry(sk ic.PrivKey, ipfsPath) (*pb.IpnsEntry, nil) {
	fmt.Println("Creating IPNS record...")
	ipnsRecord, err := createEntry(ipfsPath, sk) //create entry and sign with privatekey
	if err != nil {
	    log.Panic(err)
		return err
    }
	fmt.Printf("IPNS value: %s\n", ipnsRecord.Value)
	return ipnsRecord, err
}

// This function takes a key name and searches for it in local node Keystore.
// returns nil if sucessfull
func exportKey(keyName string) error {
	sh := shell.NewShell(localhost)
	var err error
	rb := sh.Request("key/export", keyName) //export temp key to ds
	err = rb.Exec(context.Background(), err)
	if err != nil {
		return err
	}
	return nil
}

func badRequest(wrt http.ResponseWriter, req *http.Request) {
	keys, ok := req.URL.Query()["keyName"]

	// Query()["keyName"] will return an array of items,
	// we only want the single item.
	keyName := keys[0]
	if !ok || len(keys[0]) < 1 {
        log.Println("Url Param 'keyName' is missing")
		wrt.WriteHeader(http.StatusBadRequest)
		wrt.Header().Set("Content-Type", "application/octet-stream")
		resp := make(map[string]string)
		resp["message"] = "Status Bad Request"
		jsonResp, err:= json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error in JSON marshal. Err: %s", err)
		}
		wrt.Write(jsonResp)
    } 
}

// Grab all subfiles from a Directory
func getAllFilesFromSubmittedDir(file multipart.File, fileHeader *multipart.FileHeader, r *http.Request){
	// The argument to FormFile must match the name attribute
	// of the file input on the frontend
	file, fileHeader, err := r.FormFile("files")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		os.Exit(1)
	}
	defer file.Close()
	
	//get the *fileheaders
 	files := r.MultipartForm.File["files"]

	// Open the zip file
	file, err := files[0].Open()
	if err != nil {
		return "","", fmt.Errorf("Error opening file %g", err)
	}
	defer file.Close()

	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		return "", "", fmt.Errorf("Dir not created %g", err)
	}

	// loop through the files one by one
 	for i, _ := range files {
		fmt.Println("Looping")
 		file, err := files[i].Open()
 		defer file.Close()
 		if err != nil {
 			return "", "", fmt.Errorf(strings.Join([]string{"Error opening file", files[i].Filename}, ""), err)
 		}

 		out, err := os.Create(path + "/" + files[i].Filename)

 		defer out.Close()
 		if err != nil {
 			return "","", fmt.Errorf(strings.Join([]string{"Error creating file path", files[i].Filename}, ""), err)
 		}

 		_, err = io.Copy(out, file) // file not files[i] !

 		if err != nil {
 			return "", "", fmt.Errorf(strings.Join([]string{"Error copying file", files[i].Filename}, ""), err)
 		}
	}
}