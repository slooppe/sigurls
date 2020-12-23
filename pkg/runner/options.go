package runner

import (
	"os"
	"path"

	"github.com/drsigned/sigurls/pkg/sources"
	"gopkg.in/yaml.v3"
)

// Configuration contains the fields stored in the configuration file
type Configuration struct {
	// Version indicates the version of subfinder installed.
	Version string `yaml:"version"`
	// Sources contains a list of sources to use for enumeration
	Sources []string `yaml:"sources"`
}

// Options is a
type Options struct {
	Domain      string
	Exclude     string
	IncludeSubs bool
	Use         string

	YAMLConfig Configuration
}

// ParseOptions is a
func ParseOptions(options *Options) (*Options, error) {
	directory, err := os.UserHomeDir()
	if err != nil {
		return options, err
	}

	version := "1.2.0"
	configPath := directory + "/.config/sigurls/conf.yaml"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configuration := Configuration{
			Version: version,
			Sources: sources.All,
		}

		directory, _ := path.Split(configPath)

		err := makeDirectory(directory)
		if err != nil {
			return options, err
		}

		err = configuration.MarshalWrite(configPath)
		if err != nil {
			return options, err
		}

		options.YAMLConfig = configuration
	} else {
		configuration, err := UnmarshalRead(configPath)
		if err != nil {
			return options, err
		}

		if configuration.Version != version {
			configuration.Sources = sources.All
			configuration.Version = version

			err := configuration.MarshalWrite(configPath)
			if err != nil {
				return options, err
			}
		}

		options.YAMLConfig = configuration
	}

	return options, nil
}

func makeDirectory(directory string) error {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if directory != "" {
			err = os.MkdirAll(directory, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// MarshalWrite writes the marshaled yaml config to disk
func (config *Configuration) MarshalWrite(file string) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	// Indent the spaces too
	enc := yaml.NewEncoder(f)
	enc.SetIndent(4)
	err = enc.Encode(&config)
	f.Close()
	return err
}

// UnmarshalRead reads the unmarshalled config yaml file from disk
func UnmarshalRead(file string) (Configuration, error) {
	config := Configuration{}

	f, err := os.Open(file)
	if err != nil {
		return config, err
	}

	err = yaml.NewDecoder(f).Decode(&config)

	f.Close()

	return config, err
}
