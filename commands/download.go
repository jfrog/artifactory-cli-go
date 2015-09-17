package commands

import (
  "sync"
  "strconv"
  "encoding/json"
  "github.com/JFrogDev/artifactory-cli-go/utils"
)

func Download(url, downloadPattern, props, user, password string, recursive, flat, dryRun bool,
    minSplitSize int64, splitCount, threads int) {
    aqlUrl := url + "api/search/aql"
    data := utils.BuildAqlSearchQuery(downloadPattern, recursive, props)

    println("Searching Artifactory using AQL query: " + data)
    resp, json := utils.SendPost(aqlUrl, []byte(data), user, password)
    println("Artifactory response:", resp.Status)

    if resp.StatusCode == 200 {
        resultItems := parseAqlSearchResponse(json)
        downloadFiles(resultItems, url, user, password, flat, dryRun, minSplitSize, splitCount, threads)
        println("Downloaded " + strconv.Itoa(len(resultItems)) + " artifacts from Artifactory.")
    }
}

func downloadFiles(resultItems []AqlSearchResultItem, url, user, password string, flat bool, dryRun bool,
    minSplitSize int64, splitCount, threads int) {
    size := len(resultItems)
    var wg sync.WaitGroup
    for i := 0; i < threads; i++ {
        wg.Add(1)
        go func(threadId int) {
            for j := threadId; j < size; j += threads {
                downloadPath := buildDownloadUrl(url, resultItems[j])
                logMsgPrefix := utils.GetLogMsgPrefix(threadId, dryRun)
                println(logMsgPrefix + " Downloading " + downloadPath)
                if !dryRun {
                    downloadFile(downloadPath, resultItems[j].Path, resultItems[j].Name,
                        user, password, flat, splitCount, minSplitSize, logMsgPrefix)
                }
            }
            wg.Done()
        }(i)
    }
    wg.Wait()
}

func downloadFile(downloadPath, localPath, localFileName, user, password string, flat bool,
    splitCount int, minSplitSize int64, logMsgPrefix string) {

    details := utils.GetFileDetailsFromArtifactory(downloadPath, user, password)
    localFilePath := localPath + "/" + localFileName
    if shouldDownloadFile(localFilePath, details, user, password) {
        if splitCount == 0 || minSplitSize < 0 || minSplitSize*1000 > details.Size || !details.AcceptRanges {
            resp := utils.DownloadFile(downloadPath, localPath, localFileName, flat, user, password)
            println(logMsgPrefix + " Artifactory response:", resp.Status)
        } else {
            utils.DownloadFileConcurrently(
                downloadPath, localPath, localFileName, flat, user, password,
                    details.Size, splitCount, logMsgPrefix)
        }
    } else {
        println(logMsgPrefix + " File already exists locally.")
    }
}

func buildDownloadUrl(baseUrl string, resultItem AqlSearchResultItem) string {
    if resultItem.Path == "." {
        return baseUrl + resultItem.Repo + "/" + resultItem.Name
    }
    return baseUrl + resultItem.Repo + "/" + resultItem.Path + "/" + resultItem.Name
}

func shouldDownloadFile(localFilePath string, artifactoryFileDetails *utils.FileDetails, user string, password string) bool {
    if !utils.IsFileExists(localFilePath) {
        return true
    }
    localFileDetails := utils.GetFileDetails(localFilePath)
    if localFileDetails.Md5 != artifactoryFileDetails.Md5 || localFileDetails.Sha1 != artifactoryFileDetails.Sha1 {
       return true
    }
    return false
}

func parseAqlSearchResponse(resp []byte) []AqlSearchResultItem {
    var result AqlSearchResult
    err := json.Unmarshal(resp, &result)
    utils.CheckError(err)
    return result.Results
}

type AqlSearchResult struct {
    Results []AqlSearchResultItem
}

type AqlSearchResultItem struct {
     Repo string
     Path string
     Name string
 }