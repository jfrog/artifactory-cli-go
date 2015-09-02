package commands

import (
  "strings"
  "regexp"
  "strconv"
  "runtime"
  "net/http"
  "github.com/JFrogDev/artifactory-cli-go/utils"
)

func Upload(url string, localPath string, targetPath string, recursive bool, flat bool, props string, user string, password string, useRegExp bool, dryRun bool) {
    artifacts := getFilesToUpload(localPath, targetPath, recursive, flat, useRegExp)

    for _, artifact := range artifacts {
        target := url + artifact.TargetPath
        uploadFile(artifact.LocalPath, target, props, user, password, dryRun)
    }
}

func prepareUploadPath(path string) string {
    path = strings.Replace(path, "\\", "/", -1)
    path = strings.Replace(path, "../", "", -1)
    path = strings.Replace(path, "./", "", -1)
    return path
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
    var wildcard string
    if runtime.GOOS == "windows" {
        wildcard = ".*"
    } else {
        wildcard = ".\\*"
    }
    localpath = strings.Replace(localpath, ".", "\\.", -1)
    localpath = strings.Replace(localpath, "*", wildcard, -1)
    return localpath
}

func getFilesToUpload(localpath string, targetPath string, recursive bool, flat bool, useRegExp bool) []Artifact {
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

    var paths []string
    if recursive {
        paths = utils.ListFilesRecursive(rootPath)
    } else {
        paths = utils.ListFiles(rootPath)
    }

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
                if flat {
                    target += utils.GetFileNameFromPath(path)
                } else {
                    uploadPath := prepareUploadPath(path)
                    target += uploadPath
                }
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
    if rootPath == "" {
        return "."
    }
    return rootPath
}

func uploadFile(localPath string, targetPath string, props string, user string, password string, dryRun bool) {
    if (props != "") {
        targetPath += ";" + props
    }
    print("Uploading artifact: " + targetPath + "...")
    fileContent := utils.ReadFile(localPath)

    var deployed bool = false
    var resp *http.Response
    if len(fileContent) >= 10240 {
        resp = tryChecksumDeploy(fileContent, targetPath, user, password, dryRun)
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

func tryChecksumDeploy(fileContent []byte, targetPath string, user string, password string, dryRun bool) *http.Response {
    checksum := utils.CalcChecksum(fileContent)

    headers := make(map[string]string)
    headers["X-Checksum-Deploy"] = "true"
    headers["X-Checksum-Sha1"] = checksum.Sha1
    headers["X-Checksum-Md5"] = checksum.Md5

    return utils.PutContent(nil, headers, targetPath, user, password, dryRun)
}

type Artifact struct {
    LocalPath string
    TargetPath string
}