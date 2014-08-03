# http2ftp - An experimental protocol bridge
## Idea
http2ftp is an **experimental** bridge between the FTP and the HTTP protocol. It provides the necessary means to easily create an FTP frontend for a webapplication, to allow easy import/export of data for customers, which need FTP interfaces (ERP-Systems anyone? ;).

## State
The bridge is considered highly experimental. Even though I am loosely planning on using it for a project myself in the future, I created it as means to learn and play around with [Go](http://golang.org). The first version of this endeavor has been created with 1 1/2 days of Go knowledge. Therefore be warned it might not be production ready ;).

## Installation
After installing Go simply issue the following command to install the project into your `$GOPATH/bin` folder:

```
$ go install github.com/jakobwesthoff/http2ftp/cmd/http2ftp
```

The `http2ftp` executable should be ready for usage after that. Execute it with the `-h` flag to retrieve more information about configuration options.

## Configuration
The bridge is configured using a bunch of JSON files. Each user account has its own `.json` file within the configuration folder. The filename reflects the username. To allow a user called `jakob` to login to the FTP bridge a file named `jakob.json` needs to be created inside the configuration directory. The path of the directory itself may be specified using the `--config` flag.

The configuration file hosts the `password` for the user account, as well the directory and file structure displayed once the user is logged in. Each file or directory may specify so called `endpoint`s, which reflect their mapping to a HTML resource. Endpoints may be `read` and `write` endpoints allowing to specify what should happen once a user reads the resource or writes to it. Directories might either contain other directory and file definitions directly within the configuration or provide a read endpoint, which is queried once the directory is listed. This endpoint has to return the exact same JSON structure, which is used to configure directory `entities` within the static configuration.

An example configuration, which demonstrates the different structural possibilities can be found under `config/jakob.json` within the repository.