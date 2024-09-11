package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/invopop/yaml"
	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Connection string
	other      string
}
type ConfigStruct struct {
	Crypto struct {
		JwtSecretKey string
	}
	DataCollectors struct {
		TFfoodplanAPIURL string
	}

	DatabaseConfig DatabaseConfig
}

var Config ConfigStruct = ConfigStruct{
	DatabaseConfig: DatabaseConfig{
		Connection: "postgres://user:password@localhost:5432/db?sslmode=disable",
	},
	Crypto: struct{ JwtSecretKey string }{
		JwtSecretKey: "secret",
	},
	DataCollectors: struct{ TFfoodplanAPIURL string }{
		TFfoodplanAPIURL: "http://www.treffpunkt-fanny.de/images/stories/dokumente/Essensplaene/api/TFfoodplanAPI.php",
	},
}

// Function to create a default config file if it doesn't exist and no env vars are set
func createDefaultConfigFile(configPath string) error {
	configData, err := yaml.Marshal(&Config)
	if err != nil {
		return fmt.Errorf("error marshaling default config: %w", err)
	}

	err = os.WriteFile(configPath, configData, 0644)
	if err != nil {
		return fmt.Errorf("error writing default config file: %w", err)
	}

	fmt.Println("Default config file created at", configPath)
	return nil
}

// Function to check if all environment variables for the config struct are unset
func areEnvVarsUnset(config interface{}) bool {
	val := reflect.ValueOf(config)
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		envTag := field.Tag.Get("env")
		if envTag == "" {
			// Fallback to field name if no 'env' tag is provided
			envTag = strings.ToUpper(field.Name)
		}
		if os.Getenv(envTag) != "" {
			// If any environment variable is set, return false
			return false
		}
	}
	return true
}
func LoadConfig() error {
	// Default config file path
	configPath := "./config.yml"

	// If a command line argument is provided, use it as the config path
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	} else if os.Getenv("CONFIG_FILE") != "" {
		configPath = os.Getenv("CONFIG_FILE")
	}

	// Check if all environment variables are unset, and if config file doesn't exist, create default
	if _, err := os.Stat(configPath); os.IsNotExist(err) && areEnvVarsUnset(Config) {
		fmt.Println("No environment variables set and no config file found. Creating default config file...")
		if err := createDefaultConfigFile(configPath); err != nil {
			log.Fatalf("Error creating default config file: %v", err)
		}
	}

	viper.SetConfigFile(configPath)

	// Bind environment variables automatically using struct tags
	val := reflect.ValueOf(Config)
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		envTag := field.Tag.Get("env")
		if envTag == "" {
			// Fallback to field name if no 'env' tag is provided
			envTag = strings.ToUpper(field.Name)
		}
		viper.BindEnv(field.Name, envTag) // Bind field to corresponding env variable
	}

	// Read the config file if it exists
	if _, err := os.Stat(configPath); err == nil {
		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Parse config into the provided struct
	if err := viper.Unmarshal(&Config); err != nil {
		return fmt.Errorf("unable to decode into struct: %w", err)
	}

	return nil
}
