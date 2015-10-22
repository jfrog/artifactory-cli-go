package commands

import (
    "os"
    "fmt"
    "bytes"
    "syscall"
    "io/ioutil"
    "encoding/json"
    "golang.org/x/crypto/ssh/terminal"
    "github.com/JFrogDev/artifactory-cli-go/utils"
)

func Config(details *utils.ArtifactoryDetails, interactive, shouldEncPassword bool) {
    password := details.Password
    var bytePassword []byte
    if interactive {
        if details.Url == "" {
            print("Artifactory Url: ")
            fmt.Scanln(&details.Url)
        }
        if details.User == "" {
            print("User: ")
            fmt.Scanln(&details.User)
        }
        if details.Password == "" {
            print("Password: ")
            var err error
            bytePassword, err = terminal.ReadPassword(int(syscall.Stdin))
            password = string(bytePassword)
            utils.CheckError(err)
        }
    }
    details.Url = utils.AddTrailingSlashIfNeeded(details.Url)
    updatedArtifactoryDetails := handlePasswordEncryption(details, password, shouldEncPassword)
    writeConfFile(updatedArtifactoryDetails)
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

func handlePasswordEncryption(details *utils.ArtifactoryDetails, password string, shouldEncPassword bool) *utils.ArtifactoryDetails {
    var passwordToSave string
    if shouldEncPassword {
        response, encPassword := utils.GetEncryptedPasswordFromArtifactory( &utils.ArtifactoryDetails { details.Url, details.User, password })
        switch response.StatusCode {
            // In case Artifactory does not allow encrypted password we should query the user for the not encrypted password
            // with a warning that says that the un-encrypted password will be saved to a file
            case 409:
                utils.Exit("\nYour Artifactory server is not configured to encrypt passwords\n" +
                        "You may use \"art config --enc-password=false\"")
            case 200:
                passwordToSave = encPassword
            default:
                utils.Exit("\nArtifactory response: " + response.Status)
        }
    } else {
        // In case we do not want to encrypt password we would save the user input
        passwordToSave = password
    }
    return &utils.ArtifactoryDetails { details.Url, details.User, passwordToSave }
}

func allowUnencryptedPassword(allow string) bool{
    return allow == "yes" || allow == "y" || allow == "true"
}

func getConFilePath() string {
    userDir := utils.GetHomeDir()
    if userDir == "" {
        utils.Exit("Couldn't find home directory. Make sure your HOME environment variable is set.")
    }
    confPath := userDir + "/.jfrog/"
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