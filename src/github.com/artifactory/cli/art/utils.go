package main

import (
  "os"
)

func CheckError(err error) {
    if err != nil {
        panic(err)
    }
}

func Exit(msg string) {
  println(msg)
  os.Exit(1)
}

func BuildAqlJson(repo string, path string, name string) string {
    json :=
    "{" +
        "\"repo\": {" +
            "\"$match\":" + "\"" + repo + "\"" +
        "}," +
        "\"path\": {" +
            "\"$match\":" + "\"" +  path + "\"" +
        "}," +
        "\"name\":{" +
            "\"$match\":" + "\"" + name + "\"" +
        "}" +
    "}"

    return json
}