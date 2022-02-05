import (
	"io/ioutil"

	ic "github.com/libp2p/go-libp2p-core/crypto"
	keystore "github.com/ipfs/go-ipfs-keystore"
	fsrepo "github.com/ipfs/go-ipfs/repo/fsrepo"
	pb "github.com/ipfs/go-ipns/pb"
)

type IpnsEntry struct {
	Name  string
	Value string
}

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
// returns nil if sucessfull & stores key as file in current dir.
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