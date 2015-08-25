package main

import (
    "os"
    "io"
	"bytes"
	"encoding/hex"
	"net/http"
	"io/ioutil"
    "path/filepath"
    "crypto/md5"
    "crypto/sha1"
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

func ReadFile(filePath string) []byte {
	content, err := ioutil.ReadFile(filePath)
	CheckError(err)
	return content
}

func tryChecksumDeploy(fileContent []byte, targetPath string) *http.Response {
    checksum := CalcChecksum(fileContent)

    headers := make(map[string]string)
    headers["X-Checksum-Deploy"] = "true"
    headers["X-Checksum-Sha1"] = checksum.sha1
    headers["X-Checksum-Md5"] = checksum.md5

    return PutContent(nil, headers, targetPath, User, Password, DryRun)
}

// Sends an HTTP PUT request to the specified URL, sending the specified content.
func PutContent(content []byte, headers map[string]string, url string, user string, password string, dryRun bool) *http.Response {
	if dryRun {
        return nil
	}
	var data *bytes.Buffer = bytes.NewBufferString("")
	if content != nil {
	    data = bytes.NewBuffer(content)
	}
    req, err := http.NewRequest("PUT", url, data)
    if user != "" && password != "" {
	    req.SetBasicAuth(user, password)
    }
    for name := range headers {
println(name + " " + headers[name])
        req.Header.Set(name, headers[name])
    }
    client := &http.Client{}
    resp, err := client.Do(req)
    CheckError(err)
    defer resp.Body.Close()

    return resp
}

// Sends an HTTP PUT request to the specified URL, sending the file in the
// specified path.
func PutFile(filePath string, url string, user string, password string, dryRun bool) {
    println("Uploading " + filePath + " to " + url)
	content, err := ioutil.ReadFile(filePath)
	CheckError(err)
    PutContent(content, nil, url, user, password, dryRun)
}

func DownloadFile(downloadPath string, localPath string, fileName string, flat bool) {
    println("Downloading " + downloadPath)

    if !flat && localPath != "" {
        os.MkdirAll(localPath ,0777)
        fileName = localPath + "/" + fileName
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

func CalcChecksum(data []byte) *CheckSum {
    checksum := new(CheckSum)
    md5Res := md5.Sum(data)
    sha1Res := sha1.Sum(data)
    checksum.md5 = hex.EncodeToString(md5Res[:])
    checksum.sha1 = hex.EncodeToString(sha1Res[:])

    return checksum
}

type CheckSum struct {
    md5 string
    sha1 string
}