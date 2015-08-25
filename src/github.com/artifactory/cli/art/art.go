package main

import (
  "os"
  "strings"
  "regexp"
  "strconv"
  "net/http"
  "github.com/codegangsta/cli"
  "encoding/json"
)

var LocalPath string

var DryRun bool
var Url string
var User string
var Password string
var TargetPath string
var Flat bool
var UseRegExp bool

func main() {
    app := cli.NewApp()
    app.Name = "art"
    app.Usage = "Artifactory CLI"

    app.Commands = []cli.Command{
        {
            Name: "upload",
            Flags: GetUploadFlags(),
            Aliases: []string{"u"},
            Usage: "upload <local path> <repo name:repo path>",
            Action: func(c *cli.Context) {
                Upload(c)
            },
        },
        {
            Name: "download",
            Flags: GetDownloadFlags(),
            Aliases: []string{"d"},
            Usage: "download <repo path>",
            Action: func(c *cli.Context) {
                Download(c)
            },
        },
    }

    app.Run(os.Args)
}

func GetFlags() []cli.Flag {
    return []cli.Flag{
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
}

func GetUploadFlags() []cli.Flag {
    flags := []cli.Flag{
        nil,nil,nil,nil,nil,
    }
    copy(flags[0:3], GetFlags())
    flags[3] = cli.BoolFlag{
         Name:  "dry-run",
         Usage: "Set to true to disable communication with Artifactory",
    }
    flags[4] = cli.BoolFlag{
         Name:  "regexp",
         Usage: "Set to true to use a regular expression instead of wildcards expression to collect files to upload",
    }
    return flags
}

func GetDownloadFlags() []cli.Flag {
    flags := []cli.Flag{
        nil,nil,nil,nil,
    }
    copy(flags[0:3], GetFlags())
    flags[3] = cli.BoolFlag{
        Name:  "flat",
        Usage: "Set to true if you do not wish to have the Artifactory repository path structure created locally for your downloaded files",
    }
    return flags
}

func InitFlags(c *cli.Context) {
    Url = GetMandatoryFlag(c, "url")
    if !strings.HasSuffix(Url, "/") {
        Url += "/"
    }

    User = c.String("user")
    Password = c.String("password")
    DryRun = c.Bool("dry-run")
    Flat = c.Bool("flat")
    UseRegExp = c.Bool("regexp")
}

func GetFilesToUpload() []Artifact {
    rootPath := GetRootPath(LocalPath)
    if !IsPathExists(rootPath) {
        Exit("Path does not exist: " + rootPath)
    }
    if !UseRegExp {
        LocalPathToRegExp()
    }
    artifacts := []Artifact{}
    // If the path is a single file then return it
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
            target = strings.Replace(target, "{" + strconv.Itoa(i) + "}", groups[i], -1)
        }
        if (size > 0) {
            artifacts = append(artifacts, Artifact{path, target})
        }
    }
    return artifacts
}

func LocalPathToRegExp() {
    LocalPath = strings.Replace(LocalPath, ".", "\\.", -1)
    LocalPath = strings.Replace(LocalPath, "*", ".*", -1)
}

func Download(c *cli.Context) {
    InitFlags(c)
    if len(c.Args()) != 1 {
        Exit("Wrong number of arguments. Try 'art download --help'.")
    }

    url := Url + "api/search/aql"
    pattern := CheckAndGetRepoPathFromArg(c.Args()[0])
    data := BuildAqlSearchQuery(pattern)

    println("AQL query: " + data)

    json := SendPost(url, data, User, Password)
    resultItems := ParseAqlSearchResponse(json)
    size := len(resultItems)

    for i := 0; i < size; i++ {
        downloadPath := Url + resultItems[i].Repo + "/" + resultItems[i].Path + "/" + resultItems[i].Name
        DownloadFile(downloadPath, resultItems[i].Path, resultItems[i].Name, Flat)
    }
}

func Upload(c *cli.Context) {
    InitFlags(c)
    size := len(c.Args())
    if size != 2 {
        Exit("Wrong number of arguments. Try 'art upload --help'.")
    }
    LocalPath = c.Args()[0]
    TargetPath = CheckAndGetRepoPathFromArg(c.Args()[1])
    artifacts := GetFilesToUpload()

    for _, artifact := range artifacts {
        target := Url + artifact.targetPath
        UploadFile(artifact.localPath, target)
    }
}

func UploadFile(localPath string, targetPath string) {
    println("Uploading artifact: " + targetPath)
    fileContent := ReadFile(localPath)

    var deployed bool = false
    var resp *http.Response
    if len(fileContent) >= 10240 {
        resp = tryChecksumDeploy(fileContent, targetPath)
        deployed = (resp.StatusCode == 201 || resp.StatusCode == 200)
    }
    if !deployed {
        resp = PutContent(fileContent, nil, targetPath, User, Password, DryRun)
    }
    println("Artifactory response: " + resp.Status)
}

func ParseAqlSearchResponse(resp []byte) []AqlSearchResultItem {
    var result AqlSearchResult
    err := json.Unmarshal(resp, &result)

    CheckError(err)
    return result.Results
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

func CheckAndGetRepoPathFromArg(arg string) string {
    Sections := strings.Split(arg, ":")
    if len(Sections) != 2 || Sections[0] == "" || Sections[1] == "" {
        Exit("Invalid repo path format: '" + arg + "'. Should be [repo:path].")
    }
    path := strings.Replace(arg, ":/", "/", -1)
    return strings.Replace(path, ":", "/", -1)
}

type Artifact struct {
    localPath string
    targetPath string
}

type AqlSearchResult struct {
    Results []AqlSearchResultItem
}

type AqlSearchResultItem struct {
     Repo string
     Path string
     Name string
 }