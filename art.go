package main

import (
  "strings"
  "os"
  "github.com/codegangsta/cli"
  "github.com/JFrogDev/artifactory-cli-go/commands"
  "github.com/JFrogDev/artifactory-cli-go/utils"
)

var dryRun bool
var url string
var user string
var password string
var props string
var recursive bool
var flat bool
var useRegExp bool

func main() {
    app := cli.NewApp()
    app.Name = "Artifactory CLI"
    app.Usage = "See https://github.com/JFrogDev/artifactory-cli-go for usage instructions."

    app.Commands = []cli.Command{
        {
            Name: "upload",
            Flags: GetUploadFlags(),
            Aliases: []string{"u"},
            Usage: "upload <local path> <repo path>",
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
         Usage: "[Mandatiry] Artifactory URL",
        },
        cli.StringFlag{
         Name:  "user",
         Usage: "[Optional] Artifactory user",
        },
        cli.StringFlag{
         Name:  "password",
         Usage: "[Optional] Artifactory password",
        },
    }
}

func GetUploadFlags() []cli.Flag {
    flags := []cli.Flag{
        nil,nil,nil,nil,nil,nil,nil,nil,
    }
    copy(flags[0:3], GetFlags())
    flags[3] = cli.StringFlag{
         Name:  "props",
         Usage: "[Optional] List of properties in the form of key1=value1;key2=value2,... to be attached to the uploaded artifacts.",
    }
    flags[4] = cli.StringFlag{
        Name:  "recursive",
        Value:  "",
        Usage: "[Default: true] Set to false if you do not wish to collect artifacts in sub-folders to be uploaded to Artifactory.",
    }
    flags[5] = cli.BoolFlag{
        Name:  "flat",
        Usage: "[Default: false] If not set to true, and the upload path ends with a slash, files are uploaded according to their file system hierarchy.",
    }
    flags[6] = cli.BoolFlag{
         Name:  "regexp",
         Usage: "[Default: false] Set to true to use a regular expression instead of wildcards expression to collect files to upload.",
    }
    flags[7] = cli.BoolFlag{
         Name:  "dry-run",
         Usage: "[Default: false] Set to true to disable communication with Artifactory.",
    }
    return flags
}

func GetDownloadFlags() []cli.Flag {
    flags := []cli.Flag{
        nil,nil,nil,nil,nil,nil,
    }
    copy(flags[0:3], GetFlags())
    flags[3] = cli.StringFlag{
         Name:  "props",
         Usage: "[Optional] List of properties in the form of key1=value1;key2=value2,... Only artifacts with these properties will be downloaded.",
    }
    flags[4] = cli.StringFlag{
        Name:  "recursive",
        Value:  "",
        Usage: "[Default: true] Set to false if you do not wish to include the download of artifacts inside sub-folders in Artifactory.",
    }
    flags[5] = cli.BoolFlag{
        Name:  "flat",
        Usage: "[Default: false] Set to true if you do not wish to have the Artifactory repository path structure created locally for your downloaded files.",
    }
    return flags
}

func InitFlags(c *cli.Context) {
    url = GetMandatoryFlag(c, "url")
    if !strings.HasSuffix(url, "/") {
        url += "/"
    }

    user = c.String("user")
    password = c.String("password")
    props = c.String("props")
    dryRun = c.Bool("dry-run")
    flat = c.Bool("flat")
    useRegExp = c.Bool("regexp")

    if c.String("recursive") == "" {
        recursive = true
    } else {
        recursive = c.Bool("recursive")
    }
}

func Download(c *cli.Context) {
    InitFlags(c)
    if len(c.Args()) != 1 {
        utils.Exit("Wrong number of arguments. Try 'art download --help'.")
    }
    pattern := c.Args()[0]
    commands.Download(url, pattern, recursive, props, user, password, flat, dryRun)
}

func Upload(c *cli.Context) {
    InitFlags(c)
    size := len(c.Args())
    if size != 2 {
        utils.Exit("Wrong number of arguments. Try 'art upload --help'.")
    }
    localPath := c.Args()[0]
    targetPath := c.Args()[1]
    commands.Upload(url, localPath, targetPath, recursive, flat, props, user, password, useRegExp, dryRun)
}

// Get a CLI flagg. If the flag does not exist, utils.Exit with a message.
func GetMandatoryFlag(c *cli.Context, flag string) string {
    value := c.String(flag)
    if value == "" {
        utils.Exit("The --" + flag + " flag is mandatory")
    }
    return value
}