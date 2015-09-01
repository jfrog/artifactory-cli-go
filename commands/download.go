package commands

import (
  "strings"
  "encoding/json"
  "github.com/JFrogDev/artifactory-cli-go/utils"
)

func Download(url string, downloadPattern string, props string, user string, password string, flat bool, dryRun bool) {
    aqlUrl := url + "api/search/aql"
    if strings.HasSuffix(downloadPattern, "/") {
        downloadPattern += "*"
    }

    data := utils.BuildAqlSearchQuery(downloadPattern, props)

    println("AQL query: " + data)

    json := utils.SendPost(aqlUrl, data, user, password)
    resultItems := parseAqlSearchResponse(json)
    size := len(resultItems)

    for i := 0; i < size; i++ {
        downloadPath := url + resultItems[i].Repo + "/" + resultItems[i].Path + "/" + resultItems[i].Name
        print("Downloading " + downloadPath + "...")

        localFilePath := resultItems[i].Path + "/" + resultItems[i].Name
        if utils.ShouldDownloadFile(localFilePath, downloadPath, user, password) {
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