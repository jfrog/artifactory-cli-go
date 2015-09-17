package commands

import (
  "os"
  "strings"
  "regexp"
  "strconv"
  "sync"
  "net/http"
  "github.com/JFrogDev/artifactory-cli-go/utils"
)

func Upload(url, localPath, targetPath string, recursive bool, flat bool, props, user, password string, threads int, useRegExp, dryRun bool) {
    // Get the list of artifacts to be uploaded to Artifactory:
    artifacts := getFilesToUpload(localPath, targetPath, recursive, flat, useRegExp)
    size := len(artifacts)

    // Start a thread for each channel and start uploading:
    var wg sync.WaitGroup
    for i := 0; i < threads; i++ {
        wg.Add(1)
        go func(threadId int) {
            for j := threadId; j < size; j += threads {
                target := url + artifacts[j].TargetPath
                uploadFile(artifacts[j].LocalPath, target, props, user, password, dryRun, utils.GetLogMsgPrefix(threadId, dryRun))
            }
            wg.Done()
        }(i)
    }
    wg.Wait()
    println("Uploaded " + strconv.Itoa(size) + " artifacts to Artifactory.")
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
    wildcard := ".*"
    localpath = strings.Replace(localpath, ".", "\\.", -1)
    localpath = strings.Replace(localpath, "*", wildcard, -1)
    if strings.HasSuffix(localpath, "/") || strings.HasSuffix(localpath, "\\") {
        localpath += wildcard
    }
    localpath = "^" + localpath + "$"
    return localpath
}

func getFilesToUpload(localpath string, targetPath string, recursive bool, flat bool, useRegExp bool) []Artifact {
    if strings.Index(targetPath, "/") < 0 {
        targetPath += "/"
    }

    rootPath := getRootPath(localpath, useRegExp)
    if !utils.IsPathExists(rootPath) {
        utils.Exit("Path does not exist: " + rootPath)
    }
    localpath = prepareLocalPath(localpath, useRegExp)
    artifacts := []Artifact{}
    // If the path is a single file then return it
    if !utils.IsDir(rootPath) {
        targetPath := prepareUploadPath(targetPath + rootPath)
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

func uploadFile(localPath string, targetPath string, props string, user string, password string,
    dryRun bool, logMsgPrefix string) {
    if (props != "") {
        targetPath += ";" + props
    }
    println(logMsgPrefix + " Uploading artifact: " + targetPath)
    file, err := os.Open(localPath)
    utils.CheckError(err)
    defer file.Close()
    fileInfo, err := file.Stat()
    utils.CheckError(err)

    var checksumDeployed bool = false
    var resp *http.Response
    var details *utils.FileDetails
    if fileInfo.Size() >= 10240 {
        resp, details = tryChecksumDeploy(localPath, targetPath, user, password, dryRun)
        checksumDeployed = !dryRun && (resp.StatusCode == 201 || resp.StatusCode == 200)
    }
    if !dryRun && !checksumDeployed {
        resp = utils.UploadFile(file, targetPath, user, password, details)
    }
    if !dryRun {
        var strChecksumDeployed string
        if checksumDeployed {
            strChecksumDeployed = " (Checksum deploy)"
        } else {
            strChecksumDeployed = ""
        }
        println(logMsgPrefix + " Artifactory response: " + resp.Status + strChecksumDeployed)
    }
}

func tryChecksumDeploy(filePath string, targetPath, user, password string,
    dryRun bool) (*http.Response, *utils.FileDetails) {

    details := utils.GetFileDetails(filePath)
    headers := make(map[string]string)
    headers["X-Checksum-Deploy"] = "true"
    headers["X-Checksum-Sha1"] = details.Sha1
    headers["X-Checksum-Md5"] = details.Md5

    if dryRun {
        return nil, details
    }
    resp, _ := utils.SendPut(targetPath, nil, headers, user, password)
    return resp, details
}

type Artifact struct {
    LocalPath string
    TargetPath string
}