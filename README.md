[![Build Status](https://travis-ci.org/JFrogDev/artifactory-cli-go.svg)](https://travis-ci.org/JFrogDev/artifactory-cli-go)

## Artifactory CLI

Artifactory CLI provides a command line interface for uploading and downloading artifacts to and from Artifactory.
It supports Artifactory version 3.5.3 and above.

### Getting Started

#### Downloading the executables from Bintray

If you do not wish to install/use go, you can [download executables](https://bintray.com/jfrog/artifactory-cli-go)
for Linux, Mac and Windows from JFrog Bintray. Select the architecture you want, download the art executable, and place it in your path.

#### Building the command line executable

If you prefer, you may instead build the client in go.

##### Setup GO on your machine

* Make sure you have a working Go environment. [See the install instructions](http://golang.org/doc/install).
* Navigate to the directory where you want to create the *artifactory-cli-go* project.
* Set the value of the GOPATH environment variable to the full path of this directory.

##### Download Artifactory CLI from GitHub

Run the following command to create the *artifactory-cli-go* project:
```console
$ go get github.com/JFrogDev/artifactory-cli-go/...
```

Navigate to the following directory
```console
$ cd $GOPATH/bin
```

You'll find the art (or art.exe on Windows) executable there.

#### Usage

You can copy the *art* executable to any location on your file-system as long as you add it to your *PATH* environment variable,
so that you can access it from any path.

##### Command syntax

```console
$ art command-name options arguments
```

The sections below specify the available commands, their respective options and additional arguments that may be needed.
*art* should be followed by a command name (for example, upload), a list of options (for example, --url=http://...)
and the list of arguments for the command.

##### The *upload* command

###### Function
Used to upload artifacts to Artifactory.

###### Command options
```console
   --url          [Mandatory] Artifactory URL.
   --user         [Optional] Artifactory user.
   --password     [Optional] Artifactory password.
   --props        [Optional] List of properties in the form of "key1=value1;key2=value2,..." to be attached to the uploaded artifacts.
   --deb          [Optional] Used for Debian packages in the form of distribution/component/architecture.
   --flat         [Default: true] If not set to true, and the upload path ends with a slash, artifacts are uploaded according to their file system hierarchy.
   --recursive    [Default: true] Set to false if you do not wish to collect artifacts in sub-folders to be uploaded to Artifactory.
   --regexp       [Default: false] Set to true to use a regular expression instead of wildcards expression to collect artifacts to upload.
   --threads      [Default: 3] Number of artifacts to upload in parallel.
   --dry-run      [Default: false] Set to true to disable communication with Artifactory.
```
###### Arguments
* The first argument is the local file-system path to the artifacts to be uploaded to Artifactory.
The path can include a single file or multiple artifacts, by using the * wildcard.
**Important:** If the path is provided as a regular expression (with the --regexp=true option) then
the first regular expression appearing as part of the argument must be enclosed in parenthesis.

* The second argument is the upload path in Artifactory.
The argument should have the following format: [repository name]/[repository path]
The path can include symbols in the form of {1}, {2}, ...
These symbols are replaced with the sections enclosed with parenthesis in the first argument.

###### Examples

This example uploads the *froggy.tgz* file to the root of the *my-local-repo* repository
```console
$ art upload "froggy.tgz" "my-local-repo/" --url=http://domain/artifactory --user=admin --password=password
```

This example collects all the zip artifacts located under the build directory (including sub-directories).
and uploads them to the *my-local-repo* repository, under the zipFiles folder, while keeping the artifacts original names.
```console
$ art upload build/*.zip libs-release-local/zipFiles/ --url=http://domain/artifactory --user=admin --password=password
```
And on Windows:
```console
$ art upload "build\\*.zip" "libs-release-local/zipFiles/" --url=http://domain/artifactory --user=admin --password=password
```

##### The *download* command

###### Function
Used to download artifacts from Artifactory.

###### Command options
```console
   --url          [Mandatory] Artifactory URL
   --user         [Optional] Artifactory user
   --password     [Optional] Artifactory password
   --props        [Optional] List of properties in the form of "key1=value1;key2=value2,..." Only artifacts with these properties will be downloaded.
   --flat         [Default: false] Set to true if you do not wish to have the Artifactory repository path structure created locally for your downloaded artifacts
   --recursive    [Default: true] Set to false if you do not wish to include the download of artifacts inside sub-directories in Artifactory.
   --min-split    [Default: 5120] Minimum file size in KB to split into ranges. Set to -1 for no splits.
   --split-count  [Default: 3] Number of parts to split a file when downloading. Set to 0 for no splits.
   --threads      [Default: 3] Number of artifacts to download in parallel.
```

###### Arguments
The command expects one argument - the path of artifacts to be downloaded from Artifactory.
The argument should have the following format: [repository name]/[repository path]
The path can include a single artifact or multiple artifacts, by using the * wildcard.
The artifacts are downloaded and saved to the current directory, while saving their folder structure.

###### Examples

This example downloads the *cool-froggy.zip* artifact located at the root of the *my-local-repo* repository to current directory.
```console
$ art download "my-local-repo/cool-froggy.zip" --url=http://domain/artifactory --user=admin --password=password
```

This example downloads all artifacts located in the *my-local-repo* repository under the *all-my-frogs* folder to the *all-my-frog* directory located unde the current directory.
```console
$ art download "my-local-repo/all-my-frogs/" --url=http://domain/artifactory --user=admin --password=password
```

##### The *config* command

###### Function
Used to configure the Artifactory URL, user and passwords, so that you don't have to send them as options
for the *upload* and *download* commands.
The configuration is saved at ~/.jfrog/art-cli.conf

###### Command options
```console
   --interactive  [Default: true] Set to false if you do not wish the config command to be interactive. If true, the --url option becomes optional.
   --enc-password [Default: true] If set to false then the configured password will not be encrypted using Artifatory's encryption API.
   --url          [Optional] Artifactory URL.
   --user         [Optional] Artifactory user.
   --password     [Optional] Artifactory password.
```

###### Arguments
* If no arguments are sent, the command will configure the Artifactory URL, user and password sent through the command options
or through the command's interactive prompt.
* The *show* argument will make the command show the stored configuration.
* The *clear* argument will make the command clear the stored configuration.

###### Important Note

if your Artifactory server has [encrypted password set to required](https://www.jfrog.com/confluence/display/RTF/Configuring+Security#ConfiguringSecurity-PasswordEncryptionPolicy) you should use your API Key as your password.

###### Examples

Configure the Artifactory details through an interactive propmp.
```console
$ art config
```

Configure the Artifactory details through the command options.

```console
$ art config --url=http://domain/artifactory --user=admin --password=password
```

Show the configured Artifactory details.
```console
$ art config show
```

Clear the configured Artifactory details.
```console
$ art config clear
```
