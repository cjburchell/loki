package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/satori/go.uuid"
)

// Endpoint configuration
type Endpoint struct {
	ID           string `json:"id"`
	Description  string `json:"description"`
	Path         string `json:"path"`
	Method       string `json:"method"`
	ResponseBody string `json:"response_body"`
	ContentType  string `json:"content_type"`
	Response     int    `json:"response"`
}

// GetEndpoints configuration
func GetEndpoints() ([]Endpoint, error) {
	results, err := load()
	if err != nil {
		return nil, err
	}

	endpoints := make([]Endpoint, len(results))
	index := 0
	for _, value := range results {
		endpoints[index] = value
		index++
	}

	return endpoints, nil
}

// GetEndpoint with given ID
func GetEndpoint(id string) (*Endpoint, error) {
	results, err := load()
	if err != nil {
		return nil, err
	}

	if item, ok := results[id]; ok {
		return &item, nil
	}

	return nil, nil
}

var lock = &sync.Mutex{}

// AddEndpoint in configuration
func AddEndpoint(endpoint Endpoint) (string, error) {
	lock.Lock()
	defer lock.Unlock()
	endpoint.ID = uuid.Must(uuid.NewV4()).String()
	endpoints, err := load()
	if err != nil {
		return "", err
	}

	endpoints[endpoint.ID] = endpoint
	return endpoint.ID, save(endpoints)
}

// UpdateEndpoint in configuration
func UpdateEndpoint(endpoint Endpoint) error {
	lock.Lock()
	defer lock.Unlock()
	endpoints, err := load()
	if err != nil {
		return err
	}

	if _, ok := endpoints[endpoint.ID]; ok {
		endpoints[endpoint.ID] = endpoint
		return save(endpoints)
	}

	return fmt.Errorf("unable to find logger with Id %s", endpoint.ID)
}

// DeleteEndpoint in configuration
func DeleteEndpoint(id string) error {
	lock.Lock()
	defer lock.Unlock()
	endpoints, err := load()
	if err != nil {
		return err
	}

	if _, ok := endpoints[id]; ok {
		delete(endpoints, id)
		return save(endpoints)
	}

	return fmt.Errorf("unable to find logger with Id %s", id)
}

// Setup the configuration
func Setup(file string) error {
	configFileName = file
	return nil
}

var configFileName string

func load() (map[string]Endpoint, error) {
	loggers := make(map[string]Endpoint)
	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		return loggers, nil
	}

	fileData, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return loggers, err
	}

	err = json.Unmarshal(fileData, &loggers)
	return loggers, err
}

func save(config map[string]Endpoint) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configFileName, configJSON, 0644)
}
