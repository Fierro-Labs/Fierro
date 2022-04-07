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
	"path"
	"path/filepath"
	"mime/multipart"
	"errors"

	shell "github.com/ipfs/go-ipfs-api"

)

// This function will set the appropiate headers for when there is an error.
func writeJSONError(w http.ResponseWriter, msg string, err error) {
	resp := make(map[string]interface{})
	resp["message"] = msg
	resp["error"] = err
	jsonResp, err:= json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error in JSON marshal. Err: %s", err.Error())
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

// Takes a folder that is in .zip format and saves it to specified dir
// returns the location and name of file
func saveDir(file multipart.File, fileHeader *multipart.FileHeader, dir string, maxUploadSize int64) (string, error) {
	pattern := fmt.Sprintf("%d", time.Now().UnixNano())
	path := strings.Join([]string{dir, "/user", pattern}, "")
	// fmt.Println("path: ", path)

	// Create the uploads folder if it doesn't
	// already exist
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", err
	}

	out, err := os.Create(strings.Join([]string{path,"/", cleanFileName(fileHeader.Filename)}, ""))
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	return path, nil
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

// Unzip directory at source location, iteratively calls unzipFile to unzip sub structures.
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

// Unzips a file to specified dir. 
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

// Removes the file extension of a given file name using slices. 
// returns the altered name as a string
func removeExtenstion(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

// returns the last element of the path. 
func cleanFileName(fileName string) string {
	return path.Base(fileName)
}

// This function will check and restrict the file types submitted
// returns custom error message along with err/nil
func checkFileType(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// Read 512 bytes of the file
	buff := make([]byte, 512)
	_, err := file.Read(buff)
	if err != nil {
		return "Error in reading Dir", err
	}

	// check content/mime type for zip folders
	fileType := http.DetectContentType(buff)
	// fmt.Println("File type:", fileType)
	switch fileType {
		case "application/zip":
			break
		case "application/x-gzip":
			fmt.Println("File is compressed with gzip")
		default:
			fmt.Println("File is not compressed")
			return "The provided file format is not allowed. Please upload a compressed/zip folder", errors.New("Error with DetectContentType")
	}
	
	// Move request body pointer to start of body
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "Error returning request pointer to beginning", err
	}
	return "Success", nil
}