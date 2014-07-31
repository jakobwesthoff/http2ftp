package http2ftp

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/yob/graval"
)

type httpDriver struct {
	Configurations    map[string]Configuration
	AuthenticatedUser interface{}
	UserConfiguration *Configuration
}

func (driver *httpDriver) Authenticate(username string, password string) bool {
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

func (driver *httpDriver) Bytes(filepath string) int {
	log.Printf("Size: %v", filepath)
	virtualEntity, virtualEntityExists := driver.UserConfiguration.FilePathMap[filepath]

	if !virtualEntityExists || virtualEntity.File == nil {
		return -1
	}

	return virtualEntity.File.Size
}

func (driver *httpDriver) ModifiedTime(filepath string) (time.Time, error) {
	log.Printf("Last modified: %v", filepath)
	return time.Now(), nil
}

func (driver *httpDriver) ChangeDir(path string) bool {
	log.Printf("Change Directory: %v", path)

	if path == "/" {
		return true
	}

	virtualEntity, virtualEntityExists := driver.UserConfiguration.FilePathMap[path]
	if (!virtualEntityExists || virtualEntity.Folder == nil) {
		return false
	}

	if (virtualEntity.Folder.Endpoint != nil) {
		log.Printf("Fetching dynamic folder: %s", path)
		// The folder contents is dynamic and needs to be fetched
		body, fetchError := fetchHTTPResourceBody(
			virtualEntity.Folder.Endpoint.Method,
			virtualEntity.Folder.Endpoint.URL,
		)

		if (fetchError != nil) {
			log.Printf("Failed to fetch dynamic folder contents: %v", fetchError)
			return false
		}

		integrationError := IntegrateNewVirtualFolderEntity(
			driver.UserConfiguration,
			path,
			body,
		)

		if (integrationError != nil) {
			log.Printf("Failed to integrate dynamic folder contents: %v", integrationError)
			return false
		}
	}

	return true
}

func fillDirContentsInfo(entities []virtualEntity) []os.FileInfo {
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

func (driver *httpDriver) DirContents(path string) []os.FileInfo {
	log.Printf("Listing: %v", path)

	if path == "/" {
		return fillDirContentsInfo(driver.UserConfiguration.Entities)
	}

	virtualEntity, virtualEntityExists := driver.UserConfiguration.FilePathMap[path]
	if !virtualEntityExists {
		// If it does not exist return an empty list
		return []os.FileInfo{}
	}

	return fillDirContentsInfo(virtualEntity.Folder.Entities)
}

func (driver *httpDriver) DeleteDir(path string) bool {
	return false
}

func (driver *httpDriver) DeleteFile(filepath string) bool {
	return false
}

func (driver *httpDriver) Rename(soure, target string) bool {
	return false
}

func (driver *httpDriver) MakeDir(path string) bool {
	return false
}

func (driver *httpDriver) GetFile(filepath string) (string, error) {
	log.Printf("Transmit file request: %v", filepath)

	virtualEntity, virtualEntityExists := driver.UserConfiguration.FilePathMap[filepath]
	if !virtualEntityExists || virtualEntity.File == nil {
		return "", fmt.Errorf("Invalid File requested %v", filepath)
	}

	endpoint := virtualEntity.File.Endpoint

	log.Printf("Requesting HTTP: %s %s", endpoint.Method, endpoint.URL)
	body, requestError := fetchHTTPResourceBody(endpoint.Method, endpoint.URL)

	if requestError != nil {
		return "", requestError
	}

	log.Printf("Transmit file: %v", filepath)
	return body, nil
}

func (driver *httpDriver) PutFile(filepath string, reader io.Reader) bool {
	return false
}
