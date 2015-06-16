package main

import (
    "os"
    "io"
    "strings"
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

// Sends an HTTP PUT request to the specified URL, sending the file in the
// specified path.
func PutFile(filePath string, url string, user string, password string, dryRun bool) {
    println("Uploading " + filePath + " to " + url)
	fileContent, err := ioutil.ReadFile(filePath)
	CheckError(err)
	if dryRun {
        return
	}

    req, err := http.NewRequest("PUT", url, bytes.NewBuffer(fileContent))
    if user != "" && password != "" {
	    req.SetBasicAuth(user, password)
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    CheckError(err)
    defer resp.Body.Close()

    println("Artifactory response:", resp.Status)
}

func DownloadFile(downloadPath string, localPath string, flat bool) {
    println("Downloading " + downloadPath)
    index := strings.LastIndex(localPath, "/")
    length := len(localPath)
    var fileName string
    var dir string

    if index != -1 && index != length-1 {
        dir = localPath[: index]
        fileName = localPath[index+1 : length-1]
    } else {
        fileName = localPath
    }
    if !flat && dir != "" {
        os.MkdirAll(dir ,0777)
        fileName = dir + "/" + fileName
    }

    out, err := os.Create(fileName)
    CheckError(err)
    defer out.Close()
    resp, err := http.Get(downloadPath)
    CheckError(err)
    defer resp.Body.Close()
    _, err = io.Copy(out, resp.Body)
    CheckError(err)

    println("Artifactory response:", resp.Status)
}

func SendPost(url string, data string, user string, password string) []byte {
    return Send("POST", url, data, user, password)
}

func SendGet(url string, user string, password string) []byte {
    return Send("GET", url, "", user, password)
}

func Send(method string, url string, data string, user string, password string) []byte {
    var req *http.Request
    var err error
    if data != "" {
        req, err = http.NewRequest(method, url, bytes.NewBuffer([]byte(data)))
    } else {
        req, err = http.NewRequest(method, url, nil)
    }
    CheckError(err)

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