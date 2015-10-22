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
    var password []byte
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
            password, err = terminal.ReadPassword(int(syscall.Stdin))
            utils.CheckError(err)
        }
    }
    updatedArtifactoryDetails := handlePasswordEncryption(details, string(password), shouldEncPassword)
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

/** 
 ** This function should handle all casses for encryption password:
 ** 1. User didn't disable encryption by using --enc-password=false:
 **     1.1. Artifactory is online and we are sending a request to get the encrypted password we want to save to config file
 **     1.2. Artifactory is online but configured to reject requests with encrypted password
 **         1.2.1. In this case user will be asked to allow us to save the plain text password inside the file
 ** 2. User used the --enc--password=false flag so the password is saved to the config file as is. 
 **/
func handlePasswordEncryption(details *utils.ArtifactoryDetails, password string, shouldEncPassword bool) *utils.ArtifactoryDetails { 
        passwordToSave := password
        if shouldEncPassword {
            encPassword, err, tryUnencrypted := utils.GetArtifactoryEncryptedPassword(details.Url, details.User, password)
            passwordToSave = encPassword
            if tryUnencrypted {
                passwordToSave = password
                var saveUnEncPasswrod string 
                printUnEncryptedPasswordWarningsToUser(err)
                fmt.Scanln(&saveUnEncPasswrod)
                if !userAllowedUnEncPassword(saveUnEncPasswrod) {
                    passwordToSave = ""
                }
            } else {
                utils.CheckError(err)        
            }
        }
        return &utils.ArtifactoryDetails { details.Url, details.User, passwordToSave }
}

func userAllowedUnEncPassword(allow string) bool{
    return allow == "yes" || allow == "y" || allow == "true"
}

func printUnEncryptedPasswordWarningsToUser(err error){
    println("\n" + err.Error())
    println("To avoid saving password to the config file you should enable encrption in Artifactory or use --password with each command")
    println("Save unencrypted password to config file? [yes/no]")
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