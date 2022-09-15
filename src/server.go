package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/gammazero/deque"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/robfig/cron"
)

/* TODO: I need to posibly implement a new endpoint within Kubo to support the following endpoint?
// Current progress is to just add support within the Add fx
// With this I should be able to use  --follow and construct a request to my server

// Another problem is the fact that if I don't want to reuse the command format, then I need to think about how to change lines 155-164 in remotepin.go
// to add in a string to the object through Fiero/go-pinning and send it.

// Make sure remotepin.go uses the correct endpoint in Add
*/

const FILE = "Hello"
const MAX_UPLOAD_SIZE = 3072 * 1024 //3kib * 1 kib = 3MiB
const localhost = "localhost:5001"

var abs, _ = filepath.Abs("../")
var q deque.Deque

var MLTRADRS = []string{"/ip4/10.40.2.219/tcp/4001/p2p/12D3KooWJXVZCQzCB28qyDmSGwPLo3Gk2aN9QWctnQkXSc1KCTw2"}
var users = make(map[string]PinResults)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: HomePage")
}

func main() {
	users["testauth"] = PinResults{Count: 0}

	connStr := "postgresql://<username>:<password>@<database_ip>/todos?sslmode=disable"
	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Init a cron job that runs every x minutes
	c := cron.New()
	c.AddFunc("@every 3m", func() {
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
	router.HandleFunc("/pins/", GetRecord).Methods("GET")
	router.HandleFunc("/pins/", StartFollowing).Methods("POST")
	router.HandleFunc("/pins/{requestid}", StopFollowing).Methods("Delete")

	fs := http.FileServer(http.Dir(abs + "/static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	fmt.Printf("Starting server at http://localhost:8082\n")
	log.Fatal(http.ListenAndServe(":8082", router))
}
