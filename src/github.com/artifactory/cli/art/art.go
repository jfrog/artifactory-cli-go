package main

import (
  "os"
  "strings"
  "regexp"
  "strconv"
  "github.com/codegangsta/cli"
  "encoding/json"
)

var LocalPath string

var DryRun bool
var Url string
var User string
var Password string
var TargetPath string

func main() {
    app := cli.NewApp()
    app.Name = "art"
    app.Usage = "Artifactory CLI"

    flags := []cli.Flag{
        cli.BoolFlag{
         Name:  "dry-run",
         Usage: "Set to true to disable communication with Artifactory",
        },
        cli.StringFlag{
         Name:  "url",
         Usage: "Artifactory URL",
        },
        cli.StringFlag{
         Name:  "user",
         Usage: "Artifactory user",
        },
        cli.StringFlag{
         Name:  "password",
         Usage: "Artifactory password",
        },
    }

    app.Commands = []cli.Command{
        {
            Name: "upload",
            Flags: flags,
            Aliases: []string{"u"},
            Usage: "upload <local path> <repo name:repo path>",
            Action: func(c *cli.Context) {
                Upload(c)
            },
        },
        {
            Name: "download",
            Flags: flags,
            Aliases: []string{"d"},
            Usage: "download <repo path>",
            Action: func(c *cli.Context) {
                Download(c)
            },
        },
    }

    app.Run(os.Args)
}

func InitFlags(c *cli.Context) {
    Url = GetMandatoryFlag(c, "url")
    if !strings.HasSuffix(Url, "/") {
        Url += "/"
    }

    User = c.String("user")
    Password = c.String("password")
    DryRun = c.Bool("dry-run")
}

func GetFilesToUpload() []Artifact {
    rootPath := GetRootPath(LocalPath)
    if !IsPathExists(rootPath) {
        Exit("Path does not exist: " + rootPath)
    }
    artifacts := []Artifact{}
    if !IsDir(rootPath) {
        artifacts = append(artifacts, Artifact{rootPath, TargetPath})
        return artifacts
    }
    r, err := regexp.Compile(LocalPath)
    CheckError(err)

    paths := ListFiles(rootPath)
    for _, path := range paths {
        groups := r.FindStringSubmatch(path)
        size := len(groups)
        target := TargetPath
        for i := 1; i < size; i++ {
            target = strings.Replace(target, "$" + strconv.Itoa(i), groups[i], -1)
        }
        if ( size > 0) {
            artifacts = append(artifacts, Artifact{path, target})
        }
    }
    return artifacts
}

func GetPathsFromArtifactory(repo string, path string) {
    aqlJson := BuildAqlJson(repo, path, "*")
    data := "items.find(" + aqlJson + ")"
    responseJson := SendPost(Url + "api/search/aql", data, User, Password)
    println(string(responseJson))

	var f interface{}
	err := json.Unmarshal(responseJson, &f)
	CheckError(err)
	m := f.(map[string]interface{})
	println(m)
}

func Download(c *cli.Context) {
    InitFlags(c)
    size := len(c.Args())
    if size != 1 {
        Exit("Wrong number of arguments")
    }

    CheckAndGetRepoPathFromArg(c.Args()[0])
    split := strings.Split(c.Args()[0], ":")
    repoName := split[0]
    repoPath := split[1]
    GetPathsFromArtifactory(repoName, repoPath)
}

func Upload(c *cli.Context) {
    InitFlags(c)
    size := len(c.Args())
    if size != 2 {
        Exit("Wrong number of arguments")
    }
    LocalPath = c.Args()[0]
    TargetPath = CheckAndGetRepoPathFromArg(c.Args()[1])
    artifacts := GetFilesToUpload()

    for _, artifact := range artifacts {
        target := Url + artifact.targetPath
        println("Uploading artifact " + artifact.localPath + " to " + target)
        if !DryRun {
            PutFile(artifact.localPath, target, User, Password)
        }
    }
}

// Get a CLI flagg. If the flag does not exist, exit with a message.
func GetMandatoryFlag(c *cli.Context, flag string) string {
    value := c.String(flag)
    if value == "" {
        Exit("The --" + flag + " flag is mandatory")
    }
    return value
}

// Get the local root path, from which to start collecting artifacts to be uploaded to Artifactory.
func GetRootPath(path string) string {
    index := strings.Index(path, "(")
    if index == -1 {
        return path
    }
    return path[0:index]
}

// Gets the Artifactory target path for artifacts deployment.
func CheckAndGetRepoPathFromArg(arg string) string {
    if strings.Index(arg, ":") == -1 {
        Exit("Invalid repo path format: '" + arg + "'. Should be [repo:path].")
    }
    return strings.Replace(arg, ":", "/", -1)
}

type Artifact struct {
    localPath string
    targetPath string
}