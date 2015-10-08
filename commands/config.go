package commands

import (
    "os"
    "fmt"
    "bytes"
    "os/user"
    "io/ioutil"
    "encoding/json"
    "github.com/JFrogDev/artifactory-cli-go/utils"
)

func Config(details *utils.ArtifactoryDetails, interactive bool) {
    if interactive {
        if details.Url == "" {
            print("Url: ")
            fmt.Scanln(&details.Url)
        }
        if details.User == "" {
            print("User: ")
            fmt.Scanln(&details.User)
        }
        if details.Password == "" {
            print("Password: ")
            fmt.Scanln(&details.Password)
        }
    }
    writeConfFile(details)
}

func ShowConfig() {
    details := readConfFile()
    if details.Url != "" {
        println("Url: " + details.Url)
    }
    if details.User != "" {
        println("User: " + details.User)
    }
    if details.Password != "" {
        println("Password: " + details.Password)
    }
}

func ClearConfig() {
    writeConfFile(new(utils.ArtifactoryDetails))
}

func GetConfig() *utils.ArtifactoryDetails {
    return readConfFile()
}

func getConFilePath() string {
    userDir, err := user.Current()
    utils.CheckError(err)
    confPath := userDir.HomeDir + "/.jfrog/"
    os.MkdirAll(confPath ,0777)
    return confPath + "art-cli.conf"
}

func writeConfFile(details *utils.ArtifactoryDetails) {
    confFilePath := getConFilePath()
    if !utils.IsFileExists(confFilePath) {
        out, err := os.Create(confFilePath)
        utils.CheckError(err)
        defer out.Close()
    }

    b, err := json.Marshal(&details)
    utils.CheckError(err)
    var content bytes.Buffer
    err = json.Indent(&content, b, "", "  ")
    utils.CheckError(err)

    ioutil.WriteFile(confFilePath,[]byte(content.String()), 0x777)
}

func readConfFile() *utils.ArtifactoryDetails {
    confFilePath := getConFilePath()
    details := new(utils.ArtifactoryDetails)
    if !utils.IsFileExists(confFilePath) {
        return details
    }
    content := utils.ReadFile(confFilePath)
    json.Unmarshal(content, &details)

    return details
}