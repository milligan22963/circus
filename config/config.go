// Package config is used to define any configuration that isn't passed in from the command line
// or is default options that can be overridden
package config

import (
	"io/ioutil"
	"path/filepath"

	"github.com/milligan22963/afmlog"
	"github.com/milligan22963/circus/pkg/management"
	"gopkg.in/yaml.v2"
	"tinygo.org/x/bluetooth"
)

type WebServerSettings struct {
	Host     string `yaml:"ws_host"`
	FileRoot string `yaml:"ws_root"`
	Port     int    `yaml:"ws_port"`
}

type CircusConfiguration struct {
	WebServerSettings WebServerSettings    `yaml:"website"`
	LogSettings       afmlog.Configuration `yaml:"log"`
}

// DefaultConfigPath to our default config
const DefaultConfigPath = "settings.yaml"

// AppConfiguration is configuration
type AppConfiguration struct {
	CircusConfiguration CircusConfiguration
	AppActive           chan struct{}
	Skull               chan *management.Skull
	Adapter             *bluetooth.Adapter
}

func (configuration *CircusConfiguration) LoadConfiguration(filename string) error {
	fileContents, err := ioutil.ReadFile(filepath.Clean(filename))

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(fileContents, configuration)
	if err != nil {
		return err
	}
	return configuration.LogSettings.LoadConfiguration()
}

func (appConfig *AppConfiguration) GetLogger() *afmlog.Log {
	return appConfig.CircusConfiguration.LogSettings.UserLog
}

// NewSiteConfiguration creates an instance of the site configuration struct
func NewSiteConfiguration(configFile string) *AppConfiguration {
	appConfig := &AppConfiguration{
		CircusConfiguration: CircusConfiguration{},
		AppActive:           make(chan struct{}),
		Skull:               make(chan *management.Skull),
		Adapter:             bluetooth.DefaultAdapter,
	}

	err := appConfig.CircusConfiguration.LoadConfiguration(configFile)
	if err != nil {
		panic(err)
	}

	return appConfig
}
