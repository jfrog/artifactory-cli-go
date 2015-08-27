package utils

import (
	"net/http"
    "crypto/md5"
    "crypto/sha1"
    "encoding/hex"
)

func TryChecksumDeploy(fileContent []byte, targetPath string, user string, password string, dryRun bool) *http.Response {
    checksum := CalcChecksum(fileContent)

    headers := make(map[string]string)
    headers["X-Checksum-Deploy"] = "true"
    headers["X-Checksum-Sha1"] = checksum.Sha1
    headers["X-Checksum-Md5"] = checksum.Md5

    return PutContent(nil, headers, targetPath, user, password, dryRun)
}

func ShouldDownloadFile(localFilePath string, downloadPath string, user string, password string, dryRun bool) bool {
    if !IsFileExists(localFilePath) {
        return true
    }
    localChecksum := CalcChecksum(ReadFile(localFilePath))
    artifactoryChecksum := FetchChecksumFromArtifactory(downloadPath, user, password)
    if localChecksum.Md5 != artifactoryChecksum.Md5 || localChecksum.Sha1 != artifactoryChecksum.Sha1 {
       return true
    }
    return false
}

func CalcChecksum(data []byte) *CheckSum {
    checksum := new(CheckSum)
    md5Res := md5.Sum(data)
    sha1Res := sha1.Sum(data)
    checksum.Md5 = hex.EncodeToString(md5Res[:])
    checksum.Sha1 = hex.EncodeToString(sha1Res[:])

    return checksum
}

func FetchChecksumFromArtifactory(downloadUrl string, user string, password string) *CheckSum {
    resp, _ := SendHead(downloadUrl, user, password)
    checksum := new(CheckSum)
    checksum.Md5 = resp.Header.Get("X-Checksum-Md5")
    checksum.Sha1 = resp.Header.Get("X-Checksum-Sha1")
    return checksum
}

type CheckSum struct {
    Md5 string
    Sha1 string
}