package main

import (
    "github.com/yob/graval"
    "io"
    "io/ioutil"
    "log"
    "time"
    "os"
    "fmt"
    "encoding/json"
)

type RequestEndpoint struct {
    Method string
    Url string
}

type VirtualFile struct {
    Size int
    Endpoint RequestEndpoint
}

type VirtualFolder struct {
    Entities []VirtualEntity
}

type VirtualEntity struct {
    Name string
    File *VirtualFile
    Folder *VirtualFolder
}

type Configuration struct {
    Username string
    Password string
    Entities []VirtualEntity
    FilePathMap map[string]VirtualEntity
}

func LoadConfiguration(path string) (map[string]Configuration, error) {
    pathInfo, pathErr := os.Stat(path)

    if os.IsNotExist(pathErr) {
        return nil, fmt.Errorf("The configuration path does not exist: %v", path)
    }

    if !pathInfo.IsDir() {
        return nil, fmt.Errorf("The configuration path is not a directory: %v", path)
    }

    configurations := make(map[string]Configuration)

    dirContents, readDirErr := ioutil.ReadDir(path)

    if readDirErr != nil {
        return nil, readDirErr
    }

    for _, fileInfo := range dirContents {
        if fileInfo.IsDir() {
            continue
        }

        fileContents, readFileErr := ioutil.ReadFile(path + "/" + fileInfo.Name())
        if readFileErr != nil {
            return nil, readFileErr
        }

        var configuration Configuration

        var unmarshalError = json.Unmarshal(fileContents, &configuration)
        if unmarshalError != nil {
            return nil, unmarshalError
        }

        configuration.FilePathMap = make(map[string]VirtualEntity)
        CreateFilePathMap("", configuration.Entities, configuration.FilePathMap)

        configurations[configuration.Username] = configuration
    }

    return configurations, nil
}

func CreateFilePathMap(path string, entities []VirtualEntity, filePathMap map[string]VirtualEntity) {
    for _, entity := range entities {
        entityPath := path + "/" + entity.Name
        switch {
            case entity.File != nil:
                filePathMap[entityPath] = entity
            case entity.Folder != nil:
                filePathMap[entityPath] = entity
                CreateFilePathMap(entityPath, entity.Folder.Entities, filePathMap)
        }
    }
}

/**
 * Struct implementing the needed graval protocol to support "http -> ftp" mapping
 */
type HttpDriver struct {
    Configurations map[string]Configuration
    AuthenticatedUser interface{}
    UserConfiguration *Configuration
}

func (driver* HttpDriver) Authenticate(username string, password string) bool {
    log.Printf("Authenticated: Username: %s, Password: %s", username, password)
    userConfiguration, userExists := driver.Configurations[username]

    if !userExists {
        log.Printf("No configuration for user: %v", username)
        return false
    }

    driver.UserConfiguration = &userConfiguration

    if userConfiguration.Password != password {
        log.Printf("Invalid password for user given: %v", username)
        return false
    }

    driver.AuthenticatedUser = username

    return true
}

func (driver* HttpDriver) Bytes(filepath string) int {
    log.Printf("Size: %v", filepath)
    virtualEntity, virtualEntityExists := driver.UserConfiguration.FilePathMap[filepath]

    if !virtualEntityExists || virtualEntity.File == nil {
        return -1
    }

    return virtualEntity.File.Size
}

func (driver* HttpDriver) ModifiedTime(filepath string) (time.Time, error) {
    log.Printf("Last modified: %v", filepath)
    return time.Now(), nil
}

func (driver* HttpDriver) ChangeDir(path string) bool {
    log.Printf("Change Directory: %v", path)
    virtualEntity, virtualEntityExists := driver.UserConfiguration.FilePathMap[path]

    return virtualEntityExists && virtualEntity.File != nil
}

func FillDirContentsInfo(entities []VirtualEntity) []os.FileInfo {
    contents := []os.FileInfo{}

    for _, entity := range entities {
        switch {
            case entity.File != nil:
                contents = append(contents, graval.NewFileItem(entity.Name, entity.File.Size))
            case entity.Folder != nil:
                contents = append(contents, graval.NewDirItem(entity.Name))
        }
    }

    return contents
}

func (driver* HttpDriver) DirContents(path string) []os.FileInfo {
    log.Printf("Listing: %v", path)

    if path == "/" {
        return FillDirContentsInfo(driver.UserConfiguration.Entities)
    }

    virtualEntity, virtualEntityExists := driver.UserConfiguration.FilePathMap[path]
    if !virtualEntityExists {
        // If it does not exist return an empty list
        return []os.FileInfo{}
    }

    return FillDirContentsInfo(virtualEntity.Folder.Entities)
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

    return "", fmt.Errorf("Invalid File requested %v", filepath)

}

func (driver* HttpDriver) PutFile(filepath string, reader io.Reader) bool {
    return false
}

/**
 * Factory to create and give graval the Http Driver
 */
type HttpDriverFactory struct{
    Configurations map[string]Configuration
}

func (factory *HttpDriverFactory) NewDriver() (graval.FTPDriver, error) {
	return &HttpDriver{Configurations: factory.Configurations,
         AuthenticatedUser: nil,
         UserConfiguration: nil},
     nil
}

/**
 * Run the thing :)
 */
func main() {
	configurations, err := LoadConfiguration("./foo")
    if err != nil {
        log.Printf("Foo")
        log.Fatal(err)
    }

    log.Printf("Starting up server")

    httpFactory := &HttpDriverFactory{configurations}
    server := graval.NewFTPServer(&graval.FTPServerOpts{Factory: httpFactory})

	err = server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
