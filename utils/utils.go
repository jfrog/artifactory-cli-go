package utils

import (
  "os"
  "fmt"
  "strings"
  "os/user"
  "strconv"
  "runtime"
)

var ExitCodeError ExitCode = ExitCode{1}
var ExitCodeWarning ExitCode = ExitCode{2}

type ExitCode struct {
    Code int
}

func GetVersion() string {
    return "1.2.1"
}

func CheckError(err error) {
    if err != nil {
        panic(err)
    }
}

func Exit(exitCode ExitCode, msg string) {
    if msg != "" {
        fmt.Println(msg)
    }
    os.Exit(exitCode.Code)
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

func AddTrailingSlashIfNeeded(url string) string {
    if url != "" && !strings.HasSuffix(url, "/") {
        url += "/"
    }
    return url
}

type Flags struct {
    ArtDetails *ArtifactoryDetails
    DryRun bool
    Props string
    Deb string
    Recursive bool
    Flat bool
    UseRegExp bool
    Threads int
    MinSplitSize int64
    SplitCount int
    Interactive bool
    EncPassword bool
}

type ArtifactoryDetails struct {
    Url string
    User string
    Password string
    SshKeyPath string
    SshAuthHeaders map[string]string
}