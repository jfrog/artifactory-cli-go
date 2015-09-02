package commands

import (
  "strings"
  "encoding/json"
  "github.com/JFrogDev/artifactory-cli-go/utils"
)

func Download(url string, downloadPattern string, recursive bool, props string, user string, password string, flat bool, dryRun bool) {
    aqlUrl := url + "api/search/aql"
    if strings.HasSuffix(downloadPattern, "/") {
        downloadPattern += "*"
    }

    data := utils.BuildAqlSearchQuery(downloadPattern, recursive, props)

    println("AQL query: " + data)

    json := utils.SendPost(aqlUrl, data, user, password)
    resultItems := parseAqlSearchResponse(json)
    size := len(resultItems)

    for i := 0; i < size; i++ {
        downloadPath := buildDownloadUrl(url, resultItems[i])
        print("Downloading " + downloadPath + "...")

        localFilePath := resultItems[i].Path + "/" + resultItems[i].Name
        if shouldDownloadFile(localFilePath, downloadPath, user, password) {
            resp := utils.DownloadFile(downloadPath, resultItems[i].Path, resultItems[i].Name, flat, user, password, dryRun)
            if !dryRun {
                println("Artifactory response:", resp.Status)
            } else {
                println()
            }
        } else {
            println("File already exists locally.")
        }
    }
}

func buildDownloadUrl(baseUrl string, resultItem AqlSearchResultItem) string {
    if resultItem.Path == "." {
        return baseUrl + resultItem.Repo + "/" + resultItem.Name
    }
    return baseUrl + resultItem.Repo + "/" + resultItem.Path + "/" + resultItem.Name
}

func shouldDownloadFile(localFilePath string, downloadPath string, user string, password string) bool {
    if !utils.IsFileExists(localFilePath) {
        return true
    }
    localChecksum := utils.CalcChecksum(utils.ReadFile(localFilePath))
    artifactoryChecksum := utils.FetchChecksumFromArtifactory(downloadPath, user, password)
    if localChecksum.Md5 != artifactoryChecksum.Md5 || localChecksum.Sha1 != artifactoryChecksum.Sha1 {
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