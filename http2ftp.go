package main

import (
    "github.com/yob/graval"
    "io"
    "log"
    "time"
    "os"
    "fmt"
)

/**
 * Struct implementing the needed graval protocol to support "http -> ftp" mapping
 */
type HttpDriver struct {}

func (driver* HttpDriver) Authenticate(username string, password string) bool {
    log.Printf("Authenticated: Username: %s, Password: %s", username, password)
    // For now everyone is allowed to authenticate
    return true
}

func (driver* HttpDriver) Bytes(filepath string) int {
    log.Printf("Size: %v", filepath)
    switch filepath {
        case "/foo.txt":
            return 100
        case "/bar.txt":
            return 200
        default:
            return -1
    }
}

func (driver* HttpDriver) ModifiedTime(filepath string) (time.Time, error) {
    log.Printf("Last modified: %v", filepath)
    return time.Now(), nil
}

func (driver* HttpDriver) ChangeDir(path string) bool {
    log.Printf("Change Directory: %v", path)
    return path == "/"
}

func (driver* HttpDriver) DirContents(path string) []os.FileInfo {
    log.Printf("Listing: %v", path)
    contents := []os.FileInfo{}

    if path != "/" {
        return contents
    }

    contents = append(contents, graval.NewFileItem("foo.txt", 100))
    contents = append(contents, graval.NewFileItem("bar.txt", 200))
    return contents
}

func (driver* HttpDriver) DeleteDir(path string) bool {
    return false
}

func (driver* HttpDriver) DeleteFile(filepath string) bool {
    return false
}

func (driver* HttpDriver) Rename(soure, target string) bool {
    return false
}

func (driver* HttpDriver) MakeDir(path string) bool {
    return false
}

func (driver* HttpDriver) GetFile(filepath string) (string, error) {
    log.Printf("Transmitting: %v", filepath)

    switch filepath {
        case "/foo.txt":
            return "FOO und so!", nil
        case "/bar.txt":
            return "BAR und so!!!1elf!", nil
        default:
            return "", fmt.Errorf("Invalid File requested %v", filepath)
    }
}

func (driver* HttpDriver) PutFile(filepath string, reader io.Reader) bool {
    return false
}

/**
 * Factory to create and give graval the Http Driver
 */
type HttpDriverFactory struct{}

func (factory *HttpDriverFactory) NewDriver() (graval.FTPDriver, error) {
	return &HttpDriver{}, nil
}

/**
 * Run the thing :)
 */
func main() {
	log.Printf("Starting up server")

    httpFactory := &HttpDriverFactory{}

    server := graval.NewFTPServer(&graval.FTPServerOpts{Factory: httpFactory})

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
