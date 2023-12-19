package utils

import (
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
)

func ReadYAMLFiles(folderPath string) (map[string]*Scenario, error) {
	result := make(map[string]*Scenario)

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".yaml") {
			continue
		}

		filePath := filepath.Join(folderPath, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		var yamlMap Scenario
		err = yaml.Unmarshal(data, &yamlMap)
		if err != nil {
			return nil, err
		}

		result[file.Name()] = &yamlMap
	}

	return result, nil
}
