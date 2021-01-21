package config

import (
	"encoding/json"
	"github.com/cjburchell/loki/models"
	"io/ioutil"
	"os"

	log "github.com/cjburchell/uatu-go"

	"github.com/pkg/errors"
)

// GetEndpoints configuration
func GetEndpoints(log log.ILog) ([]models.Endpoint, error) {
	results, err := load(log)
	if err != nil {
		return nil, err
	}

	endpoints := make([]models.Endpoint, len(results))
	index := 0
	for name, value := range results {
		value.Name = name
		endpoints[index] = value
		index++
	}

	return endpoints, nil
}

// Setup the configuration
func Setup(file string) error {
	configFileName = file
	return nil
}

var configFileName string

func load(log log.ILog) (map[string]models.Endpoint, error) {
	loggers := make(map[string]models.Endpoint)
	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		log.Warnf("Config file %s not found", configFileName)
		return loggers, nil
	}

	log.Printf("loading config file %s", configFileName)
	fileData, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return loggers, errors.WithStack(err)
	}

	err = json.Unmarshal(fileData, &loggers)
	return loggers, errors.WithStack(err)
}
