package config // Separate package for app configuration

import (
	"flag" // For CLI flags
	"log"  // For logging errors
	"os"   // For environment variables & file checks

	"github.com/ilyakaznacheev/cleanenv" // Library to read config into struct
)

// HTTPServer holds HTTP-related config
type HTTPServer struct {
	Addr string `yaml:"addr"` // Server address (e.g. ":8080")
}

// Config is the main application config
type Config struct {
	Env         string     `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string     `yaml:"storage_path" env-required:"true"`
	HTTPServer  HTTPServer `yaml:"http_server"`
}

// MustLoad loads config OR stops the app if config is invalid
func MustLoad() *Config {

	// Step 1: Try reading config path from ENV
	configPath := os.Getenv("CONFIG_PATH")

	// Step 2: If ENV not set, try CLI flag
	if configPath == "" {
		flagPath := flag.String("config", "", "path to the configuration file")
		flag.Parse()
		configPath = *flagPath

		// If still empty → crash
		if configPath == "" {
			log.Fatal("CONFIG_PATH is not set")
		}
	}

	// Step 3: Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	// Step 4: Read config into struct
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config file: %v", err)
	}

	// Step 5: Return loaded config
	return &cfg
}
