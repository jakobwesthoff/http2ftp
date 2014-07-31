package http2ftp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type requestEndpoint struct {
	Method string
	URL    string
}

type virtualFile struct {
	Size     int
	Endpoint requestEndpoint
}

type virtualFolder struct {
	Entities []virtualEntity
}

type virtualEntity struct {
	Name   string
	File   *virtualFile
	Folder *virtualFolder
}

// Configuration representing one of the UserConfigurations
type Configuration struct {
	Username    string
	Password    string
	Entities    []virtualEntity
	FilePathMap map[string]virtualEntity
}

// LoadConfiguration parses all configurations from wihtin the given directory
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
		fullPath := path + "/" + fileInfo.Name()

		if fileInfo.IsDir() || strings.ToLower(filepath.Ext(fullPath)) != ".json" {
			continue
		}

		fileContents, readFileErr := ioutil.ReadFile(fullPath)
		if readFileErr != nil {
			return nil, readFileErr
		}

		var configuration Configuration

		var unmarshalError = json.Unmarshal(fileContents, &configuration)
		if unmarshalError != nil {
			return nil, unmarshalError
		}

		configuration.Username = strings.TrimSuffix(filepath.Base(fullPath), filepath.Ext(fullPath))

		configuration.FilePathMap = make(map[string]virtualEntity)
		createFilePathMap("", configuration.Entities, configuration.FilePathMap)

		configurations[configuration.Username] = configuration
	}

	return configurations, nil
}

func createFilePathMap(path string, entities []virtualEntity, filePathMap map[string]virtualEntity) {
	for _, entity := range entities {
		entityPath := path + "/" + entity.Name
		switch {
		case entity.File != nil:
			filePathMap[entityPath] = entity
		case entity.Folder != nil:
			filePathMap[entityPath] = entity
			createFilePathMap(entityPath, entity.Folder.Entities, filePathMap)
		}
	}
}
