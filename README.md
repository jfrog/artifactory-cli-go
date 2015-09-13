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
$ go build art.go
```

An executable file named *art* (the Artifactory CLI executable) was created in the current directory.

### Usage

You can copy the *art* executable to a different location on your file-system and add it
to your *PATH* environment variable, so that you can access it from any path.

#### General command structure
*art* should be followed by a command name (for example, upload), a list of options (for example, --url=http://...)
and the list of arguments for the command.
```console
$ art command-name options arguments
```

#### The Upload command

##### Function
Used to upload artifacts to Artifactory.

##### Options
```console
   --url          [Mandatory] Artifactory URL.
   --user         [Optional] Artifactory user.
   --password     [Optional] Artifactory password.
   --props        [Optional] List of properties in the form of key1=value1;key2=value2,... to be attached to the uploaded artifacts.
   --flat         [Default: false] If not set to true, and the upload path ends with a slash, artifacts are uploaded according to their file system hierarchy.
   --recursive    [Default: true] Set to false if you do not wish to collect artifacts in sub-folders to be uploaded to Artifactory.
   --regexp       [Default: false] Set to true to use a regular expression instead of wildcards expression to collect artifacts to upload.
   --threads      [Default: 3] Number of artifacts to upload in parallel.
   --dry-run      [Default: false] Set to true to disable communication with Artifactory.
```
##### Arguments
* The first argument is the path to the artifacts to be uploaded to Artifactory.
The path can include a single file or multiple artifacts, by using the * wildcard.
**Important:** If the path is provided as a regular expression (with the --regexp=true option) then
the first regular expression appearing as part of the argument must be enclosed in parenthesis.

* The second argument is the upload path in Artifactory.
The argument should have the following format: [repository name]:[repository path]
The path can include symbols in the form of {1}, {2}, ...
These symbols are replaced with the sections enclosed with parenthesis in the first argument.

##### Examples

This example uploads the *froggy.tgz* file to the root of the *my-local-repo* repository
```console
$ art upload froggy.tgz my-local-repo/ --url=http://domain/artifactory --user=admin --password=password
```


This example collects all the zip artifacts located under the build directory (including sub-directories)
   and uploads them to the *my-local-repo* repository, under the zipFiles folder, while keeping the artifacts original names.
   ```console
$ art upload build/*.zip libs-release-local/zipFiles/ --url=http://domain/artifactory --user=admin --password=password
   ```

#### The Download command

##### Function
Used to download artifacts from Artifactory.


##### Options
```console
   --url          [Mandatory] Artifactory URL
   --user         [Optional] Artifactory user
   --password     [Optional] Artifactory password
   --props        [Optional] List of properties in the form of key1=value1;key2=value2,... Only artifacts with these properties will be downloaded.
   --flat         [Default: false] Set to true if you do not wish to have the Artifactory repository path structure created locally for your downloaded artifacts
   --recursive    [Default: true] Set to false if you do not wish to include the download of artifacts inside sub-directories in Artifactory.
   --min-split    [Default: 5120] Minimum file size in KB to split into ranges. Set to -1 for no splits.
   --split-count  [Default: 3] Number of parts to split a file when downloading. Set to 0 for no splits.
   --threads      [Default: 3] Number of artifacts to download in parallel.
```

##### Arguments
The command expects one argument - the path of artifacts to be downloaded from Artifactory.
The argument should have the following format: [repository name]:[repository path]
The path can include a single artifact or multiple artifacts, by using the * wildcard.
The artifacts are downloaded and saved to the current directory, while saving their folder structure.

##### Examples

This example downloads the *cool-froggy.zip* artifact located at the root of the *my-local-repo* repository to current directory.
```console
$ art download my-local-repo/cool-froggy.zip --url=http://domain/artifactory --user=admin --password=password
```

This example downloads all artifacts located in the *my-local-repo* repository under the *all-my-frogs* folder to the *all-my-frog* directory located unde the current directory.
```console
$ art download my-local-repo/all-my-frogs/ --url=http://domain/artifactory --user=admin --password=password
```