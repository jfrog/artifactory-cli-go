package commands

import (
  "strings"
  "regexp"
  "strconv"
  "net/http"
  "github.com/JFrogDev/artifactory-cli-go/utils"
)

func Upload(url string, localPath string, targetPath string, user string, password string, useRegExp bool, dryRun bool) {
    artifacts := getFilesToUpload(localPath, targetPath, useRegExp)

    for _, artifact := range artifacts {
        target := url + artifact.TargetPath
        uploadFile(artifact.LocalPath, target, user, password, dryRun)
    }
}

func prepareLocalPath(localpath string, useRegExp bool) string {
    if strings.HasPrefix(localpath, "./") {
        localpath = localpath[2:]
    } else
    if strings.HasPrefix(localpath, ".\\") {
        localpath = localpath[3:]
    }
    if !useRegExp {
        localpath = localPathToRegExp(localpath)
    }
    return localpath
}

func localPathToRegExp(localpath string) string {
    localpath = strings.Replace(localpath, ".", "\\.", -1)
    localpath = strings.Replace(localpath, "*", ".*", -1)
    return localpath
}

func getFilesToUpload(localpath string, targetPath string, useRegExp bool) []Artifact {
    rootPath := getRootPath(localpath, useRegExp)
    if !utils.IsPathExists(rootPath) {
        utils.Exit("Path does not exist: " + rootPath)
    }
    localpath = prepareLocalPath(localpath, useRegExp)
    artifacts := []Artifact{}
    // If the path is a single file then return it
    if !utils.IsDir(rootPath) {
        artifacts = append(artifacts, Artifact{rootPath, targetPath})
        return artifacts
    }

    r, err := regexp.Compile(localpath)
    utils.CheckError(err)

    paths := utils.ListFiles(rootPath)
    for _, path := range paths {
        if utils.IsDir(path) {
            continue
        }

        groups := r.FindStringSubmatch(path)
        size := len(groups)
        target := targetPath
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

// Get the local root path, from which to start collecting artifacts to be uploaded to Artifactory.
func getRootPath(path string, useRegExp bool) string {
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
        if useRegExp {
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
        if useRegExp {
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

func uploadFile(localPath string, targetPath string, user string, password string, dryRun bool) {
    print("Uploading artifact: " + targetPath + "...")
    fileContent := utils.ReadFile(localPath)

    var deployed bool = false
    var resp *http.Response
    if len(fileContent) >= 10240 {
        resp = utils.TryChecksumDeploy(fileContent, targetPath, user, password, dryRun)
        deployed = !dryRun && (resp.StatusCode == 201 || resp.StatusCode == 200)
    }
    if !deployed {
        resp = utils.PutContent(fileContent, nil, targetPath, user, password, dryRun)
    }
    if !dryRun {
        println("Artifactory response: " + resp.Status)
    } else {
        println()
    }
}

type Artifact struct {
    LocalPath string
    TargetPath string
}