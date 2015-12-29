package utils

import (
    "io"
    "fmt"
    "bytes"
    "regexp"
    "strconv"
    "io/ioutil"
    "encoding/json"
    "github.com/JFrogDev/artifactory-cli-go/Godeps/_workspace/src/golang.org/x/crypto/ssh"
)

func SshAuthentication(details *ArtifactoryDetails) {
    _, host, port := parseUrl(details.Url)

    fmt.Println("Performing SSH authentication...")
    if details.SshKeyPath == "" {
        Exit(ExitCodeError, "Cannot invoke the SshAuthentication function with no SSH key path. ")
    }

    buffer, err := ioutil.ReadFile(details.SshKeyPath)
    CheckError(err)
	key, err := ssh.ParsePrivateKey(buffer)
	CheckError(err)
    sshConfig := &ssh.ClientConfig{
        User: "admin",
        Auth: []ssh.AuthMethod{
            ssh.PublicKeys(key),
        },
    }

    hostAndPort := host + ":" + strconv.Itoa(port)
    connection, err := ssh.Dial("tcp", hostAndPort, sshConfig)
    CheckError(err)
    defer connection.Close()

    session, err := connection.NewSession()
    CheckError(err)
    defer session.Close()

    stdout, err := session.StdoutPipe()
    CheckError(err)

    var buf bytes.Buffer
    go io.Copy(&buf, stdout)

    session.Run("cli-authenticate")

    var result SshAuthResult
    err = json.Unmarshal(buf.Bytes(), &result)
    CheckError(err)
    details.Url = AddTrailingSlashIfNeeded(result.Href)
    details.SshAuthHeaders = result.Headers
    fmt.Println("SSH authentication successful.")
}

func parseUrl(url string) (protocol, host string, port int) {
    pattern1 := "^(.+)://(.+):([0-9].+)/$"
    pattern2 := "^(.+)://(.+)$"

    r, err := regexp.Compile(pattern1)
    CheckError(err)
    groups := r.FindStringSubmatch(url)
    if len(groups) == 4 {
        protocol = groups[1]
        host = groups[2]
        port, err = strconv.Atoi(groups[3])
        if err != nil {
            Exit(ExitCodeError, "URL: " + url + " is invalid. Expecting ssh://<host>:<port> or http(s)://...")
        }
        return
    }

    r, err = regexp.Compile(pattern2)
    CheckError(err)
    groups = r.FindStringSubmatch(url)
    if len(groups) == 3 {
        protocol = groups[1]
        host = groups[2]
        port = 80
        return
    }
    return
}

type SshAuthResult struct {
    Href string
    Headers map[string]string
}