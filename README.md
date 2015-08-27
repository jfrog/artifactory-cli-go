## Artifactory CLI

Artifactory CLI provides a command line interface for uploading and downloading artifacts to and from Artifactory.

### Installation

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
