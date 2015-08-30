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
**Important:** The first wildcard in the expression must be enclosed in parenthesis.

* The second argument is the upload path in Artifactory.
The argument should have the following format: [repository name]:[repository path]
The path can include symbols in the form of {1}, {2}, ...
These symbols are replaced with the sections enclosed with parenthesis in the first argument.

##### Examples

The following command uploads 'froggy.tgz' to the root of 'my-local-repo' repository

```console
$ art upload froggy.tgz my-local-repo --url=http://domain/artifactory --user=admin --password=password
```

Upload all the files from the current directory to another directory in 'my-local-repo'

```console
$ art upload * my-local-repo/uploaded/ --url=http://domain/artifactory --user=admin --password=password
```

The following command collects all the zip files located under the build directory (including sub-directories)
and uploads them to the libs-release-local repository, under the zipFiles folder, while keeping the files original names.

```console
$ artifactory-cli-go upload build/(*.zip) libs-release-local:zipFiles/{1} --url=http://localhost:8081/artifactory --user=admin --password=password
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
