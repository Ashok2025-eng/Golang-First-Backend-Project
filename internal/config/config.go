package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"addr"` // Add YAML tag
}

// Config is the main application config
type Config struct {
	Env         string     `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string     `yaml:"storage_path" env-required:"true"`
	HTTPServer  HTTPServer `yaml:"http_server"` // struct embedding
}

// MustLoad loads config or stops the app if config is invalid
func MustLoad() *Config {
	var configPath string

	// Step 1: Check ENV variable
	configPath = os.Getenv("CONFIG_PATH") //check enc variable

	// Step 2: If ENV not set, check CLI flag
	if configPath == "" {
		flags := flag.String("config", "config/local.yaml", "path to the configuration file")
		flag.Parse()
		configPath = *flags

		if configPath == "" {
			log.Fatal("config path is not set")
		}
	}

	// Step 3: Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) { // make sure the file actually exists if not prints an erro
		log.Fatalf("config file does not exist: %s", configPath)
	}

	// Step 4: Read config into struct
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil { // force YAML
		log.Fatalf("cannot read config file: %v", err) // stops the program immediately
	}

	// ✅ Print which config file is being loaded (optional)
	log.Printf("Loaded config from: %s", configPath)

	return &cfg
}
