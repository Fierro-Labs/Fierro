package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

type PublishResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ResolvedPath struct {
	Path string
}

// Grab uploaded file and add to ipfs
// returns ipfsPath (ipfsURI + CID)
func AddFile(w http.ResponseWriter, r *http.Request) {
	var dir = abs + "/Uploads"

	FilePath, err := saveFile(r, dir, MAX_UPLOAD_SIZE) // grab uploaded content & save to disk
	if err != nil {
		writeJSONError(w, "Error in saveContent", err)
		return
	}

	ipfsPath, err := addToIPFS(FilePath, "") // add content file to IPFS
	if err != nil {
		writeJSONError(w, "Error in addToIPFS", err)
		return
	}
	writeJSONSuccess(w, "Success - addFile", ipfsPath)
}

// Grab uploaded file and add to ipfs
// returns ipfsPath (ipfsURI + CID)
func addFolder(w http.ResponseWriter, r *http.Request) {
	var dir = abs + "/Uploads"

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

	cleanedFileName := cleanFileName(fileHeader.Filename)
	fileName := removeExtenstion(cleanedFileName)
	err = unzipSource(dirPath+"/"+cleanedFileName, dirPath+"/unzip/")
	if err != nil {
		writeJSONError(w, "Error in unzipSource", err)
		return
	}

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

	writeJSONSuccess(w, "Success - addFolder", ipfsPath)
}

// Adds path to IPFS, given as a string
// returns CID
func addToIPFS(path string, option string) (string, error) {
	sh := shell.NewShell(localhost)
	var ipfsPath string

	if option == "r" {
		cid, err := sh.AddDir(path)
		if err != nil {
			fmt.Println("Error in adding content to ipfs: ", err)
			return "", err
		}
		ipfsPath = ipfsURI + cid

	} else {
		fileReader, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s ", err)
			return "", err
		}
		defer fileReader.Close()

		cid, err := sh.Add(fileReader)
		if err != nil {
			fmt.Println("Error in adding content to ipfs: ", err)
			return "", err
		}
		ipfsPath = ipfsURI + cid
	}

	return ipfsPath, nil
}

// Custom publishing function returns the response object and error
// Kept as close as possible to Publish method found at gh.com/go-ipfs-api
func PublishToIPFS(contentHash string, key string) (*PublishResponse, error) {
	var pubResp PublishResponse
	sh := shell.NewShell(localhost)
	req := sh.Request("name/publish", contentHash).Option("key", key)
	req.Option("resolve", true)
	err := req.Exec(context.Background(), &pubResp)
	if err != nil {
		return nil, err
	}

	return &pubResp, nil
}

// This function is needed to let the world know your Record exists.
func publishToIPNS(ipfsPath string, KeyName string) (*PublishResponse, error) {
	pubResp, err := PublishToIPFS(ipfsPath, KeyName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in Publish: %s ", err)
		return nil, err
	}

	return pubResp, nil
}

// This function will resolve/download the content pointed to by the record.
func resolve(ipnsKey string) (string, error) {
	var path ResolvedPath
	sh := shell.NewShell(localhost)
	req := sh.Request("name/resolve", ipnsKey).Option("dht-timeout", "180s") // timeout after 3 minutes
	err := req.Exec(context.Background(), &path)
	if err != nil {
		return "", err
	}
	return path.Path, nil
}
