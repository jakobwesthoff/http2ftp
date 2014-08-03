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
	Read *readEndpoint
	Write *writeEndpoint
}

type readEndpoint struct {
	Method string
	URL    string
}

type writeEndpoint struct {
	URL string
	Parameter string
}

type virtualFile struct {
	Size     int
	Endpoint requestEndpoint
}

type virtualFolder struct {
	Entities []virtualEntity
	Endpoint requestEndpoint
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
		updateFilePathMap("", configuration.Entities, configuration.FilePathMap)

		configurations[configuration.Username] = configuration
	}

	return configurations, nil
}

// UnmarshalFolderConfiguration reads the JSON entity, which defines a VirtualFolder
func unmarshalVirtualFolderEntities(virtualFolderEntities string) ([]virtualEntity, error) {
    var virtualEntities []virtualEntity
    var unmarshalError = json.Unmarshal([]byte(virtualFolderEntities), &virtualEntities)
    return virtualEntities, unmarshalError
}

// IntegrateNewVirtualFolderEntity combines a given VirtualEntity (Folder) in JSON format with an
// existing Configuration.
//
// All subpaths of the given virtualPath are removed from the configuration and replaced with
// the new node
func IntegrateNewVirtualFolderEntity(configuration *Configuration, virtualPath string, virtualFolderEntities string) error {
    oldVirtualEntity, oldVirtualEntityExists := configuration.FilePathMap[virtualPath]

	if (!oldVirtualEntityExists) {
		return fmt.Errorf(
			"Request integration of dynamic virtual folder, which has no matching entity: %s",
			virtualPath,
		)
	}

	entities, unmarshalError := unmarshalVirtualFolderEntities(virtualFolderEntities)

	if (unmarshalError != nil) {
		return unmarshalError
	}

	removeAllEntitiesBelowPath(configuration, virtualPath)
	updateFilePathMap(virtualPath, entities, configuration.FilePathMap)
	oldVirtualEntity.Folder.Entities = entities

	return nil
}

func removeAllEntitiesBelowPath(configuration *Configuration, virtualPath string) {
	for path, _ := range configuration.FilePathMap {
		if (strings.Index(path, virtualPath + "/") != 0) {
			continue;
		}
		delete(configuration.FilePathMap, path)
	}
}

func updateFilePathMap(path string, entities []virtualEntity, filePathMap map[string]virtualEntity) {
	for _, entity := range entities {
		entityPath := path + "/" + entity.Name
		switch {
		case entity.File != nil:
			filePathMap[entityPath] = entity
		case entity.Folder != nil:
			filePathMap[entityPath] = entity
			updateFilePathMap(entityPath, entity.Folder.Entities, filePathMap)
		}
	}
}
