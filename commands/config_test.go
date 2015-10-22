package commands

import (
    "testing"
    "reflect"
    "encoding/json"
    "github.com/JFrogDev/artifactory-cli-go/utils"
)

func TestConfig(t *testing.T){
    inputDetails := utils.ArtifactoryDetails { "http://localhost:8081/artifactory", "username", "password" }
    Config(&inputDetails, false, false)
    outputConfig := GetConfig()
    printConfigStruct(&inputDetails)
    printConfigStruct(outputConfig)
    if !reflect.DeepEqual(inputDetails, *outputConfig) {
        t.Error("Unexpected configuration was saved to file. Expected: " + configStructToString(&inputDetails) + " Got " + configStructToString(outputConfig))
    }
}

func configStructToString(artConfig *utils.ArtifactoryDetails) string {
    marshaledStruct, _ := json.Marshal(*artConfig)
    return string(marshaledStruct)
}

func printConfigStruct(artConfig *utils.ArtifactoryDetails){
    stringSturct := configStructToString(artConfig)
    println(stringSturct)
}