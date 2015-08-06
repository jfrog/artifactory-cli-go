package main

import (
  "strings"
)

func BuildAqlSearchQuery(searchPattern string) string {
    index := strings.Index(searchPattern, "/")
    if index == -1 {
        Exit("Invalid search pattern: " + searchPattern)
    }

    repo := searchPattern[:index]
    searchPattern = searchPattern[index+1:]

    Pairs := CreatePathFilePairs(searchPattern)
    Size := len(Pairs)

    json :=
        "{" +
            "\"repo\": \"" + repo + "\"," +
            "\"$or\": ["

    for i := 0; i < Size; i++ {
        json +=
            "{" +
                BuildInnerQuery(repo, Pairs[i].Path, Pairs[i].File) +
            "}"

        if (i+1 < Size) {
            json += ","
        }
    }

    json +=
            "]" +
        "}"


    return "items.find(" + json + ")"
}

func BuildInnerQuery(repo string, path string, name string) string {
    query :=
        "\"$and\": [{" +
            "\"path\": {" +
                "\"$match\":" + "\"" +  path + "\"" +
            "}," +
            "\"name\":{" +
                "\"$match\":" + "\"" + name + "\"" +
            "}" +
        "}]"

    return query
}

func CreatePathFilePairs(pattern string) []PathFilePair {
    Pairs := []PathFilePair{}
    if (pattern == "*") {
        Pairs = append(Pairs, PathFilePair{"*", "*"})
        return Pairs
    }

    Index := strings.LastIndex(pattern, "/")
    Path := ""
    if (Index > 0) {
        Path = pattern[0:Index]
        Name := pattern[Index+1:]
        Pairs = append(Pairs, PathFilePair{Path, Name})
        pattern = Name
    }

    Sections := strings.Split(pattern, "*")
    Size := len(Sections)

    for i := 0; i < Size; i++ {
        if (Sections[i] == "") {
            continue
        }

        Options := []string{}
        if (i > 0) {
            Options = append(Options, "/" + Sections[i])
        }
        if (i+1 < Size) {
            Options = append(Options, Sections[i] + "/")
        }

        for _, Option := range Options {
            Str := ""
            for j := 0; j < Size; j++ {
                if (j > 0) {
                    Str += "*"
                }
                if (j == i) {
                    Str += Option
                } else {
                    Str += Sections[j]
                }
            }
            Split := strings.Split(Str, "/")

            if (Path != "") {
                Path += "/"
            }
            Pairs = append(Pairs, PathFilePair{Path + Split[0], Split[1]})
        }
    }

    return Pairs
}

type PathFilePair struct {
    Path string
    File string
}