package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"net/http"

    "github.com/gorilla/mux"
	"github.com/gammazero/deque"
	"github.com/robfig/cron"
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
	fmt.Println("ipfs Path: ", ipfsPath)
	writeJSONSuccess(w, "Success - addFile", ipfsPath)
}

func main() {
	// Used to test if keys need to be passed as objects or ints?
	// testFunctions(ipfsPath)
	c := cron.New()
	c.AddFunc("@every 2m", func() {
		ipfsPath := follow()
		fmt.Printf("%s\n",ipfsPath)	
	})
	c.Start()
	
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func(){
		<-channel
		c.Stop()
		os.Exit(1)
	}()

	q.PushBack("k51qzi5uqu5dm876hw4kh2mn58rnajofhoohohymt9bui38q6ogsa0rrct6fnh")
	// handles api/website routes.
	router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", index)
	router.HandleFunc("/getKey", GetKey).Methods("GET")
	router.HandleFunc("/postKey", PostKey).Methods("POST")
	router.HandleFunc("/deleteKey", DeleteKey).Methods("DELETE")
	router.HandleFunc("/postRecord", PostRecord).Methods("POST")
	router.HandleFunc("/putRecord", PutRecord).Methods("PUT")
	router.HandleFunc("/getRecord", GetRecord).Methods("GET")
	router.HandleFunc("/startFollowing", StartFollowing).Methods("POST")
	router.HandleFunc("/stopFollowing", StopFollowing).Methods("Delete")
	router.HandleFunc("/addFile", AddFile).Methods("POST")


	fmt.Printf("Starting server at port 8082\n")
	log.Fatal(http.ListenAndServe(":8082", router))
}