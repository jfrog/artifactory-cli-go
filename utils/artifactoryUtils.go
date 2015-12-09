package utils

import (
    "os"
    "io"
    "fmt"
    "sync"
    "strconv"
    "net/http"
    "crypto/md5"
    "crypto/sha1"
    "encoding/hex"
)

func GetFileDetails(filePath string) *FileDetails {
    details := new(FileDetails)
    details.Md5 = calcMd5(filePath)
    details.Sha1 = calcSha1(filePath)

    file, err := os.Open(filePath)
    CheckError(err)
    defer file.Close()

    fileInfo, err := file.Stat()
    CheckError(err)
    details.Size = fileInfo.Size()

    return details
}

func calcSha1(filePath string) string {
    file, err := os.Open(filePath)
    CheckError(err)
    defer file.Close()

    var resSha1 []byte
    hashSha1 := sha1.New()
    _, err = io.Copy(hashSha1, file)
    CheckError(err)
    return hex.EncodeToString(hashSha1.Sum(resSha1))
}

func calcMd5(filePath string) string {
    file, err := os.Open(filePath)
    CheckError(err)
    defer file.Close()

    var resMd5 []byte
    hashMd5 := md5.New()
    _, err = io.Copy(hashMd5, file)
    CheckError(err)
    return hex.EncodeToString(hashMd5.Sum(resMd5))
}

func GetFileDetailsFromArtifactory(downloadUrl string, user string, password string) *FileDetails {
    resp, _ := SendHead(downloadUrl, user, password)
    fileSize, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
    CheckError(err)

    fileDetails := new(FileDetails)

    fileDetails.Md5 = resp.Header.Get("X-Checksum-Md5")
    fileDetails.Sha1 = resp.Header.Get("X-Checksum-Sha1")
    fileDetails.Size = fileSize
    fileDetails.AcceptRanges = resp.Header.Get("Accept-Ranges") == "bytes"
    return fileDetails
}

func GetEncryptedPasswordFromArtifactory(artifactoryDetails *ArtifactoryDetails) (*http.Response, string) {
	apiUrl := artifactoryDetails.Url + "api/security/encryptedPassword"
	resp, body := SendGet(apiUrl, nil, artifactoryDetails.User, artifactoryDetails.Password)
	return resp, string(body)
}

func DownloadFileConcurrently(downloadPath, localPath, fileName, logMsgPrefix string, fileSize int64, flags *Flags) {
    tempLoclPath := GetTempDirPath() + "/" + localPath

    var wg sync.WaitGroup
    chunkSize := fileSize / int64(flags.SplitCount)
    mod := fileSize % int64(flags.SplitCount)

    for i := 0; i < flags.SplitCount ; i++ {
        wg.Add(1)
        start := chunkSize * int64(i)
        end := chunkSize * (int64(i) + 1)
        if i == flags.SplitCount-1 {
            end += mod
        }
        go func(start, end int64, i int) {
            headers := make(map[string]string)
            headers["Range"] = "bytes=" + strconv.FormatInt(start, 10) +"-" + strconv.FormatInt(end-1, 10)
            resp, body := SendGet(downloadPath, headers, flags.ArtDetails.User, flags.ArtDetails.Password)

            fmt.Println(logMsgPrefix + " [" + strconv.Itoa(i) + "]:", resp.Status + "...")

            os.MkdirAll(tempLoclPath ,0777)
            filePath := tempLoclPath + "/" + fileName + "_" + strconv.Itoa(i)

            createFileWithContent(filePath, body)
            wg.Done()
        }(start, end, i)
    }
    wg.Wait()

    if !flags.Flat && localPath != "" {
        os.MkdirAll(localPath ,0777)
        fileName = localPath + "/" + fileName
    }

    if IsPathExists(fileName) {
        err := os.Remove(fileName)
        CheckError(err)
    }

    destFile, err := os.Create(fileName)
    CheckError(err)
    defer destFile.Close()
    for i := 0; i < flags.SplitCount; i++ {
        tempFilePath := GetTempDirPath() + "/" + fileName + "_" + strconv.Itoa(i)
        AppendFile(tempFilePath, destFile)
    }
    fmt.Println(logMsgPrefix + " Done downloading.")
}

func createFileWithContent(filePath string, content []byte) {
    out, err := os.Create(filePath)
    CheckError(err)
    defer out.Close()
    out.Write(content)
}

type FileDetails struct {
    Md5 string
    Sha1 string
    Size int64
    AcceptRanges bool
}