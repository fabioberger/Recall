package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v1"
)

var DatabaseYAMLPath = os.Getenv("GOPATH") + "/src/github.com/fabioberger/recall/db/dbconf.yml"

var Env string = os.Getenv("GO_ENV")

var Host string
var Port string
var Database dbConfig

type config struct {
	Host         string
	Port         string
	AllowOrigins []string
	Database     dbConfig
}

type dbConfig struct {
	Driver string
	Open   string
}

var Prod config = config{
	Host: "", // TODO: Set this to our api domain
	Port: "8080",
}

var Dev config = config{
	Host: "localhost",
	Port: "3000",
}

var Test config = config{
	Host: "localhost",
	Port: "4000",
}

func Init() {
	if os.Getenv("RECALL_PSQL_HOST") == "" {
		msg := `The environment variable RECALL_PSQL_HOST is not set.
		If you are on your own laptop, you should probably set it to localhost.
		If you are on an ubuntu aws server, set it to /var/run/postgresql.`
		panic(msg)
	}
	if Env == "production" {
		Use(Prod)
		ParseDatabaseYAML("production")
	} else if Env == "development" || Env == "" {
		Use(Dev)
		ParseDatabaseYAML("development")
	} else if Env == "test" {
		Use(Test)
		ParseDatabaseYAML("test")
	} else {
		panic("Unkown environment. Don't know what configuration to use!")
	}
}

func Use(c config) {
	Host = c.Host
	Port = c.Port
}

func ParseDatabaseYAML(env string) {
	// Read all data from the file and unmarshall it
	var data map[string]struct {
		Driver string
		Open   string
	}
	content, err := ioutil.ReadFile(DatabaseYAMLPath)
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(content, &data); err != nil {
		panic(err)
	}
	envData := data[env]

	// parse the env variable and set it properly
	if strings.Contains(envData.Open, "$POSTGRES_PORT_5432_TCP_ADDR") {
		envData.Open = strings.Replace(envData.Open, "$POSTGRES_PORT_5432_TCP_ADDR", os.Getenv("POSTGRES_PORT_5432_TCP_ADDR"), -1)
	}
	if strings.Contains(envData.Open, "$POSTGRES_USER") {
		envData.Open = strings.Replace(envData.Open, "$POSTGRES_USER", os.Getenv("POSTGRES_USER"), -1)
	}
	if strings.Contains(envData.Open, "$POSTGRES_PASSWORD") {
		envData.Open = strings.Replace(envData.Open, "$POSTGRES_PASSWORD", os.Getenv("POSTGRES_PASSWORD"), -1)
	}
	fmt.Println("[database] Using psql paramaters:", envData.Open)

	// Construct a dbConfig object from envData
	Database = dbConfig{
		Driver: envData.Driver,
		Open:   envData.Open,
	}
}
