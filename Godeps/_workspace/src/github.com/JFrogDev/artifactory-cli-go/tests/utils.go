package tests

import (
    "runtime"
    "github.com/JFrogDev/artifactory-cli-go/utils"
)

func GetFlags() *utils.Flags {
    flags := new(utils.Flags)
    flags.ArtDetails = new(utils.ArtifactoryDetails)
    flags.DryRun = true
    flags.EncPassword = true
    flags.Threads = 3

    return flags
}

func GetFileSeperator() string {
    if runtime.GOOS == "windows" {
        return "\\\\"
    }
    return "/"
}
