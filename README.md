## Artifactory CLI

Artifactory CLI provides a command line interface for uploading and downloading artifacts to and from Artifactory.

### Build the command line executable

Make sure you have a working Go environment. [See the install instructions](http://golang.org/doc/install).

CD to the directory where you want to create the *artifactory-cli-go* project.
Set the value of the *GOPATH* environment variable to the full path of this  directory.

Run the following command to create the *artifactory-cli-go* project:
```console
$ go get github.com/JFrogDev/artifactory-cli-go
```

CD into the following directory
```console
$ cd $GOPATH/src/github.com/JFrogDev/artifactory-cli-go
```

Create the Artifactory CLI executable by running:
```console
$ go install
```

The Artifactory CLI executable was created in $GOPATH/bin and is ready to be used.

### Usage

You can add the path of the CLI executable to your *PATH* environment variable, so that you can access it from any path.

#### General command structure
artifactory-cli-go should be followed by a command name (for example, upload), a list of options (for example, --url=http://...)
and the list of arguments for the command.
```console
$ artifactory-cli-go command-name options arguments
```

#### The Upload command

##### Function
Used to upload artifacts to Artifactory.

##### Options
```console
   --url        Artifactory URL
   --user       Artifactory user
   --password   Artifactory password
   --dry-run    Set to true to disable communication with Artifactory
   --regexp     Set to true to use a regular expression instead of wildcards expression to collect files to upload
```
##### Arguments
* The first argument is the path to the files to be uploaded to Artifactory.
The path can include a single file or multiple files, by using the * wildcard.
**Important:** If the path is provided as a regular expression (with the --regexp=true option) then
the first regular expression appearing as part of the argument must be enclosed in parenthesis.

* The second argument is the upload path in Artifactory.
The argument should have the following format: [repository name]:[repository path]
The path can include symbols in the form of {1}, {2}, ...
These symbols are replaced with the sections enclosed with parenthesis in the first argument.

##### Examples

This example uploads the 'froggy.tgz' file to the root of the *my-local-repo* repository
```console
$ artifactory-cli-go upload froggy.tgz my-local-repo:/ --url=http://domain/artifactory --user=admin --password=password
```


This example collects all the zip files located under the build directory (including sub-directories)
   and uploads them to the *my-local-repo* repository, under the zipFiles folder, while keeping the files original names.
   ```console
   $ artifactory-cli-go upload build/*.zip libs-release-local:zipFiles/ --url=http://domain/artifactory --user=admin --password=password
   ```

#### The Download command

##### Function
Used to download artifacts from Artifactory.

##### Options
```console
   --url        Artifactory URL
   --user       Artifactory user
   --password   Artifactory password
   --flat       Set to true if you do not wish to have the Artifactory repository path structure created locally for your downloaded files
```

##### Arguments
The command expects one argument - the path of files to be downloaded from Artifactory.
The argument should have the following format: [repository name]:[repository path]
The path can include a single file or multiple files, by using the * wildcard.
The artifacts are downloaded and saved to the current directory, while saving their folder structure.

##### Examples

This example downloads the *cool-froggy.zip* file located at the root of the *my-local-repo* repository to current directory.
```console
$ artifactory-cli-go download my-local-repo:cool-froggy.zip --url=http://domain/artifactory --user=admin --password=password
```

This example downloads all files located in the *my-local-repo* repository under the *all-my-frogs* folder to the *all-my-frog* directory located unde the current directory.
```console
$ artifactory-cli-go download my-local-repo:all-my-frogs/ --url=http://domain/artifactory --user=admin --password=password
```