package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"os/signal"
	"syscall"
	"net/http"

    "github.com/gorilla/mux"
	"github.com/gammazero/deque"
	"github.com/robfig/cron"
)

const FILE = "Hello"
const MAX_UPLOAD_SIZE = 3072 * 1024 //3kib * 1 kib = 3MiB
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

	FileName, err := saveFile(r, dir, MAX_UPLOAD_SIZE) // grab uploaded content & save to disk
	if err != nil {
		writeJSONError(w, "Error in saveContent", err)
		return
	}

	fmt.Println("Adding to IPFS...")
	ipfsPath, err := addToIPFS(dir+"/"+FileName, "") // add content file to IPFS
	if err != nil {
		writeJSONError(w, "Error in addToIPFS", err)
		return
	}
	fmt.Println("ipfs Path: ", ipfsPath)
	writeJSONSuccess(w, "Success - addFile", ipfsPath)
}

// Grab uploaded file and add to ipfs
// returns ipfsPath (ipfsURI + CID)
func addFolder(w http.ResponseWriter, r *http.Request) {
	const dir = "Uploads"

	// Check uploaded files/Dir is not bigger than MAX_UPLOAD_SIZE
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		writeJSONError(w, "File too big. Max size is 3MB", err)
		return
	}

	// The argument to FormFile must match the name attribute
	// of the file input on the frontend
	file, fileHeader, err := r.FormFile("files")
	if err != nil {
		writeJSONError(w, "Error accessing FormFile", err)
		return
	}
	defer file.Close()

	if msg, err := checkFileType(file, fileHeader); err != nil {
		writeJSONError(w, msg, err)
		return
	}

	dirPath, err := saveDir(file, fileHeader, dir, MAX_UPLOAD_SIZE) // grab uploaded content & save to disk
	if err != nil {
		writeJSONError(w, "Error in saveDir", err)
		return
	}

	// fmt.Println("Dir created: ", dirPath)
	cleanedFileName := cleanFileName(fileHeader.Filename)
	fileName := removeExtenstion(cleanedFileName)
	err = unzipSource(dirPath+"/"+cleanedFileName, dirPath+"/unzip/")
	if err != nil {
		writeJSONError(w, "Error in unzipSource", err)
		return
	}

	// fmt.Println("fileName: ", fileName)
	fmt.Println("Adding to IPFS...")
	ipfsPath, err := addToIPFS(strings.Join([]string{dirPath, "/unzip/", fileName}, ""), "r") // add content file to IPFS
	if err != nil {
		writeJSONError(w, "Error in addToIPFS", err)
		return
	}

	 // Remove all the directories and files
    // Using RemoveAll() function
    err = os.RemoveAll(dirPath)
    if err != nil {
		writeJSONError(w, "Error in deleting directory", err)
        return
    }

	fmt.Println("ipfs Path: ", ipfsPath)
	writeJSONSuccess(w, "Success - addFolder", ipfsPath)
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
	router.HandleFunc("/addFolder", addFolder).Methods("POST")

	fs := http.FileServer(http.Dir("./static/"))
    router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	fmt.Printf("Starting server at http://localhost:8082\n")
	log.Fatal(http.ListenAndServe(":8082", router))
}