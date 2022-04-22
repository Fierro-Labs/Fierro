package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/gammazero/deque"
	"github.com/gorilla/mux"
	"github.com/robfig/cron"
)

const FILE = "Hello"
const MAX_UPLOAD_SIZE = 3072 * 1024 //3kib * 1 kib = 3MiB
const localhost = "localhost:5001"

const ipfsURI = "/ipfs/"
const ipnsURI = "/ipns/"

var abs, _ = filepath.Abs("../")
var q deque.Deque

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: HomePage")
}

func main() {
	// Init a cron job that runs every x minutes
	c := cron.New()
	c.AddFunc("@every 1m", func() {
		ipfsPath := follow()
		fmt.Printf("%s\n", ipfsPath)
	})
	c.Start()

	// Catch ctrl+c signal to kill threads
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-channel
		c.Stop()
		os.Exit(1)
	}()

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

	fs := http.FileServer(http.Dir(abs + "/static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	fmt.Printf("Starting server at http://localhost:8082\n")
	log.Fatal(http.ListenAndServe(":8082", router))
}
