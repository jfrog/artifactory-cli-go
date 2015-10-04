package commands

import (
  "sync"
  "strconv"
  "encoding/json"
  "github.com/JFrogDev/artifactory-cli-go/utils"
)

func Download(downloadPattern string, flags *utils.Flags) {
    aqlUrl := flags.ArtDetails.Url + "api/search/aql"
    data := utils.BuildAqlSearchQuery(downloadPattern, flags.Recursive, flags.Props)

    println("Searching Artifactory using AQL query: " + data)
    resp, json := utils.SendPost(aqlUrl, []byte(data), flags.ArtDetails.User, flags.ArtDetails.Password)
    println("Artifactory response:", resp.Status)

    if resp.StatusCode == 200 {
        resultItems := parseAqlSearchResponse(json)
        downloadFiles(resultItems, flags)
        println("Downloaded " + strconv.Itoa(len(resultItems)) + " artifacts from Artifactory.")
    }
}

func downloadFiles(resultItems []AqlSearchResultItem, flags *utils.Flags) {
    size := len(resultItems)
    var wg sync.WaitGroup
    for i := 0; i < flags.Threads; i++ {
        wg.Add(1)
        go func(threadId int) {
            for j := threadId; j < size; j += flags.Threads {
                downloadPath := buildDownloadUrl(flags.ArtDetails.Url, resultItems[j])
                logMsgPrefix := utils.GetLogMsgPrefix(threadId, flags.DryRun)
                println(logMsgPrefix + " Downloading " + downloadPath)
                if !flags.DryRun {
                    downloadFile(downloadPath, resultItems[j].Path, resultItems[j].Name, logMsgPrefix, flags)
                }
            }
            wg.Done()
        }(i)
    }
    wg.Wait()
}

func downloadFile(downloadPath, localPath, localFileName, logMsgPrefix string, flags *utils.Flags) {
    details := utils.GetFileDetailsFromArtifactory(downloadPath, flags.ArtDetails.User, flags.ArtDetails.Password)
    localFilePath := localPath + "/" + localFileName
    if shouldDownloadFile(localFilePath, details, flags.ArtDetails.User, flags.ArtDetails.Password) {
        if flags.SplitCount == 0 || flags.MinSplitSize < 0 || flags.MinSplitSize*1000 > details.Size || !details.AcceptRanges {
            resp := utils.DownloadFile(downloadPath, localPath, localFileName, flags.Flat, flags.ArtDetails.User, flags.ArtDetails.Password)
            println(logMsgPrefix + " Artifactory response:", resp.Status)
        } else {
            utils.DownloadFileConcurrently(
                downloadPath, localPath, localFileName, logMsgPrefix, details.Size, flags)
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