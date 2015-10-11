package main

import (
	"github.com/JFrogDev/artifactory-cli-go/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/JFrogDev/artifactory-cli-go/commands"
	"github.com/JFrogDev/artifactory-cli-go/utils"
	"os"
	"strconv"
	"strings"
)

var flags = new(utils.Flags)

func main() {
    defer utils.RemoveTempDir()

    app := cli.NewApp()
    app.Name = "art"
    app.Usage = "See https://github.com/JFrogDev/artifactory-cli-go for usage instructions."
    app.Version = "1.0.0"

    app.Commands = []cli.Command{
        {
            Name: "config",
            Flags: getConfigFlags(),
            Aliases: []string{"c"},
            Usage: "config",
            Action: func(c *cli.Context) {
                config(c)
            },
        },
        {
            Name: "upload",
            Flags: getUploadFlags(),
            Aliases: []string{"u"},
            Usage: "upload <local path> <repo path>",
            Action: func(c *cli.Context) {
                upload(c)
            },
        },
        {
            Name: "download",
            Flags: getDownloadFlags(),
            Aliases: []string{"d"},
            Usage: "download <repo path>",
            Action: func(c *cli.Context) {
                download(c)
            },
        },
    }

    app.Run(os.Args)
}

func getFlags() []cli.Flag {
    return []cli.Flag{
        cli.StringFlag{
         Name:  "url",
         Usage: "[Mandatory] Artifactory URL",
        },
        cli.StringFlag{
         Name:  "user",
         Usage: "[Optional] Artifactory username",
        },
        cli.StringFlag{
         Name:  "password",
         Usage: "[Optional] Artifactory password",
        },
    }
}

func getUploadFlags() []cli.Flag {
    flags := []cli.Flag{
        nil,nil,nil,nil,nil,nil,nil,nil,nil,
    }
    copy(flags[0:3], getFlags())
    flags[3] = cli.StringFlag{
         Name:  "props",
         Usage: "[Optional] List of properties in the form of key1=value1;key2=value2,... to be attached to the uploaded artifacts.",
    }
    flags[4] = cli.StringFlag{
        Name:  "recursive",
        Value:  "",
        Usage: "[Default: true] Set to false if you do not wish to collect artifacts in sub-folders to be uploaded to Artifactory.",
    }
    flags[5] = cli.StringFlag{
        Name:  "flat",
        Value:  "",
        Usage: "[Default: true] If not set to true, and the upload path ends with a slash, files are uploaded according to their file system hierarchy.",
    }
    flags[6] = cli.BoolFlag{
         Name:  "regexp",
         Usage: "[Default: false] Set to true to use a regular expression instead of wildcards expression to collect files to upload.",
    }
    flags[7] = cli.StringFlag{
         Name:  "threads",
         Value:  "",
         Usage: "[Default: 3] Number of artifacts to upload in parallel.",
    }
    flags[8] = cli.BoolFlag{
         Name:  "dry-run",
         Usage: "[Default: false] Set to true to disable communication with Artifactory.",
    }
    return flags
}

func getDownloadFlags() []cli.Flag {
    flags := []cli.Flag{
        nil,nil,nil,nil,nil,nil,nil,nil,nil,
    }
    copy(flags[0:3], getFlags())
    flags[3] = cli.StringFlag{
         Name:  "props",
         Usage: "[Optional] List of properties in the form of key1=value1;key2=value2,... Only artifacts with these properties will be downloaded.",
    }
    flags[4] = cli.StringFlag{
        Name:  "recursive",
        Value:  "",
        Usage: "[Default: true] Set to false if you do not wish to include the download of artifacts inside sub-folders in Artifactory.",
    }
    flags[5] = cli.StringFlag{
        Name:  "flat",
        Value:  "",
        Usage: "[Default: false] Set to true if you do not wish to have the Artifactory repository path structure created locally for your downloaded files.",
    }
    flags[6] = cli.StringFlag{
        Name:  "min-split",
        Value:  "",
        Usage: "[Default: 5120] Minimum file size in KB to split into ranges when downloading. Set to -1 for no splits.",
    }
    flags[7] = cli.StringFlag{
        Name:  "split-count",
        Value:  "",
        Usage: "[Default: 3] Number of parts to split a file when downloading. Set to 0 for no splits.",
    }
    flags[8] = cli.StringFlag{
         Name:  "threads",
         Value:  "",
         Usage: "[Default: 3] Number of artifacts to download in parallel.",
    }
    return flags
}

func getConfigFlags() []cli.Flag {
    flags := []cli.Flag{
        nil,nil,nil,nil,
    }
    flags[0] = cli.StringFlag{
         Name:  "interactive",
         Usage: "[Default: true] Set to false if you do not want the config command to be interactive. If true, the --url option becomes optional.",
    }
    copy(flags[1:4], getFlags())
    return flags
}

func initFlags(c *cli.Context, cmd string) {
    if c.String("recursive") == "" {
        flags.Recursive = true
    } else {
        flags.Recursive = c.Bool("recursive")
    }
    if c.String("interactive") == "" {
        flags.Interactive = true
    } else {
        flags.Interactive = c.Bool("interactive")
    }

    if cmd == "config" {
        flags.ArtDetails = getArtifactoryDetails(c, false)
        if !flags.Interactive && flags.ArtDetails.Url == "" {
            utils.Exit("The --url option is mandatory when the --interactive option is set to false")
        }
    } else {
        flags.ArtDetails = getArtifactoryDetails(c, true)
        if flags.ArtDetails.Url == "" {
            utils.Exit("The --url option is mandatory")
        }
    }

    strFlat := c.String("flat")
    if cmd == "upload" {
        if strFlat == "" {
            flags.Flat = true
        } else {
            flags.Flat, _ = strconv.ParseBool(strFlat)
        }
    } else {
        if strFlat == "" {
            flags.Flat = false
        } else {
             flags.Flat, _ = strconv.ParseBool(strFlat)
         }
    }

    flags.Props = c.String("props")
    flags.DryRun = c.Bool("dry-run")
    flags.UseRegExp = c.Bool("regexp")
    var err error
    if c.String("threads") == "" {
        flags.Threads = 3
    } else {
        flags.Threads, err = strconv.Atoi(c.String("threads"))
        if err != nil || flags.Threads < 1 {
            utils.Exit("The '--threads' option should have a numeric positive value.")
        }
    }
    if c.String("min-split") == "" {
        flags.MinSplitSize = 5120
    } else {
        flags.MinSplitSize, err = strconv.ParseInt(c.String("min-split"), 10, 64)
        if err != nil {
            utils.Exit("The '--min-split' option should have a numeric value. Try 'art download --help'.")
        }
    }
    if c.String("split-count") == "" {
        flags.SplitCount = 3
    } else {
        flags.SplitCount, err = strconv.Atoi(c.String("split-count"))
        if err != nil {
            utils.Exit("The '--split-count' option should have a numeric value. Try 'art download --help'.")
        }
        if flags.SplitCount > 15 {
            utils.Exit("The '--split-count' option value is limitted to a maximum of 15.")
        }
        if flags.SplitCount < 0 {
            utils.Exit("The '--split-count' option cannot have a negative value.")
        }
    }
}

func config(c *cli.Context) {
    if len(c.Args()) > 1 {
        utils.Exit("Wrong number of arguments. Try 'art config --help'.")
    } else
    if len(c.Args()) == 1 {
        if c.Args()[0] == "show" {
            commands.ShowConfig()
        } else
        if c.Args()[0] == "clear" {
            commands.ClearConfig()
        } else {
            utils.Exit("Unknown argument '" + c.Args()[0] + "'. Available arguments are 'show' and 'clear'.")
        }
    } else {
        initFlags(c, "config")
        commands.Config(flags.ArtDetails, flags.Interactive)
    }
}

func download(c *cli.Context) {
    initFlags(c, "download")
    if len(c.Args()) != 1 {
        utils.Exit("Wrong number of arguments. Try 'art download --help'.")
    }
    pattern := c.Args()[0]
    commands.Download(pattern, flags)
}

func upload(c *cli.Context) {
    initFlags(c, "upload")
    size := len(c.Args())
    if size != 2 {
        utils.Exit("Wrong number of arguments. Try 'art upload --help'.")
    }
    localPath := c.Args()[0]
    targetPath := c.Args()[1]
    commands.Upload(localPath, targetPath, flags)
}

func getArtifactoryDetails(c *cli.Context, includeConfig bool) *utils.ArtifactoryDetails {
    details := new(utils.ArtifactoryDetails)
    details.Url = c.String("url")
    details.User = c.String("user")
    details.Password = c.String("password")

    if includeConfig {
        if details.Url == "" || details.User == "" || details.Password == "" {
            confDetails := commands.GetConfig()
            if details.Url == "" {
                details.Url = confDetails.Url
            }
            if details.User == "" {
                details.User = confDetails.User
            }
            if details.Password == "" {
                details.Password = confDetails.Password
            }
        }
    }
    if details.Url != "" && !strings.HasSuffix(details.Url, "/") {
        details.Url += "/"
    }
    return details
}