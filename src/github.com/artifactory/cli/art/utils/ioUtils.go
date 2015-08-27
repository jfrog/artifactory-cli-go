package utils

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

func IsFileExists(path string) bool {
    if !IsPathExists(path) {
        return false
    }
    f, err := os.Stat(path)
    CheckError(err)
    return !f.IsDir()
}

func ReadFile(filePath string) []byte {
	content, err := ioutil.ReadFile(filePath)
	CheckError(err)
	return content
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
	content, err := ioutil.ReadFile(filePath)
	CheckError(err)
    PutContent(content, nil, url, user, password, dryRun)
}

func DownloadFile(downloadPath string, localPath string, fileName string, flat bool, user string, password string) *http.Response {
    if !flat && localPath != "" {
        os.MkdirAll(localPath ,0777)
        fileName = localPath + "/" + fileName
    }

    out, err := os.Create(fileName)
    CheckError(err)
    defer out.Close()
    resp, body := SendGet(downloadPath, user, password)
    out.Write(body)
    CheckError(err)

    return resp
}

func SendPost(url string, data string, user string, password string) []byte {
    _, body := Send("POST", url, data, user, password)
    return body
}

func SendGet(url string, user string, password string) (*http.Response, []byte) {
    return Send("GET", url, "", user, password)
}

func SendHead(url string, user string, password string) (*http.Response, []byte) {
    return Send("HEAD", url, "", user, password)
}

func Send(method string, url string, data string, user string, password string) (*http.Response, []byte) {
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
    body, _ := ioutil.ReadAll(resp.Body)
    return resp, body
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