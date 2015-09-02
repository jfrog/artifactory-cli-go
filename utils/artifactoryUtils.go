package utils

import (
    "crypto/md5"
    "crypto/sha1"
    "encoding/hex"
)

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