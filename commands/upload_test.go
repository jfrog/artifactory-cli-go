package commands

import (
    "testing"
    "strconv"
    "github.com/JFrogDev/artifactory-cli-go/utils"
    "github.com/JFrogDev/artifactory-cli-go/tests"
)

func TestSingleFileUpload(t *testing.T) {
    flags := tests.GetFlags()
    uploadCount1 := Upload("testdata/a.txt", "repo-local", flags)
    uploadCount2 := Upload("testdata/aa.txt", "repo-local", flags)
    if uploadCount1 != 1 {
        t.Error("Expected 1 file to be uploaded. Got " + strconv.Itoa(uploadCount1) + ".")
    }
    if uploadCount2 != 1 {
        t.Error("Expected 1 file to be uploaded. Got " + strconv.Itoa(uploadCount2) + ".")
    }
}

func TestPatternRecursiveUpload(t *testing.T) {
    flags := tests.GetFlags()
    flags.Recursive = true
    testPatternUpload(t, flags)
}

func TestPatternNonRecursiveUpload(t *testing.T) {
    flags := tests.GetFlags()
    flags.Recursive = false
    testPatternUpload(t, flags)
}

func testPatternUpload(t *testing.T, flags *utils.Flags) {
    sep := tests.GetFileSeperator()
    uploadCount1 := Upload("testdata" + sep + "*", "repo-local", flags)
    uploadCount2 := Upload("testdata" + sep + "a*", "repo-local", flags)
    uploadCount3 := Upload("testdata" + sep + "b*", "repo-local", flags)

    if uploadCount1 != 3 {
        t.Error("Expected 3 file to be uploaded. Got " + strconv.Itoa(uploadCount1) + ".")
    }
    if uploadCount2 != 2 {
        t.Error("Expected 2 file to be uploaded. Got " + strconv.Itoa(uploadCount2) + ".")
    }
    if uploadCount3 != 1 {
        t.Error("Expected 1 file to be uploaded. Got " + strconv.Itoa(uploadCount3) + ".")
    }
}
