package utils

import (
    "os"
	"bytes"
	"strings"
	"strconv"
	"sync"
	"net/http"
	"io/ioutil"
    "path/filepath"
)

func GetFileNameFromPath(path string) string {
    index := strings.LastIndex(path, "/")
    if index != -1 {
        return path[index+1:]
    }
    index = strings.LastIndex(path, "\\")
    if index != -1 {
        return path[index+1:]
    }
    return path
}

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

// Sends an HTTP PUT request to the specified URL, sending the file in the
// specified path.
func PutFile(filePath string, url string, user string, password string, dryRun bool) {
	content, err := ioutil.ReadFile(filePath)
	CheckError(err)
	if !dryRun {
        SendPut(url, content, nil, user, password)
	}
}

func DownloadFile(downloadPath string, localPath string, fileName string, flat bool, user string, password string) *http.Response {
    if !flat && localPath != "" {
        os.MkdirAll(localPath ,0777)
        fileName = localPath + "/" + fileName
    }

    out, err := os.Create(fileName)
    CheckError(err)
    defer out.Close()
    resp, body := SendGet(downloadPath, nil, user, password)
    out.Write(body)
    CheckError(err)
    return resp
    return nil
}

func DownloadFileConcurrently(downloadPath string, localPath string, fileName string, flat bool, user string, password string, fileSize int) {
    Threads := 3
    TempDirName := "C:\\temp\\art_cli_temp"
    tempLoclPath := TempDirName + "/" + localPath

    var wg sync.WaitGroup
    chunkSize := fileSize / Threads
    mod := fileSize % Threads

    for i := 0; i < Threads ; i++ {
        wg.Add(1)
        start := chunkSize * i
        end := chunkSize * (i + 1)
        if i == Threads-1 {
            end += mod
        }
        go func(start int, end int, i int) {
            headers := make(map[string]string)
            headers["Range"] = "bytes=" + strconv.Itoa(start) +"-" + strconv.Itoa(end-1)
            resp, body := SendGet(downloadPath, headers, user, password)

            print("[" + strconv.Itoa(i) + "]:", resp.Status + "...")

            os.MkdirAll(tempLoclPath ,0777)
            filePath := tempLoclPath + "/" + fileName + "_" + strconv.Itoa(i)

            out, err := os.Create(filePath)
            CheckError(err)
            defer out.Close()

            out.Write(body)
            CheckError(err)
            wg.Done()
        }(start, end, i)
    }
    wg.Wait()

    if !flat && localPath != "" {
        os.MkdirAll(localPath ,0777)
        fileName = localPath + "/" + fileName
    }

    out, err := os.Create(fileName)
    CheckError(err)
    defer out.Close()

    for i := 0; i < Threads ; i++ {
        tempFilePath := TempDirName + "/" + fileName + "_" + strconv.Itoa(i)
        content := ReadFile(tempFilePath)
        out.Write(content)
        CheckError(err)
    }
    println("Done downloading.")
}

func SendPut(url string, content []byte, headers map[string]string, user string, password string) (*http.Response, []byte) {
    return Send("PUT", url, content, headers, user, password)
}

func SendPost(url string, content []byte, user string, password string) []byte {
    _, body := Send("POST", url, content, nil, user, password)
    return body
}

func SendGet(url string, headers map[string]string, user string, password string) (*http.Response, []byte) {
    return Send("GET", url, nil, headers, user, password)
}

func SendHead(url string, user string, password string) (*http.Response, []byte) {
    return Send("HEAD", url, nil, nil, user, password)
}

func Send(method string, url string, content []byte, headers map[string]string, user string, password string) (*http.Response, []byte) {
    var req *http.Request
    var err error

    if content != nil {
        req, err = http.NewRequest(method, url, bytes.NewBuffer(content))
    } else {
        req, err = http.NewRequest(method, url, nil)
    }
    CheckError(err)

    if user != "" && password != "" {
	    req.SetBasicAuth(user, password)
    }
    if headers != nil {
        for name := range headers {
            req.Header.Set(name, headers[name])
        }
    }
    client := &http.Client{}
    resp, err := client.Do(req)
    CheckError(err)
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    return resp, body
}

// Return the recursive list of files and directories in the specified path
func ListFilesRecursive(path string) []string {
    fileList := []string{}
    err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
        fileList = append(fileList, path)
        return nil
    })
    CheckError(err)
    return fileList
}

// Return the list of files and directories in the specified path
func ListFiles(path string) []string {
    if !strings.HasSuffix(path, "/") {
        path += "/"
    }
    fileList := []string{}
    files, _ := ioutil.ReadDir("./")
    for _, f := range files {
        fileList = append(fileList, path + f.Name())
    }
    return fileList
}