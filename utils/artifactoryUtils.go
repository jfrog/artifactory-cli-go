package utils

import (
    "strconv"
    "crypto/md5"
    "crypto/sha1"
    "encoding/hex"
)

func GetFileDetails(data []byte) *FileDetails {
    details := new(FileDetails)
    md5Res := md5.Sum(data)
    sha1Res := sha1.Sum(data)
    details.Md5 = hex.EncodeToString(md5Res[:])
    details.Sha1 = hex.EncodeToString(sha1Res[:])
    details.Size = len(data)

    return details
}

func GetFileDetailsFromArtifactory(downloadUrl string, user string, password string) *FileDetails {
    resp, _ := SendHead(downloadUrl, user, password)
    fileSize, err := strconv.Atoi(resp.Header.Get("Content-Length"))
    CheckError(err)

    fileDetails := new(FileDetails)
    fileDetails.Md5 = resp.Header.Get("X-Checksum-Md5")
    fileDetails.Sha1 = resp.Header.Get("X-Checksum-Sha1")
    fileDetails.Size = fileSize
    return fileDetails
}

type FileDetails struct {
    Md5 string
    Sha1 string
    Size int
}