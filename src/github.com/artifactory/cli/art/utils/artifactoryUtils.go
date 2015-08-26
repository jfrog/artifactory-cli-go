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
    headers["X-Checksum-Sha1"] = checksum.sha1
    headers["X-Checksum-Md5"] = checksum.md5

    return PutContent(nil, headers, targetPath, user, password, dryRun)
}

func CalcChecksum(data []byte) *CheckSum {
    checksum := new(CheckSum)
    md5Res := md5.Sum(data)
    sha1Res := sha1.Sum(data)
    checksum.md5 = hex.EncodeToString(md5Res[:])
    checksum.sha1 = hex.EncodeToString(sha1Res[:])

    return checksum
}

func FetchChecksumFromArtifactory(downloadUrl string) {

}

type CheckSum struct {
    md5 string
    sha1 string
}