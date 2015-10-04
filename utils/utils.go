package utils

import (
  "os"
  "strconv"
)

func CheckError(err error) {
    if err != nil {
        panic(err)
    }
}

func Exit(msg string) {
    println(msg)
    os.Exit(1)
}

func GetLogMsgPrefix(threadId int, dryRun bool) string {
    var strDryRun string
    if dryRun {
        strDryRun = " [Dry run]"
    } else {
        strDryRun = ""
    }
    return "[Thread " + strconv.Itoa(threadId) + "]" + strDryRun
}

type Flags struct {
    ArtDetails ArtifactpryDetails
    DryRun bool
    Props string
    Recursive bool
    Flat bool
    UseRegExp bool
    Threads int
    MinSplitSize int64
    SplitCount int
}

type ArtifactpryDetails struct {
    Url string
    User string
    Password string
}