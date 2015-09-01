package main

import (
  "os"
  "strings"
  "regexp"
  "strconv"
  "net/http"
  "github.com/codegangsta/cli"
  "github.com/JFrogDev/artifactory-cli-go/utils"
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
    app.Name = "Artifactory CLI"
    app.Usage = "See https://github.com/JFrogDev/artifactory-cli-go for usage instructions."

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

func prepareLocalPath() {
    if strings.HasPrefix(LocalPath, "./") {
        LocalPath = LocalPath[2:]
    } else
    if strings.HasPrefix(LocalPath, ".\\") {
        LocalPath = LocalPath[3:]
    }
    if !UseRegExp {
        LocalPathToRegExp()
    }
}

func GetFilesToUpload() []Artifact {
    rootPath := GetRootPath(LocalPath)
    if !utils.IsPathExists(rootPath) {
        utils.Exit("Path does not exist: " + rootPath)
    }
    prepareLocalPath()
    artifacts := []Artifact{}
    // If the path is a single file then return it
    if !utils.IsDir(rootPath) {
        artifacts = append(artifacts, Artifact{rootPath, TargetPath})
        return artifacts
    }

    r, err := regexp.Compile(LocalPath)
    utils.CheckError(err)

    paths := utils.ListFiles(rootPath)
    for _, path := range paths {
        if utils.IsDir(path) {
            continue
        }

        groups := r.FindStringSubmatch(path)
        size := len(groups)
        target := TargetPath
        if (size > 0) {
            for i := 1; i < size; i++ {
                group := strings.Replace(groups[i], "\\", "/", -1)
                target = strings.Replace(target, "{" + strconv.Itoa(i) + "}", group, -1)
            }
            if strings.HasSuffix(target, "/") {
                target += utils.GetFileNameFromPath(path)
            }
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
        utils.Exit("Wrong number of arguments. Try 'art download --help'.")
    }

    url := Url + "api/search/aql"
    pattern := CheckAndGetRepoPathFromArg(c.Args()[0])
    if strings.HasSuffix(pattern, "/") {
        pattern += "*"
    }

    data := BuildAqlSearchQuery(pattern)

    println("AQL query: " + data)

    json := utils.SendPost(url, data, User, Password)
    resultItems := ParseAqlSearchResponse(json)
    size := len(resultItems)

    for i := 0; i < size; i++ {
        downloadPath := Url + resultItems[i].Repo + "/" + resultItems[i].Path + "/" + resultItems[i].Name
        print("Downloading " + downloadPath + "...")

        localFilePath := resultItems[i].Path + "/" + resultItems[i].Name
        if utils.ShouldDownloadFile(localFilePath, downloadPath, User, Password) {
            resp := utils.DownloadFile(downloadPath, resultItems[i].Path, resultItems[i].Name, Flat, User, Password, DryRun)
            if !DryRun {
                println("Artifactory response:", resp.Status)
            } else {
                println()
            }
        } else {
            println("File already exists locally.")
        }
    }
}

func Upload(c *cli.Context) {
    InitFlags(c)
    size := len(c.Args())
    if size != 2 {
        utils.Exit("Wrong number of arguments. Try 'art upload --help'.")
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
    print("Uploading artifact: " + targetPath + "...")
    fileContent := utils.ReadFile(localPath)

    var deployed bool = false
    var resp *http.Response
    if len(fileContent) >= 10240 {
        resp = utils.TryChecksumDeploy(fileContent, targetPath, User, Password, DryRun)
        deployed = !DryRun && (resp.StatusCode == 201 || resp.StatusCode == 200)
    }
    if !deployed {
        resp = utils.PutContent(fileContent, nil, targetPath, User, Password, DryRun)
    }
    if !DryRun {
        println("Artifactory response: " + resp.Status)
    } else {
        println()
    }
}

func ParseAqlSearchResponse(resp []byte) []AqlSearchResultItem {
    var result AqlSearchResult
    err := json.Unmarshal(resp, &result)

    utils.CheckError(err)
    return result.Results
}

// Get a CLI flagg. If the flag does not exist, utils.Exit with a message.
func GetMandatoryFlag(c *cli.Context, flag string) string {
    value := c.String(flag)
    if value == "" {
        utils.Exit("The --" + flag + " flag is mandatory")
    }
    return value
}

// Get the local root path, from which to start collecting artifacts to be uploaded to Artifactory.
func GetRootPath(path string) string {
    // The first step is to split the local path pattern into sections, by the file seperator.
    seperator := "/"
    sections := strings.Split(path, seperator)
    if len(sections) == 1 {
        seperator = "\\"
        sections = strings.Split(path, seperator)
    }

    // If we got only one section, meaning no file seperators, we should return "." if it is a pattern
    // or if it is not a pattern, return it as is.
    if len(sections) == 1 {
        if UseRegExp {
            if strings.Index(path, "(") != -1 {
                return "."
            }
        } else {
            if strings.Index(path, "*") != -1 {
                return "."
            }
        }
        return path
    }

    // Now we start building the root path, making sure to leave out the sub-directory that includes the pattern.
    rootPath := ""
    for _, section := range sections {
        if section == "" {
            continue
        }
        if UseRegExp {
            if strings.Index(section, "(") != -1 {
                break
            }
        } else {
            if strings.Index(section, "*") != -1 {
                break
            }
        }
        if rootPath != "" {
            rootPath += seperator
        }
        rootPath += section
    }
    return rootPath
}

func CheckAndGetRepoPathFromArg(arg string) string {
    Sections := strings.Split(arg, ":")
    if len(Sections) != 2 || Sections[0] == "" || Sections[1] == "" {
        utils.Exit("Invalid repo path format: '" + arg + "'. Should be [repo:path].")
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