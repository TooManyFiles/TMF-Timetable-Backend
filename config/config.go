package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	goUntisAPIstructs "github.com/Mr-Comand/goUntisAPI/structs"
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
		Untis        struct {
			FixedIV string // Must be 16 bytes for AES
			Salt    string
		}
	}
	DataCollectors struct {
		TFfoodplanAPIURL string
		UntisApiConfig   goUntisAPIstructs.ApiConfig
	}

	DatabaseConfig DatabaseConfig
	CanSignUp      bool
}

var Config ConfigStruct = ConfigStruct{
	DatabaseConfig: DatabaseConfig{
		Connection: "postgres://user:password@localhost:5432/db?sslmode=disable",
	},
	Crypto: struct {
		JwtSecretKey string
		Untis        struct {
			FixedIV string // Must be 16 bytes for AES
			Salt    string
		}
	}{
		JwtSecretKey: "secret",
		Untis: struct {
			FixedIV string // Must be 16 bytes for AES
			Salt    string
		}{
			FixedIV: "example_iv123456",
			Salt:    "example_salt",
		},
	},
	DataCollectors: struct {
		TFfoodplanAPIURL string
		UntisApiConfig   goUntisAPIstructs.ApiConfig
	}{
		TFfoodplanAPIURL: "http://www.treffpunkt-fanny.de/images/stories/dokumente/Essensplaene/api/TFfoodplanAPI.php",
		UntisApiConfig: goUntisAPIstructs.ApiConfig{
			Server:    "school.server.domain",
			User:      "username",
			Password:  "password",
			Useragent: "client",
			School:    "school",
		},
	},
	CanSignUp: true,
}

// Function to create a default config file if it doesn't exist and no env vars are set
func createDefaultConfigFile(configPath string) error {
	// Get the directory from the config path
	dir := filepath.Dir(configPath)

	// Create the necessary directories
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}
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
func areEnvVarsUnset(config interface{}, prefix string) bool {
	val := reflect.ValueOf(config)

	// If the value is not a struct, return true (base case)
	if val.Kind() != reflect.Struct {
		return true
	}
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldValue := val.Field(i)

		// Determine the environment variable name
		var envTag string
		if tag := field.Tag.Get("env"); tag != "" {
			envTag = tag
		} else {
			envTag = strings.ToUpper(field.Name)
		}

		// If there's a prefix, join it with the field name
		if prefix != "" {
			envTag = prefix + "_" + envTag
		}

		// Check if the environment variable is set
		if os.Getenv(envTag) != "" {
			// If any environment variable is set, return false
			return false
		}

		// If the field is a struct, recursively check its fields
		if fieldValue.Kind() == reflect.Struct {
			if !areEnvVarsUnset(fieldValue.Interface(), envTag) {
				return false
			}
		}
	}
	return true
}

// Recursive function to bind environment variables
func getEnvVars(config interface{}, prefix string) {
	val := reflect.ValueOf(config)

	// Ensure we are dealing with a struct
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldValue := val.Field(i)

		// Get the 'env' tag or fallback to the uppercase field name
		Tag := field.Tag.Get("env")
		if Tag == "" {
			Tag = strings.ToUpper(field.Name)
		}

		// If there's a prefix (for nested structs), prepend it to the env variable name
		if prefix != "" {
			Tag = prefix + "." + Tag
		}
		// Bind the environment variable
		if os.Getenv(strings.ReplaceAll(Tag, ".", "_")) != "" {
			v.Set(strings.ToLower(Tag), os.Getenv(strings.ReplaceAll(Tag, ".", "_")))
		}
		// If the field is a struct, recursively bind its fields
		if fieldValue.Kind() == reflect.Struct {
			getEnvVars(fieldValue.Interface(), Tag)
		}
	}
}

var v *viper.Viper

func LoadConfig() error {
	v = viper.New()

	// Default config file path
	configPath := "./config.yml"

	// If a command line argument is provided, use it as the config path
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	} else if os.Getenv("CONFIG_FILE") != "" {
		configPath = os.Getenv("CONFIG_FILE")
	}

	// Check if all environment variables are unset, and if config file doesn't exist, create default
	if _, err := os.Stat(configPath); os.IsNotExist(err) && areEnvVarsUnset(Config, "") {
		fmt.Println("No environment variables set and no config file found. Creating default config file...")
		if err := createDefaultConfigFile(configPath); err != nil {
			log.Fatalf("Error creating default config file: %v", err)
		}
	}

	v.SetEnvPrefix("")
	v.SetConfigFile(configPath)

	// Read the config file if it exists
	if _, err := os.Stat(configPath); err == nil {
		if err := v.ReadInConfig(); err != nil {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Get environment variables automatically using struct tags
	getEnvVars(Config, "")

	// Parse config into the provided struct
	if err := v.Unmarshal(&Config); err != nil {
		return fmt.Errorf("unable to decode into struct: %w", err)
	}
	return nil
}
