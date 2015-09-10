package utils

import (
    "io"
    "os"
	"bytes"
	"strings"
	"strconv"
	"sync"
	"bufio"
	"net/http"
	"io/ioutil"
    "path/filepath"
)

var tempDirPath string

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

func IsDirExists(path string) bool {
    if !IsPathExists(path) {
        return false
    }
    f, err := os.Stat(path)
    CheckError(err)
    return f.IsDir()
}

func ReadFile(filePath string) []byte {
	content, err := ioutil.ReadFile(filePath)
	CheckError(err)
	return content
}

func UploadFile(f *os.File, url, user, password string) *http.Response {
    req, err := http.NewRequest("PUT", url, f)
    CheckError(err)
    if user != "" && password != "" {
	    req.SetBasicAuth(user, password)
    }
    client := &http.Client{}
    res, err := client.Do(req)
    CheckError(err)
    return res
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

func DownloadFileConcurrently(downloadPath string, localPath string, fileName string,
    flat bool, user string, password string, fileSize int64, splitCount int) {

    tempLoclPath := GetTempDirPath() + "/" + localPath

    var wg sync.WaitGroup
    chunkSize := fileSize / int64(splitCount)
    mod := fileSize % int64(splitCount)

    for i := 0; i < splitCount ; i++ {
        wg.Add(1)
        start := chunkSize * int64(i)
        end := chunkSize * (int64(i) + 1)
        if i == splitCount-1 {
            end += mod
        }
        go func(start, end int64, i int) {
            headers := make(map[string]string)
            headers["Range"] = "bytes=" + strconv.FormatInt(start, 10) +"-" + strconv.FormatInt(end-1, 10)
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

    if IsPathExists(fileName) {
        err := os.Remove(fileName)
        CheckError(err)
    }
    for i := 0; i < splitCount; i++ {
        tempFilePath := GetTempDirPath() + "/" + fileName + "_" + strconv.Itoa(i)
        AppendFile(tempFilePath, fileName)
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

func GetTempDirPath() string {
    if tempDirPath == "" {
        path, err := ioutil.TempDir("", "artifactory.cli.")
        CheckError(err)
        tempDirPath = path
    }
    return tempDirPath
}

func RemoveTempDir() {
    if IsDirExists(tempDirPath) {
        os.RemoveAll(tempDirPath)
    }
}

// Reads the content of the file in the source path and appends it to
// the file in the destination path.
func AppendFile(srcPath, destPath string) {
    srcFile, err := os.Open(srcPath)
    CheckError(err)

    defer func() {
        err := srcFile.Close();
        CheckError(err)
    }()

    reader := bufio.NewReader(srcFile)

    var destFile *os.File
    if IsPathExists(destPath) {
        destFile, err = os.OpenFile(destPath, os.O_APPEND, 0666)
    } else {
        destFile, err = os.Create(destPath)
    }
    CheckError(err)
    defer func() {
        err := destFile.Close();
        CheckError(err)
    }()

    writer := bufio.NewWriter(destFile)
    buf := make([]byte, 1024)
    for {
        n, err := reader.Read(buf)
        if err != io.EOF {
            CheckError(err)
        }
        if n == 0 {
            break
        }
        _, err = writer.Write(buf[:n])
        CheckError(err)
    }
    err = writer.Flush()
    CheckError(err)
}