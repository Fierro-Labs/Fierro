package main

import (
	"fmt"
	"log"
	"archive/zip"
	"net/http"
	"os"
	"time"
	"context"
	"encoding/json"
	"io"
	"strings"
	"math/rand"
	"strconv"
	"path/filepath"


	shell "github.com/ipfs/go-ipfs-api"

)

// This function will set the appropiate headers for when there is an error.
func writeJSONError(w http.ResponseWriter, msg string, err error) {
	resp := make(map[string]interface{})
	resp["message"] = msg
	resp["error"] = err
	jsonResp, err:= json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error in JSON marshal. Err: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

// This function will set the appropiate headers for when there is an error.
func writeJSONSuccess(w http.ResponseWriter, msg string, val string) {
	resp := make(map[string]string)
	resp["message"] = msg
	resp["value"] = val
	jsonResp, err:= json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error in JSON marshal. Err: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

// This function grabs the specified parameter value out the URL
func GetParam(r *http.Request, parameter string) (string, bool) {
	params, ok := r.URL.Query()[parameter]
	// Query()[parameter] will return an array of items,
	// we only want the single item.
	if !ok || len(params[0]) < 1 {
		fmt.Println("Error: Missing " + parameter)
		return "Missing "+parameter+" parameter", !ok
	} 
	return params[0], ok
}

// Grabs file with specified file size and save to specified dir
// returns name of file
func saveFile(r *http.Request, dir string, size int64) (string, error) {
	var FileName string

	fmt.Println("Saving file...")
	r.ParseMultipartForm(size) // limit max input length
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s ", err)
		return "", err
	}
	defer file.Close()
	FileName = header.Filename
	
	f, err := os.OpenFile(dir+"/"+FileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Error in OpenFile: ", err)
		return "", err
	}
	io.Copy(f, file)

	return FileName, nil
}

// Grabs a set of files and saves them to specified users dir
// returns the directory name
func saveDir(r *http.Request, dir string, size int64) (string, string, error) {
	pattern := strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(10000))
	path := strings.Join([]string{dir, "/user", pattern}, "")

	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		return "Dir not created", "", err
	}
	err = r.ParseMultipartForm(size) // grab the multipart form
 	if err != nil {
 		return "Problem parsing Form", "", err
 	}

 	formdata := r.MultipartForm // ok, no problem so far, read the Form data

 	//get the *fileheaders
 	files := formdata.File["files"] // grab the filenames

 	for i, _ := range files { // loop through the files one by one
 		file, err := files[i].Open()
 		defer file.Close()
 		if err != nil {
 			return "", "", err
 		}

 		out, err := os.Create(path + "/" + files[i].Filename)

 		defer out.Close()
 		if err != nil {
 			return "No such file or directory", files[i].Filename, err
 		}

 		_, err = io.Copy(out, file) // file not files[i] !

 		if err != nil {
 			return "Error with io.copy", files[i].Filename, err
 		}
	}
	return path, files[0].Filename, nil
}

// Generate keys and embed records. Meant to test how keys are needed to be passed.
func testFunctions(ipfsPath string) {
	const KeyName = "temp"
	sh := shell.NewShell(localhost)

	key, err := sh.KeyGen(context.Background(), KeyName) //generate temp key to local node
	if err != nil {
		panic(err)
	}

	// publishToIPNS(ipfsPath, key.Name) // publish ipfsPath to ipfs under temp key
	fmt.Println("Exporting key...")
	err = exportKey(KeyName)
	if err != nil {
		panic(err)
	}

	fmt.Println("Deleting key...")
	err = deleteKey(key.Name) // delete temp key
	if err != nil {
		panic(err)
	}

	err = importKey(key.Name, key.Name + ".key")
	if err != nil {
		panic(err)
	}
	fmt.Printf("functions individually work.\n\n")
}

func unzipSource(source, destination string) error {
    // 1. Open the zip file
    reader, err := zip.OpenReader(source)
    if err != nil {
        return err
    }
    defer reader.Close()

    // 2. Get the absolute destination path
    destination, err = filepath.Abs(destination)
    if err != nil {
        return err
    }

	fmt.Println("Absolute path: ", destination)
    // 3. Iterate over zip files inside the archive and unzip each of them
    for _, f := range reader.File {
        err := unzipFile(f, destination)
        if err != nil {
            return err
        }
    }

    return nil
}

func unzipFile(f *zip.File, destination string) error {
    // 4. Check if file paths are not vulnerable to Zip Slip
    filePath := filepath.Join(destination, f.Name)
    if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
        return fmt.Errorf("invalid file path: %s", filePath)
    }

    // 5. Create directory tree
    if f.FileInfo().IsDir() {
        if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
            return err
        }
        return nil
    }

    if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
        return err
    }

    // 6. Create a destination file for unzipped content
    destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
    if err != nil {
        return err
    }
    defer destinationFile.Close()

    // 7. Unzip the content of a file and copy it to the destination file
    zippedFile, err := f.Open()
    if err != nil {
        return err
    }
    defer zippedFile.Close()

    if _, err := io.Copy(destinationFile, zippedFile); err != nil {
        return err
    }
    return nil
}

func removeExtenstion(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}






// Leaving to go eat. I am trying to accept a zip folder (works) and unzip it to a different directory. So I made an unzip dir to unzip the stuff too and then I will add to IPFS from that dir.