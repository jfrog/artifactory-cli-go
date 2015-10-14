package utils

import (
  "os"
  "os/user"
  "strconv"
  "runtime"
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

func GetFileSeperator() string {
    if runtime.GOOS == "windows" {
        return "\\"
    }
    return "/"
}

func GetHomeDir() string {
    user, err := user.Current()
    if err == nil {
        return user.HomeDir
    }
    home := os.Getenv("HOME")
    if home != "" {
        return home
    }
    return "";
}

type Flags struct {
    ArtDetails *ArtifactoryDetails
    DryRun bool
    Props string
    Recursive bool
    Flat bool
    UseRegExp bool
    Threads int
    MinSplitSize int64
    SplitCount int
    Interactive bool
}

type ArtifactoryDetails struct {
    Url string
    User string
    Password string
}