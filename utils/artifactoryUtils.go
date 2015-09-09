package utils

import (
    "os"
    "io"
    "strconv"
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

type FileDetails struct {
    Md5 string
    Sha1 string
    Size int64
    AcceptRanges bool
}