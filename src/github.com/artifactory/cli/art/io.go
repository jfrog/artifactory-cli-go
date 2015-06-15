package main

import (
    "os"
	"bytes"
	"net/http"
	"io/ioutil"
    "path/filepath"
)

func IsDir(path string) bool {
    if !IsPathExists(path) {
        return false
    }
    f, err := os.Stat(path)
    CheckError(err)
    return f.IsDir()
}

func IsPathExists(path string) bool {
    _, err := os.Stat(path)
    return !os.IsNotExist(err)
}

// Sends an HTTP PUT request to specified URL, sending the file in the
// specified path.
func PutFile(filePath string, url string, user string, password string) {
	fileContent, err := ioutil.ReadFile(filePath)
	CheckError(err)

    req, err := http.NewRequest("PUT", url, bytes.NewBuffer(fileContent))
    if user != "" && password != "" {
	    req.SetBasicAuth(user, password)
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    CheckError(err)
    defer resp.Body.Close()

    println("Response status:", resp.Status)
}

func SendPost(url string, data string, user string, password string) []byte {
    req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
    if user != "" && password != "" {
	    req.SetBasicAuth(user, password)
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    CheckError(err)
    defer resp.Body.Close()

    println("Response status:", resp.Status)
    body, _ := ioutil.ReadAll(resp.Body)
    return body
}

// Return the list of all files and directories (recursive) in the specified path
func ListFiles(path string) []string {
    fileList := []string{}
    err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
        fileList = append(fileList, path)
        return nil
    })
    CheckError(err)
    return fileList
}