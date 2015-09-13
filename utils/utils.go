package utils

import (
  "os"
  "strconv"
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

func GetLogMsgPrefix(threadId int) string {
    return "[thread " + strconv.Itoa(threadId) + "]"
}