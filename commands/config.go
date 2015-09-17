package commands

import (
    "bytes"
    "encoding/json"
    //"os/user"
    "github.com/JFrogDev/artifactory-cli-go/utils"
)

func Config() {
    props := make(map[string]string)
    props["key1"] = "val1"
    props["key2"] = "val2"

    user := &props
    b, err := json.Marshal(user)
    utils.CheckError(err)

    var out bytes.Buffer
    err = json.Indent(&out, b, "", "  ")
    utils.CheckError(err)
    println(string(out.Bytes()))

    /*
    usr, err := user.Current()
    utils.CheckError(err)
    confFile = usr.HomeDir + "/.jfrog/cli.conf"
    println("Looking for config file '" + confFile + "'")
    if len(c.Args()) == 0 {
        if !utils.IsPathExists(confFile) {
            println("CLI conf file does not exists")
        } else {
            println("CLI conf file content:")
            // TODO: Read the flags from the conf and display
        }
    } else {
        key := c.Args()[0]
        val := c.Args()[1]
        println("Adding " + key + "=" + val + " to the CLI conf file")
        // TODO: Add or modify the entry in the conf file and create the file if needed
    }
    */
}