package config

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/DmitryVesenniy/go-rest-framework/common"
	"github.com/joho/godotenv"
)

var (
	config         Config
	once           sync.Once
	DEFAULT_VALUES map[string]string = map[string]string{
		"HOST":             "localhost:3000",
		"DEBUG_HOST":       "localhost:8080",
		"REDIS_HOST":       "localhost:6379",
		"REDIS_PASSW":      "",
		"REDIS_DB":         "0",
		"LOG_FILE":         "./log.txt",
		"SUPER_USER_EMAIL": "admin@aisog.ru",
		"SUPER_USER_PASSW": "Adm1n",
		"MEDIA_PATH":       "media",
	}

	TokenExpires         time.Duration = time.Minute * 30
	RefreshTokenExpaires time.Duration = time.Hour * 24 * 365
)

type Config struct {
	BasePath    string
	Debug       bool   `env:"DEBUG"`
	SecretKey   string `env:"SECRET_KEY"`
	Host        string `env:"HOST"`
	DebugHost   string `env:"DEBUG_HOST"`
	StartWorker bool   `env:"START_WORKER"`

	PSQLServerHost string `env:"PSQLDB_HOST"`
	PSQLPort       string `env:"PSQLDB_PORT"`
	PSQLDatabase   string `env:"PSQLDB_NAME"`
	PSQLUserName   string `env:"PSQLDB_USER"`
	PSQLPasswd     string `env:"PSQLDB_PASSWORD"`

	RedisHost   string `env:"REDIS_HOST"`
	RedisPasswd string `env:"REDIS_PASSW"`
	RedisDB     int    `env:"REDIS_DB"`

	TokenExpires         time.Duration
	RefreshTokenExpaires time.Duration

	EmailHost         string `env:"EMAIL_HOST"`
	EmailPort         int    `env:"EMAIL_PORT"`
	EmailHostUser     string `env:"EMAIL_HOST_USER"`
	EmailHostPassword string `env:"EMAIL_HOST_PASSWORD"`

	SuperUserEmail string `env:"SUPER_USER_EMAIL"`
	SuperUserPassw string `env:"SUPER_USER_PASSW"`

	MediaPath string `env:"MEDIA_PATH"`

	LogFile string `env:"LOG_FILE"`
}

func (c *Config) GetMediaPath() string {
	return filepath.Join(c.BasePath, c.MediaPath)
}

func Get() *Config {
	once.Do(func() {
		dir, _ := os.Getwd()

		filesEnv := []string{
			filepath.Join(dir, "configs", ".env"),
			filepath.Join(dir, "configs", ".env.local"),
		}
		envFile := ""

		for _, _file := range filesEnv {
			if isExists, _ := common.Exists(_file); isExists {
				envFile = _file
			}
		}
		if err := godotenv.Load(envFile); err != nil {
			log.Fatal("[!] could not find environment variable file [./configs/.env]")
		}
		config = loadEnv()
		config.BasePath = dir
		config.TokenExpires = TokenExpires
		config.RefreshTokenExpaires = RefreshTokenExpaires
	})
	return &config
}

func loadEnv() Config {
	config = Config{}
	val := reflect.ValueOf(&config).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag.Get("env")

		valueEnv, _ := os.LookupEnv(tag)
		if valueEnv == "" {
			valueEnv = DEFAULT_VALUES[tag]
		}

		switch typeField.Type.Kind() {
		case reflect.Int:
			value, _ := strconv.Atoi(valueEnv)
			valueField.Set(reflect.ValueOf(value))
		case reflect.String:
			valueField.SetString(valueEnv)
		case reflect.Bool:
			value := false
			if valueEnv == "true" {
				value = true
			}
			valueField.SetBool(value)
		}
	}

	return config
}
