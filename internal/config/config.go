package config

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	DocsPath string `yaml:"docs_path" env-default:"./docs"`
	SwaggerPath string `yaml:"swagger_path" env-default:"./docs/swagger.json"`
	HTTPServer `yaml:"http_server"`
	Database   `yaml:"database"`
}

type HTTPServer struct {
	Host        string        `yaml:"host" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"30"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"30"`
}
type Database struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"5432"`
	DBName   string `yaml:"dbname" env-default:"auth"`
	User     string `yaml:"user" env-default:"postgres"`
	Password string `yaml:"password" env-default:"postgres"`
}

func ReadConfig() (*Config, error) {
	configPath := "./config/config.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	return &cfg, nil
}

var (
	JWTSecret          = []byte("your-secret-key") // для access token
	RefreshJWTSecret   = []byte("refresh-secret")  // для refresh token (должен отличаться)
	AccessTokenExpiry  = 15 * time.Minute          // короткое время жизни access token
	RefreshTokenExpiry = 7 * 24 * time.Hour        // длительное время жизни refresh token
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	UserID uint `json:"user_id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}
