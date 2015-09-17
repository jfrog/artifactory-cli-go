package commands

import (
    "os"
    "bytes"
    "os/user"
    "io/ioutil"
    "encoding/json"
    "github.com/JFrogDev/artifactory-cli-go/utils"
)

func Config(props map[string]string) {
    oldProps := readConfFile()
    println(len(oldProps))
    writeConfFile(props)
}

func getConFilePath() string {
    userDir, err := user.Current()
    utils.CheckError(err)
    confPath := userDir.HomeDir + "/.jfrog/"
    os.MkdirAll(confPath ,0777)
    return confPath + "art-cli.conf"
}

func writeConfFile(props map[string]string) {
    confFilePath := getConFilePath()
    if !utils.IsFileExists(confFilePath) {
        out, err := os.Create(confFilePath)
        utils.CheckError(err)
        defer out.Close()
    }

    b, err := json.Marshal(&props)
    utils.CheckError(err)
    var content bytes.Buffer
    err = json.Indent(&content, b, "", "  ")
    utils.CheckError(err)

    ioutil.WriteFile(confFilePath,[]byte(content.String()), 0x777)
}

func readConfFile() map[string]string {
    confFilePath := getConFilePath()
    if !utils.IsFileExists(confFilePath) {
        return make(map[string]string)
    }
    content := utils.ReadFile(confFilePath)
    props := make(map[string]string)
    json.Unmarshal(content, &props)
    return props
}